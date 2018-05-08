package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const AverageRatesInterval = time.Minute * 10

type Storage interface {
	Save(interface{}) error
	Fetch(duration time.Duration) ([]interface{}, error)
}

type Route struct {
	Method string
	Path   string
}

type Handler struct {
	routeHandlers map[Route]http.HandlerFunc
}

func (h *Handler) setRoute(method string, path string, f http.HandlerFunc) error {
	route := Route{
		Method: method,
		Path:   path,
	}
	if h.routeHandlers[route] != nil {
		return errors.New(`api: Handler already set for route ` + route.Method + ` ` + route.Path)
	}
	h.routeHandlers[route] = f
	return nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := Route{
		Method: r.Method,
		Path:   r.URL.Path,
	}
	fmt.Fprintf(os.Stdout, "serve route:\n%#v\n", route)
	if handlerFunc, ok := h.routeHandlers[route]; ok {
		handlerFunc(w, r)
	} else {
		http.Error(w, `Bad request`, http.StatusBadRequest)
	}
}

type Api struct {
	handler *Handler
	server  *http.Server
	Storage Storage
	logger  *log.Logger
}

func NewApi(listenAddr string, apiPath string, storage Storage, logger *log.Logger) (*Api, error) {
	handler := Handler{
		routeHandlers: make(map[Route]http.HandlerFunc, 0),
	}
	api := Api{
		handler: &handler,
		server: &http.Server{
			Addr:    listenAddr,
			Handler: &handler,
		},
		Storage: storage,
		logger:  logger,
	}
	err := api.setRoute(`GET`, apiPath+`/avg`, api.GetAverage)
	if err != nil {
		return nil, err
	}
	return &api, nil
}

func (a *Api) setRoute(method string, path string, f http.HandlerFunc) error {
	return a.handler.setRoute(method, path, f)
}

func (a *Api) Start() error {
	return a.server.ListenAndServe()
}

func (a *Api) Stop() error {
	return a.server.Shutdown(context.Background())
}

func (a *Api) GetAverage(w http.ResponseWriter, req *http.Request) {
	result, err := a.Storage.Fetch(AverageRatesInterval)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	jsonResult, err := json.Marshal(result)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	n, err := w.Write(jsonResult)
	if err != nil {
		a.logger.Println(`api: Response write error`, err.Error())
	}
	l := len(jsonResult)
	if n != l {
		a.logger.Printf("api: Incorrect bytes lenght written, %#v expected, %#v written\n", l, n)
	}
}

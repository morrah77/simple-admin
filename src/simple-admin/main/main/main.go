package main

import (
	"flag"
	"log"
	"os"

	"simple-admin/api"
	"simple-admin/storage"

	"github.com/pkg/errors"
)

type Conf struct {
	ListenAddr string
	ApiPath    string
	Storage    string
	Dsn        string
}

var (
	conf   Conf
	logger *log.Logger
)

func init() {
	flag.StringVar(&conf.ListenAddr, `listen-addr`, `:8080`, `Address to listen`)
	flag.StringVar(&conf.ApiPath, `api-path`, `/simple-admin/v0`, `API path to handle rate requests`)
	flag.StringVar(&conf.Storage, `storage`, `postgres`, `Storage type`)
	flag.StringVar(&conf.Dsn, `dsn`, ``, `Storage DSN`)
	logger = log.New(os.Stdout, ``, log.Flags())
}

func main() {
	flag.Parse()
	stor, err := storage.NewStorage(conf.Storage, conf.Dsn, logger)
	if err != nil {
		panic(errors.Wrap(err, `Could not create storage!`))
	}
	defer func() {
		stor.Stop()
	}()

	apiInstance, err := api.NewApi(conf.ListenAddr, conf.ApiPath, stor, logger)
	if err != nil {
		panic(errors.Wrap(err, `Could not create API!`))
	}
	defer func() {
		apiInstance.Stop()
	}()
	err = apiInstance.Start()
	if err != nil {
		panic(errors.Wrap(err, `Could not start API!`))
	}
}

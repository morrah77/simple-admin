FROM golang:1.9.1

ENV GOPATH=/go:/go/src/github.com/morrah77/simple-admin
WORKDIR /go/src/github.com/morrah77/simple-admin
COPY . .

RUN go get -u github.com/golang/dep/cmd/dep

RUN cd src/simple-admin/ && dep ensure && cd ../..
RUN go install simple-admin/main
CMD ./bin/main --listen-addr :8080

#Simple microservices based data entities admin system example - backend part

Purposed to manage data entities (users && user subscriptions ones implemented) via REST API and/or web SPA

Contains several services:

 - database
 - backend allowing CRUD operations after  authentification
 - TODO consider to implement dedicated authentification service
 - frontend exposing web interface

###Environment variables

`export GOPATH=pwd && echo $GOPATH`

###Dependencies injection

`cd src/simple-admin/ && dep ensure && cd ../..`

###Build

  - locally `go install simple-admin/main`


  - with docker `docker build -f Dockerfile . -t simple-admin-backend`

###Run

  - locally
  
  `docker run --name=postgres --rm -d -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=changeit -e POSTGRES_DB=simple-admin -v simple-admin_pgdata:/var/lib/postgresql/data postgres:9.6.8`
  
  `docker inspect postgres | grep \"IPAddress\":` to know postgres container IP
  
  `bin/main --listen-addr :8080 --api-path /api/v0 --storage postgres --dsn "host=<IP_address> port=5432 user=postgres password=changeit dbname=simple-admin sslmode=disable"`
  
  please remember that you should have built frontend Docker image to run web SPA, i. e., `cd <your_frontend_project_path> && docker build -f Dockerfile . -t simple-admin-frontend`
  
  `docker run --name=frontend --rm -d -e BACKEND_URL=http://localhost:8080 -e PORT=8081 simple-admin-frontend`


  - with docker `docker-compose -f docker-compose.yml -p simple-admin up -d`

###Stop

  - locally `docker stop postgres:9.6`
  
  - with docker `docker-compose -f docker-compose.yml -p simple-admin down`
  
###Logs

 - with docker
 
 `docker logs -f simple-admin_simple-admin-service_1`
 
 `docker logs -f simple-admin_postgres_1`

###Test

####Unit tests for packages

`go test -race .src/simple-admin/api && go test -race .src/simple-admin/storage`

####Manually

`curl -iv -X GET http://localhost:8080/simple-admin/v0/avg`
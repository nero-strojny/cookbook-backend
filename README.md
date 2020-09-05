# Cookbook

## Running Locally:

- install go
- run `go get go.mongodb.org/mongo-driver/mongo`
- run `go get -u github.com/gorilla/mux`
- run `go get golang.org/x/crypto/bcrypt`
- add config file with the proper access to the mongodb
- run `go run main.go -DB_STRING "<DB_STRING>"`

## Testing
- run `go get github.com/stretchr/testify/assert`
- run `go test -v ./test -DB_STRING "<DB_STRING>" -ENV "dev"`

## Development:
- learned from [ToDo App Tutorial](https://levelup.gitconnected.com/build-a-todo-app-in-golang-mongodb-and-react-e1357b4690a6)

## Deployment:
- scp over the files: `scp -r -i <pem file here> <path to cookbook backend> <ec2 instance>`
- ssh into the ec2 instance `ssh -i "<pem file here>" <ec2 instance>`
- start docker if it is not already started: `sudo service docker start`
- build the docker container: `docker build -t docker-example .`
- run the docker container: `docker run --detach --publish 8080:8080 docker-example`
- follow the logs if needed: `docker logs follow <container id>`

## Authors
- Jake Strojny @jstrojny
- Alexandra Nero @alexandra-nero
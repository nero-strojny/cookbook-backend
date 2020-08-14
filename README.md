# Cookbook

## Running Locally:

- install go
- run `go get go.mongodb.org/mongo-driver/mongo`
- run `go get -u github.com/gorilla/mux`
- add config file with the proper access to the mongodb
- run `go run main.go`

## Development:
- learned from [ToDo App Tutorial](https://levelup.gitconnected.com/build-a-todo-app-in-golang-mongodb-and-react-e1357b4690a6)

## Deployment:
- ssh into the ec2 instance `ssh -i "<pem file here>" ec2-user@ec2-54-145-81-149.compute-1.amazonaws.com`
- start docker if it is not already started: `sudo service docker start`
- run the docker container: `docker run --detach --publish 8080:8080 docker-example`
- follow the logs if needed: `docker logs follow <container id>`

## Authors
- Jake Strojny @jstrojny
- Alexandra Nero @alexandra-nero
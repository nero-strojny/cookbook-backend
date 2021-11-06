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
- Infrastructure is fargate based, see cloud formation 

aws cloudformation create-stack --stack-name tastyboi-server-fargate --template-body file://cloudformation.yaml --capabilities "CAPABILITY_NAMED_IAM" --region <region> --parameters "ParameterKey=packageName,ParameterValue=tastyboi-server" "ParameterKey=version,ParameterValue=test-1" "ParameterKey=hostedZoneCertificateArn,ParameterValue=arn:aws:acm:<region>:<account id>:certificate/<certificate id>" "ParameterKey=maximumCount,ParameterValue=4" "ParameterKey=cpu,ParameterValue=1024" "ParameterKey=memory,ParameterValue=4096"

## Authors
- Jake Strojny @jstrojny
- Alexandra Nero @alexandra-nero


aws cloudformation create-stack --stack-name tastyboi-server-fargate --template-body file://cloudformation.yaml --capabilities "CAPABILITY_NAMED_IAM" --region us-east-1 --parameters "ParameterKey=packageName,ParameterValue=tastyboi-server" "ParameterKey=version,ParameterValue=test-1" "ParameterKey=hostedZoneCertificateArn,ParameterValue=arn:aws:acm:us-east-1:401221454868:certificate/ff10a345-1e6d-412c-9a1d-214c62adefe4" "ParameterKey=maximumCount,ParameterValue=4" "ParameterKey=cpu,ParameterValue=1024" "ParameterKey=memory,ParameterValue=4096"
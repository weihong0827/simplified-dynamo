gen-proto:
	protoc --proto_path=proto proto/*.proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative

webc:
	go run ./webclient

first-node:
	go run ./server --addr="127.0.0.1:50051" --webclient="http://127.0.0.1:8080/addNode?port=50051"
second-node:
	go run ./server --addr="127.0.0.1:50052" --webclient="http://127.0.0.1:8080/addNode?port=50052"
third-node:
	go run ./server --addr="127.0.0.1:50053" --webclient="http://127.0.0.1:8080/addNode?port=50053"

fourth-node:
	go run ./server --addr="127.0.0.1:50054" --webclient="http://127.0.0.1:8080/addNode?port=50054" 

fifth-node:
	go run ./server --addr="127.0.0.1:50055" --webclient="http://127.0.0.1:8080/addNode?port=50055" 

docker-build:
	docker build -t node .

run-second-node:docker-build
	docker rm second-node
	docker run --network host --name second-node node ./bin/server --addr="127.0.0.1:50052" --webclient="http://127.0.0.1:8080/addNode?port=50052"


run-sixth-node:docker-build
	docker rm sixth-node
	docker run --network host --name sixth-node node ./bin/server --addr="127.0.0.1:50056" --webclient="http://webclient:8080/addNode?port=50056"

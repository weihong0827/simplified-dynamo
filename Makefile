gen-proto:
	protoc --proto_path=proto proto/*.proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative

webclient:
	go run ./webclient

first-node:
	go run ./server --addr="127.0.0.1:50051" --webclient="http://127.0.0.1:8080/addNode?port=50051"
second-node:
	go run ./server --addr="127.0.0.1:50052" --webclient="http://127.0.0.1:8080/addNode?port=50052"
third-node:
	go run ./server --addr="127.0.0.1:50053" --webclient="http://127.0.0.1:8080/addNode?port=50053"
gen-proto:
	protoc --proto_path=proto proto/*.proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative

seed:
	go run ./server --addr="127.0.0.1:50051"
second-node:
	go run ./server --addr="127.0.0.1:50052" --seed_addr="127.0.0.1:50051"
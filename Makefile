proto:
	protoc -I=./api/dolorosa --go_out=. --go-grpc_out=. online-control.proto
	protoc -I=./api/kafka --go_out=. --go-grpc_out=. log.proto

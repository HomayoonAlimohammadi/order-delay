.PHONY: gen-proto
gen-proto:
	protoc -I=./proto --go_out=./proto --go_opt=paths=source_relative \
	 --go-grpc_out=./proto --go-grpc_opt=paths=source_relative order.proto

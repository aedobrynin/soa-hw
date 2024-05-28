grpc-generate:
	protoc \
	-I protos \
	--go_out=./internal/grpcadapter/gen/ \
	--go_opt=paths=source_relative \
	--go-grpc_out=./internal/grpcadapter/gen/ \
	--go-grpc_opt=paths=source_relative \
	--experimental_allow_proto3_optional \
	protos/statistics.proto

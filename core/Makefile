oapi-generate:
	oapi-codegen --config ./codegen/core.yaml ./api/core.yaml

posts-client-generate:
	protoc \
	-I ../posts/protos/ \
	--go_out=./internal/clients/postsclient/gen/ \
	--go_opt=paths=source_relative \
	--go-grpc_out=./internal/clients/postsclient/gen/ \
	--go-grpc_opt=paths=source_relative \
	../posts/protos/posts.proto

statistics-client-generate:
	protoc \
	-I ../statistics/protos/ \
	--go_out=./internal/clients/statisticsclient/gen/ \
	--go_opt=paths=source_relative \
	--go-grpc_out=./internal/clients/statisticsclient/gen/ \
	--go-grpc_opt=paths=source_relative \
	--experimental_allow_proto3_optional \
	../statistics/protos/statistics.proto

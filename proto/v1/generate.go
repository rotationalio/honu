package pb

//go:generate protoc -I . --go_out=. --go_opt=module=github.com/rotationalio/honu/proto/v1 --go-grpc_out=. --go-grpc_opt=module=github.com/rotationalio/honu/proto/v1 honu.proto

package pb

//go:generate protoc -I . --go_out=.. --go_opt=module=github.com/rotationalio/honu --go-grpc_out=.. --go-grpc_opt=module=github.com/rotationalio/honu object/v1/object.proto replica/v1/replica.proto

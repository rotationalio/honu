package pb

//go:generate protoc -I . --go_out=.. --go_opt=module=github.com/rotationalio/honu --go-grpc_out=.. --go-grpc_opt=module=github.com/rotationalio/honu object/v1/object.proto replica/v1/replica.proto pagination/v1/pagination.proto peers/v1/peers.proto honu/v1/honu.proto


## Use

- install plug
```go
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

```

- clone this repo
```bash
cd pb && git clone https://github.com/googleapis/api-common-protos.git
```

- tree
```text
├── README.md
├── client
├── go.mod
├── go.sum
├── pb
│   ├── api-common-protos
│   │    ├── google
│   │    |    ├── api
│   │    |    ├── ...
│   │    │    │   ├── annotations.proto
│   │    │    │   ├── http.proto
│   │    │    │   ....
│   ├── demo.pb.go
│   ├── demo.pb.gw.go
│   ├── demo.proto
│   └── demo_grpc.pb.go
└── server
    └── main.go
```

- generate proto
```bash
cd pb && protoc -I./api-common-protos -I. --go_out=. --go-grpc_out=. --grpc-gateway_out=. demo.proto
```
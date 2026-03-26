module tribecart/api-gateway

go 1.24.6

replace github.com/tribecart/proto => ../../proto

require (
	github.com/gorilla/mux v1.8.1
	github.com/tribecart/proto v0.0.0-00010101000000-000000000000
	github.com/rs/cors v1.11.1
	google.golang.org/grpc v1.75.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/golang/protobuf v1.5.4 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
)
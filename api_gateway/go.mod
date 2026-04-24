module elex_storage/api_gateway

go 1.24.0

require (
	elex_storage/pkg v0.0.0
	github.com/joho/godotenv v1.5.1
	github.com/julienschmidt/httprouter v1.3.0
	go.uber.org/fx v1.24.0
	google.golang.org/grpc v1.79.3
)

replace elex_storage/pkg => ../pkg

require (
	github.com/docker/distribution v2.8.3+incompatible // indirect
	github.com/jinzhu/copier v0.4.0 // indirect
	go.uber.org/dig v1.19.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.1 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)

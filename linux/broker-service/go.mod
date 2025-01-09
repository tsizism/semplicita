module broker

go 1.23.3

replace github.com/tsizism/semplicita/linux/shared => ./../shared

require (
	github.com/rs/cors v1.11.1
	github.com/tsizism/semplicita/linux/shared v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.69.2
	google.golang.org/protobuf v1.36.1
)

require (
	github.com/rabbitmq/amqp091-go v1.10.0
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241015192408-796eee8c2d53 // indirect
)

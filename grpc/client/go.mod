module github.com/lucaspere/grpc/client

go 1.20

require (
	github.com/lucaspere/grpc/service v0.0.0
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.5.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/grpc v1.37.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

replace github.com/lucaspere/grpc/service => ../service

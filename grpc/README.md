# Annotations

## Chapter 8 - Building RPC Application using gRPC

RPC (Remote Procedure Call) is a style of service communication where one calls a function from another service over a network like calling in the same process.
in Go, you can create a RPC service using [net/rpc](https://pkg.go.dev/net/rpc) and (net/rpc/jsonpc)[https://pkg.go.dev/net/rpc/jsonrpc] native packages.
RPC frameworks have to deal with two things: how the function call gets converted to a network request and how it gets transmitted.

### gRPC in Golang

gRPC (Google Remote Procedure Call) is a _universal_ framework that allows services to communicate with others using HTTP2. It is universal because you can use it in different languages. For example, you can create a gRPC server in Go using [grpc-go](https://github.com/grpc/grpc-go) and a client in Node.js using [grpc-js](https://github.com/grpc/grpc-node/tree/master). As it uses the HTTP2 protocol to communicate, the message transmitted over the network must be in binary format as required by HTTP2. To allow you to Encode and Decode your data, you can use the Protubuf IDL to define your service and to use it as data serialization.

### Protobuf

Has two main uses: as Inteface Description Language (IDL) and as Data Serialization:

- IDL: to create a gRPC service, the first step is to define its _contract_ using Protobuf specification. With this, you can structure the service schema of data to be serialized or deserialized. Once you have the message interface, you can use a code generator ([Go](https://pkg.go.dev/github.com/golang/protobuf/protoc-gen-go), [Nodejs](https://github.com/improbable-eng/ts-protoc-gen)) to create language-specific code for you;
- Data Serialization: it is a binary format that efficiently encode and decode structured data in a compact and platform-independent way. It is compact because the protobuf decrease the size of message transmitted over a network, and it is platform-independent because it binary data is a machine format, thus all computer or platform could understand it.

#### IDL

As other languages, the protobuf has rules that must respect to be able to compile. To create a **_service_**, use `service` keyword following the name of it in _PascalCase_; for a **_method_** use the `rpc` keywork following the name with its **_inputs_** and the **_output_** with the same style of service. The input and output must be a `message`, each **_message_** has one or more **_fields_** which have a **_type_**, **_name_**, and a **_number_**. Example a protobuf service:

```proto
syntax = "proto3";

service Cal {
    rpc Add(AddRequest) returns int32 {}
}

message AddRequest {
    int32 num_1 = 1;
    int32 num_2 = 2;
}
```

Once you have your service structured, you can use the compile _protoc_ with one generator plugin to the language you want. In summary, Protobuf is a powerful and efficient data serialization format that is well-suited for use of distributed system and APIs that need to be fast, compact, and plataform-independent.

### gRPC

With your generated codes, you can create your client or server using the grpc-go package. This package is same like http package: create a server over a _transport protocol_ normally the _TCP_. The gRPC does not enforce you to implement all the methods defined in the protobuf file, so you can implement only what expect to use it.

#### Testing

To test a gRPC, you can use the [bufconn](https://pkg.go.dev/google.golang.org/grpc/test/bufconn) package. It simulates a gRPC network running in-memory. So, it makes easy to test without the complexity of managing a real network service.

#### Json Serialization

To intoperability with _json_ data format, the grpc-go package offers [protojson](https://pkg.go.dev/google.golang.org/protobuf/encoding/protojson) package for that. This packages allows you to unsmarshall Protobuf struct to valid json format and vice-versa.

#### Error Handling

The reliability principle is one the most important in API Design because it offers the consumer services to know what is going on with the server. The grpc-go has the [status](https://pkg.go.dev/google.golang.org/grpc/status) and [code](https://pkg.go.dev/google.golang.org/grpc/codes) to allow you a response with a property error status.

### Summary

1. gRPC is a _frameworks_ created by Google in 2015 that use the HTTP2 communication protocol and protobuf as the data serialization format.
2. The first step to create a gRPC service is define a schema using the Protobuf IDL specification.
3. Once you have the schema, you can use the compile with the property language plugin to generate the codes.
4. gRPC do not enforces you to implement all the services and methods defined in Protobuf schema.
5. grpc-go package offers bufconn and status to test a grpc service and handler errors.
6. It also offers the protojson package to serialization json format and vice-versa.

### References

[gRPC Website](https://grpc.io/)
[Protobuf Website](https://protobuf.dev/)
[gRPC package](https://pkg.go.dev/google.golang.org/grpc)
[gRPC status_codes](https://developers.google.com/maps-booking/reference/grpc-api/status_codes?hl=pt-br)
[]

# Distributed Services with Go

## Part. Get Started

### Chapter 2 - Let's Go

This first chapter was an introduction to set up the initial project that we gonna use throughout the book.
We learned: How to build a simple JSON/HTTP service that accepts and responds with JSON and stores the request's data in in-memory.

#### References

1. https://github.com/golang/go/wiki/Modules

### Chapter 2 - Structure Data with Protocol Buffers (Protobuf)

#### Characteristics

###### _Consistent schemas_

As it's an IDL, you must respect the semantics and grammar of the language. With that, all your services remain consistent with the data model throughout your whole system.

###### _Versioning for free_

Each field in a message has one number to the compile manager for the versioning. You can add/remove fields and the compile will check and will warn the services the changes;

###### _Less boilerPlate_

The Protobuf libraries handle encoding and decoding your messages to your structures.

###### _Extensibility_

You can add plugins to extend your Protobuf functionalities. For example, the [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)(https://github.com/grpc-ecosystem/grpc-gateway) creates a REST service in Go according to your Protobuf description.

###### _Language agnosticism_

The Protobuf compiler translates the specification into different languages.

###### _Performance_

Protobuf is highly performant, and has smaller payloads and serializes up to six times faster than JSON.

#### Steps to use Protobuf

1. Install the compiler: you have to go the [repository](https://github.com/protocolbuffers/protobuf) and follows the installation instructions.
2. Define the domain Types: the IDL has version 2 and 3. Go to the [site](https://protobuf.dev/programming-guides/proto3/) and define your service domain accordding the language specification.
3. Compile your Protobuf file to your languange using the [plugins](https://protobuf.dev/reference/)
4. Implement your service using the generated codes made by the compiler.

### Chapter 3 - Write a Log package

#### What can I learned from building a log myself?

1. Solve problems using logs and discover how they can make hard problems easier.
2. Change existing log-based systems to fit your needs and build your own log-based systems.
3. Write and read data efficiently when building storage engines.
4. protect against data loss caused by system failures.
5. Encode data to persist it to a disk or write your own wire protocols and send data between applications.

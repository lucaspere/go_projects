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

#### The powerful of Log

Log systems has several utilities and implication in diferents softwares, as: Database enginers, Filesystems, and States managements.

- Database Enginers: this enginers uses a log functionality to make the writen of data durable. For example, when a new commit for write some data coming, the enginer first write this to a called _write-ahead log_ (WAL), and later process the WAL to aplly the changes to their database's data files. It also uses this tool to allow replication: instead to use the data in files, it uses the WAL to send the data to their replicas.
- Filesystems: the _ext_ filesystem uses a _journal_ algorithm to log changes instead of directly changing the disk's data file. Once the filesystem has safely written the changes to the journal, it then applies those changes to the data file.
- State Managements: Front-end libraries like Redux uses a type o log to track the UI changes and allowing the user to undo the state to previous one. This allow more UI interation and more control for the users.

#### How Logs Work

A **log** is normally saved as a file of **record**. A log is an append-only sequence of records, where are is appended to the end of the file. When a record is append to a log, it assigns the record a unique and sequential **offset** number that acts like the ID. The user not have disks with infinite space, which means that the log can't append to the same file forever, so it split the log into a list of **segments**. When the log grows too big, it free up disk space by deleting old segments.

There're one special segment called **active segment**. This kind is segment where the log is currently using to append records. Each segment comprises a **store file** and a **index file**. The store file is where we store the record data, and the index file is where the index for each record in the store file is saved. The index is used to map record offsets to their position to speed up the process of reading the records. Index files are small enought that we can memory-map them and make operations on the file as fast as operating on in-memory data.

##### log system's terms:

- Record the data stored in our log.
- Store the file we store records in.
- Index the file we store index entries in.
- Segment the abstraction that ties a store and an index together.
- Log the abstraction that ties all the segments together.

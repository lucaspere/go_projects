# Building Distributed Application in Gin

Distributed Applications are the new style of software development in these days.
My goals to read this book are the following:

- How to build a RESTful application using GIN and how to structure it.
- Learning the MongoDB drive in Go;
- How to scale a Web Application using Docker;
- How to deploy a Web Application to the cloud; and
- How to set up a CI/CD for a Web Application.

## Section 1 - Inside the Gin Framework

Gin is a WEB framework with a focus on creating a fast Restful Go application. It offers a set of tools and features that enables developers to build a high-performance web application like Routing, Middlewares, Error Handling, and JSON serialization.

### Chapter 1 - Getting Started with Gin

Gin is one of the fastest web frameworks in the Go language. To create and listen for HTTP requests, just use the `gin.Default()` method that returns an `*Engine` pointer to use as HTTP routing and middleware.

### References

[Gin Framework](https://github.com/gin-gonic/gin)
[What is HTTP Middleware?](https://www.moesif.com/blog/engineering/middleware/What-Is-HTTP-Middleware/)
[Routing in Go](https://benhoyt.com/writings/go-routing/)

## Section 2 - Distributed Microservices

Distributed microservices is an architectural style for building applications that consists of multiple smaller services, each running independently and communicating with each other over a network. Each microservice typically has its API and Database; and communicates with other microservices using protocols such as REST or messaging systems like RabbitMQ.

In this section, we gonna build our Gin microservices with their API using Rest and Database with MongoDB. And we'll integrate them using RabbitMQ as the message system.

### Chapter 2 - Setting Up API Endpoints

We set up our API interface using REST style with Gin because REST protocol is a good fit for **_north-south_** traffic communication. This traffic is when a _consumer_ of the API is out of API's boundary, that is, is out-of-process.

The steps to create a RESTful service are the following:

1. Define the structure of the Model

```go
type Recipe struct {
	ID           primitive.ObjectID `json:"id"`
	Name         string             `json:"name"`
	Tags         []string           `json:"tags"`
	Ingredients  []string           `json:"ingredients"`
	Instructions []string           `json:"instructions"`
	PublishedAt  time.Time          `json:"publishedAt"`
}
```

2. Define the route path name based on the model name. This is the famous _CRUD_ operation and the _resource_ is the name of the model structure.

| HTTP Method | Resource      | Description                       |
| ----------- | ------------- | --------------------------------- |
| GET         | /recipes      | Returns a list of recipes         |
| GET         | /recipes/{id} | Returns a single recipe by its id |
| POST        | /recipes      | Creates a new recipe              |
| PUT         | /recipes{id}  | Updates an existing recipe        |
| DELETE      | /recipes{id}  | Deletes an existing recipe        |

3. Implements the API using Gin framework routing.
   A _routing_ library is just a simple _map table_ where each key is a resource path and the value is an **_HTTP Handler_**.
   ```go
   router.POST("/recipes", recipientHandlers.NewRecipeHandler)
   router.GET("/recipes", recipientHandlers.ListRecipesHandler)
   router.PUT("/recipes/:id", recipientHandlers.UpdateRecipeHandler)
   router.DELETE("/recipes/:id", recipientHandlers.DeleteRecipeHandler)
   router.GET("/recipes/:id", recipientHandlers.GetOneRecipeHandler)
   ```
4. Documents the Rest API using the OpenAPI specification.
   The OpenAPI is a set of patterns that facilitates the documentation of API and allows consumers to understand the interface of the API. We can use the [gin-swagger](https://github.com/swaggo/gin-swagger) tool to generate the OpenAPI specs based on the code implementation.

### Chapter 3 - Managing Data Persistent with MongoDB

#### MongoDB

With the API interface and documentation, the user can now interact with the API to integrate with her application. However, all data generated by the interaction of the user will be lost when the server goes off because all data is saved in memory. For that reason, we need to persist the data with persistent service. We gonna use MongoDB for that purpose.

MongoDB is a _NoSQL_ database that uses the _document_ paradigm to structure and organize the data. The MongoDB Driver saves a document in _BSON_ format which is a binary representation of a MongoDB document. To use the driver methods (Insert, Update, Find), we must unmarshal Go struct types into BSON and marshall it to Go struct

```go
// To create a new client instance.
client, err := mongo.Connect(context.Background(), options.CLient().ApplyURI(URI))

// Verify if the server is online.
if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
	log.Fatal(err)
}

// Connect with a collection.
collection := client.Database(DATABASE).Collection("recipes")
```

#### Redis

Redis is also another _NoSQL_ but it uses the _key/value_ type. This type is like a hash table where the key is the identifier of the value. As the hash table has the _O(1)_ complexity, it is very fast in querying the data. For that reason, Redis is a good choice to use as the Cache system.

```go
// Creates a new Redis Client instance.
redisClient = redis.NewClient(&redis.Options{
	Addr:     ADDR,
	Password: PASSWORD,
	DB:       0,
})

// Verify the status of client instance.
fmt.Println(redisClient.Ping())
```

### Chapter 4 - Building API Authentication

Authentication is a core feature that all applications should have if the security concern is important. Authentication allows you to know who is the user interaction with your app. This is important because it offers to you create actions for a specific user like metrics and authorization, which is limit the user actions. There are several types of Authentication, one of them being API-KEY, JSON Web Token (JWT), and thirty-party OAuth2.0.

- API-KEY: the most basic authentication feature. It consists in sending a base64 value at _Header_ request (X-API-KEY or Authorization). The advantage of this approach is the simplest way to authenticate. The disadvantage is that it can authenticate the project, but not the user. Also, anyone with a key can use the app.
- JWT: most used in the Web context. This type allows you to identify who is interacting with your app and allows you to set an expiration time for that token. With that, you can have more control over the user's interaction with your app. The advantage is that is easy to set up this type of authentication because it is a Web standard by RFC and has several thirty-party libraries in different languages.
- OAuth 2.0: a SaaS, like Firebase Authentication. It offers you to delegate all the responsibility of managing and maintaining an authentication service to it. The advantage is the reducing the cost of maintaining and managing a dedicated auth service.

#### API-KEY

To implements an API-KEY type, you can generate a base64 value using the _OpenSSL_ library with the command `openssl rand -base64 16`. The user will send this value in a Header request and, for each router, you verify the Header value.

```go
 if c.GetHeader("X-API-KEY") != os.Getenv("X_API_KEY") {
       c.JSON(http.StatusUnauthorized, gin.H{
          "error": "API key not provided or invalid"})
       return
   }
```

#### JWT

JSON Web Token is a secure way to transmit information between parties as a JSON Object because all its information is digitally signed with a secret, so it can be verified and trusted. A JWT Token has three parts separated by dots:

- HEADER: H5256 Algorithm, type "JWT";
- PAYLOAD: the payload;
- SIGNATURE: HMACSHA256 (Base64(header) + "." + Base64(Payload), SECRET)

Example: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Imx1Y2FzIiwiZXhwIjoxNjgyOTkyOTUyfQ.4m6QpjcM5D3aYTC5QmT7Dt9TE_dcwOjLPF6C-RGxkes`

To implement a JWT authentication, you can follow these steps:

1. Install the [jwt](https://github.com/golang-jwt/jwt) package.
2. Using the library, generate the token with your _secret_

```go
et := time.Now().Add(10 * time.Minute)
claims := &Claims{
	Username: user.Username,
	StandardClaims: jwt.StandardClaims{
		ExpiresAt: et.Unix(),
	},
}
t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
ts, err := t.SignedString([]byte(os.Getenv("JWT_SECRET")))
```

3. Create a logic to validate the Token receive in the request

```go
tokenValue := c.GetHeader("Authorization")
claims := &Claims{}
tkn, err := jwt.ParseWithClaims(
	tokenValue, claims,
	func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
```

### Chapter 6: Scaling a Gin Application

Scale an application means increasing its capability to handle more user interactions. It can be done by increasing the hardware capability (RAM, CPU, Disk) or by creating clones (different instances of the application). The first one is called _Vertical Scaling_, the second is _Horizontal Scaling_. In this chapter, we'll scale horizontally with **Docker** replicas and scaling workloads with **RabbitMQ** message broker. Also, we use the **nginx** to serve as a _reverse_ _proxy_ to our replicas.

#### Scaling Workloads using RabbitMQ

When a Monolithic architecture is migrated to microservices architecture where each service is responsible for handler one single business logic, and the way each service is communicating with other changes as well. In monolithic architecture, the communication is in the process, where the SO orchestrates it. In microservices, this type up to the network, with a framework doing the job, like HTTP Request, Message System, RPC, GraphQL and others.

RabbitMQ is a Message System framework focused on Resilient and Safaly message delivery used by many large companies. It acts as Message Broker in the system and offers different event patterns like Pub/Sub and FanIn/FanOut.

For this application, we gonna use the Pub/Sub pattern. This pattern is like a magazine subscription: people assign to receive a magazine paper according to an event, for example, a period (weekly). All the subscribers will receive a clone of the magazine edition, if one person subscribes after an edition, he will receive the edition followers, never the old. In this example, the Magazine is a Publisher and the People is the Subscriber. The RabbitMQ is the delivery of a message.

##### Steps to create a Pub/Sub in RabbitMQ

First, you create a TCP connection with the RabbitMQ URI;

```go
conn, err := amqp.Dial(os.GetEnv("RABBITMQ_URI"))
if err != nil {
	log.Fatal(err)
}
```

2. Then you open a server channel to process the messages in the connection

```go
channel, err := conn.Channel()
if err != nil {
	log.Fatal(err)
}
```

3. With the channel ready to process data, you can publish or consume the messages
   3.1. The interface to Publish a message is:

   ```go
    Publish(exchange, key string, mandatory, immediate bool, msg Publishing) error
   ```

   the `msg` parameter is a `Publishing` struct and it has a `Body []byte` field.
   With that, you can deliver your messages, as bytes, to your subscribers.

   3.2. The interface to Consume messages is:

   ```go
   Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args Table) (<-chan Delivery, error)
   ```

   Observe that the messages are delivered via a `chan` type. With this, we can structure the code to keep listening for messages published using the `for...loop`

   ```go
   for msg := range messages {
      fmt.Println(msg)
   }
   ```

#### Scaling Horizontally using Docker

Scaling an application is about giving more power to the application to be able to handle more user interaction. Docker and Docker Compose allow you to scale your application quickly and with less effort. The command `docker compose --scale app=5` will create 5 `app`'s containers sharing the same host.

#### Add a Load Balance using Nginx

With the app running with multiple clones, they are splitting the consumer workloads between them. However, when a clone is busy handling a request, and others are idle, how can you know where to send more requests? To resolve that problem, you can use a strategy called **_Load Balancing_** which is a _reverse proxy_ to sit between your client and the instances of your app.

The most famous Load Balancing is **Nginx**, but there is another one that is gaining prominence called _[traefik](https://github.com/traefik/traefik)_. To use the Nginx as load balacing, first you must use the Nginx docker image and create a `nginx.conf` file at the root of your application source code.

### References

[Microservices and Distributed Systems](https://cleancommit.io/blog/are-microservices-distributed-systems/)
[REST Paper](https://www.ics.uci.edu/~fielding/pubs/dissertation/rest_arch_style.htm)
[OpenAPI Specification](https://swagger.io/specification/)
[NoSQL Databases](https://www.ibm.com/topics/nosql-databases)
[MongoDB BSON](https://www.mongodb.com/docs/drivers/go/current/fundamentals/bson/)
[Redis Caching](https://redis.io/docs/manual/client-side-caching/)
[Why is API Key?](https://cloud.google.com/endpoints/docs/openapi/when-why-api-key?hl=pt-br)
[JWT RFC](https://datatracker.ietf.org/doc/html/rfc7518)
[export local environment online](https://ngrok.com/)
[Comparing the different Scaling](https://www.section.io/blog/scaling-horizontally-vs-vertically/)
[OAuth 2.0](https://docs.sensedia.com/pt/api-platform-guide/4.3.x.x/other-info/oauth20.html)
[RabbitMQ explained](https://www.javaguides.net/2018/12/how-rabbitmq-works-and-rabbitmq-core-concepts.html)
[What is Pub-Sub Pattern](https://learn.microsoft.com/en-us/azure/architecture/patterns/publisher-subscriber)
[What is Docker?](https://learn.microsoft.com/en-us/dotnet/architecture/microservices/container-docker-introduction/docker-defined)
[Nginx Reverse Proxy](https://medium.com/globant/understanding-nginx-as-a-reverse-proxy-564f76e856b2)
[Load Balancing](https://samwho.dev/load-balancing/)

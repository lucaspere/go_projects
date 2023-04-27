package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/lucaspere/go_projects/recipes-api/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var ctx context.Context
var err error
var client *mongo.Client
var collection *mongo.Collection
var redisClient *redis.Client

func init() {
	ctx = context.Background()
	client, err = mongo.Connect(ctx,
		options.Client().ApplyURI(os.Getenv("MONGO_URI")),
	)
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	status := redisClient.Ping()
	fmt.Println(status)
	log.Println("Connected to MongoDB")

	collection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
}

func main() {
	router := gin.Default()
	recipientHandlers := handlers.NewRecipesHandler(ctx, collection, redisClient)
	router.POST("/recipes", recipientHandlers.NewRecipeHandler)
	router.GET("/recipes", recipientHandlers.ListRecipesHandler)
	router.PUT("/recipes/:id", recipientHandlers.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", recipientHandlers.DeleteRecipeHandler)

	router.Run()
}

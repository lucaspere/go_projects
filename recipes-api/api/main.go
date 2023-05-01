package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/lucaspere/go_projects/recipes-api/api/handlers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var ctx context.Context
var err error
var client *mongo.Client
var collectionUsers *mongo.Collection
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

	collectionUsers = client.Database(os.Getenv("MONGO_DATABASE")).Collection("users")
	collection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
}

func main() {
	router := gin.Default()
	recipientHandlers := handlers.NewRecipesHandler(ctx, collection, redisClient)
	usersHandlers := handlers.NewAuthHandler(ctx, collectionUsers)

	router.POST("/recipes", usersHandlers.AuthMiddleware(), recipientHandlers.NewRecipeHandler)
	router.GET("/recipes", recipientHandlers.ListRecipesHandler)
	router.PUT("/recipes/:id", usersHandlers.AuthMiddleware(), recipientHandlers.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", usersHandlers.AuthMiddleware(), recipientHandlers.DeleteRecipeHandler)
	router.GET("/recipes/:id", recipientHandlers.GetOneRecipeHandler)
	router.GET("/dashboard", DashboardHandler)

	router.POST("/signin", usersHandlers.SignInHandler)
	router.POST("/signup", usersHandlers.SignUpHanlder)
	router.POST("/refresh", usersHandlers.RefreshHandler)

	if enviroment := os.Getenv("GIN_ENV"); enviroment == "production" {
		router.RunTLS(":8080", "./certs/localhost.crt", "./certs/localhost.key")
	} else {
		router.Run(":8000")
	}
}

type Recipe struct {
	Title     string `json:"title" bson:"title"`
	Thumbnail string `json:"thumbnail" bson:"thumbnail"`
	URL       string `json:"url" bson:"url"`
}

func DashboardHandler(c *gin.Context) {
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(ctx)
	recipes := make([]Recipe, 0)
	for cur.Next(ctx) {
		var recipe Recipe
		cur.Decode(&recipe)
		recipes = append(recipes, recipe)
	}

	c.JSON(http.StatusOK, gin.H{
		"recipes": recipes,
	})
}

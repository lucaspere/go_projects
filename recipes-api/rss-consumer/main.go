package main

import (
	"context"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Feed struct {
	Entries []Entry `xml:"entry"`
}
type Entry struct {
	Link struct {
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Thumbnail struct {
		URL string `xml:"url,attr"`
	} `xml:"thumbnail"`
	Title string `xml:"title"`
}

var channelAmqp *amqp.Channel
var client *mongo.Client
var ctx context.Context

func init() {
	ctx = context.Background()
	client, _ = mongo.Connect(ctx,
		options.Client().ApplyURI(os.Getenv("MONGO_URI")),
	)
	amqpConnection, err := amqp.Dial(os.Getenv("RABBITMQ_URI"))
	if err != nil {
		log.Fatal(err)
	}

	channelAmqp, _ = amqpConnection.Channel()
}

func GetFeedEntries(url string) ([]Entry, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	byteValue, _ := io.ReadAll(resp.Body)
	var feed Feed
	xml.Unmarshal(byteValue, &feed)

	return feed.Entries, nil
}
func main() {
	c := make(chan os.Signal, 1)
	go ConsumeAMQ(c)
	s := <-c
	log.Printf("Got signal: %v\n", s)

}

type Request struct {
	URL string `json:"url"`
}

func ConsumeAMQ(c chan<- os.Signal) {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URI"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	channel, _ := conn.Channel()
	defer channel.Close()

	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	msgs, err := channel.Consume(
		os.Getenv("RABBITMQ_QUEUE"),
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	for msg := range msgs {
		log.Printf("Received a message: %s", msg.Body)
		entries, err := GetFeedEntries(string(msg.Body))
		if err != nil {
			log.Println(err)
			continue
		}
		go SaveEntries(&entries, collection)
	}
}

func SaveEntries(entries *[]Entry, c *mongo.Collection) error {
	log.Println("Saving bson data")
	for _, entry := range (*entries)[2:] {
		data := bson.M{
			"title":     entry.Title,
			"thumbnail": entry.Thumbnail.URL,
			"url":       entry.Link.Href,
		}
		c.InsertOne(ctx, data)
	}

	return nil
}

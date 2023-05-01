package main

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/streadway/amqp"
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

func init() {
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
	sm := http.NewServeMux()
	sm.HandleFunc("/parse", ParserHandler)

	log.Fatal(http.ListenAndServe(":5000", sm))
}

type Request struct {
	URL string `json:"url"`
}

func ParserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var request Request
	payload, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(payload, &request)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	go func() {
		err = channelAmqp.Publish(
			"",
			os.Getenv("RABBITMQ_QUEUE"),
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        []byte(request.URL),
			},
		)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()
	w.WriteHeader(http.StatusOK)
}

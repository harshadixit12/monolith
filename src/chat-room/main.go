package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
)

type App struct {
	DB *dynamodb.Client
}

type Client struct {
	ID   string
	Chan chan Message
}

type Broker struct {
	mu      sync.RWMutex
	clients map[string]*Client
}

func NewBroker() *Broker {
	return &Broker{
		clients: make(map[string]*Client),
	}
}

func (b *Broker) AddClient(id string) *Client {
	client := &Client{
		ID:   id,
		Chan: make(chan Message, 10),
	}

	b.mu.Lock()
	b.clients[id] = client
	b.mu.Unlock()

	return client
}

func (b *Broker) Remove(id string) {
	b.mu.Lock()
	if client, ok := b.clients[id]; ok {
		close(client.Chan)
		delete(b.clients, id)
	}
	b.mu.Unlock()
}

func (b *Broker) Broadcast(message Message) {
	b.mu.RLock()
	for _, client := range b.clients {
		select {
		case client.Chan <- message:
		default:
		}
	}
	b.mu.RUnlock()
}

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "streaming not supported :(", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	id := uuid.New().String()
	client := b.AddClient(id)
	defer b.Remove(id)

	ctx := r.Context()

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-client.Chan:
			fmt.Println("flusing...", msg)
			if !ok {
				return
			}

			fmt.Fprintf(w, "data: Message: %s, sent by %s at %s\n\n", msg.Content, msg.Sender, msg.CreatedAt)
			flusher.Flush()
		}
	}
}

const defaultChatroom string = "common"

func subscribe(w http.ResponseWriter, r *http.Request) {
	// Create a connection, add it to manager, return a connection ID, store to DB
	username := r.Header.Get("username")
	if username == "" {
		w.WriteHeader(500)
		return
	}

}

type Message struct {
	ID        string
	Sender    string
	Chatroom  string
	Content   string
	CreatedAt string
}

func wrappedSendHandler(app *App) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		body, err := io.ReadAll(r.Body)
		message := &Message{
			ID:        uuid.New().String(),
			Sender:    r.Header.Get("sender"),
			Content:   string(body),
			Chatroom:  defaultChatroom,
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
		}

		av, err := attributevalue.MarshalMap(message)

		if err != nil {
			http.Error(w, "marshal error", 500)
			return
		}

		_, err = app.DB.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String("Messages"),
			Item:      av,
		})
		if err != nil {
			http.Error(w, "put item error", 500)
			fmt.Println(err)
			return
		}
		broker.Broadcast(*message)
		fmt.Println(message)
		fmt.Fprintln(w, "Message stored")
	})
}

var broker = NewBroker()

func main() {
	customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           "http://localhost:8000",
			SigningRegion: "us-west-2",
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithEndpointResolver(customResolver), config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
		Value: aws.Credentials{
			AccessKeyID: "abcd", SecretAccessKey: "a1b2c3", SessionToken: "",
			Source: "Mock credentials used above for local instance",
		},
	}))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	db := dynamodb.NewFromConfig(cfg)
	app := &App{DB: db}

	http.HandleFunc("POST /messages", wrappedSendHandler(app))
	http.Handle("GET /listen", broker)

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}
}

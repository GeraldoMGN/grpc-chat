package main

import (
	"context"
	"fmt"
	"io"
	"math"
	"strings"
	"time"

	"github.com/thesayfulla/grpc-chat-example/chat"
	"google.golang.org/grpc"
)

var (
	requestStartTime time.Time
	requestTimes     []time.Duration
)

func handleInput(stream chat.Chat_ChatClient, clientName string) {
	clientRequest := "message"
	clientRequest = strings.TrimSpace(clientRequest)

	requestStartTime = time.Now()

	err := stream.Send(&chat.ChatMessage{
		User:    clientName,
		Message: clientRequest + "\n",
	})
	if err != nil {
		panic(err)
	}
}

func handleResponse(stream chat.Chat_ChatClient) {
	_, err := stream.Recv()
	elapsed := time.Since(requestStartTime)
	requestTimes = append(requestTimes, elapsed)

	if err == io.EOF {
		return
	} else if err != nil {
		panic(err)
	}
	/*
		msg := data.GetMessage()
		println((msg))
	*/
}

func calculateMetrics() {
	sum := 0.0
	for _, requestTime := range requestTimes {
		sum += float64(requestTime)
	}
	mean := sum / float64(len(requestTimes))

	sd := 0.0
	for _, requestTime := range requestTimes {
		sd += math.Pow(float64(requestTime)-mean, 2)
	}
	sd = math.Sqrt(sd / float64(len(requestTimes)))

	println("Mean: ", mean/float64(time.Millisecond))
	println("SD: ", sd/float64(time.Millisecond))
}

func main() {
	fmt.Printf("Olá, digite seu nome: ")
	clientName := "user"
	fmt.Printf("Olá, pode começar a falar!\n")

	ctx := context.Background()

	conn, err := grpc.Dial("localhost:9997", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := chat.NewChatClient(conn)
	stream, err := c.Chat(ctx)
	if err != nil {
		panic(err)
	}

	for len(requestTimes) < 10000 {
		handleInput(stream, clientName)
		handleResponse(stream)
	}

	calculateMetrics()
}

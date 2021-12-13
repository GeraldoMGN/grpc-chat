package main

import (
	"bufio"
	"github.com/GeraldoMGN/grpc-chat/chat"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
)

var (
	requestStartTime time.Time
	requestTimes     []time.Duration
)

func handleInput(stream chat.Chat_ChatClient, clientReader *bufio.Reader, clientName string) {
	clientRequest, err := clientReader.ReadString('\n')
	clientRequest = strings.TrimSpace(clientRequest)

	requestStartTime = time.Now()

	err = stream.Send(&chat.ChatMessage{
		User:    clientName,
		Message: clientRequest + "\n",
	})
	if err != nil {
		panic(err)
	}
}

func handleResponse(stream chat.Chat_ChatClient) {
	message, err := stream.Recv()
	serverResponse := message.GetMessage()

	elapsed := time.Since(requestStartTime)
	requestTimes = append(requestTimes, elapsed)

	switch err {
	case nil:
		fmt.Print("\033[H\033[2J")
		fmt.Println(strings.TrimSpace(serverResponse))
	case io.EOF:
		log.Println("Servidor fechou a conexão")
		return
	default:
		log.Printf("Erro no servidor: %v\n", err)
		return
	}
}

func main() {
	clientReader := bufio.NewReader(os.Stdin)

	fmt.Printf("Olá, digite seu nome: ")
	clientName, err := clientReader.ReadString('\n')
	clientName = strings.TrimSpace(clientName)
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

	for {
		handleInput(stream, clientReader, clientName)
		handleResponse(stream)
	}
}

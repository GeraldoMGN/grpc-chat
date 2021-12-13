package main

import (
	"io"
	"net"
	"strings"

	"github.com/GeraldoMGN/grpc-chat/chat"

	"google.golang.org/grpc"
)

var (
	chatHistory string
)

type Connection struct {
	conn chat.Chat_ChatServer
	send chan *chat.ChatMessage
}

func (c *Connection) GetMessages() error {
	for {
		data, err := c.conn.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		msg := data.GetUser() + ": " + data.GetMessage()
		if strings.Split(msg, ": ")[1] != "\n" && len(chatHistory) < 1024 {
			chatHistory += msg
		}
		c.conn.Send(&chat.ChatMessage{
			User:    "",
			Message: chatHistory + "\b",
		})
	}
}

type ChatServer struct{}

func (c *ChatServer) Chat(stream chat.Chat_ChatServer) error {
	conn := &Connection{
		conn: stream,
		send: make(chan *chat.ChatMessage),
	}

	err := conn.GetMessages()

	return err
}

func main() {
	lst, err := net.Listen("tcp", "0.0.0.0:9997")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	srv := &ChatServer{}
	chat.RegisterChatServer(s, srv)

	err = s.Serve(lst)
	if err != nil {
		panic(err)
	}
}

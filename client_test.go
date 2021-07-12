package gocord

import (
	"fmt"
	"os"
)

var (
	client *Client
)

func init() {
	var err error
	client, err = New(Options{
		Token: os.Getenv("BOT_TOKEN"),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(client.Options.Token)
	client.On("ready", func(ctx *Context) {
		fmt.Println("Client ready")
	})
	client.Connect()
}

package examples

import (
	"fmt"

	"github.com/gocord/gocord"
)

func main() {
	client, err := gocord.New(gocord.Options{
		Token:   "YOUR TOKEN HERE",
		Intents: gocord.Intents.ALL,
	})
	/* client creation has errored if nil */
	if err != nil {
		panic(err)
	}

	client.On("ready", func(ctx *gocord.Context) {
		fmt.Printf("%s ready", ctx.Client.User.Username)
	})

	client.On("messageCreate", func(ctx *gocord.Context) {
		/* check if context is empty , very unlikely */
		if ctx == nil {
			return
		}
		fmt.Println(ctx.Message.Content)
	})

	defer client.Close()
	client.Connect()
}

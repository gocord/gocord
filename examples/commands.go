package examples

import (
	"gitlab.com/weary/gocord"
	"gitlab.com/weary/gocord/commands"
)

func main() {
	client, err := gocord.New(gocord.Options{
		Token:   "abc",
		Intents: gocord.Intents.ALL,
	})
	defer client.Close()
	if err != nil {
		panic(err)
	}

	client.On("messageCreate", func(ctx *gocord.Context) {
		cmds := commands.New(commands.Options{
			Prefix: "!",
		})

		cmds.On("abc", func(c *gocord.Context) {

		})

		cmds.Handle(ctx)
	})

	client.Connect()
}

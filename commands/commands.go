package commands

import (
	"strings"

	"github.com/gocord/gocord"
)

// TODO add custom context that inherits gocord.Context?
type Command func(ctx *gocord.Context)

type CommandHandler struct {
	Prefix   string
	Commands []CommandRoute
}

type Options struct {
	Prefix string
}

type CommandRoute struct {
	Name    string
	Aliases []string
	Handler Command
}

func New(opts Options) *CommandHandler {
	return &CommandHandler{
		Prefix: opts.Prefix,
	}
}

func (c *CommandHandler) Handle(ctx *gocord.Context) {
	if ctx == nil ||
		ctx.Type != gocord.EVENTS.MESSAGE_CREATE ||
		ctx.Message.Author.Bot {
		return
	}
	args := strings.Split(ctx.Message.Content, " ")
	command := args[0][len(c.Prefix):]
	if len(args) > 1 {
		args = args[1:]
	} else {
		args = []string{}
	}
	for _, cmd := range c.Commands {
		if cmd.Name == command {
			cmd.Handler(ctx)
			return
		}
		for _, alias := range cmd.Aliases {
			if alias == command {
				cmd.Handler(ctx)
				return
			}
		}
	}
}

func (c *CommandHandler) On(command string, handler Command) *CommandRoute {
	route := CommandRoute{
		Name:    command,
		Handler: handler,
	}
	c.Commands = append(c.Commands, route)
	return &c.Commands[len(c.Commands)-1]
}

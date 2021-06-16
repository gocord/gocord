package gocord

type EventFunc func(context *Context)

type Client struct {

	// General Values
	Options Options
	ws      *Websocket

	// Events
	events map[string][]EventFunc

	// JSON
	User User `json:"user"`

	// Cache
	Guilds   *GuildCache
	Channels *ChannelCache
	Users    *UserCache
}

type Options struct {
	Token   string
	Intents int
}

// Creates the client and allows us to access methods to interact with Discord.
func New(options Options) (c *Client, err error) {
	c = &Client{
		Options: options,
		events:  make(map[string][]EventFunc),
	}

	if c == nil {
		err = ErrClientCreate
	}
	return
}

// Opens a websocket connection with Discord.
func (c *Client) Connect() error {
	websocket := newWebsocket(c)
	c.ws = websocket

	if err := websocket.connect(); err != nil {
		return err
	}

	<-websocket.listening
	return nil
}

// Closes the current websocket connection with Discord.
func (c *Client) Close() {
	c.ws.conn.Close()
}

// Creates an event listener to listen for the specified event.
func (c *Client) On(event string, ev EventFunc) {
	c.events[event] = append(c.events[event], ev)
}

// Internal call function.
func (c *Client) call(event string, ctx *Context) {
	for _, fn := range c.events[event] {
		fn(ctx)
	}
}

// Exported call function. User can run their own events.
func (c *Client) Call(event string, ctx *Context) {
	if protectedEvents[event] {
		return
	}
	for _, fn := range c.events[event] {
		fn(ctx)
	}
}

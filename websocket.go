package gocord

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/gobwas/ws"
)

type Websocket struct {
	conn          net.Conn
	gateway       string
	lastHeartbeat time.Time
	seq           int64
	interval      time.Duration
	client        *Client
	listening     chan interface{}
}

func newWebsocket(c *Client) *Websocket {
	ws := Websocket{
		client: c,
	}
	return &ws
}

func (w *Websocket) getGateway() {
	res, err := http.Get("https://discordapp.com/api/v9/gateway")
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	var t struct {
		URL string `json:"url"`
	}

	if json.Unmarshal(body, &t) != nil {
		return
	}

	w.gateway = t.URL
}

func (w *Websocket) connect() error {
	/* get gateway endpoint from api */
	w.getGateway()

	/* err if gateway is empty */
	if w.gateway == "" {
		return ErrNoGateway
	}

	/* init gateway conn */
	header := http.Header{}
	header.Add("accept-encoding", "zlib")
	var err error
	w.conn, _, _, err = ws.Dial(context.Background(), w.gateway)
	if err != nil {
		return ErrConnFailed
	}
	fmt.Println("connected to ws")
	w.listening = make(chan interface{})

	// read first message
	m, err := w.readMessage()
	// t, m, err := w.conn.Read(context.Background())
	if err != nil {
		return ErrCannotRead
	}
	w.handleEvent(m)

	w.identify()

	m, err = w.readMessage()
	if err != nil {
		return ErrCannotRead
	}

	// Ready Event
	w.handleEvent(m)

	// Initaliase Cache
	{
		w.client.Guilds = &GuildCache{}
		w.client.Guilds.cache.Init()
		w.client.fetchGuilds()
	}

	// TODO: Check for other websocket inital messages for things like gateway resume

	// Make calls to this channel and add more listeners
	w.listening = make(chan interface{})

	go w.heartbeat(w.interval)
	go w.events()

	return nil
}

type identifyDataProperties struct {
	OS      string `json:"$os"`
	Browser string `json:"$browser"`
	Device  string `json:"$device"`
}

type identifyData struct {
	Token      string                 `json:"token"`
	Intents    int                    `json:"intents"`
	Properties identifyDataProperties `json:"properties"`
}

type identify struct {
	Op   int          `json:"op"`
	Data identifyData `json:"d"`
}

func (w *Websocket) writeJSON(data interface{}) error {
	bytes, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))
	return ws.WriteFrame(w.conn, ws.NewTextFrame(bytes))
	// json.NewEncoder(w.conn).Encode(&data)
}

func (w *Websocket) identify() {
	w.writeJSON(identify{
		Op: 2,
		Data: identifyData{
			Token:   w.client.Options.Token,
			Intents: w.client.Options.Intents,
			Properties: identifyDataProperties{
				"linux", "gocord/client 0.1", "gocord/client 0.1",
			},
		},
	})
}

func (w *Websocket) heartbeat(interval time.Duration) {
	ticker := time.NewTicker(interval * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.lastHeartbeat = time.Now()
			w.writeJSON(heartbeatOp{1, w.seq})
		case <-w.listening:
			return
		}
	}
}

func (w *Websocket) readMessage() ([]byte, error) {
	frame, err := ws.ReadFrame(w.conn)
	if err != nil {
		return nil, err
	}
	fmt.Println(frame.Header.OpCode)
	// if _, err := io.Copy(&buf, w.conn); err != nil {
	// 	return nil, err
	// }
	return frame.Payload, nil
}

func (w *Websocket) events() {
	for {
		m, err := w.readMessage()
		if err != nil {
			return
		}
		select {
		case <-w.listening:
			return
		default:
			if err := w.handleEvent(m); err != nil {
				return
			}
		}
	}
}

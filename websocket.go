package gocord

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

func (w *Websocket) Close() {
	w.conn.Close()
	w.listening <- nil
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
	w.listening = make(chan interface{})

	// read first message
	f, err := w.readMessage()
	// t, m, err := w.conn.Read(context.Background())
	if err != nil {
		return ErrCannotRead
	}
	w.handleEvent(f)

	w.identify()

	f, err = w.readMessage()
	if err != nil {
		fmt.Println(err)
		return ErrCannotRead
	}
	w.handleEvent(f)

	// Initaliase Caches
	{
		w.client.Guilds = &GuildCache{}
		w.client.Guilds.cache.Init()
		w.client.fetchGuilds()
	}

	// TODO: Check for other websocket inital messages for things like gateway resume

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
	fmt.Printf("Writing: %s\n", string(bytes))
	frame := ws.NewTextFrame(bytes)
	fmt.Println(frame.Header.Length)
	return ws.WriteFrame(w.conn, frame)
}

func (w *Websocket) mustWriteJSON(data interface{}) {
	if err := w.writeJSON(data); err != nil {
		panic(err)
	}
}

func (w *Websocket) identify() {
	w.mustWriteJSON(identify{
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
			w.sendHeartbeat()
		case <-w.listening:
			return
		}
	}
}

func (w *Websocket) sendHeartbeat() {
	w.writeJSON(heartbeatOp{1, w.seq})
}

func (w *Websocket) readMessage() (ws.Frame, error) {
	// var err error
	// var frame ws.Frame
	frame, err := ws.ReadFrame(w.conn)
	// for err == io.EOF {
	// 	// fmt.Println(err)
	// 	frame, err = ws.ReadFrame(w.conn)
	// }
	if frame.Header.OpCode == ws.OpClose {
		fmt.Println("closing...")
		w.Close()
		return frame, ErrConnClose
	}
	fmt.Println(string(frame.Payload))
	return frame, err
}

func (w *Websocket) events() {
	for {
		m, err := w.readMessage()
		if err != nil {
			if err == io.EOF {
				continue
			}
			fmt.Println(err)
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

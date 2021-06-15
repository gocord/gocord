package gocord

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Websocket struct {
	conn          *websocket.Conn
	gateway       string
	lastHeartbeat time.Time
	seq           int64
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
	w.conn, _, err = websocket.DefaultDialer.Dial(w.gateway, header)
	if err != nil {
		return ErrConnFailed
	}

	w.conn.SetCloseHandler(func(code int, text string) error {
		return nil
	})

	t, m, err := w.conn.ReadMessage()
	if err != nil {
		return ErrCannotRead
	}

	ev, err := w.readEvent(t, m)
	if err != nil {
		return err
	}

	if ev.Op != 10 {
		return ErrExpectedHello
	}

	var h helloOp
	if err := json.Unmarshal(ev.Data, &h); err != nil {
		return ErrEventDecode
	}

	w.identify()

	t, m, err = w.conn.ReadMessage()
	if err != nil {
		return ErrCannotRead
	}
	ev, err = w.readEvent(t, m)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(ev.Data, &w.client); err != nil {
		return ErrEventDecode
	}

	// Ready Event
	ctx, eType := acquireContext(ev, w.client)
	w.client.call(eType, ctx)

	// Initaliase Caches
	w.client.Guilds.cache.Init()

	//
	w.listening = make(chan interface{})

	go w.heartbeat(h.Interval)
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

func (w *Websocket) identify() {
	w.conn.WriteJSON(identify{
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
			w.conn.WriteJSON(heartbeatOp{1, w.seq})
		case <-w.listening:
			return
		}
	}
}

func (w *Websocket) events() {
	for {
		mt, m, err := w.conn.ReadMessage()
		if err != nil {
			return
		}
		select {
		case <-w.listening:
			return
		default:
			ev, err := w.readEvent(mt, m)
			if err != nil {
				return
			}

			ctx, eType := acquireContext(ev, w.client)
			if eType == "" || ctx == nil {
				return
			}
			w.client.call(eType, ctx)
		}
	}
}

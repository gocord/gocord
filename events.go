package gocord

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"io"
	"time"

	"github.com/gorilla/websocket"
)

type Event struct {
	Type string          `json:"t"`
	Seq  int64           `json:"s"`
	Op   int             `json:"op"`
	Data json.RawMessage `json:"d"`

	Other interface{} `json:"-"`
}

type helloOp struct {
	Interval time.Duration `json:"heartbeat_interval"`
}

type heartbeatOp struct {
	Op   int   `json:"op"`
	Data int64 `json:"d"`
}

// TODO: migrate everything to read raw response and handle all in one function here

func (w *Websocket) readEvent(mType int, data []byte) (*Event, error) {

	var reader io.Reader
	reader = bytes.NewBuffer(data)

	if mType == websocket.BinaryMessage {
		z, err := zlib.NewReader(reader)
		if err != nil {
			return nil, ErrDecompressEvent
		}
		defer z.Close()

		reader = z
	}

	var e *Event
	dec := json.NewDecoder(reader)
	if err := dec.Decode(&e); err != nil {
		return e, ErrEventDecode
	}

	if e.Op == 1 {
		err := w.conn.WriteJSON(heartbeatOp{1, w.seq})
		if err != nil {
			return e, ErrHeartbeat
		}

		return e, nil
	}

	w.seq = e.Seq

	return e, nil
}

package gocord

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"errors"
	"io"
	"time"

	"nhooyr.io/websocket"
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

// TODO: add hello event checking
func (w *Websocket) handleEvent(mType websocket.MessageType, data []byte) error {

	// Define as io.Reader for zlib
	var reader io.Reader
	reader = bytes.NewBuffer(data)

	if mType == websocket.MessageBinary {
		zl, err := zlib.NewReader(reader)
		if err != nil {
			return err
		}
		defer zl.Close()
		reader = zl
	}

	// Unmarshal websocket message into event
	var ev *Event
	dec := json.NewDecoder(reader)
	if err := dec.Decode(&ev); err != nil {
		return err
	}

	switch ev.Op {
	case 10:
		var h helloOp
		if err := json.Unmarshal(ev.Data, &h); err != nil {
			return err
		}

		w.client.ReadyAt = time.Now()
		w.interval = h.Interval
	}

	// Set websocket sequence
	w.seq = ev.Seq

	var resp struct {
		Type int `json:"type"`
	}
	if err := json.Unmarshal(ev.Data, &resp); err != nil {
		return err
	}
	if resp.Type != 0 {
		return errors.New("invalid data type")
	}

	var ctx Context
	ctx.Type = ev.Type

	switch ev.Type {
	case EVENTS.MESSAGE_CREATE:
		var message Message
		json.Unmarshal([]byte(ev.Data), &message)
		message.Channel = w.client.getChannel(message.ChannelID)
		ctx.Message = &message
	/* this should be temporary lol , have this switch every event */
	case EVENTS.CHANNEL_CREATE, EVENTS.CHANNEL_DELETE, EVENTS.CHANNEL_UPDATE:
		var channel Channel
		json.Unmarshal([]byte(ev.Data), &channel)
		ctx.Channel = &channel
	case EVENTS.GUILD_BAN_ADD, EVENTS.GUILD_BAN_REMOVE:
		var member Member
		json.Unmarshal([]byte(ev.Data), &member)
		ctx.Member = &member
	case EVENTS.GUILD_MEMBER_ADD, EVENTS.GUILD_MEMBER_REMOVE, EVENTS.GUILD_MEMBER_UPDATE:
		var member Member
		json.Unmarshal([]byte(ev.Data), &member)
		ctx.Member = &member
	}

	w.client.call(getEventName(ev.Type), &ctx)
	return nil
}

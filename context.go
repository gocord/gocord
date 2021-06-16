package gocord

import (
	"bytes"
	"encoding/json"
)

// TODO: finish labeling this struct

type Context struct {
	// Event Type
	Type string
	// Only exists on messageCreate
	Message *Message
	// Always Exists
	Client *Client
	// Exists in messageCreate, channelCreate, channelDelete and channelUpdate
	Channel *Channel
	// Exists on guildMemberAdd
	Member *Member
}

func acquireContext(ev *Event, client *Client) (*Context, string) {
	reader := bytes.NewBuffer(ev.Data)
	var response struct {
		Type int `json:"type"`
	}
	dec := json.NewDecoder(reader)
	dec.Decode(&response)

	var ctx Context
	var eventType string

	ctx.Client = client

	if response.Type != 0 {
		return &ctx, eventType
	}
	eventType = GetEventName(ev.Type)
	ctx.Type = ev.Type

	switch ev.Type {
	case EVENTS.MESSAGE_CREATE:
		var message Message
		json.Unmarshal([]byte(ev.Data), &message)
		message.Channel = client.getChannel(message.ChannelID)

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
	return &ctx, eventType
}

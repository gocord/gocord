package gocord

import (
	"encoding/json"
	"errors"
)

type Message struct {
	// Client
	Client *Client

	// JSON
	Content   string            `json:"content"`
	ID        *Snowflake        `json:"id"`
	ChannelID string            `json:"channel_id"`
	Member    *Member           `json:"member"`
	Embeds    []*Embed          `json:"embeds"`
	Author    *User             `json:"author"`
	Mentions  []*MessageMention `json:"mentions"`
	Partial   bool              `json:"partial"`
	Pinned    bool              `json:"pinned"`
	Thread    *MessageThread    `json:"thread"`
	TTS       bool              `json:"tts"`
	Nonce     string            `json:"nonce"`
	Type__    int               `json:"type"`

	// General
	ApplicationID *Snowflake
	Deleted       bool
	System        bool
	Type          *MessageType

	// Fetched data
	Channel *Channel
	Guild   *Guild
}

type MessageThread struct {
}

type MessageType string

var MessageTypes = StringArray{
	"DEFAULT", "RECIPIENT_ADD", "RECIPIENT_REMOVE",
	"CALL", "CHANNEL_NAME_CHANGE", "CHANNEL_ICON_CHANGE",
	"PINS_ADD", "GUILD_MEMBER_JOIN", "USER_PREMIUM_GUILD_SUBSCRIPTION",
	"USER_PREMIUM_GUILD_SUBSCRIPTION_TIER_1", "USER_PREMIUM_GUILD_SUBSCRIPTION_TIER_2",
	"USER_PREMIUM_GUILD_SUBSCRIPTION_TIER_3", "CHANNEL_FOLLOW_ADD", "",
	"GUILD_DISCOVERY_DISQUALIFIED", "GUILD_DISCOVERY_REQUALIFIED", "GUILD_DISCOVERY_GRACE_PERIOD_INITIAL_WARNING",
	"GUILD_DISCOVERY_GRACE_PERIOD_FINAL_WARNING", "THREAD_CREATED", "REPLY",
	"APPLICATION_COMMAND", "THREAD_STARTER_MESSAGE", "GUILD_INVITE_REMINDER",
}

var SystemMessageTypes = StringArray{
	"RECIPIENT_ADD", "RECIPIENT_REMOVE",
	"CALL", "CHANNEL_NAME_CHANGE", "CHANNEL_ICON_CHANGE",
	"PINS_ADD", "GUILD_MEMBER_JOIN", "USER_PREMIUM_GUILD_SUBSCRIPTION",
	"USER_PREMIUM_GUILD_SUBSCRIPTION_TIER_1", "USER_PREMIUM_GUILD_SUBSCRIPTION_TIER_2",
	"USER_PREMIUM_GUILD_SUBSCRIPTION_TIER_3", "CHANNEL_FOLLOW_ADD", "",
	"GUILD_DISCOVERY_DISQUALIFIED", "GUILD_DISCOVERY_REQUALIFIED", "GUILD_DISCOVERY_GRACE_PERIOD_INITIAL_WARNING",
	"GUILD_DISCOVERY_GRACE_PERIOD_FINAL_WARNING", "THREAD_CREATED", "REPLY",
	"APPLICATION_COMMAND", "THREAD_STARTER_MESSAGE", "GUILD_INVITE_REMINDER",
}

func (mt *MessageType) check() (bool, bool) {
	return MessageTypes.Includes(string(*mt)), SystemMessageTypes.Includes(string(*mt))
}

type MessageMention struct {
}

func newMessage(client *Client, data json.RawMessage) *Message {
	msg := Message{
		Client:  client,
		Deleted: false,
	}
	if err := msg.patch(data); err != nil {
		return nil
	}

	return &msg
}

func (m *Message) patch(data json.RawMessage) error {
	if err := json.Unmarshal(data, m); err != nil {
		return err
	}

	m.Type = (*MessageType)(&MessageTypes[m.Type__])
	mt, smt := m.Type.check()
	if !mt {
		return errors.New("invalid message type")
	}
	m.System = smt

	if m.Author != nil {
		m.Client.Users.cache.add(m.Author)
	}

	m.Channel = m.Client.getChannel(m.ChannelID)

	return nil
}

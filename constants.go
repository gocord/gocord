package gocord

import "fmt"

var ErrClientCreate = fmt.Errorf("an error occurred while instantiating the client")

var ErrNoGateway = fmt.Errorf("failed to get gateway uri")
var ErrConnFailed = fmt.Errorf("failed to connect to gateway")
var ErrCannotRead = fmt.Errorf("failed to read gateway message")

var ErrDecompressEvent = fmt.Errorf("failed to decompress event message")
var ErrEventDecode = fmt.Errorf("failed to decode even message")
var ErrHeartbeat = fmt.Errorf("failed to send heartbeat")
var ErrExpectedHello = fmt.Errorf("gateway returned non-hello response")

var Intents = struct {
	GUILDS                   int
	GUILD_MEMBERS            int
	GUILD_BANS               int
	GUILD_EMOJIS             int
	GUILD_INTEGRATIONS       int
	GUILD_WEBHOOKS           int
	GUILD_INVITES            int
	GUILD_VOICE_STATES       int
	GUILD_PRESENCES          int
	GUILD_MESSAGES           int
	GUILD_MESSAGE_REACTIONS  int
	GUILD_MESSAGE_TYPING     int
	DIRECT_MESSAGES          int
	DIRECT_MESSAGE_REACTIONS int
	DIRECT_MESSAGE_TYPING    int
	ALL                      int
}{
	1 << 0, 1 << 1, 1 << 2,
	1 << 3, 1 << 4, 1 << 5,
	1 << 6, 1 << 7, 1 << 8,
	1 << 9, 1 << 10, 1 << 11,
	1 << 12, 1 << 13, 1 << 14,
	1 << 0 & 1 << 1 & 1 << 2 & 1 << 3 & 1 << 4 & 1 << 5 & 1 << 6 & 1 << 7 & 1 << 8 & 1 << 9 & 1 << 10 & 1 << 11 & 1 << 12 & 1 << 13 & 1 << 14,
}

var EVENTS = struct {
	READY               string
	MESSAGE_CREATE      string
	CHANNEL_CREATE      string
	CHANNEL_UPDATE      string
	CHANNEL_DELETE      string
	GUILD_BAN_ADD       string
	GUILD_BAN_REMOVE    string
	GUILD_MEMBER_ADD    string
	GUILD_MEMBER_REMOVE string
	GUILD_MEMBER_UPDATE string
}{
	"READY",
	"MESSAGE_CREATE",
	"CHANNEL_CREATE",
	"CHANNEL_UPDATE",
	"CHANNEL_DELETE",
	"GUILD_BAN_ADD",
	"GUILD_BAN_REMOVE",
	"GUILD_MEMBER_ADD",
	"GUILD_MEMBER_REMOVE",
	"GUILD_MEMBER_UPDATE",
}

var protectedEvents = map[string]bool{
	GetEventName(EVENTS.READY):               true,
	GetEventName(EVENTS.MESSAGE_CREATE):      true,
	GetEventName(EVENTS.CHANNEL_CREATE):      true,
	GetEventName(EVENTS.CHANNEL_UPDATE):      true,
	GetEventName(EVENTS.CHANNEL_DELETE):      true,
	GetEventName(EVENTS.GUILD_BAN_ADD):       true,
	GetEventName(EVENTS.GUILD_BAN_REMOVE):    true,
	GetEventName(EVENTS.GUILD_MEMBER_ADD):    true,
	GetEventName(EVENTS.GUILD_MEMBER_REMOVE): true,
	GetEventName(EVENTS.GUILD_MEMBER_UPDATE): true,
}

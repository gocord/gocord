package gocord

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

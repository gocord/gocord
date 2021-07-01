package gocord

type Message struct {
	// Client
	Client *Client

	// JSON
	Content   string           `json:"content"`
	ID        Snowflake        `json:"id"`
	ChannelID string           `json:"channel_id"`
	Member    Member           `json:"member"`
	Embeds    []Embed          `json:"embeds"`
	Author    User             `json:"author"`
	Mentions  []MessageMention `json:"mentions"`
	Partial   bool             `json:"partial"`

	// Fetched data
	Channel *Channel
	Guild   *Guild
}

type MessageMention struct {
}

package gocord

type Message struct {
	Client    *Client
	Content   string `json:"content"`
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	Member    Member `json:"member"`
	User      User   `json:"user"`
	Author    Author `json:"author"`
	Channel   *Channel
	Guild     *Guild
}

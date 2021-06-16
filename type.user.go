package gocord

// TODO this is all a fucking mess , clean it up

type User struct {
	// Client
	client *Client

	// JSON
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Bot           bool   `json:"bot"`
	Partial       bool   `json:"partial"`
}

type Author struct {
	// JSON
	Bot bool `json:"bot"`

	// Inheritance
	User
}

// idk if member should inherit user
type Member struct {
	// Inheritance
	User

	// General
	Guild *Guild
}

func (m *Member) Ban(reason string) error {
	return m.Guild.BanMember(m.User.ID, reason)
}

// Caching related

type UserCache struct {
	cache Cache
}

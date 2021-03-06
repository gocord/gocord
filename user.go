package gocord

type User struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Bot           bool   `json:"bot"`
}

type Author struct {
	Bot bool `json:"bot"`
	User
}

// idk if member should inherit user
type Member struct {
	User
	Guild *Guild
}

func (m *Member) Ban(reason string) error {
	return m.Guild.BanMember(m.User.ID, reason)
}

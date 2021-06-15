package gocord

import "fmt"

type Guild struct {
	ID       string
	Channels []*Channel
	Members  []*Member
	Client   *Client
}

func (g *Guild) BanMember(memberid, reason string) error {
	_, err := g.Client.sendRequest(
		fmt.Sprintf("/guilds/%s/bans/%s", g.ID, memberid),
		"PUT",
		fmt.Sprintf(`{"reason":"%s"}`, reason),
	)
	return err
}

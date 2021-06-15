package gocord

import "fmt"

type Guild struct {
	ID       string
	Channels []*Channel
	Members  []*Member
	Client   *Client
}

type GuildCache struct {
	// don't export cache xd
	cache Cache
}

func (c *GuildCache) Get(snowflake string) Guild {
	return c.cache.Get(snowflake).(Guild)
}

// Bans the member from the guild
func (g *Guild) BanMember(memberid, reason string) error {
	_, err := g.Client.sendRequest(
		fmt.Sprintf("/guilds/%s/bans/%s", g.ID, memberid),
		"PUT",
		fmt.Sprintf(`{"reason":"%s"}`, reason),
	)
	return err
}

// Unbans the member from the guild
func (g *Guild) UnbanMember(memberId, reason string) error {
	_, err := g.Client.sendRequest(
		fmt.Sprintf("/guilds/%s/bans/%s", g.ID, memberId),
		"DELETE",
		"", // TODO: Add reason support
	)
	return err
}

// Kicks (or removes) a member from the guild
func (g *Guild) KickMember(memberId, reason string) error {
	_, err := g.Client.sendRequest(
		fmt.Sprintf("/guilds/%s/members/%s", g.ID, memberId),
		"DELETE",
		"", // TODO: Add reason support
	)
	return err
}

// Creates a role
func (g *Guild) CreateRole(name, permissions string, color int, hoist, mentionable bool) error {
	_, err := g.Client.sendRequest(
		fmt.Sprintf("/guilds/%s/roles", g.ID),
		"POST",
		fmt.Sprintf(
			`{"name":"%s", "permissions":"%s", "color":%d, "hoist":%t, "mentionable":%t}`,
			name,
			permissions,
			color,
			hoist,
			mentionable,
		),
	)
	return err
}

// Deletes a role
func (g *Guild) DeleteRole(roleId string) error {
	_, err := g.Client.sendRequest(
		fmt.Sprintf("/guilds/%s/roles/%s", g.ID, roleId),
		"DELETE",
		"", // TODO: Add reason support
	)
	return err
}

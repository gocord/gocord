package gocord

import (
	"encoding/json"
	"fmt"
)

type Guild struct {
	// Client
	client *Client

	// JSON
	Partial bool      `json:"partial"`
	ID      Snowflake `json:"id"`

	// Fetched
	Channels []*Channel
	Members  []*Member
}

// Bans the member from the guild
func (g *Guild) BanMember(memberid, reason string) error {
	_, err := g.client.sendRequest(
		fmt.Sprintf("/guilds/%s/bans/%s", g.ID, memberid),
		"PUT",
		fmt.Sprintf(`{"reason":"%s"}`, reason),
	)
	return err
}

// Unbans the member from the guild
func (g *Guild) UnbanMember(memberId, reason string) error {
	_, err := g.client.sendRequest(
		fmt.Sprintf("/guilds/%s/bans/%s", g.ID, memberId),
		"DELETE",
		"", // TODO: Add reason support
	)
	return err
}

// Kicks (or removes) a member from the guild
func (g *Guild) KickMember(memberId, reason string) error {
	_, err := g.client.sendRequest(
		fmt.Sprintf("/guilds/%s/members/%s", g.ID, memberId),
		"DELETE",
		"", // TODO: Add reason support
	)
	return err
}

// Creates a role
func (g *Guild) CreateRole(name, permissions string, color int, hoist, mentionable bool) error {
	_, err := g.client.sendRequest(
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
	_, err := g.client.sendRequest(
		fmt.Sprintf("/guilds/%s/roles/%s", g.ID, roleId),
		"DELETE",
		"", // TODO: Add reason support
	)
	return err
}

// Caching related

type GuildCache struct {
	cache Cache
}

func (c *GuildCache) Get(snowflake string) Guild {
	return c.cache.get(snowflake).(Guild)
}

func (c *Client) fetchGuilds() error {
	g, err := c.sendRequest("/users/@me/guilds", "GET", "")
	if err != nil {
		return err
	}
	var guilds []Guild
	json.Unmarshal([]byte(g), &g)
	for _, guild := range guilds {
		c.Guilds.cache.set(guild.ID.string, guild)
	}
	return nil
}

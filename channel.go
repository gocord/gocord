package gocord

import (
	"encoding/json"
	"fmt"
)

type Channel struct {
	// Client
	client *Client

	ID   string `json:"id"`
	Name string `json:"name"`
	NSFW bool   `json:"nsfw"`
}

type ChannelCache struct {
	cache Cache
}

func (c *ChannelCache) Get(snowflake string) Channel {
	return c.cache.get(snowflake).(Channel)
}

// Fetches a channel.
func (c *Client) getChannel(id string) *Channel {
	resp, err := c.sendRequest(
		fmt.Sprintf("/channels/%s", id),
		"GET",
		"",
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	var channel Channel
	json.Unmarshal([]byte(resp), &channel)
	channel.client = c
	return &channel
}

// Sends a message to a specified channel.
func (c *Channel) SendMessage(message string) error {
	/* need something better then c.Client */
	_, err := c.client.sendRequest(
		fmt.Sprintf("/channels/%s/messages", c.ID),
		"POST",
		fmt.Sprintf(`{"content": "%s"}`, message),
	)
	/* add checking for response */
	return err
}

func (c *Channel) SendEmbed(embed Embed) error {
	body, err := json.Marshal(embed)
	if err != nil {
		return err
	}
	_, err = c.client.sendRequest(
		fmt.Sprintf("/channels/%s/messages", c.ID),
		"POST",
		fmt.Sprintf(`{"embed": %s}`, body),
	)
	return err
}

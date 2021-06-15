package gocord

type Cache struct {
	cache map[string]interface{}
}

func (c *Cache) Init(data ...map[string]interface{}) {
	// make cache mutable
	c.cache = make(map[string]interface{})
	// if data is supplied , enter it into cache
	if len(data) > 0 {
		for _, d := range data {
			for k, v := range d {
				c.cache[k] = v
			}
		}
	}
}

func (c *Cache) Get(snowflake string) interface{} {
	return c.cache[snowflake]
}

func (c *Cache) Set(snowflake string, value interface{}) {
	c.cache[snowflake] = value
}

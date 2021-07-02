package gocord

type Cache struct {
	cache map[string]interface{}
}

// Initalise cache
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

func (c *Cache) get(snowflake string) interface{} {
	return c.cache[snowflake]
}

func (c *Cache) set(snowflake string, value interface{}) {
	c.cache[snowflake] = value
}

func (c *Cache) add(data interface{}) {
	d, ok := data.(struct{ ID string })
	if !ok {
		return
	}
	c.cache[d.ID] = d
}

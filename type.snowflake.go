package gocord

import (
	"strconv"
	"time"
)

type Snowflake string

func (s *Snowflake) ReadTimestamp() (time.Time, error) {
	n, err := strconv.ParseInt(string(*s), 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return (time.Unix((n/4194304)+1420070400000, 0)), nil
}

package gocord

type StringArray []string

func (s StringArray) Includes(val string) bool {
	for _, v := range s {
		if val == v {
			return true
		}
	}
	return false
}

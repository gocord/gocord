package gocord

import "strings"

// really fucking unreadable :D
func GetEventName(eName string) string {
	eName = strings.ToLower(eName)
	for i, rn := range eName {
		if rn == '_' {
			eName = eName[:i] + eName[i+1:]
			eName = eName[:i] + string(rn-28) + eName[i+1:]
		}
	}
	return eName
}

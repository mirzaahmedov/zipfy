package generate

import (
	"strings"
)

const allowedChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func UUIDFromInt(id int) string {
	var result string
	for id > 0 {
		result = string(allowedChars[id%62]) + result
		id = id / 62
	}
	return result
}

func ParseIntFromUUID(id string) int {
	var result int
	for _, r := range id {
		result = result*62 + strings.Index(allowedChars, string(r))
	}
	return result
}

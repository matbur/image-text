package server

import (
	"fmt"
	"strings"
)

func isIn(s string, ss []string) bool {
	for _, i := range ss {
		if s == i {
			return true
		}
	}
	return false
}

func parsePath(s string) (string, string, string, error) {
	s = strings.Trim(s, "/")
	ss := strings.Split(s, "/")
	if len(ss) != 3 {
		return "", "", "", fmt.Errorf("malformed path '%s'", s)
	}
	return ss[0], ss[1], ss[2], nil
}

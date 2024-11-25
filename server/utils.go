package server

import (
	"fmt"
	"strings"
)

func parsePath(s string) (string, string, string, error) {
	s = strings.Trim(s, "/")
	ss := strings.Split(s, "/")
	if len(ss) != 3 {
		return "", "", "", fmt.Errorf("malformed path '%s'", s)
	}
	return ss[0], ss[1], ss[2], nil
}

package templates

import "fmt"

func envTitleSuffix() string {
	if EnvLabel == "" {
		return ""
	}
	return fmt.Sprintf(" (%s)", EnvLabel)
}

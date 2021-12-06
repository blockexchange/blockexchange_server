package core

import "regexp"

var nameregex = regexp.MustCompile("^[a-zA-Z0-9_.-]*$")

func ValidateName(name string) bool {
	return name != "" && nameregex.Match([]byte(name))
}

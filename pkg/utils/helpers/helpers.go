package helpers

import (
	"fmt"
	"strings"
)

type Flags string

type Flag interface {
	SetFlag(flag string, value string) string
}

func (f *Flags) SetFlag(flag, value string) string {
	*f = Flags(strings.TrimSpace(string(*f) + " " + fmt.Sprintf("-%s=%s", flag, value)))
	return string(*f)
}

func GetFlagIndex(flags []string) int {
	for i, arg := range flags {
		if arg[0] == '-' {
			return i
		}
	}
	return 0
}

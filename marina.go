package marina

import (
	"fmt"
	"strings"
)

const (
	MajorVersion = 0
	MinorVersion = 1
)

func SimpleTypename(x interface{}) string {
	t := fmt.Sprintf("%T", x)
	t = strings.TrimPrefix(t, "*")
	t = strings.TrimPrefix(t, "netdl.")
	return t
}

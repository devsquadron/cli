//go:build !prod
// +build !prod

package constants

import (
	"fmt"
)

const (
	PORT              = ":4000"
	GLOBAL_CONFIG     = "devsquadron_test.yml"
	CACHE_COMPLETIONS = "completions_test.yml"
)

var (
	ENDPOINT = fmt.Sprintf("http://127.0.0.1%s/api/", PORT)
)

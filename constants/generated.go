package constants

import (
	_ "embed"
)

//go:generate bash version.sh
//go:embed Version.txt

var Version string

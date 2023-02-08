package exception

import (
	"fmt"
	"os"

	"github.com/devsquadron/ds/message"
)

func CheckErr(err error) {
	if err != nil {
		fmt.Println(message.Red("ERROR", err.Error()))
		os.Exit(1)
	}
}

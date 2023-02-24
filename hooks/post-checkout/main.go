package main

import (
	"fmt"

	"github.com/devsquadron/cli/hooks"
)

func main() {
	var (
		err error
		msg string
	)
	msg, err = hooks.PostCheckout()
	fmt.Println(msg)

	// TODO: this should only be in the verbose case
	if err != nil {
		fmt.Println(err)
	}
}

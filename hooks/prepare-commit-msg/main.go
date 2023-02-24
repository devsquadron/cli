package main

import (
	"fmt"
	"os"

	"github.com/devsquadron/cli/hooks"

	"github.com/devsquadron/cli/message"

	"github.com/devsquadron/models"
)

func appendToCommit(args []string) error {
	var (
		err         error
		cmtF        *os.File
		tskId       uint64
		tsk         *models.Task
		percentCmnt string
	)

	cmtF, err = os.OpenFile(os.Args[1], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	tskId, err = hooks.Cfg.TaskNumber()
	if err != nil {
		return err
	}

	tsk, err = hooks.TaskClient.GetTaskById(tskId, hooks.Cfg.Team())
	if err != nil {
		return err
	}
	percentCmnt = fmt.Sprintf("# %d%s", tsk.Percent, "%")

	_, err = cmtF.WriteString(percentCmnt)
	if err != nil {
		return err
	}

	return cmtF.Close()
}

func main() {
	var (
		err error
	)
	err = appendToCommit(os.Args)
	if err != nil {
		fmt.Println(
			message.Yellow("WARNING", "unable to populate task percent in commit"),
		)
		fmt.Println(err.Error())
	}
}

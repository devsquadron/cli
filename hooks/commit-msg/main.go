package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/devsquadron/cli/hooks"

	"github.com/devsquadron/cli/message"

	"github.com/devsquadron/models"
)

func commitMsg() error {
	var (
		err                       error
		commitFileBit             []byte
		percCmtLn, percNumStr, tm string
		tskId                     uint64
		tsk                       *models.Task
		newPrct                   int64
	)
	commitFileBit, err = os.ReadFile(os.Args[1])
	if err != nil {
		return err
	}
	lns := strings.Split(string(commitFileBit), "\n")
	percCmtLn = lns[len(lns)-2]
	percNumStr = strings.Trim(percCmtLn, "# %")
	newPrct, err = strconv.ParseInt(percNumStr, 10, 0)
	tm = hooks.Cfg.Team()

	tskId, err = hooks.Cfg.TaskNumber()
	if err != nil {
		return err
	}

	// TODO: remove unessecary get and replace with request method to only update task percentage
	tsk, err = hooks.TaskClient.GetTaskById(tskId, tm)
	if err != nil {
		return err
	}

	tsk.Percent = int(newPrct)
	if newPrct == 100 {
		tsk.Status = models.Status.Review
	}
	err = hooks.TaskClient.UpdateTask(tsk, tm)
	if err != nil {
		return err
	}
	fmt.Println(message.Green("Task Percent", percNumStr))
	return nil
}

func main() {
	var err error
	err = commitMsg()
	if err != nil {
		fmt.Println(
			message.Yellow("WARNING", "unable to update task percent"),
		)
		fmt.Println(err.Error())
	}
}

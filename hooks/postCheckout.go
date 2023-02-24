package hooks

import (
	"fmt"

	"github.com/devsquadron/cli/message"

	"github.com/devsquadron/models"
)

var (
	PostCheckoutSuccessMessage = func(title string) string {
		return message.Green("Tracking Task", fmt.Sprintf("Found task with title \"%s\" ", title))
	}
)

func PostCheckout() (string, error) {
	var (
		err     error
		dsUsrNm string
		taskNum uint64
		task    *models.Task
	)

	taskNum, err = Cfg.TaskNumber()
	if err != nil {
		return message.Yellow(
			"Warning",
			"Could not get a task number from your branch name.",
		), err
	}

	tm := Cfg.Team()
	task, err = TaskClient.GetTaskById(taskNum, tm)
	if err != nil {
		return message.Yellow(
			"Not Tracking",
			"You can proceed on this branch, but your progress is not visible",
		), err
	}

	dsUsrNm, err = Cfg.UsernameE()
	if err != nil {
		return message.Red(
			"No Username",
			"we tried to get your username from 'git config user.name'",
		), err
	}
	task.Developer = dsUsrNm
	task.Status = models.Status.Developing
	err = TaskClient.UpdateTask(task, tm)
	if err != nil {
		return message.Yellow(
			"Warning",
			"failed to update task with your details.",
		), err
	}
	return PostCheckoutSuccessMessage(task.Title), err
}

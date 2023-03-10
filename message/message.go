package message

import (
	"fmt"

	"github.com/devsquadron/models"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorViolet = "\033[35m"
	colorBlue   = "\033[34m"
	colorWhite  = "\033[37m"
)

func colorTagPrint(tag string, message string, color string) string {
	return fmt.Sprintf("%s[%s]%s %s", color, tag, ColorReset, message)
}

func Green(tag string, message string) string {
	return colorTagPrint(tag, message, ColorGreen)
}

func Yellow(tag string, message string) string {
	return colorTagPrint(tag, message, ColorYellow)
}

func Red(tag string, message string) string {
	return colorTagPrint(tag, message, ColorRed)
}

// TODO: find best way to prettyprint tasks
func printTask(tsk *models.Task, full bool) {
	fmt.Printf("%s[%d]%s %s\n", ColorYellow, tsk.ID, ColorReset, tsk.Title)
	if tsk.Tags != nil && len(tsk.Tags) != 0 {
		fmt.Printf("%s\t  [tags]%s %s\n", ColorGreen, ColorReset, tsk.Tags)
	}
	if tsk.Status == models.Status.Developing {
		fmt.Printf("%s\t[status]%s %s %s%s\n", colorViolet, ColorReset, tsk.Status, fmt.Sprintf("%d", tsk.Percent), "%")
	} else {
		fmt.Printf("%s\t[status]%s %s\n", colorViolet, ColorReset, tsk.Status)
	}
	if tsk.Developer != "" {
		fmt.Printf("%s\t   [dev]%s %s\n", colorBlue, ColorReset, tsk.Developer)
	}
	if full {
		fmt.Println(tsk.Criterion)
	}
}

func Task(tsk *models.Task) {
	printTask(tsk, true)
}

func TaskAbb(tsk *models.Task) {
	printTask(tsk, false)
}

func Header() {
	fmt.Println("________________________________________________")
}

func ListTasksByContext(allTsks *[]models.Task, listTagFlag string, listDevFlag string) {
	fmt.Println(fmt.Sprintf("tag %s | dev %s", listTagFlag, listDevFlag))
	var curTsks []models.Task

	for _, sts := range models.Statuses {
		curTsks = []models.Task{}
		for _, t := range *allTsks {
			if t.Status == sts {
				curTsks = append(curTsks, t)
			}
		}
		if len(curTsks) > 0 {
			fmt.Println("________________________________________________")
			fmt.Println(fmt.Sprintf("status %s", sts))
			for _, st := range curTsks {
				TaskAbb(&st)
			}
		}
	}
}

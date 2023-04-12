package message

import (
	"fmt"
	"os"

	"github.com/charmbracelet/glamour"
	"github.com/devsquadron/models"
	"golang.org/x/term"
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

func printTask(tsk *models.Task, full bool, includeStatus bool) {
	fmt.Printf("%s[%d]%s %s\n", ColorYellow, tsk.ID, ColorReset, tsk.Title)
	if tsk.Tags != nil && len(tsk.Tags) != 0 {
		fmt.Printf("%s\t  [tags]%s %s\n", ColorGreen, ColorReset, tsk.Tags)
	}
	if includeStatus {
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
	printTask(tsk, true, true)
}

func TaskAbb(tsk *models.Task, includeStatus bool) {
	printTask(tsk, false, includeStatus)
}

func PrettyTask(tsk *models.Task) error {
	var (
		err             error
		prettyCriterion string
	)
	printTask(tsk, false, true)
	prettyCriterion, err = glamour.Render(tsk.Criterion, "dark")
	if err != nil {
		return err
	}
	fmt.Print(prettyCriterion)
	return nil
}

func Header() {
	var (
		err error
		w   int
	)
	w, _, err = term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("_________________________________________")
		return
	}
	for i := 0; i < w; i++ {
		fmt.Print("_")
	}
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
			Header()
			fmt.Println(sts)
			for _, st := range curTsks {
				TaskAbb(&st, false)
			}
		}
	}
}

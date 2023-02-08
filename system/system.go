package system

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/devsquadron/ds/configuration"
	"github.com/devsquadron/ds/constants"
	"github.com/devsquadron/ds/exception"
	"github.com/devsquadron/requests"

	"github.com/devsquadron/models"

	"github.com/nanvenomous/exfs"
	"golang.org/x/term"
)

const (
	tempFileName = "devsquadron-task-*.md"
)

type SystemType interface {
	GetArg(args []string, exp string) string
	GetTwoArgs(args []string, exp string) (string, string)
	GetTaskId(idS string) (uint64, error)
	GetPassword() (string, error)
	GitUsername() (string, error)
	GitEmail() (string, error)
	GetText(string, string) (string, error)
	GetConfirmedPassword() (string, string, error)
	EditTempMarkdownFileDS(txt string) (string, error)
	CheckoutBranch(tsk *models.Task) error
}

type System struct {
	FS *exfs.FileSystem
}

func NewSystem(fs *exfs.FileSystem) SystemType {
	return &System{FS: fs}
}

func checkArgs(args []string, l int, exp string) {
	if len(args) != l {
		exception.CheckErr(errors.New(fmt.Sprintf("expected arguments %s", exp)))
	}
}

func (sys *System) GetArg(args []string, exp string) string {
	checkArgs(args, 1, exp)
	return args[0]
}

func (sys *System) GetTwoArgs(args []string, exp string) (string, string) {
	checkArgs(args, 2, exp)
	return args[0], args[1]
}

func (sys *System) GetTaskId(idS string) (uint64, error) {
	num, err := strconv.ParseUint(idS, 10, 64)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func (sys *System) GetPassword() (string, error) {
	var (
		err     error
		pswdByt []byte
	)
	fmt.Printf("Enter Password: ")
	pswdByt, err = term.ReadPassword(int(syscall.Stdin))
	fmt.Println("")
	if err != nil {
		return "", err
	}

	return string(pswdByt), nil
}

func (sys *System) GitUsername() (string, error) {
	var (
		err  error
		outs string
	)
	outs, _, err = sys.FS.Capture("git", []string{"config", "user.name"})
	return strings.TrimSpace(outs), err
}

func (sys *System) GitEmail() (string, error) {
	var (
		err  error
		outs string
	)
	outs, _, err = sys.FS.Capture("git", []string{"config", "user.email"})
	return strings.TrimSpace(outs), err
}

func (sys *System) GetText(prmpt string, dflt string) (string, error) {
	var (
		err   error
		input string
		rdr   *bufio.Reader
	)
	if dflt != "" {
		fmt.Print(
			fmt.Sprintf("%s (%s): ", prmpt, dflt),
		)
	} else {
		fmt.Print(
			fmt.Sprintf("%s: ", prmpt),
		)
	}
	rdr = bufio.NewReader(os.Stdin)
	input, err = rdr.ReadString('\n')
	if err != nil {
		return "", err
	}

	input = strings.TrimSuffix(input, "\n")
	if input != "" {
		return input, nil
	} else if dflt != "" {
		return dflt, nil
	}
	return "", errors.New(fmt.Sprintf("you must enter a %s", prmpt))
}

func (sys *System) GetConfirmedPassword() (string, string, error) {
	var (
		err         error
		pswdByt     []byte
		cnfmPswdByt []byte
	)
	fmt.Printf("Enter Password: ")
	pswdByt, err = term.ReadPassword(int(syscall.Stdin))
	fmt.Println("")
	if err != nil {
		return "", "", err
	}

	fmt.Printf("Confirm Password: ")
	cnfmPswdByt, err = term.ReadPassword(int(syscall.Stdin))
	fmt.Println("")
	if err != nil {
		return "", "", err
	}

	if string(pswdByt) != string(cnfmPswdByt) {
		return "", "", errors.New("Passwords did not match.")
	}
	return string(pswdByt), string(cnfmPswdByt), nil
}

func (sys *System) EditTempMarkdownFileDS(txt string) (string, error) {
	var editor string
	editor = os.Getenv(exfs.EDITOR_ENV_VAR)
	if editor == "" {
		return "", errors.New("devsquadron requires the EDITOR environment variable to be set\n for more info see https://developersquadron.com/troubleshooting/editor/")
	}

	return sys.FS.EditTemporaryFile(editor, tempFileName, txt)
}

func (sys *System) CheckoutBranch(tsk *models.Task) error {
	return sys.FS.Execute(
		"git",
		[]string{
			"checkout",
			"-b",
			fmt.Sprintf("%d-%s", tsk.ID, strings.Replace(tsk.Title, " ", "_", -1)),
		},
	)
}

func GlobalSetup() (
	SystemType,
	configuration.ConfigurationType,
	*requests.TaskClient,
	*requests.DeveloperClient,
	*requests.TeamClient,
) {
	fs := exfs.NewFileSystem()
	sys := NewSystem(fs)

	cfg := configuration.NewConfiguration(fs)
	tsk := requests.NewTaskClient(constants.ENDPOINT)
	dev := requests.NewDeveloperClient(constants.ENDPOINT)
	tm := requests.NewTeamClient(constants.ENDPOINT)
	return sys, cfg, tsk, dev, tm
}

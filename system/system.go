package system

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/devsquadron/cli/configuration"
	"github.com/devsquadron/cli/constants"
	"github.com/devsquadron/cli/exception"
	"github.com/devsquadron/requests"

	"github.com/devsquadron/models"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

type ConfigPrompt struct {
	Prompt      string
	CheckFunc   func() (string, error)
	DefaultFunc func() (string, error)
	SetFunc     func(string)
}

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type ConfigurationTextInputModel struct {
	focusIndex int
	Inputs     []textinput.Model
	prompts    []string
	cursorMode textinput.CursorMode
}

func InitialConfigurationTextInputModel(cps []ConfigPrompt, noCacheFlag bool) ConfigurationTextInputModel {
	var (
		gitCfg, exstCfg string
	)
	m := ConfigurationTextInputModel{
		Inputs: []textinput.Model{},
	}

	var t textinput.Model
	for i := range cps {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 32
		t.Prompt = cps[i].Prompt + ": "
		gitCfg, _ = cps[i].DefaultFunc()
		exstCfg, _ = cps[i].CheckFunc()
		if exstCfg == "" {
			exstCfg = gitCfg
		}
		t.Placeholder = exstCfg

		switch i {
		case 0:
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.CharLimit = 64
		}

		m.Inputs = append(m.Inputs, t)
	}
	m.Inputs[0].Focus()

	return m
}

func (m ConfigurationTextInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ConfigurationTextInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			for i := range m.Inputs {
				m.Inputs[i].SetValue(m.Inputs[i].Placeholder)
			}
			return m, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.Inputs) {
				return m, tea.Quit
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.Inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.Inputs)
			}

			cmds := make([]tea.Cmd, len(m.Inputs))
			for i := 0; i <= len(m.Inputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.Inputs[i].Focus()
					m.Inputs[i].PromptStyle = focusedStyle
					m.Inputs[i].TextStyle = focusedStyle
					continue
				}
				m.Inputs[i].Blur()
				m.Inputs[i].PromptStyle = noStyle
				m.Inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *ConfigurationTextInputModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Inputs))

	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m ConfigurationTextInputModel) View() string {
	var b strings.Builder

	for i := range m.Inputs {
		b.WriteString(m.Inputs[i].View())
		if i < len(m.Inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.Inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("Press enter to select the defaults. Press Esc or Ctrl+c to quit"))

	return b.String()
}

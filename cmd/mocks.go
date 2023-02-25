/*
Copyright © 2022 nanvenomous mrgarelli@gmail.com
*/
package cmd

import (
	"errors"

	"github.com/devsquadron/models"
	"github.com/devsquadron/project-manager/database"
	"github.com/devsquadron/project-manager/security"
	"github.com/devsquadron/project-manager/services"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

var (
	mockDevs  = []string{"nanvenomous", "developerinthemaking", "dirtymoney", "theliteraldev"}
	mockTasks = []models.Task{
		{
			Title:     "get apple developer account",
			Tags:      []string{"business"},
			Developer: "dirtymoney",
			Criterion: `- create an [apple developer account](https://developer.apple.com/programs/enroll/)
- use email "developersquadron@gmail.com" (reach out to nanvenomous for password)
- do research on getting our app "verified" so our users don't see:
    - ![example error](https://support.apple.com/library/content/dam/edam/applecare/images/en_US/macos/Big-Sur/macos-big-sur-alert-unverified-developer.png)`,
		},
		{
			Title: "pursue promotion avenues",
			Tags:  []string{"business"},
			Criterion: `# this is a list to help promote devsquadron (once the time comes)
- [maas lalani](https://github.com/maaslalani)
    - [draw](https://github.com/maaslalani/draw)
    - [slides](https://github.com/maaslalani/slides)
    - [helix](https://www.youtube.com/watch?v=tGYvUXYN-c0)
- [the primagen](https://www.youtube.com/@ThePrimeagen)
    - [burnout](https://www.youtube.com/watch?v=jTmFW1J-KLc&t=2s)
- [tj](https://www.youtube.com/@teej_dv)
- [distro tube](https://www.youtube.com/@DistroTube)
    - would need to be FOSS
- [charmbracelet](https://github.com/charmbracelet)
- [corey butler](https://github.com/coreybutler)
- [y-combinator](https://www.ycombinator.com/)
- [Aaron Jack](https://www.youtube.com/channel/UCRLEADhMcb8WUdnQ5_Alk7g)
    - is a go developer`,
		},
		{
			Title:     "how to use tag",
			Tags:      []string{"documentation"},
			Developer: "developerinthemaking",
			Criterion: "- add documentation on how to use tag\n- you can run `ds tag --help` to get more info on using tags\n- the general idea is\n\t- add a tag with `ds tag 12 dev`\n\t- delete a tag with `ds tag --delete dev 12`",
		},
		{
			Title:     "list by status and developer",
			Tags:      []string{"documentation"},
			Criterion: "- we need to update the documentation for [listing tasks](https://developersquadron.com/getting-started/listing-tasks/) with the additional flags\n\t- try it out! use the command `ds list --status Review`\n- in addition you can now list by developer as well\n\t- try `ds list --developer developerinthemaking`\n\t- `ds info` generates shell completions for this",
		},
		{
			Title:     "ask to join existing team",
			Tags:      []string{"documentation"},
			Criterion: "- on [this page](https://developersquadron.com/install/getting-started/)\n- make it clear if a team is already created then you will not need to run this step\n- for now we are just sending a message to a teammate to be added via command `ds setup team grow <username>`",
		},
		{
			Title:     "user interface",
			Tags:      []string{"development"},
			Criterion: "### Determine library\n\nwe want a library with support for text editor plugins & that can reach the most people the fastest\n\n### desktop & mobile\n- [fyne](https://github.com/fyne-io/fyne)\n- desktop & mobile w/ one codebase\n- can use our current request library from cli\n- terminal emulator [guide](https://ishuah.com/2021/03/10/build-a-terminal-emulator-in-100-lines-of-go/)\n- fyne markdown renderer [pull req](https://github.com/fyne-io/fyne/pull/2301)\n\n### web\n- [beego](https://github.com/beego/beego)\n\t- frontend web framework in go\n- [react](https://github.com/facebook/react)\n- [vim for node](https://www.npmjs.com/package/js-vim)\n- see [vim for wasm](https://github.com/rhysd/vim.wasm)",
		},
		{
			Title:     "indexed search by title",
			Tags:      []string{"development"},
			Developer: "nanvenomous",
			Criterion: "",
		},
		{
			Title:     "refactor flags",
			Tags:      []string{"development"},
			Criterion: "#### similar flags\n- notice the similarity between\n\t- `ds-cli/dscmd/list.go` `listDevFlag`\n\t- `ds-cli/dscmd/assign.go` `assignDevFlag`\n\n#### directory `ds-cli/dscmd`\n- run `ls ds-cli/dscmd`\n- get all flags from those commands\n#### make a new `system/flags.go`\n- move all the flag logic there\n- then then commands themselves will just implement the flags\n- this should help with flag consistency & less typing :)",
		},
		{
			Title:     "ask to join team",
			Tags:      []string{"development"},
			Criterion: "- add the text editor error link\n- from devsqudron troubleshot page",
		},
		{
			Title:     "context based listing",
			Tags:      []string{"development"},
			Criterion: "# `ds-cli/dscmd/list.go`\n\n####  `for _, tsk = range *allTsks {`\n- here we are listing all the tasks we got from query\n- the come in creation order by default\n- if we are listing by `tag` or `developer` then create some functionality\n\t- in `message/message.go` to list in status order\n\t- & that selectively removes listing things on a per task basis if they are all the same\n#### `fmt.Println(________________________________________________)`\n- move the page break functionality to `message/message.go`",
		},
		{
			Title:     "command completions for developer",
			Tags:      []string{"development"},
			Developer: "nanvenomous",
			Criterion: "#### `ds-cli/dscmd/assign.go`\n- flag `assignDevFlag` needs completions similar to:\n\t- `ds-cli/dscmd/list.go` `RegisterFlagCompletionFunc`\n#### `ds-cli/dscmd/list.go`\n- the flag `listDevFlag` also needs the same completions",
		},
		{
			Title:     "info command",
			Tags:      []string{"development"},
			Developer: "nanvenomous",
			Criterion: "- create a `ds info` command\n- to start it should list:\n\t- devsquadron username\n\t- devsquadron email\n\t- team\n\t- teammates\n\t- public / private\n\t- current task count",
		},
		{
			Title:     "convert release ci to go code",
			Tags:      []string{"development"},
			Developer: "nanvenomous",
			Criterion: "- convert this method to go code\n- put in `scripts/release/main.go`",
		},
	}
)

func AddAllMockData() error {
	var (
		err    error
		clnt   *mongo.Client
		dtb    *mongo.Database
		ctx    = context.Background()
		devDtb *database.DeveloperDatabase
		tskDtb *database.TaskDatabase
		tmDtb  *database.TeamDatabase
	)

	clnt, dtb, tskDtb, devDtb, tmDtb, err = services.GetDatabases(ctx, prod, quiet, security.NewSecurity())
	if err != nil {
		return err
	}

	if !prod {
		err = dtb.Drop(ctx)
		if err != nil {
			return err
		}
	} else {
		return errors.New("cannot add mock data to production!")
	}

	_, err = tmDtb.CreateTeam(&models.Team{Name: "devsquadron", Developers: []string{}})
	if err != nil {
		return err
	}

	for _, dv := range mockDevs {
		_, _, err = devDtb.CreateDeveloper(&models.Developer{Name: dv, Password: "superawesomepassword", ConfirmPassword: "superawesomepassword"})
		if err != nil {
			return err
		}
	}

	for _, dv := range mockDevs {
		growDev, err := devDtb.GetDevFromName(dv)
		if err != nil {
			return err
		}

		growDev.Teams = append(growDev.Teams, "devsquadron")
		_, err = devDtb.UpdateDeveloper(growDev)
		if err != nil {
			return err
		}

		_, err = tmDtb.GrowTeam("devsquadron", dv)
		if err != nil {
			return err
		}
	}

	for i, tk := range mockTasks {
		res, err := tskDtb.CreateTask(&tk, "devsquadron")
		if err != nil {
			return err
		}
		mockTasks[i].ID = res
	}

	_, err = tskDtb.GetTasks("devsquadron", "", []string{}, "")
	if err != nil {
		return err
	}

	mockTasks[10].Status = models.Status.Review
	mockTasks[11].Status = models.Status.Developing
	_, err = tskDtb.UpdateTask("devsquadron", &mockTasks[10])
	if err != nil {
		return err
	}
	_, err = tskDtb.UpdateTask("devsquadron", &mockTasks[11])
	if err != nil {
		return err
	}

	return clnt.Disconnect(ctx)
}

var mocksCmd = &cobra.Command{
	Use:   "mocks",
	Short: "creates the mocks for full integration tests",
	RunE: func(cmd *cobra.Command, args []string) error {
		return AddAllMockData()
	},
}

func init() {
	rootCmd.AddCommand(mocksCmd)
}

package configuration

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/devsquadron/ds/constants"

	"github.com/devsquadron/ds/exception"

	"github.com/nanvenomous/exfs"
	"gopkg.in/yaml.v3"
)

const (
	EDITOR_ENV_VAR                = "EDITOR"
	tempFileName                  = "devsquadron-task-*.md"
	PROJECT_CONFIG_FILE_NAME      = ".devsquadron.yml"
	PROJECT_CONFIG_DIRECTORY_NAME = "devsquadron"
)

var (
	errCompCache  error
	errGlobalCfg  error
	errProjectCfg error
)

func (cg *Configuration) makeDsOsFile(
	osDir string,
	flNm string,
) (string, error) {
	var (
		err       error
		osFilePth string
	)
	err = os.MkdirAll(osDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	osFilePth = path.Join(osDir, flNm)

	if _, err = os.Stat(osFilePth); os.IsNotExist(err) {
		fl, err := os.Create(osFilePth)
		if err != nil {
			return "", err
		}
		err = fl.Close()
		if err != nil {
			return "", err
		}
	}
	return osFilePth, nil
}

type UserDirectories struct {
	Config string
	Cache  string
}

var userDirs UserDirectories

func (cg *Configuration) getDefaultGlobalConfig() (string, string, error) {
	var (
		err                error
		cfgFile, cacheFile string
	)

	// getOsConfig := func() error {
	// 	userDirs.Config, err = os.UserConfigDir()
	// 	return err
	// }

	// err = exfs.RunOn(&exfs.OperatingSystemRoute{
	// 	Windows: func() error {
	// hm, err := os.UserHomeDir()
	// if err != nil {
	// return err
	// }
	// userDirs.Config = filepath.Join(hm, ".config")
	// 		return nil
	// 	},
	// 	Linux: getOsConfig,
	// 	Mac:   getOsConfig,
	// })
	userDirs.Config, err = os.UserConfigDir()
	if err != nil {
		return "", "", err
	}

	cfgFile, err = cg.makeDsOsFile(
		filepath.Join(userDirs.Config, PROJECT_CONFIG_DIRECTORY_NAME),
		constants.GLOBAL_CONFIG,
	)
	if err != nil {
		return "", "", err
	}

	userDirs.Cache, err = os.UserCacheDir()
	cacheFile, err = cg.makeDsOsFile(
		filepath.Join(userDirs.Cache, PROJECT_CONFIG_DIRECTORY_NAME, cg.Project.Team),
		constants.CACHE_COMPLETIONS,
	)
	if err != nil {
		return "", "", err
	}

	return cfgFile, cacheFile, nil
}

type CompletionCache struct {
	Tags       []string `yaml:"tags"`
	Developers []string `yaml:"developers"`
}

type ProjectConfig struct {
	Team string `yaml:"team"`
}

type GlobalConfig struct {
	Token    string `yaml:"token"`
	Username string `yaml:"username"`
	Email    string `yaml:"email"`
}

type ConfigurationType interface {
	Setup(string)
	Write() error
	Tags() []string
	SetTags([]string)
	Developers() []string
	SetDevelopers([]string)
	Token() string
	Team() string
	TeamE() (string, error)
	SetTeam(string)
	SetToken(string)
	TaskNumber() (uint64, error)
	UsernameE() (string, error)
	Username() string
	SetUsername(string)
	EmailE() (string, error)
	Email() string
	SetEmail(string)
	PrintAllPaths()
}

type Configuration struct {
	FS                  *exfs.FileSystem
	Project             *ProjectConfig
	Global              *GlobalConfig
	CompletionCache     *CompletionCache
	PathGlobal          string
	PathCompletionCache string
	PathProject         string
}

func NewConfiguration(fs *exfs.FileSystem) ConfigurationType {
	return &Configuration{
		FS:              fs,
		Project:         &ProjectConfig{},
		Global:          &GlobalConfig{},
		CompletionCache: &CompletionCache{},
	}
}

func (cg *Configuration) hydrateCompletionCache(pth string) error {
	var (
		err      error
		fileByte []byte
	)

	fileByte, err = ioutil.ReadFile(pth)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(fileByte, cg.CompletionCache)
	if err != nil {
		return err
	}
	cg.PathCompletionCache = pth
	return nil
}

func (cg *Configuration) hydrateGlobalConfig(cfgPth string) error {
	var (
		err      error
		fileByte []byte
	)

	fileByte, err = ioutil.ReadFile(cfgPth)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(fileByte, cg.Global)
	if err != nil {
		return err
	}

	cg.PathGlobal = cfgPth
	return nil
}

func (cg *Configuration) hydrateProjectConfig() error {
	var (
		err      error
		filePath string
		fileByte []byte
		cfg      ProjectConfig
	)
	filePath, err = cg.FS.FindFileInAboveCurDir(PROJECT_CONFIG_FILE_NAME)
	if err != nil {
		return errors.New(
			fmt.Sprintf(
				"%s\n%s",
				err.Error(),
				"for more info see https://developersquadron.com/troubleshooting/project-config/",
			),
		)
	}

	fileByte, err = ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(fileByte, &cfg)
	if err != nil {
		return err
	}

	cg.Project = &cfg
	cg.PathProject = filePath
	return nil
}

func (cg *Configuration) Setup(cfgFile string) {
	var (
		defCfgFl, compCacheFl string
		err                   error
	)

	err = cg.hydrateProjectConfig()
	if err != nil {
		errProjectCfg = err
	}

	defCfgFl, compCacheFl, err = cg.getDefaultGlobalConfig()
	if err != nil {
		errCompCache = err
		errGlobalCfg = err
	}
	if cfgFile == "" {
		cfgFile = defCfgFl
	}

	err = cg.hydrateCompletionCache(compCacheFl)
	if err != nil {
		errCompCache = err
	}

	err = cg.hydrateGlobalConfig(cfgFile)
	if err != nil {
		errGlobalCfg = err
	}
}

func (cg *Configuration) Write() error {
	var (
		err                                          error
		newGlbCfgByt, newCompCacheByt, newProjCfgByt []byte
	)
	if errCompCache == nil {
		newCompCacheByt, err = yaml.Marshal(&cg.CompletionCache)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(cg.PathCompletionCache, newCompCacheByt, 0777)
		if err != nil {
			return err
		}
	} else {
		return errCompCache
	}

	if errGlobalCfg == nil {
		newGlbCfgByt, err = yaml.Marshal(&cg.Global)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(cg.PathGlobal, newGlbCfgByt, 0777)
		if err != nil {
			return err
		}
	} else {
		return errGlobalCfg
	}

	if cg.PathProject != "" {
		newProjCfgByt, err = yaml.Marshal(&cg.Project)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(cg.PathProject, newProjCfgByt, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cg *Configuration) Developers() []string {
	exception.CheckErr(errCompCache)
	return cg.CompletionCache.Developers
}

func (cg *Configuration) SetDevelopers(devs []string) {
	exception.CheckErr(errCompCache)
	cg.CompletionCache.Developers = devs
}

func (cg *Configuration) Tags() []string {
	exception.CheckErr(errCompCache)
	return cg.CompletionCache.Tags
}

func (cg *Configuration) SetTags(tgs []string) {
	exception.CheckErr(errCompCache)
	cg.CompletionCache.Tags = tgs
}

func (cg *Configuration) Token() string {
	exception.CheckErr(errGlobalCfg)
	return cg.Global.Token
}

func (cg *Configuration) SetToken(tkn string) {
	exception.CheckErr(errGlobalCfg)
	cg.Global.Token = tkn
}

func (cg *Configuration) Team() string {
	exception.CheckErr(errProjectCfg)
	return cg.Project.Team
}

func (cg *Configuration) TeamE() (string, error) {
	if errProjectCfg != nil {
		return "", errProjectCfg
	}
	return cg.Project.Team, nil
}

func (cg *Configuration) SetTeam(tm string) {
	if errProjectCfg != nil {
		fl, err := os.Create(PROJECT_CONFIG_FILE_NAME)
		exception.CheckErr(err)
		exception.CheckErr(fl.Close())
		cg.PathProject = PROJECT_CONFIG_FILE_NAME
		errProjectCfg = nil
		fmt.Println()
		fmt.Println(
			fmt.Sprintf("...created '%s' in current directory.", PROJECT_CONFIG_FILE_NAME),
		)
	}
	cg.Project.Team = tm
}

func (cg *Configuration) TaskNumber() (uint64, error) {
	var (
		err           error
		outStr, errs  string
		taskNumberReg *regexp.Regexp
		taskNumberStr string
	)
	outStr, errs, err = cg.FS.Capture("git", []string{"rev-parse", "--abbrev-ref", "HEAD"})
	if err != nil {
		fmt.Println(errs)
		return 0, err
	}

	taskNumberReg, err = regexp.Compile("^\\d+")
	if err != nil {
		return 0, err
	}

	taskNumberStr = taskNumberReg.FindString(strings.TrimSpace(outStr))
	if taskNumberStr == "" {
		return 0, errors.New("Warning, no task number associated in the story. Not automating your progress.")
	}

	return strconv.ParseUint(taskNumberStr, 10, 64)
}

func (cg *Configuration) UsernameE() (string, error) {
	if errGlobalCfg != nil {
		return "", errGlobalCfg
	}
	if cg.Global.Username == "" {
		return "", errors.New("username is not set, run 'ds init' to set")
	}
	return cg.Global.Username, nil
}

func (cg *Configuration) Username() string {
	exception.CheckErr(errGlobalCfg)
	if cg.Global.Username == "" {
		exception.CheckErr(errors.New("username is not set, run 'ds init' to set"))
	}
	return cg.Global.Username
}

func (cg *Configuration) SetUsername(unm string) {
	exception.CheckErr(errGlobalCfg)
	cg.Global.Username = unm
}

func (cg *Configuration) EmailE() (string, error) {
	if errGlobalCfg != nil {
		return "", errGlobalCfg
	}
	if cg.Global.Email == "" {
		return "", errors.New("email is not set, run 'ds init'")
	}
	return cg.Global.Email, nil
}

func (cg *Configuration) Email() string {
	exception.CheckErr(errGlobalCfg)
	if cg.Global.Email == "" {
		exception.CheckErr(errors.New("email is not set, run 'ds init'"))
	}
	return cg.Global.Email
}

func (cg *Configuration) SetEmail(unm string) {
	exception.CheckErr(errGlobalCfg)
	cg.Global.Email = unm
}

func (cg *Configuration) PrintAllPaths() {
	fmt.Println("Global: ", cg.PathGlobal)
	fmt.Println("Cache: ", cg.PathCompletionCache)
	fmt.Println("Project: ", cg.PathProject)
}

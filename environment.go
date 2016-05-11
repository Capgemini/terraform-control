package main

import (
	"path/filepath"
	"github.com/mitchellh/cli"
	"custom/terraform-control/terraform"
	"github.com/hashicorp/otto/ui"
	"os"
	"fmt"
	"log"
	"sync"
	"io/ioutil"
	"strconv"
	"os/exec"
	)

const (
	ErrorPrefix  = "e:"
	OutputPrefix = "o:"
	)

var (
	safeEnvironments = make(map[int]*SafeEnvironment)
	ouputFile 	=	"output" 
	stateFile 	=	"state"
	)

type Environment struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Repo      string    `json:"repo"`
	Branch 	  string    `json:"branch"`
	Path      string 	`json:"path"`
	AutoApply bool		`json:"autoApply"`
	//TODO: Handle variables dynamically
	Var1 	string		`json:"var1"`
	Val1 	string		`json:"val1"`
	Var2 	string		`json:"var2"`
	Val2 	string		`json:"val2"`
	Changes	  []*Change
}

type Environments []Environment

type SafeEnvironment struct {
	sync.Mutex
	Id int
}

func NewSafeEnvironment(id int) (*SafeEnvironment){
	return &SafeEnvironment{
		Id: id,
	}
}

func GetSingletonSafeEnvironment(id int)(*SafeEnvironment){
    if (safeEnvironments[id] == nil) {
		safeEnvironments[id] = NewSafeEnvironment(id)
    } 
    return safeEnvironments[id]
}

func (e *Environment) GetPathToRepo()(string) {
	return filepath.Join(config.RootFolder, e.Name)
}

func (e *Environment) GetPathToFiles()(string) {
	return filepath.Join(e.GetPathToRepo(), e.Path)
}

func (e *Environment) GetPathToOuput()(string) {
	return filepath.Join(e.GetPathToFiles(), ouputFile)
}

func (e *Environment) GetPathToState()(string) {
	return filepath.Join(e.GetPathToFiles(), stateFile)
}

func (e *Environment) createUi() (ui.Ui) {
	cliUi := &cli.ColoredUi{
		OutputColor: cli.UiColorNone,
		InfoColor:   cli.UiColorNone,
		ErrorColor:  cli.UiColorRed,
		WarnColor:   cli.UiColorNone,
		Ui: &cli.PrefixedUi{
			AskPrefix:    OutputPrefix,
			OutputPrefix: OutputPrefix,
			InfoPrefix:   OutputPrefix,
			ErrorPrefix:  ErrorPrefix,
			Ui:           &cli.BasicUi{Writer: os.Stdout},
		},
	}

	tfUi := NewUi(cliUi, e)
	return tfUi
}

func (se *SafeEnvironment) Execute(change *Change, action *Action) (error) {
    // Agressive locking as we want the same environment to be manipulated only once at a time 
    se.Lock()
    defer se.Unlock()

	env := RepoFindEnvironment(se.Id)
	command := action.Command
	pathToRepo := env.GetPathToRepo()
	pathToOuput := env.GetPathToOuput()
	pathToState := env.GetPathToState()

	if command == "plan" {
		env.Changes = append(env.Changes, change)
		derr := config.Persistence.PutEnvironment(env)
		if derr != nil {
			log.Fatal(derr)
		}

		changesChannel := getChangesChannel()
		changesChannel <- se.Id
	}

	// TODO: Think about allowing apply any change/rollback.
	// Hacky Hacky, If running apply assume only last change can be applied
	if change == nil {
		change = env.Changes[len(env.Changes)-1]
	}

    if err := env.Execute(change, command); err != nil {
		change.Status = action.FailCode
	} else {
		change.Status = action.SuccessCode
	}

	// TODO: consider a better way of doing this by buffering or something
	// I cant be bothered today as I'm feeling so sick :O
	planOuputContent, err := ioutil.ReadFile(pathToOuput)
		if err != nil {
		    log.Fatal(err)
		}
	change.PlanOutput = string(planOuputContent)

	if command == "apply" {
		stateFileContent, err := ioutil.ReadFile(pathToState)
		if err != nil {
		    log.Fatal(err)
		}
		change.State = string(stateFileContent)		
	}

	// Override last change with new info
	env.Changes[len(env.Changes)-1] = change

	derr := config.Persistence.PutEnvironment(env)
	if derr != nil {
		log.Fatal(derr)
	}
	changesChannel <- se.Id

	os.RemoveAll(pathToRepo)
	return nil
}

func (e *Environment) Execute(change *Change, command ...string) error {

	pathToRepo := e.GetPathToRepo()
    pathToFiles := e.GetPathToFiles()

	//TODO: handle variables dynamically 
	vars := make(map[string]string)
	vars[e.Var1] = e.Val1
	vars[e.Var2] = e.Val2

	// TODO: Us git2go here
	// Clone the project
	cmd := exec.Command("git", "clone", pathToRepo)
	cmd.Dir = config.RootFolder
	err := cmd.Run()
	if err != nil {
		log.Printf("Error cloning: %v", err)
	}

	// Checkout to Headcommit
	commit := (change.HeadCommit.(map[string]interface{})["id"]).(string)
	cmd = exec.Command("git", "checkout", commit)
	cmd.Dir = pathToRepo
	err = cmd.Run()
	if err != nil {
		log.Printf("Error checking out: %v", err)
	}

	tfUi := e.createUi()

	tf := &terraform.Terraform{
		Path:      "",
		Dir:       pathToFiles,
		Ui:        tfUi,
		Variables: vars,
		Directory: config.Persistence,
		StateId:   "env-" + strconv.Itoa(e.Id),
	}

	tfUi.Header("Executing Terraform to manage infrastructure...")
	tfUi.Message("Raw Terraform output will begin streaming in below.")

	// Start the Terraform command
	err = tf.Execute(command...)
	if err != nil {
		err = fmt.Errorf("Error running Terraform: %s", err)
		log.Printf("Error running terraform: %v", err)
		return err
	}

	tfUi.Header("Terraform execution complete. Saving results...")

	return nil
}
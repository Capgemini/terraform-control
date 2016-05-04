package main

import (
	"time"
	"path/filepath"
	"github.com/mitchellh/go-homedir"
	"github.com/hashicorp/otto/directory"
	"github.com/mitchellh/cli"
	"custom/terraform-control/terraform"
	"github.com/libgit2/git2go"
	"os"
	"fmt"
	"log"
	"sync"
	"io/ioutil"
	"strconv"
	)

var safeEnvironments = make(map[int]*SafeEnvironment)

type Environment struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Repo      string    `json:"repo"`
	Branch 	  string    `json:"branch"`
	Path      string 	`json:"path"`
	modified  time.Time `json:"modified"`
	AutoApply bool		`json:"autoApply"`
	//TODO: number of variables dynamically
	Var1 string		`json:"var1"`
	Val1 string		`json:"val1"`
	Var2 string		`json:"var2"`
	Val2 string		`json:"val2"`
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

func (se *SafeEnvironment) Execute(change *Change, command string, status int) (error) {
    se.Lock()
	env := RepoFindEnvironment(se.Id)
    pathToFiles := filepath.Join(GetDataFolder(), "/repo-" + env.Name, env.Path)
	//TODO: Think about allowing apply any change/rollback.
	// If running apply assume only last change can be applied
	if (change == nil) {
		change = env.Changes[len(env.Changes)-1]
	}

    if err := env.Execute(change, command); err != nil {
		change.Status = 100
	} else {
		change.Status = status
	}

	planOutputFile := filepath.Join(pathToFiles, "/planOutput")
	planOuputContent, err := ioutil.ReadFile(planOutputFile)
	if err != nil {
		log.Printf("No output file found: %v", planOutputFile)
	    log.Fatal(err)
	}
	// TODO: consider a better way of doing this by buffering or something
	// I cant be bothered today as I'm feeling so sick :O
	change.PlanOutput = string(planOuputContent)

	if command == "apply" {
		stateFile := filepath.Join(pathToFiles, "/state")
		stateFileContent, err := ioutil.ReadFile(stateFile)
		if err != nil {
		    log.Fatal(err)
		}
		change.State = string(stateFileContent)		
	}

	db := &BoltBackend{
		Dir: filepath.Join(GetDataFolder(), "data"),
	}
	env.Changes = append(env.Changes, change)
	derr := db.PutEnvironment(env)
	if derr != nil {
		log.Fatal(derr)
	}

	os.RemoveAll(GetDataFolder()+ "/repo-" + env.Name)
	time.Sleep(5*time.Second)
	se.Unlock()
	return nil
}

func GetDataFolder()(string) {
	dataFolder, err := homedir.Expand("~/.terraform-control")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Something went wrong!!!!!: %s", err)
	}
	return dataFolder
}

func (e *Environment) Execute(change *Change, command ...string) error {

	// Build the variables
	dataFolder, err := homedir.Expand("~/.terraform-control")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Something went wrong!!!!!: %s", err)
	}

	//TODO: handle variables dynamically 
	vars := make(map[string]string)
	vars[e.Var1] = e.Val1
	vars[e.Var2] = e.Val2

	gitRepo := e.Repo
	repoPath := dataFolder + "/repo-" + e.Name
	//credentials
	repo, err := git.OpenRepository(repoPath)
 	if err != nil {
        fmt.Println("repo does not exist: creating one ")
        repo, err = git.Clone(gitRepo, repoPath, &git.CloneOptions{})
        if err != nil {
            panic(err)
        }
	}
	defer repo.Free()

	commit := (change.HeadCommit.(map[string]interface{})["id"]).(string)

	oid, err := git.NewOid(commit)
	if err != nil {
		panic(err)
	}

	changeCommit, err := repo.LookupCommit(oid)
	if err != nil {
		panic(err)
	}

	err = repo.ResetToCommit(changeCommit, git.ResetSoft, &git.CheckoutOpts{})
	if err != nil {
		panic(err)
	}

	// Build the context
	dataDir := dataFolder
	directory := &directory.BoltBackend{
		Dir: filepath.Join(dataDir, "data"),
	}

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

	// Build our executor
	tf := &terraform.Terraform{
		Path:      "",
		Dir:       filepath.Join(GetDataFolder(), "/repo-" + e.Name, e.Path),
		Ui:        tfUi,
		Variables: vars,
		Directory: directory,
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
		//infra.State = directory.InfraStatePartial
	}

	tfUi.Header("Terraform execution complete. Saving results...")

	return nil
}
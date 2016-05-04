package main

// import "fmt"
import "path/filepath"
import "log"

var currentId int

var environments Environments


// Give us some seed data
func init() {
	// RepoCreateTodo(Todo{Name: "Write presentation"})
	// RepoCreateTodo(Todo{Name: "Host meetup"})
}

func RepoIndexEnvironments() []*Environment {
	db := &BoltBackend{
	Dir: filepath.Join(GetDataFolder(), "data"),
	}
	e, derr := db.GetAllEnvironments()
	if derr != nil {
		log.Fatal(derr)
	}
	return e
}

func RepoFindEnvironment(id int) *Environment {
	db := &BoltBackend{
	Dir: filepath.Join(GetDataFolder(), "data"),
	}
	e, derr := db.GetEnvironment(id)
	if derr != nil {
		log.Fatal(derr)
	}
	return e
}

//this is bad, I don't think it passes race condtions
func RepoCreateEnvironment(e Environment) Environment {
	db := &BoltBackend{
		Dir: filepath.Join(GetDataFolder(), "data"),
	}
	derr := db.PutEnvironment(&e)
	if derr != nil {
		log.Fatal(derr)
	}
	return e
}

//this is bad, I don't think it passes race condtions
func RepoCreateChange(c Change) Change {
	db := &BoltBackend{
		Dir: filepath.Join(GetDataFolder(), "data"),
	}
	derr := db.PutChange(&c)
	if derr != nil {
		log.Fatal(derr)
	}
	return c
}

//this is bad, I don't think it passes race condtions
func RepoTerraformAction(action Action) error {
	safeEnvironment := GetSingletonSafeEnvironment(action.Id)
	// TODO: Consider similar approach to http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html
	go safeEnvironment.Execute(nil, action.Action, 2)
	return nil
}

// func RepoDestroyTodo(id int) error {
// 	for i, t := range todos {
// 		if t.Id == id {
// 			todos = append(todos[:i], todos[i+1:]...)
// 			return nil
// 		}
// 	}
// 	return fmt.Errorf("Could not find Todo with id of %d to delete", id)
// }
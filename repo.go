package main

import "log"

func RepoIndexEnvironments() []*Environment {
	e, derr := config.Persistence.GetAllEnvironments()
	if derr != nil {
		log.Fatal(derr)
	}
	return e
}

func RepoFindEnvironment(id int) *Environment {
	e, derr := config.Persistence.GetEnvironment(id)
	if derr != nil {
		log.Fatal(derr)
	}
	return e
}

//RepoCreateEnvironment - this is bad, I don't think it passes race condtions
func RepoCreateEnvironment(e Environment) Environment {
	derr := config.Persistence.PutEnvironment(&e)
	if derr != nil {
		log.Fatal(derr)
	}
	return e
}

func RepoHookHandler(c Change) Change {
	derr := c.handleHook(config.Persistence)

	if derr != nil {
		log.Fatal(derr)
	}
	return c
}

func RepoTerraformAction(action Action) error {
	safeEnvironment := GetSingletonSafeEnvironment(action.Id)
	// TODO: Consider similar approach to http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html
	go safeEnvironment.Execute(nil, action.SetExitCodes())
	return nil
}

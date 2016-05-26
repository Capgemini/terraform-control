package main

import "testing"
import "github.com/mitchellh/go-homedir"
import "fmt"
import "os"

func TestGetPathToRepo(t *testing.T) {
	e := &Environment{
		Name: "testEnv",
	}

	want, err := homedir.Expand("~/.terraform-control/testEnv")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Something went wrong when retrieving the data folder!!!: %s", err)
	}

	got := e.GetPathToRepo()

	if got != want {
		t.Errorf("GetPathToRepo() == %q, want %q", got, want)
	}
}

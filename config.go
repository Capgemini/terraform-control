package main

import (
	"path/filepath"
	"github.com/mitchellh/go-homedir"
	"os"
	"fmt"
	)

const (
	rootFolder = "~/.terraform-control"
	dataFolder = "data"
	)

type Config struct {
	Persistence   *BoltBackend
	RootFolder    string	
}

var c *Config

func init() {
	c = &Config{
		Persistence: getPersistenceBackend(),
		RootFolder: getRootFolder(),
	}	
}

func getRootFolder()(string) {
	rootFolder, err := homedir.Expand(rootFolder)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Something went wrong when retrieving the data folder!!!: %s", err)
	}
	return rootFolder
}

func getPersistenceBackend()(*BoltBackend) {
	db := &BoltBackend{
		Dir: filepath.Join(getRootFolder(), dataFolder),
	}
	return db
}

func GetConfig()(*Config) {
    return c
}


package main

import (
	"github.com/boltdb/bolt"
	"log"
	)

type Change struct {
	Repository	map[string]interface{}	`json:"repository"`
	Commits		[]interface{}	`json:"commits"`
	HeadCommit	interface{}	`json:"head_commit"`
	PlanOutput	string
	State		string
	Status		int
}

type Changes []Change

func (change *Change) handleHook(b *BoltBackend) error {
	db, err := b.db()
	if err != nil {
		return err
	}
	defer db.Close()

	return db.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket(boltEnvironmentsBucket)
		bucket, err = bucket.CreateBucketIfNotExists([]byte(
			boltEnvironmentsBucket))
		if err != nil {
			return err
		}

	    c := bucket.Cursor()
		var env *Environment
	    for k, v := c.First(); k != nil; k, v = c.Next() {
	        env = &Environment{}
	        err := b.structRead(env, v)
	        if err != nil {
				return err
	        }

	        if (env.Repo == change.Repository["ssh_url"] || env.Repo == change.Repository["git_url"] || env.Repo == change.Repository["git_url"] || env.Repo == change.Repository["html_url"]) {
                log.Printf("Triggering environment changes for repo: %v", env.Repo)
		        safeEnvironment := GetSingletonSafeEnvironment(env.Id)
				// TODO: Consider similar approach to http://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html
                action := &Action {
                	Command: "plan",

                }
                go safeEnvironment.Execute(change, action.SetExitCodes())
	        }
	    }
		return nil
	})
}

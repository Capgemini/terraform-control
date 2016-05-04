package main

type Change struct {
	Repository	map[string]interface{}	`json:"repository"`
	Commits		[]interface{}	`json:"commits"`
	HeadCommit	interface{}	`json:"head_commit"`
	PlanOutput	string
	State		string
	Status		int
}

type Changes []Change
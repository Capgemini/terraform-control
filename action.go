package main

const (
	PLAN_SUCCESS   	= 1
	PLAN_FAIL   	= 100
	APPLY_SUCCESS   = 2
	APPLY_FAIL		= 200
	REFRESH = 3
)

type Action struct {
	Id 				int       `json:"id"`
	Command      	string    `json:"action"`
	SuccessCode		int
	FailCode		int
}

func (a *Action) SetExitCodes()(*Action) {
	switch a.Command {
	case "apply":
		a.SuccessCode = APPLY_SUCCESS
		a.SuccessCode = APPLY_FAIL
	case "plan":
		a.SuccessCode = PLAN_SUCCESS
		a.SuccessCode = PLAN_FAIL
	}

	return a
}
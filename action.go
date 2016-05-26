package main

// Exported constants
const (
	PlanSuccess  = 1
	PlanFail     = 100
	ApplySuccess = 2
	ApplyFail    = 200
	Refresh      = 3
)

// Action exporting exitCodes for use later
type Action struct {
	ID          int    `json:"id"`
	Command     string `json:"action"`
	SuccessCode int
	FailCode    int
}

// SetExitCodes to use the Action for apply or plan
func (a *Action) SetExitCodes() *Action {
	switch a.Command {
	case "apply":
		a.SuccessCode = ApplySuccess
		a.SuccessCode = ApplyFail
	case "plan":
		a.SuccessCode = PlanSuccess
		a.SuccessCode = PlanFail
	}

	return a
}

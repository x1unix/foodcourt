package auth

const (
	ActionRequestCode = "REQUEST_CODE"
	ActionSubmitCode  = "SUBMIT_CODE"
	ActionSetPassword = "RESET_PASSWD"
)

type RecoveryAction struct {
	Action string            `json:"action"`
	Data   map[string]string `json:"data"`
}

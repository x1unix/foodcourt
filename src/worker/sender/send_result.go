package sender

type SendResult struct {
	GoRoutine int
	Email string
	Success bool
	Error string
}

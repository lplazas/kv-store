package errs

type NodeNotAvailable struct {
	message string
}

func (e NodeNotAvailable) Error() string {
	return e.message
}

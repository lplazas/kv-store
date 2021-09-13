package errs

type ValueNotFound struct {
	message string
}

func (e ValueNotFound) Error() string {
	return e.message
}

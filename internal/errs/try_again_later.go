package errs

type TryAgainLater struct {
	message string
}

func TryAgainLaterError(message string) error {
	return TryAgainLater{message: message}
}

func (e TryAgainLater) Error() string {
	return e.message
}

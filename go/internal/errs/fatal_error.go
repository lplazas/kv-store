package errs

import "fmt"

type fatalError struct {
	message    string
	backingErr error
}

func FatalError(message string, err error) error {
	return fatalError{
		message:    message,
		backingErr: err,
	}
}

func (e fatalError) Error() string {
	return fmt.Errorf("%s, err: %w", e.message, e.backingErr).Error()
}

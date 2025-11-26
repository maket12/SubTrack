package uc_errors

type WrappedError struct {
	Public error
	Reason error
}

func (e *WrappedError) Error() string {
	return e.Public.Error()
}

func Wrap(public, reason error) error {
	return &WrappedError{
		Public: public,
		Reason: reason,
	}
}

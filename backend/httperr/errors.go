package httperr

type BadRequestError struct {
	message string
}

func (e *BadRequestError) Error() string {
	return e.message
}

func NewBadRequestError(msg string) error {
	return &BadRequestError{message: msg}
}

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

type UnauthorizedError struct {
	message string
}

func (e *UnauthorizedError) Error() string {
	return e.message
}

func NewUnauthorizedError(msg string) error {
	return &UnauthorizedError{message: msg}
}

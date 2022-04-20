package errors

type Error struct {
	Code    int    `json:"code" yaml:"code"`
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
	Err     error  `json:"-" yaml:"-"`
}

func (e *Error) Error() string {
	if e.Err != nil {
		e.Err = ErrInternalServerError
	}

	return e.Err.Error()
}

func (e *Error) Unwrap() error {
	return e.Err
}

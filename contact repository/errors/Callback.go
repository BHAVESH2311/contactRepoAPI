package errors

// CallbackError Represents callback type error.
type CallbackError struct {
	Err string `json:"error" example:"After create failed"`
}

// Error implements error interface.
func (e CallbackError) Error() string {
	return e.Err
}

// NewCallbackError returns new instance of Callback error.
func NewCallbackError(error string) *CallbackError {
	return &CallbackError{
		Err: error,
	}
}

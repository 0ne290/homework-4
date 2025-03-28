package shared

type InvariantViolationError struct {
	Message string
}

var NilOfInvariantViolationError *InvariantViolationError

func (e *InvariantViolationError) Error() string {
	return e.Message
}

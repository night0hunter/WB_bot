package myError

type ErrorType int

const (
	DateInputError ErrorType = iota + 1
	WarehouseInputError
	CoeffInputError
	SupplyTypeError
	TrackingChoiceError
	ActionChoiceError
)

type MyError struct {
	ErrType ErrorType
	Message string
}

func (e *MyError) Error() string {
	return e.Message
}

func (e *MyError) GetErrorType() ErrorType {
	return e.ErrType
}

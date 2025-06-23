package myError

type ErrorType int

const (
	DefaultError ErrorType = iota
	DateInputError
	WarehouseInputError
	CoeffInputError
	SupplyTypeError
	TrackingChoiceError
	ActionChoiceError
	SaveStatusChoiceError
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

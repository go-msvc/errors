package errors

type Validator interface {
	Validate() error
}

package validator

type Validator interface {
	Validate(val interface{}) *ValidationErrors
}

package validator

type Validator interface {
	IsValid() error
}

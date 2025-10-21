package validator

type Validable interface {
	IsValid() error
}

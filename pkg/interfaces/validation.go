package interfaces

type Validatable interface {
	Validate() error
	Transform() error
}

package cursor

type Validator[T any] func(value T) bool

func ValidateNonNil(value interface{}) bool {
	return value != nil
}

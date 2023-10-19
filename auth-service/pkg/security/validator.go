package security

import "github.com/go-playground/validator/v10"

var val *validator.Validate

func initValidator() {
	val = validator.New()
}

func (s security) Validate(src any) error {
	err := val.Struct(&src)
	if err != nil {
		return err // TODO: wrap error
	}

	return nil
}

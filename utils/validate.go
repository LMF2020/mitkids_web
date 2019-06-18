package utils

import "gopkg.in/go-playground/validator.v9"

var validate *validator.Validate

func ValidStruct(obj interface{}) error {

	err := validator.New().Struct(obj)
	if err != nil {
		return err
	}

	return nil
}

package model

import "fmt"

type ServiceError struct {
	errorCode   int
	errorString string
}

func (e ServiceError) Error() string {
	return fmt.Sprintf(e.errorString, e.errorCode)
}

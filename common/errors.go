package common

import "fmt"

type CustomError struct {
	Message string
	Code    int
	Details string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("Error %d: %s - %s", e.Code, e.Message, e.Details)
}

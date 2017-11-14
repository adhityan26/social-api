package helper

import (
	"fmt"
	"net/http"
)

func MandatoryErrorMessage(field string, messages *[]string) (int, bool)  {
	*messages = append(*messages, fmt.Sprintf("%s cannot be empty", field))
	return  http.StatusPreconditionRequired, false
}

func UndefinedErrorMessage(message string, messages *[]string) (int, bool)  {
	*messages = append(*messages, fmt.Sprintf("Undefined error: %s", message))
	return  http.StatusInternalServerError, false
}

func DuplicateErrorMessage(field string, value string, messages *[]string) (int, bool)  {
	*messages = append(*messages, fmt.Sprintf("%s %s is already exists", field, value))
	return  http.StatusConflict, false
}

func InvalidFormatMessage(field string, fieldType string, messages *[]string) (int, bool)  {
	*messages = append(*messages, fmt.Sprintf("%s must be %s", field, fieldType))
	return  http.StatusPreconditionFailed, false
}

func RecordNotFoundMessage(data string, messages *[]string) (int, bool)  {
	*messages = append(*messages, fmt.Sprintf("%s is not found", data))
	return  http.StatusNotFound, false
}

func CustomPreconditionErrorMessage(message string, messages *[]string) (int, bool)  {
	*messages = append(*messages, fmt.Sprintf("%s", message))
	return  http.StatusPreconditionFailed, false
}

func CustomPreconditionRequiredErrorMessage(message string, messages *[]string) (int, bool)  {
	*messages = append(*messages, fmt.Sprintf("%s", message))
	return  http.StatusPreconditionRequired, false
}

package errors

import "errors"

var (
	ErrorValueNotFound        = errors.New("value not found ")
	ErrorKeyValueAlreadyExist = errors.New("key/value already exist ")
)

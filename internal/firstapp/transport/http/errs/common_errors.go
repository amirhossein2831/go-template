package errs

import "errors"

// MOGNO_ERROR
var (
	MongoInternalServerError = errors.New("mongo internal server error")
)

// APPLICATION_ERROR
var (
	DefaultServerError = errors.New("something went wrong")
)

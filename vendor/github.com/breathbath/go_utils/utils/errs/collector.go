package errs

import (
	"errors"
	"fmt"
	"strings"
)

func CollectErrors(sep string, errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	errorWord := ""
	for _, err := range errs {
		if err != nil {
			errorWord += err.Error() + sep
		}
	}

	if errorWord == "" {
		return nil
	}

	errorWord = strings.TrimRight(errorWord, sep)

	return errors.New(errorWord)
}

func WrapError(err error, errSeparator, msg string, args ...interface{}) error {
	cont := NewErrorContainer()
	cont.AddError(fmt.Errorf(msg, args...))
	cont.AddError(err)

	return cont.Result(errSeparator)
}

func AppendError(err error, errors *[]error) {
	if err != nil {
		*errors = append(*errors, err)
	}
}

type ErrorContainer struct {
	errors []error
}

func NewErrorContainer() *ErrorContainer {
	return &ErrorContainer{
		errors: []error{},
	}
}

func (ec *ErrorContainer) AddError(err error) {
	AppendError(err, &ec.errors)
}

func (ec *ErrorContainer) Result(sep string) error {
	return CollectErrors(sep, ec.errors...)
}

func (ec *ErrorContainer) AddErrorF(msg string, args ...interface{}) {
	err := fmt.Errorf(msg, args...)
	ec.AddError(err)
}

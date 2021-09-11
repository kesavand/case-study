package errs

import (
	"fmt"
	"log"
)

func FailOnError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func FailOnErrorF(template string, args ...interface{}) {
	errMsg := fmt.Errorf(template, args...)
	FailOnError(errMsg)
}

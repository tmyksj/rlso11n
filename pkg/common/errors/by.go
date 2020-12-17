package errors

import (
	"errors"
	"fmt"
)

func By(err error, format string, args ...interface{}) error {
	if err == nil {
		return errors.New(fmt.Sprintf("[ "+format+" ]", args...))
	} else {
		return errors.New(fmt.Sprintf("[ "+format+" ], caused by %v", append(args, err)...))
	}
}

package types

import (
	"os"
)

type ErrInvalidPermissions struct {
	File string
	Perm os.FileMode
	Err  error
}

func (e ErrInvalidPermissions) detailed() (valid bool) {
	if e.File != "" && e.Perm != 0 {
		valid = true
	}
	return
}

func (e ErrInvalidPermissions) Error() string {
	errMsg := "Invalid file permission"
	if e.detailed() {
		errMsg += " for file " + e.File + " should be " + e.Perm.String()
	}
	if e.Err != nil {
		errMsg += " \nerr: " + e.Err.Error()
	}
	return errMsg
}

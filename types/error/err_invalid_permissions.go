package types

type ErrInvalidPermissions struct {
	Err error
}

func (e ErrInvalidPermissions) Error() string {
	errMsg := "Invalid file permission"
	if e.Err != nil {
		errMsg += " err: " + e.Err.Error()
	}
	return errMsg
}

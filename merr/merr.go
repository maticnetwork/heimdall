package merr

type VErr struct {
	Field string
}

func (v VErr) Error() string {
	return "[ERROR] Validation Failure, missing " + v.Field
}

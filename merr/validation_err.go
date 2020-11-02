package merr

import "fmt"

type ValErr struct {
	Field  string
	Module string
}

func (v ValErr) Error() string {
	return fmt.Sprintf("[ERROR] Validation failure for field: %v, module: %v", v.Field, v.Module)
}

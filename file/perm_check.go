package file

import (
	"errors"
	"os"

	types "github.com/maticnetwork/heimdall/types/error"
)

const (
	// storing constant as this recommended as a security feature
	secretPerm os.FileMode = 0600
)

// PermCheck check the secret key and the keystore files.
// it verifies whether they are stored with the correct permissions.
func PermCheck(filePath string, validPerm os.FileMode) error {
	// get path to keystore files

	f, err := os.Stat(filePath)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return types.ErrInvalidPermissions{File: filePath, Perm: validPerm, Err: err}
	}

	filePerm := f.Mode()
	if filePerm != os.FileMode(validPerm) {
		return types.ErrInvalidPermissions{File: filePath, Perm: validPerm}
	}

	return nil
}

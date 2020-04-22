package cmd

import (
	"fmt"
	"os"
)

const (
	// storing constant as this recomended as a secruity feature
	secretPerm os.FileMode = 0600
)

var ErrInvalidPermissions = fmt.Errorf("Invalid file permission")

// permCheck check the secret key and the keystore files.
// it verifies whether they are stored with the correct permissions.
func permCheck(fileName string, validPerm uint32) (err error) {
	// get path to keystore files

	f, err := os.Stat(fileName)
	if err != nil {
		return err
	}

	if f.Mode() != os.FileMode(validPerm) {
		return ErrInvalidPermissions
		// lgr.Error("Invalid file permissions for file" + f.Name())
	}

	// check for keystore file permissions

	return nil
}

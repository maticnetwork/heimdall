package file

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/common"

	types "github.com/maticnetwork/heimdall/types/error"
)

func TestPermCheck(t *testing.T) {
	tc := []struct {
		filePath  string
		perm      os.FileMode
		validPerm os.FileMode
		expErr    error
		msg       string
	}{
		{
			filePath:  "/tmp/heimdall_test/test.json",
			perm:      0o777,
			validPerm: 0o600,
			expErr:    types.ErrInvalidPermissions{File: "/tmp/heimdall_test/test.json", Perm: 0o600},
			msg:       "test for invalid permission",
		},
		{
			filePath:  "/tmp/heimdall_test/test.json",
			perm:      0o600,
			validPerm: 0o600,
			msg:       "success",
		},
	}

	for i, c := range tc {
		// get path to UAT secrets file
		caseMsg := fmt.Sprintf("for i: %v, case: %v", i, c.msg)
		// set files for perm

		err := common.EnsureDir(filepath.Dir(c.filePath), 0o777)
		assert.Nil(t, err, caseMsg)
		_, err = os.OpenFile(c.filePath, os.O_CREATE, c.perm) // os.OpenFile creates the file if it is missing
		assert.Nil(t, err, caseMsg)

		// check file perm for secret file
		err = PermCheck(c.filePath, c.validPerm)
		assert.Equal(t, c.expErr, err)

		os.Remove(c.filePath) // cleen up
	}
}

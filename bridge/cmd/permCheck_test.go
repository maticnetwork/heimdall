package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/common"
)

func Test_permCheck(t *testing.T) {

	tc := []struct {
		confRoot  string
		perm      os.FileMode
		validPerm uint32
		expErr    error
		msg       string
	}{
		{
			confRoot:  "/tmp/heimdall_test",
			perm:      0644,
			validPerm: 0600,
			expErr:    ErrInvalidPermissions,
			msg:       "test for invalid permission",
		},
		{
			confRoot:  "/tmp/heimdall_test",
			perm:      0600,
			validPerm: 0600,
			msg:       "success",
		},
	}

	for i, c := range tc {
		conf := cfg.TestConfig()
		// get path to UAT secrets file
		caseMsg := fmt.Sprintf("for i: %v, case: %v", i, c.msg)
		// set files for perm
		if c.confRoot != "" {
			conf.SetRoot(c.confRoot)
			defer func(dirName string) {
				err := os.Remove(dirName)
				if err != nil {
					t.Fatal(err)
				}
			}(conf.RootDir)
		}
		err := common.EnsureDir(filepath.Dir(conf.PrivValidatorKeyFile()), 0777)
		assert.Nil(t, err, caseMsg)
		err = ioutil.WriteFile(conf.PrivValidatorKeyFile(), []byte(`{"priv_key": "test"}`), c.perm)
		assert.Nil(t, err, caseMsg)

		// check file perm for secret file
		err = permCheck(conf.PrivValidatorKeyFile(), c.validPerm)
		assert.Equal(t, c.expErr, err)
	}

}

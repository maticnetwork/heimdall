// Copyright 2018 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

// +build !windows

package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	gorand "math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/swarm/api"
	swarm "github.com/ethereum/go-ethereum/swarm/api/client"
)

// TestAccessPassword tests for the correct creation of an ACT manifest protected by a password.
// The test creates bogus content, uploads it encrypted, then creates the wrapping manifest with the Access entry
// The parties participating - node (publisher), uploads to second node then disappears. Content which was uploaded
// is then fetched through 2nd node. since the tested code is not key-aware - we can just
// fetch from the 2nd node using HTTP BasicAuth
func TestAccessPassword(t *testing.T) {
	cluster := newTestCluster(t, 1)
	defer cluster.Shutdown()
	proxyNode := cluster.Nodes[0]

	// create a tmp file
	tmp, err := ioutil.TempDir("", "swarm-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	// write data to file
	data := "notsorandomdata"
	dataFilename := filepath.Join(tmp, "data.txt")

	err = ioutil.WriteFile(dataFilename, []byte(data), 0666)
	if err != nil {
		t.Fatal(err)
	}

	hashRegexp := `[a-f\d]{128}`

	// upload the file with 'swarm up' and expect a hash
	up := runSwarm(t,
		"--bzzapi",
		proxyNode.URL, //it doesn't matter through which node we upload content
		"up",
		"--encrypt",
		dataFilename)
	_, matches := up.ExpectRegexp(hashRegexp)
	up.ExpectExit()

	if len(matches) < 1 {
		t.Fatal("no matches found")
	}

	ref := matches[0]

	password := "smth"
	passwordFilename := filepath.Join(tmp, "password.txt")

	err = ioutil.WriteFile(passwordFilename, []byte(password), 0666)
	if err != nil {
		t.Fatal(err)
	}

	up = runSwarm(t,
		"access",
		"new",
		"pass",
		"--dry-run",
		"--password",
		passwordFilename,
		ref,
	)

	_, matches = up.ExpectRegexp(".+")
	up.ExpectExit()

	if len(matches) == 0 {
		t.Fatalf("stdout not matched")
	}

	var m api.Manifest

	err = json.Unmarshal([]byte(matches[0]), &m)
	if err != nil {
		t.Fatalf("unmarshal manifest: %v", err)
	}

	if len(m.Entries) != 1 {
		t.Fatalf("expected one manifest entry, got %v", len(m.Entries))
	}

	e := m.Entries[0]

	ct := "application/bzz-manifest+json"
	if e.ContentType != ct {
		t.Errorf("expected %q content type, got %q", ct, e.ContentType)
	}

	if e.Access == nil {
		t.Fatal("manifest access is nil")
	}

	a := e.Access

	if a.Type != "pass" {
		t.Errorf(`got access type %q, expected "pass"`, a.Type)
	}
	if len(a.Salt) < 32 {
		t.Errorf(`got salt with length %v, expected not less the 32 bytes`, len(a.Salt))
	}
	if a.KdfParams == nil {
		t.Fatal("manifest access kdf params is nil")
	}

	client := swarm.NewClient(cluster.Nodes[0].URL)

	hash, err := client.UploadManifest(&m, false)
	if err != nil {
		t.Fatal(err)
	}

	httpClient := &http.Client{}

	url := cluster.Nodes[0].URL + "/" + "bzz:/" + hash
	response, err := httpClient.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusUnauthorized {
		t.Fatal("should be a 401")
	}
	authHeader := response.Header.Get("WWW-Authenticate")
	if authHeader == "" {
		t.Fatal("should be something here")
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("", password)

	response, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, response.StatusCode)
	}
	d, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(d) != data {
		t.Errorf("expected decrypted data %q, got %q", data, string(d))
	}

	wrongPasswordFilename := filepath.Join(tmp, "password-wrong.txt")

	err = ioutil.WriteFile(wrongPasswordFilename, []byte("just wr0ng"), 0666)
	if err != nil {
		t.Fatal(err)
	}

	//download file with 'swarm down' with wrong password
	up = runSwarm(t,
		"--bzzapi",
		proxyNode.URL,
		"down",
		"bzz:/"+hash,
		tmp,
		"--password",
		wrongPasswordFilename)

	_, matches = up.ExpectRegexp("unauthorized")
	if len(matches) != 1 && matches[0] != "unauthorized" {
		t.Fatal(`"unauthorized" not found in output"`)
	}
	up.ExpectExit()
}

// TestAccessPK tests for the correct creation of an ACT manifest between two parties (publisher and grantee).
// The test creates bogus content, uploads it encrypted, then creates the wrapping manifest with the Access entry
// The parties participating - node (publisher), uploads to second node (which is also the grantee) then disappears.
// Content which was uploaded is then fetched through the grantee's http proxy. Since the tested code is private-key aware,
// the test will fail if the proxy's given private key is not granted on the ACT.
func TestAccessPK(t *testing.T) {
	// Setup Swarm and upload a test file to it
	cluster := newTestCluster(t, 1)
	defer cluster.Shutdown()

	// create a tmp file
	tmp, err := ioutil.TempFile("", "swarm-test")
	if err != nil {
		t.Fatal(err)
	}
	defer tmp.Close()
	defer os.Remove(tmp.Name())

	// write data to file
	data := "notsorandomdata"
	_, err = io.WriteString(tmp, data)
	if err != nil {
		t.Fatal(err)
	}

	hashRegexp := `[a-f\d]{128}`

	// upload the file with 'swarm up' and expect a hash
	up := runSwarm(t,
		"--bzzapi",
		cluster.Nodes[0].URL,
		"up",
		"--encrypt",
		tmp.Name())
	_, matches := up.ExpectRegexp(hashRegexp)
	up.ExpectExit()

	if len(matches) < 1 {
		t.Fatal("no matches found")
	}

	ref := matches[0]

	pk := cluster.Nodes[0].PrivateKey
	granteePubKey := crypto.CompressPubkey(&pk.PublicKey)

	publisherDir, err := ioutil.TempDir("", "swarm-account-dir-temp")
	if err != nil {
		t.Fatal(err)
	}

	passFile, err := ioutil.TempFile("", "swarm-test")
	if err != nil {
		t.Fatal(err)
	}
	defer passFile.Close()
	defer os.Remove(passFile.Name())
	_, err = io.WriteString(passFile, testPassphrase)
	if err != nil {
		t.Fatal(err)
	}
	_, publisherAccount := getTestAccount(t, publisherDir)
	up = runSwarm(t,
		"--bzzaccount",
		publisherAccount.Address.String(),
		"--password",
		passFile.Name(),
		"--datadir",
		publisherDir,
		"--bzzapi",
		cluster.Nodes[0].URL,
		"access",
		"new",
		"pk",
		"--dry-run",
		"--grant-key",
		hex.EncodeToString(granteePubKey),
		ref,
	)

	_, matches = up.ExpectRegexp(".+")
	up.ExpectExit()

	if len(matches) == 0 {
		t.Fatalf("stdout not matched")
	}

	var m api.Manifest

	err = json.Unmarshal([]byte(matches[0]), &m)
	if err != nil {
		t.Fatalf("unmarshal manifest: %v", err)
	}

	if len(m.Entries) != 1 {
		t.Fatalf("expected one manifest entry, got %v", len(m.Entries))
	}

	e := m.Entries[0]

	ct := "application/bzz-manifest+json"
	if e.ContentType != ct {
		t.Errorf("expected %q content type, got %q", ct, e.ContentType)
	}

	if e.Access == nil {
		t.Fatal("manifest access is nil")
	}

	a := e.Access

	if a.Type != "pk" {
		t.Errorf(`got access type %q, expected "pk"`, a.Type)
	}
	if len(a.Salt) < 32 {
		t.Errorf(`got salt with length %v, expected not less the 32 bytes`, len(a.Salt))
	}
	if a.KdfParams != nil {
		t.Fatal("manifest access kdf params should be nil")
	}

	client := swarm.NewClient(cluster.Nodes[0].URL)

	hash, err := client.UploadManifest(&m, false)
	if err != nil {
		t.Fatal(err)
	}

	httpClient := &http.Client{}

	url := cluster.Nodes[0].URL + "/" + "bzz:/" + hash
	response, err := httpClient.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatal("should be a 200")
	}
	d, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(d) != data {
		t.Errorf("expected decrypted data %q, got %q", data, string(d))
	}
}

// TestAccessACT tests the e2e creation, uploading and downloading of an ACT type access control
// the test fires up a 3 node cluster, then randomly picks 2 nodes which will be acting as grantees to the data
// set. the third node should fail decoding the reference as it will not be granted access. the publisher uploads through
// one of the nodes then disappears.
func TestAccessACT(t *testing.T) {
	// Setup Swarm and upload a test file to it
	cluster := newTestCluster(t, 3)
	defer cluster.Shutdown()

	var uploadThroughNode = cluster.Nodes[0]
	client := swarm.NewClient(uploadThroughNode.URL)

	r1 := gorand.New(gorand.NewSource(time.Now().UnixNano()))
	nodeToSkip := r1.Intn(3) // a number between 0 and 2 (node indices in `cluster`)
	// create a tmp file
	tmp, err := ioutil.TempFile("", "swarm-test")
	if err != nil {
		t.Fatal(err)
	}
	defer tmp.Close()
	defer os.Remove(tmp.Name())

	// write data to file
	data := "notsorandomdata"
	_, err = io.WriteString(tmp, data)
	if err != nil {
		t.Fatal(err)
	}

	hashRegexp := `[a-f\d]{128}`

	// upload the file with 'swarm up' and expect a hash
	up := runSwarm(t,
		"--bzzapi",
		cluster.Nodes[0].URL,
		"up",
		"--encrypt",
		tmp.Name())
	_, matches := up.ExpectRegexp(hashRegexp)
	up.ExpectExit()

	if len(matches) < 1 {
		t.Fatal("no matches found")
	}

	ref := matches[0]
	grantees := []string{}
	for i, v := range cluster.Nodes {
		if i == nodeToSkip {
			continue
		}
		pk := v.PrivateKey
		granteePubKey := crypto.CompressPubkey(&pk.PublicKey)
		grantees = append(grantees, hex.EncodeToString(granteePubKey))
	}

	granteesPubkeyListFile, err := ioutil.TempFile("", "grantees-pubkey-list.csv")
	if err != nil {
		t.Fatal(err)
	}

	_, err = granteesPubkeyListFile.WriteString(strings.Join(grantees, "\n"))
	if err != nil {
		t.Fatal(err)
	}

	defer granteesPubkeyListFile.Close()
	defer os.Remove(granteesPubkeyListFile.Name())

	publisherDir, err := ioutil.TempDir("", "swarm-account-dir-temp")
	if err != nil {
		t.Fatal(err)
	}

	passFile, err := ioutil.TempFile("", "swarm-test")
	if err != nil {
		t.Fatal(err)
	}
	defer passFile.Close()
	defer os.Remove(passFile.Name())
	_, err = io.WriteString(passFile, testPassphrase)
	if err != nil {
		t.Fatal(err)
	}

	_, publisherAccount := getTestAccount(t, publisherDir)
	up = runSwarm(t,
		"--bzzaccount",
		publisherAccount.Address.String(),
		"--password",
		passFile.Name(),
		"--datadir",
		publisherDir,
		"--bzzapi",
		cluster.Nodes[0].URL,
		"access",
		"new",
		"act",
		"--grant-keys",
		granteesPubkeyListFile.Name(),
		ref,
	)

	_, matches = up.ExpectRegexp(`[a-f\d]{64}`)
	up.ExpectExit()

	if len(matches) == 0 {
		t.Fatalf("stdout not matched")
	}
	hash := matches[0]
	m, _, err := client.DownloadManifest(hash)
	if err != nil {
		t.Fatalf("unmarshal manifest: %v", err)
	}

	if len(m.Entries) != 1 {
		t.Fatalf("expected one manifest entry, got %v", len(m.Entries))
	}

	e := m.Entries[0]

	ct := "application/bzz-manifest+json"
	if e.ContentType != ct {
		t.Errorf("expected %q content type, got %q", ct, e.ContentType)
	}

	if e.Access == nil {
		t.Fatal("manifest access is nil")
	}

	a := e.Access

	if a.Type != "act" {
		t.Fatalf(`got access type %q, expected "act"`, a.Type)
	}
	if len(a.Salt) < 32 {
		t.Fatalf(`got salt with length %v, expected not less the 32 bytes`, len(a.Salt))
	}
	if a.KdfParams != nil {
		t.Fatal("manifest access kdf params should be nil")
	}

	httpClient := &http.Client{}

	// all nodes except the skipped node should be able to decrypt the content
	for i, node := range cluster.Nodes {
		log.Debug("trying to fetch from node", "node index", i)

		url := node.URL + "/" + "bzz:/" + hash
		response, err := httpClient.Get(url)
		if err != nil {
			t.Fatal(err)
		}
		log.Debug("got response from node", "response code", response.StatusCode)

		if i == nodeToSkip {
			log.Debug("reached node to skip", "status code", response.StatusCode)

			if response.StatusCode != http.StatusUnauthorized {
				t.Fatalf("should be a 401")
			}

			continue
		}

		if response.StatusCode != http.StatusOK {
			t.Fatal("should be a 200")
		}
		d, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(d) != data {
			t.Errorf("expected decrypted data %q, got %q", data, string(d))
		}
	}
}

// TestKeypairSanity is a sanity test for the crypto scheme for ACT. it asserts the correct shared secret according to
// the specs at https://github.com/ethersphere/swarm-docs/blob/eb857afda906c6e7bb90d37f3f334ccce5eef230/act.md
func TestKeypairSanity(t *testing.T) {
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		t.Fatalf("reading from crypto/rand failed: %v", err.Error())
	}
	sharedSecret := "a85586744a1ddd56a7ed9f33fa24f40dd745b3a941be296a0d60e329dbdb896d"

	for i, v := range []struct {
		publisherPriv string
		granteePub    string
	}{
		{
			publisherPriv: "ec5541555f3bc6376788425e9d1a62f55a82901683fd7062c5eddcc373a73459",
			granteePub:    "0226f213613e843a413ad35b40f193910d26eb35f00154afcde9ded57479a6224a",
		},
		{
			publisherPriv: "70c7a73011aa56584a0009ab874794ee7e5652fd0c6911cd02f8b6267dd82d2d",
			granteePub:    "02e6f8d5e28faaa899744972bb847b6eb805a160494690c9ee7197ae9f619181db",
		},
	} {
		b, _ := hex.DecodeString(v.granteePub)
		granteePub, _ := crypto.DecompressPubkey(b)
		publisherPrivate, _ := crypto.HexToECDSA(v.publisherPriv)

		ssKey, err := api.NewSessionKeyPK(publisherPrivate, granteePub, salt)
		if err != nil {
			t.Fatal(err)
		}

		hasher := sha3.NewKeccak256()
		hasher.Write(salt)
		shared, err := hex.DecodeString(sharedSecret)
		if err != nil {
			t.Fatal(err)
		}
		hasher.Write(shared)
		sum := hasher.Sum(nil)

		if !bytes.Equal(ssKey, sum) {
			t.Fatalf("%d: got a session key mismatch", i)
		}
	}
}

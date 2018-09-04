/* Copyright 2013 Steve Leonard <sleonard76@gmail.com>. All rights reserved.
Use of this source code is governed by the MIT license that can be found
in the LICENSE file.
*/

package merkle

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"testing"

	"github.com/stretchr/testify/assert"
)

// SimpleHash: does nothing

var SimpleHashData []byte

type SimpleHash struct{}

func NewSimpleHash() hash.Hash {
	return SimpleHash{}
}

func (self SimpleHash) Write(p []byte) (int, error) {
	size := self.Size()
	datalen := (len(p) / size) * size
	if len(p) == 0 || len(p)%size != 0 {
		datalen += size
	}
	data := make([]byte, datalen)
	copy(data, p)

	block := make([]byte, size)
	copy(block, data[:size])
	for i := 1; i < len(data)/size; i++ {
		_block := data[i*size : (i+1)*size]
		for j, c := range _block {
			block[j] += c
		}
	}

	SimpleHashData = append(SimpleHashData, block...)
	return size, nil
}
func (self SimpleHash) Sum(p []byte) []byte {
	p = append(p[:], SimpleHashData[:]...)
	return p
}
func (self SimpleHash) Reset() {
	SimpleHashData = nil
}
func (self SimpleHash) Size() int {
	return 32
}
func (self SimpleHash) BlockSize() int {
	return 32
}

type NotHash struct{}

func NewNotHash() hash.Hash {
	return NotHash{}
}
func (self NotHash) Write(p []byte) (int, error) {
	return 32, nil
}
func (self NotHash) Sum(p []byte) []byte {
	return p
}
func (self NotHash) Reset() {
}
func (self NotHash) Size() int {
	return 32
}
func (self NotHash) BlockSize() int {
	return 32
}

// FailingHash: always returns error on Write
type FailingHash struct {
	SucceedFor int
}

var failingHashWriteAttempts int = 0

func NewFailingHashAt(n int) FailingHash {
	failingHashWriteAttempts = 0
	return FailingHash{SucceedFor: n}
}

func NewFailingHash() FailingHash {
	return NewFailingHashAt(0)
}

func (self FailingHash) Write(p []byte) (int, error) {
	failingHashWriteAttempts += 1
	if failingHashWriteAttempts > self.SucceedFor {
		return 0, errors.New("Failed to write hash")
	} else {
		return 0, nil
	}
}
func (self FailingHash) Sum(p []byte) []byte {
	return p
}
func (self FailingHash) Reset() {
}
func (self FailingHash) Size() int {
	return 0
}
func (self FailingHash) BlockSize() int {
	return 0
}

func failNotEqual(t *testing.T, label string, input interface{},
	expect interface{}, result interface{}) {
	t.Errorf("%s(%v) != %v (%v, instead)", label, input, expect, result)
}

/* Utils */

func TestCalculateTreeHeight(t *testing.T) {
	inputs := [][]uint64{
		{0, 0},
		{1, 2},
		{2, 2},
		{3, 3},
		{4, 3},
		{5, 4},
		{6, 4},
		{7, 4},
		{8, 4},
		{9, 5},
		{15, 5},
		{16, 5},
		{17, 6},
		{31, 6},
		{32, 6},
		{63, 7},
		{64, 7},
		{65, 8},
	}
	for _, i := range inputs {
		r := calculateTreeHeight(i[0])
		if r != i[1] {
			failNotEqual(t, "calculateTreeHeight", i[0], i[1], r)
		}
		r, _ = CalculateHeightAndNodeCount(i[0])
		if r != i[1] {
			failNotEqual(t, "CalculateHeightAndNodeCount", i[0], i[1], r)
		}
	}
}

func TestCeilLogBaseTwo(t *testing.T) {
	inputs := [][]uint64{
		{0, 0},
		{1, 0},
		{2, 1},
		{3, 2},
		{7, 3},
		{15, 4},
		{32, 5},
		{63, 6},
		{64, 6},
	}
	for _, i := range inputs {
		r := ceilLogBaseTwo(i[0])
		if r != i[1] {
			failNotEqual(t, "ceilLogBaseTwo", i[0], i[1], r)
		}
	}
}

func TestLogBaseTwo(t *testing.T) {
	inputs := [][]uint64{
		{0, 0},
		{1, 0},
		{2, 1},
		{4, 2},
		{8, 3},
		{16, 4},
		{32, 5},
		{64, 6},
	}
	for _, i := range inputs {
		r := logBaseTwo(i[0])
		if r != i[1] {
			failNotEqual(t, "logBaseTwo", i[0], i[1], r)
		}
	}
}

func TestCalculateNodeCount(t *testing.T) {
	inputs := [][]uint64{
		{0, 0},
		{1, 1},
		{2, 3},
		{3, 6},
		{4, 7},
		{9, 20},
		{10, 21},
		{11, 23},
		{12, 24},
		{13, 27},
		{21, 44},
		{22, 45},
	}
	for _, i := range inputs {
		height := calculateTreeHeight(i[0])
		r := calculateNodeCount(height, i[0])
		if r != i[1] {
			failNotEqual(t, "calculateNodeCount", i[0], i[1], r)
		}
		_, r = CalculateHeightAndNodeCount(i[0])
		if r != i[1] {
			failNotEqual(t, "CalculateHeightAndNodeCount", i[0], i[1], r)
		}
	}
}

func TestNextPowerOfTwo(t *testing.T) {
	inputs := [][]uint64{
		{0, 1},
		{1, 1},
		{2, 2},
		{3, 4},
		{4, 4},
		{5, 8},
		{8, 8},
		{14, 16},
		{16, 16},
		{65535, 65536},
		{65536, 65536},
		{65537, 131072},
	}
	for _, i := range inputs {
		r := nextPowerOfTwo(i[0])
		if r != i[1] {
			failNotEqual(t, "nextPowerOfTwo", i[0], i[1], r)
		}
	}
}

func TestIsPowerOfTwo(t *testing.T) {
	type powerOfTwoResult struct {
		input  uint64
		output bool
	}
	inputs := []powerOfTwoResult{
		{0, false},
		{1, true},
		{2, true},
		{3, false},
		{4, true},
		{16, true},
		{65534, false},
		{65535, false},
		{65536, true},
		{65537, false},
		{2032131433, false},
	}
	for _, i := range inputs {
		r := isPowerOfTwo(i.input)
		if r != i.output {
			failNotEqual(t, "isPowerOfTwo", i.input, i.output, r)
		}
	}
}

/* Tree */

func containsNode(nodes []Node, node *Node) bool {
	/* Returns trueif a *Node is in a []Node */
	for i := 0; i < len(nodes); i++ {
		if node == &nodes[i] {
			return true
		}
	}
	return false
}

func createDummyTreeData(count, size int, use_rand bool) [][]byte {
	/* Creates an array of bytes with nonsense in them */
	data := make([][]byte, count)
	for i := 0; i < count; i++ {
		garbage := make([]byte, size)
		if use_rand {
			read := 0
			for read < size {
				n, _ := rand.Read(garbage[read:])
				read += n
			}
		} else {
			for i := 0; i < size; i++ {
				garbage[i] = byte((i + 1) % 0xFF)
			}
		}
		data[i] = garbage
	}
	return data
}

func verifyGeneratedTree(t *testing.T, tree *Tree, h hash.Hash) {
	/* Given a generated tree, confirm its state is correct */

	// Nodes should have been created
	assert.NotNil(t, tree.Nodes)
	assert.Equal(t, len(tree.Nodes), cap(tree.Nodes),
		"tree.Nodes len should equal its cap")

	// The leaves should not have children
	for _, n := range tree.Leaves() {
		assert.Nil(t, n.Left)
		assert.Nil(t, n.Right)
	}

	for i := tree.Height() - 1; i > 0; i-- {
		// All the other nodes should have children, and their children
		// should be in the deeper level
		deeper := tree.GetNodesAtHeight(i + 1)
		row := tree.GetNodesAtHeight(i)
		for j, n := range row {
			assert.NotNil(t, n.Left, "Left child should never be nil")
			assert.Equal(t, n.Left, &deeper[j*2])
			if j == len(row)-1 && len(deeper)%2 == 1 {
				// Last node in this level should have nil right child
				// if its unbalanced
				assert.Nil(t, n.Right)
				// Its hash should be the same as the left node hash
				assert.Equal(t, n.Left.Hash, n.Hash,
					"Left child hash should equal node hash when right child is nil")
			} else {
				assert.NotNil(t, n.Right)
				assert.Equal(t, n.Right, &deeper[j*2+1])
				assert.NotEqual(t, bytes.Equal(n.Right.Hash, n.Hash), true,
					"Right child hash should not equal node hash")
				assert.NotEqual(t, bytes.Equal(n.Left.Hash, n.Hash), true,
					"Left child hash should not equal node hash")
				verifyHashInNode(t, tree, n, h)
			}
		}

		// Each row should have prev/2 + prev%2 nodes
		prev := len(deeper)
		assert.Equal(t, len(row), prev/2+prev%2)
	}

	rootRow := tree.GetNodesAtHeight(1)
	// The root row should exist
	assert.NotNil(t, rootRow)

	// The root row should be of length 1
	assert.Equal(t, len(rootRow), 1,
		"The root row should contain only 1 node")

	// the Root() should be the only item in the top row
	assert.Equal(t, tree.Root(), &rootRow[0],
		"tree.Root() is not the expected node")

	// Verify Root Hash
	verifyHashInNode(t, tree, *tree.Root(), h)

	// The Leaves() should the deepest row
	assert.Equal(t, len(tree.Leaves()),
		len(tree.GetNodesAtHeight(tree.Height())),
		"tree.Leaves() is not the expected row")
}

func verifyHashInNode(t *testing.T, tree *Tree, n Node, h hash.Hash) {
	/* Given a node it verifies that the Node Hash was calculated correctly */
	nn, err := tree.generateNode(n.Left.Hash, n.Right.Hash, h)

	assert.Nil(t, err)
	assert.Equal(t, nn.Hash, n.Hash, "calculated Hash needs to match generated one")
}

func verifyInitialState(t *testing.T, tree *Tree) {
	assert.Nil(t, tree.Nodes)
	assert.Nil(t, tree.Levels)
}

func TestNewNode(t *testing.T) {
	h := NewSimpleHash()
	block := createDummyTreeData(1, h.Size(), true)[0]
	n, err := NewNode(h, block)
	assert.Nil(t, err)
	assert.Equal(t, bytes.Equal(n.Hash, block), true)

	// Passing a nil hash function should create a node with the unhashed block
	n, err = NewNode(nil, block)
	assert.Nil(t, err)
	assert.Equal(t, n.Hash, block)

	// Any nil argument should return blank node, no error
	n, err = NewNode(nil, nil)
	assert.Nil(t, err)
	assert.Nil(t, n.Hash)
	n, err = NewNode(h, nil)
	assert.Nil(t, err)
	assert.Nil(t, n.Hash)

	// Check hash error handling
	h = NewFailingHash()
	n, err = NewNode(h, block)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Failed to write hash")
}

func TestNewTree(t *testing.T) {
	tree := NewTree()
	verifyInitialState(t, &tree)
	assert.False(t, tree.Options.EnableHashSorting)
	assert.False(t, tree.Options.DisableHashLeaves)
}

func TestNewTreeWithOpts(t *testing.T) {
	tree := NewTreeWithOpts(TreeOptions{EnableHashSorting: true, DisableHashLeaves: true})
	verifyInitialState(t, &tree)
	assert.True(t, tree.Options.EnableHashSorting)
	assert.True(t, tree.Options.DisableHashLeaves)
}

func TestTreeUngenerated(t *testing.T) {
	tree := Tree{}
	// If data is nil, it should handle that:
	err := tree.Generate(nil, NewSimpleHash())
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Empty tree")
	assert.Nil(t, tree.Leaves())
	assert.Nil(t, tree.Root())
	assert.Equal(t, tree.Height(), uint64(0))
	assert.Nil(t, tree.Nodes)
}

func TestTreeGenerate(t *testing.T) {
	tree := Tree{}
	h := NewSimpleHash()
	// Setup some dummy data
	blockCount := 13
	blockSize := 16
	data := createDummyTreeData(blockCount, blockSize, true)

	// Generate the tree
	err := tree.Generate(data, h)
	assert.Nil(t, err)
	verifyGeneratedTree(t, &tree, h)

	// Generating with no blocks should return error
	err = tree.Generate(make([][]byte, 0, 1), h)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Empty tree")
}

func TestTreeGenerate_DisableHashLeaves(t *testing.T) {
	h := sha256.New()
	items := [][]byte{[]byte("alpha"), []byte("beta")}
	alpha := sha256.Sum256(items[0])
	beta := sha256.Sum256(items[1])
	items_hashed := [][]byte{alpha[:32], beta[:32]}

	treeHashedLeaves := NewTree()
	err := treeHashedLeaves.Generate(items, h)
	assert.Nil(t, err)

	tree := NewTreeWithOpts(TreeOptions{false, true})
	err = tree.Generate(items_hashed, h)
	assert.Nil(t, err)
	assert.Equal(t, tree.Root().Hash, treeHashedLeaves.Root().Hash)
}

func TestGenerateNodeHashOfUnbalance(t *testing.T) {
	tree := Tree{}
	tree.Options.EnableHashSorting = true
	h := NewSimpleHash()

	sampleLeft := []byte{203, 225, 206, 227, 57, 204, 31, 188, 40, 131, 158, 32, 174, 43, 15, 187, 176, 223, 90, 55, 162, 35, 25, 177, 219, 173, 93, 54, 138, 119, 188, 56}
	n, err := tree.generateNode(sampleLeft, nil, h)
	assert.Nil(t, err)
	assert.Equal(t, sampleLeft, n.Hash)
}

func TestGenerateNodeHashOrdered(t *testing.T) {
	tree := Tree{}
	tree.Options.EnableHashSorting = true
	h := NewSimpleHash()

	sampleLeft := []byte{203, 225, 206, 227, 57, 204, 31, 188, 40, 131, 158, 32, 174, 43, 15, 187, 176, 223, 90, 55, 162, 35, 25, 177, 219, 173, 93, 54, 138, 119, 188, 56}
	sampleRight := []byte{193, 201, 112, 48, 157, 84, 238, 81, 120, 81, 228, 112, 38, 213, 168, 50, 37, 170, 137, 211, 44, 177, 75, 68, 152, 252, 54, 145, 145, 146, 154, 136}

	data := make([]byte, h.Size()*2)
	copy(data[:h.Size()], sampleRight)
	copy(data[h.Size():], sampleLeft)

	expected, _ := NewNode(h, data)
	n, err := tree.generateNode(sampleLeft, sampleRight, h)
	assert.Nil(t, err)
	assert.Equal(t, expected.Hash, n.Hash)
}

func TestGenerateNodeHashStandard(t *testing.T) {
	tree := Tree{}
	h := NewSimpleHash()
	sampleLeft := []byte{203, 225, 206, 227, 57, 204, 31, 188, 40, 131, 158, 32, 174, 43, 15, 187, 176, 223, 90, 55, 162, 35, 25, 177, 219, 173, 93, 54, 138, 119, 188, 56}
	sampleRight := []byte{193, 201, 112, 48, 157, 84, 238, 81, 120, 81, 228, 112, 38, 213, 168, 50, 37, 170, 137, 211, 44, 177, 75, 68, 152, 252, 54, 145, 145, 146, 154, 136}

	data := make([]byte, h.Size()*2)
	copy(data[:h.Size()], sampleRight)
	copy(data[h.Size():], sampleLeft)

	expected, _ := NewNode(h, data)
	n, err := tree.generateNode(sampleLeft, sampleRight, h)
	assert.Nil(t, err)
	assert.Equal(t, expected.Hash, n.Hash)
}

func TestHashOrderedTreeGenerate(t *testing.T) {
	tree := Tree{}
	tree.Options.EnableHashSorting = true
	h := NewSimpleHash()

	// Setup some dummy data
	blockCount := 13
	blockSize := 16
	data := createDummyTreeData(blockCount, blockSize, true)

	// Generate the tree
	err := tree.Generate(data, h)
	assert.Nil(t, err)
	verifyGeneratedTree(t, &tree, h)

	// Generating with no blocks should return error
	err = tree.Generate(make([][]byte, 0, 1), NewSimpleHash())
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Empty tree")
}

func TestGenerateFailedHash(t *testing.T) {
	tree := NewTree()
	data := createDummyTreeData(16, 16, true)
	// Fail hash during the leaf generation
	err := tree.Generate(data, NewFailingHash())
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Failed to write hash")

	// Fail hash during internal node generation
	data = createDummyTreeData(16, 16, true)
	err = tree.Generate(data, NewFailingHashAt(20))
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Failed to write hash")
}

func TestGetNodesAtHeight(t *testing.T) {
	// ungenerate tree should return nil
	tree := NewTree()
	h := NewSimpleHash()
	assert.Nil(t, tree.GetNodesAtHeight(1))

	count := 15
	size := 16
	data := createDummyTreeData(count, size, true)
	tree.Generate(data, NewSimpleHash())
	verifyGeneratedTree(t, &tree, h)

	// invalid height should return nil
	assert.Nil(t, tree.GetNodesAtHeight(0))
	assert.Nil(t, tree.GetNodesAtHeight(tree.Height()+1))

	// check valid height = 1
	nodes := tree.GetNodesAtHeight(tree.Height())
	assert.Equal(t, len(nodes), count)
	expect := tree.Nodes[:count]
	for i := 0; i < len(nodes); i++ {
		assert.Equal(t, &expect[i], &nodes[i])
	}
}

// Returns the root hash for an array of hashes
func simpleMerkle(data [][]byte) []byte {
	h := sha256.New()
	// Build the leaves
	h0 := make([][]byte, len(data))
	for i, b := range data {
		h.Reset()
		h.Write(b)
		h0[i] = h.Sum(nil)
	}

	h1 := make([][]byte, (len(h0)+len(h0)%2)/2)
	for {
		for i := 0; i < len(h0); i += 2 {
			var sum []byte
			if len(h0)%2 == 1 && i == len(h0)-1 {
				sum = h0[i]
			} else {
				c := append(h0[i], h0[i+1]...)
				h.Reset()
				h.Write(c)
				sum = h.Sum(nil)
			}
			h1[i/2] = sum
		}
		if len(h1) == 1 {
			break
		}
		h0 = h1
		h1 = make([][]byte, (len(h0)+len(h0)%2)/2)
	}
	return h1[0]
}

func TestRootHashValue(t *testing.T) {
	// Check the root hash made by Tree against a simpler implementation
	// that finds only the root hash

	tree := Tree{}
	h := sha256.New()
	// Setup some dummy data
	blockCount := 16
	blockSize := 16
	data := createDummyTreeData(blockCount, blockSize, true)

	// Generate the tree
	err := tree.Generate(data, h)
	assert.Nil(t, err)
	verifyGeneratedTree(t, &tree, h)

	// Calculate the root hash with the simpler method
	merk := simpleMerkle(data)

	assert.Equal(t, bytes.Equal(tree.Root().Hash, merk), true)
}

/* Benchmarks */

func generateBenchmark(b *testing.B, data [][]byte, hashf hash.Hash) {
	tree := NewTree()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Generate(data, hashf)
	}
}

func noHashBenchmark(b *testing.B, n int) {
	data := createDummyTreeData(n, 1, false)
	generateBenchmark(b, data, NewNotHash())
}

func BenchmarkGenerate_1K_Blocks_NoHash(b *testing.B) {
	noHashBenchmark(b, 1000)
}

func BenchmarkGenerate_100K_Blocks_NoHash(b *testing.B) {
	noHashBenchmark(b, 100000)
}

func BenchmarkGenerate_1M_Blocks_NoHash(b *testing.B) {
	noHashBenchmark(b, 1000000)
}

func BenchmarkGenerate_512MB_512KB_MD5(b *testing.B) {
	mb := 512
	blockSize := 512 * 1024
	data := createDummyTreeData((mb*1024*1024)/blockSize, blockSize, false)
	generateBenchmark(b, data, md5.New())
}

func BenchmarkGenerate_512MB_512KB_SHA256(b *testing.B) {
	mb := 512
	blockSize := 512 * 1024
	data := createDummyTreeData((mb*1024*1024)/blockSize, blockSize, false)
	generateBenchmark(b, data, sha256.New())
}

func BenchmarkGenerate_1GB_2MB_MD5(b *testing.B) {
	mb := 1024
	blockSize := 2 * 1024 * 1024
	data := createDummyTreeData((mb*1024*1024)/blockSize, blockSize, false)
	generateBenchmark(b, data, md5.New())
}

func BenchmarkGenerate_1GB_2MB_SHA256(b *testing.B) {
	mb := 1024
	blockSize := 2 * 1024 * 1024
	data := createDummyTreeData((mb*1024*1024)/blockSize, blockSize, false)
	generateBenchmark(b, data, sha256.New())
}

func Example_complete() {
	items := [][]byte{[]byte("alpha"), []byte("beta"), []byte("gamma"), []byte("delta"), []byte("epsilon")}

	treeOptions := TreeOptions{
		EnableHashSorting: false,
		DisableHashLeaves: false,
	}

	tree := NewTreeWithOpts(treeOptions)
	err := tree.Generate(items, md5.New())
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Height: %d\n", tree.Height())
	fmt.Printf("Root: %v\n", tree.Root())
	fmt.Printf("N Leaves: %v\n", len(tree.Leaves()))
	fmt.Printf("Height 2: %v\n", tree.GetNodesAtHeight(2))
}

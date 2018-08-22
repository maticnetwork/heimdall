package amino_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-amino"
)

func TestMarshalBinary(t *testing.T) {
	var cdc = amino.NewCodec()

	type SimpleStruct struct {
		String string
		Bytes  []byte
		Time   time.Time
	}

	s := SimpleStruct{
		String: "hello",
		Bytes:  []byte("goodbye"),
		Time:   time.Now().UTC().Truncate(time.Millisecond), // strip monotonic and timezone.
	}

	b, err := cdc.MarshalBinary(s)
	assert.Nil(t, err)
	t.Logf("MarshalBinary(s) -> %X", b)

	var s2 SimpleStruct
	err = cdc.UnmarshalBinary(b, &s2)
	assert.Nil(t, err)
	assert.Equal(t, s, s2)
}

func TestUnmarshalBinaryReader(t *testing.T) {
	var cdc = amino.NewCodec()

	type SimpleStruct struct {
		String string
		Bytes  []byte
		Time   time.Time
	}

	s := SimpleStruct{
		String: "hello",
		Bytes:  []byte("goodbye"),
		Time:   time.Now().UTC().Truncate(time.Millisecond), // strip monotonic and timezone.
	}

	b, err := cdc.MarshalBinary(s)
	assert.Nil(t, err)
	t.Logf("MarshalBinary(s) -> %X", b)

	var s2 SimpleStruct
	_, err = cdc.UnmarshalBinaryReader(bytes.NewBuffer(b), &s2, 0)
	assert.Nil(t, err)

	assert.Equal(t, s, s2)
}

func TestUnmarshalBinaryReaderSize(t *testing.T) {
	var cdc = amino.NewCodec()

	var s1 string = "foo"
	b, err := cdc.MarshalBinary(s1)
	assert.Nil(t, err)
	t.Logf("MarshalBinary(s) -> %X", b)

	var s2 string
	var n int64
	n, err = cdc.UnmarshalBinaryReader(bytes.NewBuffer(b), &s2, 0)
	assert.Nil(t, err)
	assert.Equal(t, s1, s2)
	frameLengthBytes, msgLengthBytes := 1, 1
	assert.Equal(t, frameLengthBytes+msgLengthBytes+len(s1), int(n))
}

func TestUnmarshalBinaryReaderSizeLimit(t *testing.T) {
	var cdc = amino.NewCodec()

	var s1 string = "foo"
	b, err := cdc.MarshalBinary(s1)
	assert.Nil(t, err)
	t.Logf("MarshalBinary(s) -> %X", b)

	var s2 string
	var n int64
	n, err = cdc.UnmarshalBinaryReader(bytes.NewBuffer(b), &s2, int64(len(b)-1))
	assert.NotNil(t, err, "insufficient limit should lead to failure")
	n, err = cdc.UnmarshalBinaryReader(bytes.NewBuffer(b), &s2, int64(len(b)))
	assert.Nil(t, err, "sufficient limit should not cause failure")
	assert.Equal(t, s1, s2)
	frameLengthBytes, msgLengthBytes := 1, 1
	assert.Equal(t, frameLengthBytes+msgLengthBytes+len(s1), int(n))
}

func TestUnmarshalBinaryReaderTooLong(t *testing.T) {
	var cdc = amino.NewCodec()

	type SimpleStruct struct {
		String string
		Bytes  []byte
		Time   time.Time
	}

	s := SimpleStruct{
		String: "hello",
		Bytes:  []byte("goodbye"),
		Time:   time.Now().UTC().Truncate(time.Millisecond), // strip monotonic and timezone.
	}

	b, err := cdc.MarshalBinary(s)
	assert.Nil(t, err)
	t.Logf("MarshalBinary(s) -> %X", b)

	var s2 SimpleStruct
	_, err = cdc.UnmarshalBinaryReader(bytes.NewBuffer(b), &s2, 1) // 1 byte limit is ridiculous.
	assert.NotNil(t, err)
}

func TestUnmarshalBinaryBufferedWritesReads(t *testing.T) {
	var cdc = amino.NewCodec()
	var buf = bytes.NewBuffer(nil)

	// Write 3 times.
	var s1 string = "foo"
	_, err := cdc.MarshalBinaryWriter(buf, s1)
	assert.Nil(t, err)
	_, err = cdc.MarshalBinaryWriter(buf, s1)
	assert.Nil(t, err)
	_, err = cdc.MarshalBinaryWriter(buf, s1)
	assert.Nil(t, err)

	// Read 3 times.
	var s2 string
	_, err = cdc.UnmarshalBinaryReader(buf, &s2, 0)
	assert.Nil(t, err)
	assert.Equal(t, s1, s2)
	_, err = cdc.UnmarshalBinaryReader(buf, &s2, 0)
	assert.Nil(t, err)
	assert.Equal(t, s1, s2)
	_, err = cdc.UnmarshalBinaryReader(buf, &s2, 0)
	assert.Nil(t, err)
	assert.Equal(t, s1, s2)

	// Reading 4th time fails.
	_, err = cdc.UnmarshalBinaryReader(buf, &s2, 0)
	assert.NotNil(t, err)
}

package amino

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDecodeSkippedFieldsInTime(t *testing.T) {
	type testTime struct {
		Time time.Time
	}
	cdc := NewCodec()

	tm, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", "1970-01-01 00:00:00 +0000 UTC")
	assert.NoError(t, err)

	b, err := cdc.MarshalBinary(testTime{Time: tm})
	assert.NoError(t, err)
	var ti testTime
	err = cdc.UnmarshalBinary(b, &ti)
	assert.NoError(t, err)
	assert.Equal(t, testTime{Time: tm}, ti)

	tm2, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", "1970-01-01 00:00:01.978131102 +0000 UTC")
	assert.NoError(t, err)

	b, err = cdc.MarshalBinary(testTime{Time: tm2})
	assert.NoError(t, err)
	err = cdc.UnmarshalBinary(b, &ti)
	assert.NoError(t, err)
	assert.Equal(t, testTime{Time: tm2}, ti)

	t1, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", "1970-01-01 00:00:11.577968799 +0000 UTC")
	t2, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", "2078-07-10 15:44:58.406865636 +0000 UTC")
	t3, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", "1970-01-01 00:00:00 +0000 UTC")
	t4, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", "1970-01-01 00:00:14.48251984 +0000 UTC")

	type tArr struct {
		TimeAr [4]time.Time
	}
	st := tArr{
		TimeAr: [4]time.Time{t1, t2, t3, t4},
	}
	b, err = cdc.MarshalBinary(st)
	assert.NoError(t, err)

	var tStruct tArr
	err = cdc.UnmarshalBinary(b, &tStruct)
	assert.NoError(t, err)
	assert.Equal(t, st, tStruct)
}

package bor

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"
	"math"
	"sync"
)

const (
	seedSize           = int8(32)
	roundSize          = int8(1)
	positionWindowSize = int8(4)
	pivotViewSize      = seedSize + roundSize
	totalSize          = seedSize + roundSize + positionWindowSize
	ShuffleRoundCount  = 90
)

var (
	maxShuffleListSize uint64 = 1 << 40
	sha256Pool                = sync.Pool{New: func() interface{} {
		return sha256.New()
	}}
)

// ShuffleList returns list of shuffled indexes in a pseudorandom permutation `p` of `0...list_size - 1` with “seed“ as entropy.
// We utilize 'swap or not' shuffling in this implementation; we are allocating the memory with the seed that stays
// constant between iterations instead of reallocating it each iteration as in the spec. This implementation is based
// on the original implementation from protolambda, https://github.com/protolambda/eth2-shuffle
//
//	improvements:
//	 - seed is always the first 32 bytes of the hash input, we just copy it into the buffer one time.
//	 - add round byte to seed and hash that part of the buffer.
//	 - split up the for-loop in two:
//	  1. Handle the part from 0 (incl) to pivot (incl). This is mirrored around (pivot / 2).
//	  2. Handle the part from pivot (excl) to N (excl). This is mirrored around ((pivot / 2) + (size/2)).
//	 - hash source every 256 iterations.
//	 - change byteV every 8 iterations.
//	 - we start at the edges, and work back to the mirror point.
//	   this makes us process each pear exactly once (instead of unnecessarily twice, like in the spec).
func ShuffleList(input []uint64, seed [32]byte) ([]uint64, error) {
	return innerShuffleList(input, seed, true /* shuffle */)
}

// shuffles or unshuffles, shuffle=false to un-shuffle.
func innerShuffleList(input []uint64, seed [32]byte, shuffle bool) ([]uint64, error) {
	if len(input) <= 1 {
		return input, nil
	}

	if uint64(len(input)) > maxShuffleListSize {
		return nil, fmt.Errorf("list size %d out of bounds",
			len(input))
	}

	rounds := uint8(ShuffleRoundCount)
	if rounds == 0 {
		return input, nil
	}

	listSize := uint64(len(input))
	buf := make([]byte, totalSize)
	r := uint8(0)

	if !shuffle {
		r = rounds - 1
	}

	copy(buf[:seedSize], seed[:])

	for {
		buf[seedSize] = r
		ph := sha256Hash(buf[:pivotViewSize])
		pivot := FromBytes8(ph[:8]) % listSize
		mirror := (pivot + 1) >> 1
		if pivot>>8 > math.MaxUint32 {
			return nil, fmt.Errorf("pivot value out of range for uint32: %d", pivot>>8)
		}
		//nolint:gosec
		binary.LittleEndian.PutUint32(buf[pivotViewSize:], uint32(pivot>>8))
		source := sha256Hash(buf)

		byteV := source[(pivot&0xff)>>3]
		for i, j := uint64(0), pivot; i < mirror; i, j = i+1, j-1 {
			byteV, source = swapOrNot(buf, byteV, i, input, j, source)
		}

		// Now repeat, but for the part after the pivot.
		mirror = (pivot + listSize + 1) >> 1
		end := listSize - 1
		if end>>8 > math.MaxUint32 {
			return nil, fmt.Errorf("end value out of range for uint32: %d", end>>8)
		}
		//nolint:gosec
		binary.LittleEndian.PutUint32(buf[pivotViewSize:], uint32(end>>8))
		source = sha256Hash(buf)

		byteV = source[(end&0xff)>>3]
		for i, j := pivot+1, end; i < mirror; i, j = i+1, j-1 {
			byteV, source = swapOrNot(buf, byteV, i, input, j, source)
		}

		if shuffle {
			r++
			if r == rounds {
				break
			}
		} else {
			if r == 0 {
				break
			}
			r--
		}
	}

	return input, nil
}

// swapOrNot describes the main algorithm behind the shuffle where we swap bytes in the inputted value
// depending on if the conditions are met.
func swapOrNot(buf []byte, byteV byte, i uint64, input []uint64, j uint64, source [32]byte) (byte, [32]byte) {
	if j&0xff == 0xff {
		// just overwrite the last part of the buffer, reuse the start (seed, round)
		//nolint:gosec
		binary.LittleEndian.PutUint32(buf[pivotViewSize:], uint32(j>>8))
		source = sha256Hash(buf)
	}

	if j&0x7 == 0x7 {
		byteV = source[(j&0xff)>>3]
	}

	bitV := (byteV >> (j & 0x7)) & 0x1
	if bitV == 1 {
		input[i], input[j] = input[j], input[i]
	}

	return byteV, source
}

// FromBytes8 returns an integer which is stored in the little-endian format(8, 'little')
// from a byte array.
func FromBytes8(x []byte) uint64 {
	if len(x) < 8 {
		return 0
	}

	return binary.LittleEndian.Uint64(x)
}

// Hash defines a function that returns the sha256 checksum of the data passed in
func sha256Hash(data []byte) [32]byte {
	h, ok := sha256Pool.Get().(hash.Hash)
	if !ok {
		h = sha256.New()
	}

	defer sha256Pool.Put(h)
	h.Reset()

	var b [32]byte

	// The hash interface never returns an error, for that reason
	// we are not handling the error below. For reference, it is
	// stated here https://golang.org/pkg/hash/#Hash

	// #nosec G104
	h.Write(data)
	h.Sum(b[:0])

	return b
}

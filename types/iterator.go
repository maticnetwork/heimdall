package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// KVStorePrefixIteratorPaginated returns iterator over items in the selected page.
// Items iterated and skipped in ascending order.
func KVStorePrefixIteratorPaginated(kvs sdk.KVStore, prefix []byte, page, limit uint) sdk.Iterator {
	pi := &PaginatedIterator{
		Iterator: sdk.KVStorePrefixIterator(kvs, prefix),
		page:     page,
		limit:    limit,
	}
	pi.skip()

	return pi
}

// KVStoreReversePrefixIteratorPaginated returns iterator over items in the selected page.
// Items iterated and skipped in descending order.
func KVStoreReversePrefixIteratorPaginated(kvs sdk.KVStore, prefix []byte, page, limit uint) sdk.Iterator {
	pi := &PaginatedIterator{
		Iterator: sdk.KVStoreReversePrefixIterator(kvs, prefix),
		page:     page,
		limit:    limit,
	}
	pi.skip()

	return pi
}

// KVStorePrefixRangeIteratorPaginated returns iterator over items in the selected page and queries within a range.
// Items iterated and skipped in ascending order.
func KVStorePrefixRangeIteratorPaginated(kvs sdk.KVStore, page, limit uint, from, to []byte) sdk.Iterator {
	pi := &PaginatedIterator{
		Iterator: kvs.Iterator(from, to),
		page:     page,
		limit:    limit,
	}
	pi.skip()

	return pi
}

// PaginatedIterator is a wrapper around Iterator that iterates over values starting for given page and limit.
type PaginatedIterator struct {
	sdk.Iterator

	page, limit uint // provided during initialization
	iterated    uint // incremented in a call to Next

}

func (pi *PaginatedIterator) skip() {
	for i := (pi.page - 1) * pi.limit; i > 0 && pi.Iterator.Valid(); i-- {
		pi.Iterator.Next()
	}
}

// Next will panic after limit is reached.
func (pi *PaginatedIterator) Next() {
	if !pi.Valid() {
		panic(fmt.Sprintf("PaginatedIterator reached limit %d", pi.limit))
	}

	pi.Iterator.Next()
	pi.iterated++
}

// Valid if below limit and underlying iterator is valid.
func (pi *PaginatedIterator) Valid() bool {
	if pi.iterated >= pi.limit {
		return false
	}

	return pi.Iterator.Valid()
}

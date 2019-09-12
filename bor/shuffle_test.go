package bor

import (
	"reflect"
	"testing"
)

func TestShuffleList_OK(t *testing.T) {
	var list1 []uint64
	seed1 := [32]byte{1, 128, 12}
	seed2 := [32]byte{2, 128, 12}
	for i := 0; i < 10; i++ {
		list1 = append(list1, uint64(i))
	}

	list2 := make([]uint64, len(list1))
	copy(list2, list1)

	list1, err := ShuffleList(list1, seed1)
	if err != nil {
		t.Errorf("Shuffle failed with: %v", err)
	}

	list2, err = ShuffleList(list2, seed2)
	if err != nil {
		t.Errorf("Shuffle failed with: %v", err)
	}

	if reflect.DeepEqual(list1, list2) {
		t.Errorf("2 shuffled lists shouldn't be equal")
	}
	// if !reflect.DeepEqual(list1, []uint64{0, 7, 8, 6, 3, 9, 4, 5, 2, 1}) {
	// 	t.Errorf("list 1 was incorrectly shuffled got: %v", list1)
	// }
	// if !reflect.DeepEqual(list2, []uint64{0, 5, 2, 1, 6, 8, 7, 3, 4, 9}) {
	// 	t.Errorf("list 2 was incorrectly shuffled got: %v", list2)
	// }
}

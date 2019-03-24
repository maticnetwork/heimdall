package pier

import (
	"sync"
	"testing"
)

func setup() *MaticCheckpointer {
	checkpointer := NewMaticCheckpointer()
	return checkpointer
}
func teardown() {

}
func TestGenHeaderDetails(t *testing.T) {
	// checkpointer := setup()
	var wg sync.WaitGroup
	wg.Add(1)

	// checkpointer.genHeaderDetailContract(1000,wg,)
}

package processor

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/maticnetwork/heimdall/bridge/setu/util"
	"github.com/maticnetwork/heimdall/helper"
)

func TestUrl(t *testing.T) {
	// This is a placeholder for the actual test implementation.
	// The test should validate the functionality of the URL processor.
	SpanByIdURL := fmt.Sprintf(util.SpanByIdURL, strconv.FormatUint(1, 10))
	result := helper.GetHeimdallServerEndpoint(SpanByIdURL)
	if result != "http://localhost:1317/heimdall/bor/span/1" {
		t.Errorf("Expected URL to be 'http://localhost:1317/heimdall/bor/span/1', got '%s'", result)
	}
}

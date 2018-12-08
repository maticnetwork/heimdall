package test

import (
	"github.com/maticnetwork/heimdall/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdateAck(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	keeper.UpdateACKCount(ctx)
	ack := keeper.GetACKCount(ctx)
	require.Equal(t, uint64(1), ack, "Ack Count Not Equal")
}

func TestCheckpointBuffer(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	checkpoint := types.CheckpointBlockHeader{
		Proposer: "0x17cde2546df29E2bbE66a98Ae95A6Ed8604D6B2b",
	}
}

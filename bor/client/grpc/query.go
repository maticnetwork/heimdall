package grpc

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"

	proto "github.com/maticnetwork/polyproto/bor"
	protoutil "github.com/maticnetwork/polyproto/utils"
)

func (h *BorGRPCClient) GetRootHash(ctx context.Context, startBlock uint64, endBlock uint64) (string, error) {

	req := &proto.GetRootHashRequest{
		StartBlockNumber: startBlock,
		EndBlockNumber:   endBlock,
	}

	log.Info("Fetching bor root hash")

	res, err := h.client.GetRootHash(ctx, req)
	if err != nil {
		return "", err
	}

	log.Info("Fetched bor root hash")

	return res.RootHash, nil
}

func (h *BorGRPCClient) GetVoteOnHash(ctx context.Context, startBlock uint64, endBlock uint64, rootHash string, milestoneId string) (bool, error) {

	req := &proto.GetVoteOnHashRequest{
		StartBlockNumber: startBlock,
		EndBlockNumber:   endBlock,
		Hash:             rootHash,
		MilestoneId:      milestoneId,
	}

	log.Info("Fetching vote on hash")

	res, err := h.client.GetVoteOnHash(ctx, req)
	if err != nil {
		return false, err
	}

	log.Info("Fetched vote on hash")

	return res.Response, nil
}

func (h *BorGRPCClient) HeaderByNumber(ctx context.Context, blockID uint64) (*ethTypes.Header, error) {

	req := &proto.GetHeaderByNumberRequest{
		Number: blockID,
	}

	log.Info("Fetching header by number")

	res, err := h.client.HeaderByNumber(ctx, req)
	if err != nil {
		return &ethTypes.Header{}, err
	}

	log.Info("Fetched header by number")

	resp := &ethTypes.Header{
		Number:     big.NewInt(int64(res.Header.Number)),
		ParentHash: protoutil.ConvertH256ToHash(res.Header.ParentHash),
		Time:       res.Header.Time,
	}

	return resp, nil
}

func (h *BorGRPCClient) BlockByNumber(ctx context.Context, blockID uint64) (*ethTypes.Block, error) {

	req := &proto.GetBlockByNumberRequest{
		Number: blockID,
	}

	log.Info("Fetching block by number")

	res, err := h.client.BlockByNumber(ctx, req)
	if err != nil {
		return &ethTypes.Block{}, err
	}

	log.Info("Fetched block by number")

	header := ethTypes.Header{
		Number:     big.NewInt(int64(res.Block.Header.Number)),
		ParentHash: protoutil.ConvertH256ToHash(res.Block.Header.ParentHash),
		Time:       res.Block.Header.Time,
	}
	return ethTypes.NewBlock(&header, nil, nil, nil, nil), nil
}

func (h *BorGRPCClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*ethTypes.Receipt, error) {

	req := &proto.ReceiptRequest{
		Hash: protoutil.ConvertHashToH256(txHash),
	}

	log.Info("Fetching transaction receipt")

	res, err := h.client.TransactionReceipt(ctx, req)
	if err != nil {
		return &ethTypes.Receipt{}, err
	}

	log.Info("Fetched transaction receipt")

	return receiptResponseToTypesReceipt(res.Receipt), nil
}

func (h *BorGRPCClient) BorBlockReceipt(ctx context.Context, txHash common.Hash) (*ethTypes.Receipt, error) {

	req := &proto.ReceiptRequest{
		Hash: protoutil.ConvertHashToH256(txHash),
	}

	log.Info("Fetching bor block receipt")

	res, err := h.client.BorBlockReceipt(ctx, req)
	if err != nil {
		return &ethTypes.Receipt{}, err
	}

	log.Info("Fetched bor block receipt")

	return receiptResponseToTypesReceipt(res.Receipt), nil
}

func receiptResponseToTypesReceipt(receipt *proto.Receipt) *ethTypes.Receipt {
	// Bloom and Logs have been intentionally left out as they are not used in the current implementation
	return &ethTypes.Receipt{
		Type:              uint8(receipt.Type),
		PostState:         receipt.PostState,
		Status:            receipt.Status,
		CumulativeGasUsed: receipt.CumulativeGasUsed,
		TxHash:            protoutil.ConvertH256ToHash(receipt.TxHash),
		ContractAddress:   protoutil.ConvertH160toAddress(receipt.ContractAddress),
		GasUsed:           receipt.GasUsed,
		EffectiveGasPrice: big.NewInt(receipt.EffectiveGasPrice),
		BlobGasUsed:       receipt.BlobGasUsed,
		BlobGasPrice:      big.NewInt(receipt.BlobGasPrice),
		BlockHash:         protoutil.ConvertH256ToHash(receipt.BlockHash),
		BlockNumber:       big.NewInt(int64(receipt.BlockNumber)),
		TransactionIndex:  uint(receipt.TransactionIndex),
	}
}

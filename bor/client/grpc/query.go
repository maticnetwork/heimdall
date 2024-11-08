package grpc

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"

	proto "github.com/maticnetwork/polyproto/bor"
	protoutil "github.com/maticnetwork/polyproto/utils"
)

func (h *BorGRPCClient) GetRootHash(ctx context.Context, startBlock uint64, endBlock uint64) (string, error) {
	fmt.Printf(">>>>> Using BorGRPCClient GetRootHash\n")

	fmt.Printf(">>>>> GetRootHash prepare req\n")
	req := &proto.GetRootHashRequest{
		StartBlockNumber: startBlock,
		EndBlockNumber:   endBlock,
	}
	fmt.Printf(">>>>> GetRootHash req: %v\n", req)

	log.Info("Fetching bor root hash")

	fmt.Printf(">>>>> GetRootHash IN\n")
	res, err := h.client.GetRootHash(ctx, req)
	fmt.Printf(">>>>> GetRootHash OUT\n")
	if err != nil {
		fmt.Printf(">>>>> GetRootHash err: %v\n", err)
		return "", err
	}

	log.Info("Fetched bor root hash")
	fmt.Printf(">>>>> GetRootHash returning res: %v\n", res)

	return res.RootHash, nil
}

func (h *BorGRPCClient) GetVoteOnHash(ctx context.Context, startBlock uint64, endBlock uint64, rootHash string, milestoneId string) (bool, error) {
	fmt.Printf(">>>>> Using BorGRPCClient GetVoteOnHash\n")

	fmt.Printf(">>>>> GetVoteOnHash prepare req\n")
	req := &proto.GetVoteOnHashRequest{
		StartBlockNumber: startBlock,
		EndBlockNumber:   endBlock,
		Hash:             rootHash,
		MilestoneId:      milestoneId,
	}
	fmt.Printf(">>>>> GetVoteOnHash req: %v\n", req)

	log.Info("Fetching vote on hash")

	fmt.Printf(">>>>> GetVoteOnHash IN\n")
	res, err := h.client.GetVoteOnHash(ctx, req)
	fmt.Printf(">>>>> GetVoteOnHash OUT\n")
	if err != nil {
		fmt.Printf(">>>>> GetVoteOnHash err: %v\n", err)
		return false, err
	}

	log.Info("Fetched vote on hash")

	fmt.Printf(">>>>> GetVoteOnHash returning res: %v\n", res)
	return res.Response, nil
}

func (h *BorGRPCClient) HeaderByNumber(ctx context.Context, blockID uint64) (*ethTypes.Header, error) {
	fmt.Printf(">>>>> Using BorGRPCClient HeaderByNumber\n")

	fmt.Printf(">>>>> HeaderByNumber prepare req\n")
	req := &proto.GetHeaderByNumberRequest{
		Number: blockID,
	}
	fmt.Printf(">>>>> HeaderByNumber req: %v\n", req)

	log.Info("Fetching header by number")

	fmt.Printf(">>>>> HeaderByNumber IN\n")
	res, err := h.client.HeaderByNumber(ctx, req)
	fmt.Printf(">>>>> HeaderByNumber OUT\n")
	if err != nil {
		fmt.Printf(">>>>> HeaderByNumber err: %v\n", err)
		return &ethTypes.Header{}, err
	}

	log.Info("Fetched header by number")

	resp := &ethTypes.Header{
		Number:     big.NewInt(int64(res.Header.Number)),
		ParentHash: protoutil.ConvertH256ToHash(res.Header.ParentHash),
		Time:       res.Header.Time,
	}
	fmt.Printf(">>>>> HeaderByNumber returning res: %v\n", resp)

	return resp, nil
}

func (h *BorGRPCClient) BlockByNumber(ctx context.Context, blockID uint64) (*ethTypes.Block, error) {
	fmt.Printf(">>>>> Using BorGRPCClient BlockByNumber\n")

	fmt.Printf(">>>>> BlockByNumber prepare req\n")
	req := &proto.GetBlockByNumberRequest{
		Number: blockID,
	}
	fmt.Printf(">>>>> BlockByNumber req: %v\n", req)

	log.Info("Fetching block by number")

	fmt.Printf(">>>>> BlockByNumber IN\n")
	res, err := h.client.BlockByNumber(ctx, req)
	fmt.Printf(">>>>> BlockByNumber OUT\n")
	if err != nil {
		fmt.Printf(">>>>> BlockByNumber err: %v\n", err)
		return &ethTypes.Block{}, err
	}

	log.Info("Fetched block by number")

	header := ethTypes.Header{
		Number:     big.NewInt(int64(res.Block.Header.Number)),
		ParentHash: protoutil.ConvertH256ToHash(res.Block.Header.ParentHash),
		Time:       res.Block.Header.Time,
	}
	fmt.Printf(">>>>> BlockByNumber returning res: %v\n", header)
	return ethTypes.NewBlock(&header, nil, nil, nil, nil), nil
}

func (h *BorGRPCClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*ethTypes.Receipt, error) {
	fmt.Printf(">>>>> Using BorGRPCClient TransactionReceipt\n")

	fmt.Printf(">>>>> TransactionReceipt prepare req\n")
	req := &proto.ReceiptRequest{
		Hash: protoutil.ConvertHashToH256(txHash),
	}
	fmt.Printf(">>>>> TransactionReceipt req: %v\n", req)

	log.Info("Fetching transaction receipt")

	fmt.Printf(">>>>> TransactionReceipt IN\n")
	res, err := h.client.TransactionReceipt(ctx, req)
	fmt.Printf(">>>>> TransactionReceipt OUT\n")
	if err != nil {
		fmt.Printf(">>>>> TransactionReceipt err: %v\n", err)
		return &ethTypes.Receipt{}, err
	}

	log.Info("Fetched transaction receipt")

	fmt.Printf(">>>>> TransactionReceipt returning res: %v\n", res.Receipt)
	return receiptResponseToTypesReceipt(res.Receipt), nil
}

func (h *BorGRPCClient) BorBlockReceipt(ctx context.Context, txHash common.Hash) (*ethTypes.Receipt, error) {
	fmt.Printf(">>>>> Using BorGRPCClient BorBlockReceipt\n")

	fmt.Printf(">>>>> BorBlockReceipt prepare req\n")
	req := &proto.ReceiptRequest{
		Hash: protoutil.ConvertHashToH256(txHash),
	}
	fmt.Printf(">>>>> BorBlockReceipt req: %v\n", req)

	log.Info("Fetching bor block receipt")

	fmt.Printf(">>>>> BorBlockReceipt IN\n")
	res, err := h.client.BorBlockReceipt(ctx, req)
	fmt.Printf(">>>>> BorBlockReceipt OUT\n")
	if err != nil {
		fmt.Printf(">>>>> BorBlockReceipt err: %v\n", err)
		return &ethTypes.Receipt{}, err
	}

	log.Info("Fetched bor block receipt")

	fmt.Printf(">>>>> BorBlockReceipt returning res: %v\n", res.Receipt)
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

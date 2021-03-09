package util

type ChainParams struct {
	BorChainID            string `json:"bor_chain_id"`
	MaticTokenAddress     string `json:"matic_token_address"`
	StakingManagerAddress string `json:"staking_manager_address"`
	SlashManagerAddress   string `json:"slash_manager_address"`
	RootChainAddress      string `json:"root_chain_address"`
	StakingInfoAddress    string `json:"staking_info_address"`
	StateSenderAddress    string `json:"state_sender_address"`
	StateReceiverAddress  string `json:"state_receiver_address"`
	ValidatorSetAddress   string `json:"validator_set_address"`
}

type ChainmanageParams struct {
	MainchainTxConfirmations  uint64      `json:"mainchain_tx_confirmations,string"`
	MaticchainTxConfirmations uint64      `json:"maticchain_tx_confirmations,string"`
	ChainParams               ChainParams `json:"chain_params"`
}

type ChainmanagerParamsResponse struct {
	Params ChainmanageParams `json:"params"`
}

type Validator struct {
	ID               uint32 `protobuf:"varint,1,opt,name=ID,proto3,enum=heimdall.types.ValidatorID" json:"ID,omitempty"`
	StartEpoch       uint64 `protobuf:"varint,2,opt,name=start_epoch,json=startEpoch,proto3" json:"start_epoch,string,omitempty" yaml:"start_epoch"`
	EndEpoch         uint64 `protobuf:"varint,3,opt,name=end_epoch,json=endEpoch,proto3" json:"end_epoch,string,omitempty" yaml:"end_epoch"`
	Nonce            uint64 `protobuf:"varint,4,opt,name=nonce,proto3" json:"nonce,string,omitempty"`
	VotingPower      int64  `protobuf:"varint,5,opt,name=voting_power,json=votingPower,proto3" json:"voting_power,string,omitempty" yaml:"voting_power"`
	PubKey           string `protobuf:"bytes,6,opt,name=pub_key,json=pubKey,proto3" json:"pub_key,omitempty" yaml:"pub_key"`
	Signer           string `protobuf:"bytes,7,opt,name=signer,proto3" json:"signer,omitempty"`
	LastUpdated      string `protobuf:"bytes,8,opt,name=last_updated,json=lastUpdated,proto3" json:"last_updated,omitempty" yaml:"last_updated"`
	Jailed           bool   `protobuf:"varint,9,opt,name=jailed,proto3" json:"jailed,omitempty"`
	ProposerPriority int64  `protobuf:"varint,10,opt,name=proposer_priority,json=proposerPriority,proto3" json:"proposer_priority,string,omitempty" yaml:"proposer_priority"`
}

//// Validator set response
type ValidatorSet struct {
	Validators       []Validator `json:"validators"`
	Proposer         Validator   `json:"proposer"`
	TotalVotingPower int64       `json:"total_voting_power,string"`
}

type ValidatorSetResponse struct {
	ValidatorSet ValidatorSet `json:"validator_set"`
}

//// checkpoint params response
//type CheckPointParams struct {
//	CheckpointBufferTime time.Duration `protobuf:"bytes,1,opt,name=checkpoint_buffer_time,json=checkpointBufferTime,proto3,stdduration" json:"checkpoint_buffer_time,string" yaml:"checkpoint_buffer_time"`
//	AvgCheckpointLength  uint64        `protobuf:"varint,2,opt,name=avg_checkpoint_length,json=avgCheckpointLength,proto3" json:"avg_checkpoint_length,string,omitempty" yaml:"avg_checkpoint_length"`
//	MaxCheckpointLength  uint64        `protobuf:"varint,3,opt,name=max_checkpoint_length,json=maxCheckpointLength,proto3" json:"max_checkpoint_length,string,omitempty" yaml:"max_checkpoint_length"`
//	ChildBlockInterval   uint64        `protobuf:"varint,4,opt,name=child_block_interval,json=childBlockInterval,proto3" json:"child_block_interval,string,omitempty" yaml:"child_block_interval"`
//}
//
//type CheckPointParamsResponse struct {
//	Params CheckPointParams `json:"params"`
//}

type LastNoAckResponse struct {
	LastNoAck uint64 `json:"last_no_ack,string"`
}

type NextSpanSeedResponse struct {
	NextSpanSeed string `json:"next_span_seed"`
}

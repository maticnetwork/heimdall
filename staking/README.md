# Staking Module

## Table of Contents

* [Preliminary terminology](#preliminary-terminology)
* [Overview](#overview)
* [How does one join the network as a validator](#how-does-one-join-the-network-as-a-validator)
	* [How to propose a MsgValidatorJoin transaction](#how-to-propose-a-msgvalidatorjoin-transaction)
* [How does an existing validator exit the network](#how-does-an-existing-validator-exit-the-network)
	* [How to propose a MsgValidatorExit transaction](#how-to-propose-a-msgvalidatorexit-transaction)
* [How does a validator update its stake](#how-does-a-validator-update-its-stake)
	* [How to propose a MsgStakeUpdate transaction](#how-to-propose-a-msgstakeupdate-transaction)
* [How does a validator update its signer address](#how-does-a-validator-update-its-signer-address)
	* [How to propose a MsgSignerUpdate transaction](#how-to-propose-a-msgsignerupdate-transaction)
* [Query commands](#query-commands)

## Preliminary terminology

* An `epoch` represents the period until a checkpoint is submitted on Ethereum (i.e. one `epoch` ends when a checkpoint is committed on Ethereum and the next one begins).

## Overview

The `staking` module in Heimdall is responsible for a validator's stake related operations. It primarily aids in

* A node joining the protocol as a validator.
* An node leaving the protocol as a validator.
* Updating an existing validator's stake in the network.
* Updating the signer address of an existing validator.

## How does one join the network as a validator

The node that wants to be a validator stakes its tokens by invoking the `stakeFor` method on the `StakeManager` contract on L1 (Ethereum), which emits a `Staked` event:

```
/// @param signer validator address.
/// @param validatorId unique integer to identify a validator.
/// @param nonce to synchronize the events in heimdall.
/// @param activationEpoch validator's first epoch as proposer.
/// @param amount staking amount.
/// @param total total staking amount.
/// @param signerPubkey public key of the validator
event Staked(
    address indexed signer,
    uint256 indexed validatorId,
    uint256 nonce,
    uint256 indexed activationEpoch,
    uint256 amount,
    uint256 total,
    bytes signerPubkey
);
```

An existing validator on Heimdall catches this event and sends a `MsgValidatorJoin` transaction, which is represented by the data structure:

```
type MsgValidatorJoin struct {
	From            hmTypes.HeimdallAddress `json:"from"`
	ID              hmTypes.ValidatorID     `json:"id"`
	ActivationEpoch uint64                  `json:"activationEpoch"`
	Amount          sdk.Int                 `json:"amount"`
	SignerPubKey    hmTypes.PubKey          `json:"pub_key"`
	TxHash          hmTypes.HeimdallHash    `json:"tx_hash"`
	LogIndex        uint64                  `json:"log_index"`
	BlockNumber     uint64                  `json:"block_number"`
	Nonce           uint64                  `json:"nonce"`
}
```
where

* `From` represents the address of the validator that initiated the `MsgValidatorJoin` transaction on heimdall.
* `ID` represents the id of the new validator.
* `ActivationEpoch` is the `epoch` at which the new validator will be activated.
* `Amount` is the total staked amount.
* `SignerPubKey` is the signer public key of the new validator.
* `TxHash` is the hash of the staking transaction on L1.
* `LogIndex` is the index of the `Staked` log in the staking transaction receipt.
* `BlockNumber` is the L1 block number in which the staking transaction was included.
* `Nonce` is the the count representing all the staking related transactions performed from the new validator's account. This is meant to keep Heimdall and L1 in sync.

Upon broadcasting the message, it goes through `HandleMsgValidatorJoin` handler which checks the basic sanity of the transaction (verifying the validator isn't already existing, voting power, etc.).

The `SideHandleMsgValidatorJoin` side-handler in all the existing (honest) validators then ensures the authenticity of staking transaction on L1. It fetches the transaction receipt from L1 contract and validates it with the data provided in the `MsgValidatorJoin` transaction. Upon successful validation, `YES` is voted.

The `PostHandleMsgValidatorJoin` post-handler then initializes the new validator and persists in the state via the keeper:

```
// create new validator
newValidator := hmTypes.Validator{
	ID:          msg.ID,
	StartEpoch:  msg.ActivationEpoch,
	EndEpoch:    0,
	Nonce:       msg.Nonce,
	VotingPower: votingPower.Int64(),
	PubKey:      pubkey,
	Signer:      hmTypes.BytesToHeimdallAddress(signer.Bytes()),
	LastUpdated: "",
}

// update last updated
newValidator.LastUpdated = sequence.String()

// add validator to store
k.Logger(ctx).Debug("Adding new validator to state", "validator", newValidator.String())

if err = k.AddValidator(ctx, newValidator); err != nil {
	k.Logger(ctx).Error("Unable to add validator to state", "validator", newValidator.String(), "error", err)
	return hmCommon.ErrValidatorSave(k.Codespace()).Result()
}
```

The `EndBlocker` hook, which is executed at the end of a heimdall block, then adds the new validator and updates the validator set once activation `epoch` is completed:

```
// --- Start update to new validators
currentValidatorSet := app.StakingKeeper.GetValidatorSet(ctx)
allValidators := app.StakingKeeper.GetAllValidators(ctx)
ackCount := app.CheckpointKeeper.GetACKCount(ctx)

// get validator updates
setUpdates := helper.GetUpdatedValidators(
	&currentValidatorSet, // pointer to current validator set -- UpdateValidators will modify it
	allValidators,        // All validators
	ackCount,             // ack count
)

if len(setUpdates) > 0 {
	// create new validator set
	if err := currentValidatorSet.UpdateWithChangeSet(setUpdates); err != nil {
		// return with nothing
		logger.Error("Unable to update current validator set", "Error", err)
		return abci.ResponseEndBlock{}
	}

	// validator set change
	logger.Debug("[ENDBLOCK] Updated current validator set", "proposer", currentValidatorSet.GetProposer())

	// save set in store
	if err := app.StakingKeeper.UpdateValidatorSetInStore(ctx, currentValidatorSet); err != nil {
		// return with nothing
		logger.Error("Unable to update current validator set in state", "Error", err)
		return abci.ResponseEndBlock{}
	}

// more code
	}
```

### How to propose a MsgValidatorJoin transaction

The `bridge` service in an existing validator's heimdall process polls for `Staked` event periodically and generates and broadcasts the transaction once it detects and parses the event. An existing validator on the network can also leverage the CLI to send the transaction:

```
heimdallcli tx staking validator-join --proposer <PROPOSER_ADDRESS> --tx-hash <ETH_TX_HASH> --signer-pubkey <PUB_KEY> --staked-amount <STAKED_AMOUNT> --activation-epoch <ACTIVATION_EPOCH>
```

Or the REST server :

```
curl -X POST "localhost:1317/staking/validator-join?from=<PROPOSER_ADDRESS>&tx-hash=<ETH_TX_HASH>&signer-pubkey=<PUB_KEY>&staked-amount=<STAKED_AMOUNT>&activation-epoch=<ACTIVATION_EPOCH>"
```

## How does an existing validator exit the network

If a validator wishes to exit the network and unbond its stake, it invokes the `unstake` function on the `StakeManager` contract on L1, which emits an `UnstakeInit` event:

```
/// @param user address of the validator.
/// @param validatorId unique integer to identify a validator.
/// @param nonce to synchronize the events in heimdall.
/// @param deactivationEpoch last epoch for validator.
/// @param amount staking amount.
event UnstakeInit(
    address indexed user,
    uint256 indexed validatorId,
    uint256 nonce,
    uint256 deactivationEpoch,
    uint256 indexed amount
);
```

An existing validator on Heimdall catches and parses this event and sends a `MsgValidatorExit` transaction, which is represented by the data structure:

```
type MsgValidatorExit struct {
	From              hmTypes.HeimdallAddress `json:"from"`
	ID                hmTypes.ValidatorID     `json:"id"`
	DeactivationEpoch uint64                  `json:"deactivationEpoch"`
	TxHash            hmTypes.HeimdallHash    `json:"tx_hash"`
	LogIndex          uint64                  `json:"log_index"`
	BlockNumber       uint64                  `json:"block_number"`
	Nonce             uint64                  `json:"nonce"`
}
```
where

* `From` represents the address of the validator that initiated the `MsgValidatorExit` transaction on heimdall.
* `ID` represents the id of the validator to be unstaked.
* `DeactivationEpoch` is the last `epoch` as a validator.
* `TxHash` is the hash of the unstake transaction on L1.
* `LogIndex` is the index of the `UnstakeInit` log in the unstake transaction receipt.
* `BlockNumber` is the L1 block number in which the unstake transaction was included.
* `Nonce` is the the count representing all the staking related transactions performed from the validator's account.

Upon broadcasting the message, it goes through `HandleMsgValidatorExit` handler which checks the basic sanity of the data in the transaction.

The `SideHandleMsgValidatorExit` side-handler in all the existing (honest) validators then ensures the authenticity of unstake init transaction on L1. It fetches the transaction receipt from L1 contract and validates it with the data provided in the `MsgValidatorExit` transaction. Upon successful validation, `YES` is voted.

The `PostHandleMsgValidatorExit` post-handler then sets the deactivation epoch for the validator and persists in the state via the keeper:

```
validator, ok := k.GetValidatorFromValID(ctx, msg.ID)
if !ok {
	k.Logger(ctx).Error("Fetching of validator from store failed", "validatorID", msg.ID)
	return hmCommon.ErrNoValidator(k.Codespace()).Result()
}

// set end epoch
validator.EndEpoch = msg.DeactivationEpoch

// update last updated
validator.LastUpdated = sequence.String()

// update nonce
validator.Nonce = msg.Nonce

// Add deactivation time for validator
if err := k.AddValidator(ctx, validator); err != nil {
	k.Logger(ctx).Error("Error while setting deactivation epoch to validator", "error", err, "validatorID", validator.ID.String())
	return hmCommon.ErrValidatorNotDeactivated(k.Codespace()).Result()
}
```

The `EndBlocker` hook then updates the validator set once deactivation epoch is completed.

### How to propose a MsgValidatorExit transaction

The `bridge` service in an existing validator's heimdall process polls for `UnstakeInit` event periodically and generates and broadcasts the transaction once it detects and parses the event. An existing validator on the network can also leverage the CLI to send the transaction:

```
heimdallcli tx staking validator-exit --proposer <PROPOSER_ADDRESS> --id <VALIDATOR_ID> --tx-hash <ETH_TX_HASH> --nonce <VALIDATOR_NONCE> --log-index <LOG_INDEX> --block-number <BLOCK_NUMBER> --deactivation-epoch <DEACTIVATION_EPOCH>
```

Or the REST server :

```
curl -X POST "localhost:1317/staking/validator-exit?from=<PROPOSER_ADDRESS>&id=<VALIDATOR_ID>&tx-hash=<ETH_TX_HASH>&nonce=<VALIDATOR_NONCE>&log-index=<LOG_INDEX>&block-number=<BLOCK_NUMBER>&deactivation-epoch=<DEACTIVATION_EPOCH>"
```

## How does a validator update its stake

A validator can update its stake in the network by invoking the `restake` function on the `StakeManager` contract on L1, which emits an `StakeUpdate` event:

```
/// @param validatorId unique integer to identify a validator.
/// @param nonce to synchronize the events in heimdall.
/// @param newAmount the updated stake amount.
event StakeUpdate(
    uint256 indexed validatorId,
    uint256 indexed nonce,
    uint256 indexed newAmount
);
```

On Heimdall, this event is parsed and a `MsgStakeUpdate` transaction is broadcasted, which is represented by the data structure:

```
type MsgStakeUpdate struct {
	From        hmTypes.HeimdallAddress `json:"from"`
	ID          hmTypes.ValidatorID     `json:"id"`
	NewAmount   sdk.Int                 `json:"amount"`
	TxHash      hmTypes.HeimdallHash    `json:"tx_hash"`
	LogIndex    uint64                  `json:"log_index"`
	BlockNumber uint64                  `json:"block_number"`
	Nonce       uint64                  `json:"nonce"`
}
```
where

* `From` represents the address of the validator that initiated the `MsgStakeUpdate` transaction on heimdall.
* `ID` represents the id of the validator whose stake is to be updated.
* `NewAmount` is the new staked amount.
* `TxHash` is the hash of the stake update transaction on L1.
* `LogIndex` is the index of the `StakeUpdate` log in the stake update transaction receipt.
* `BlockNumber` is the L1 block number in which the stake update transaction was included.
* `Nonce` is the the count representing all the staking related transactions performed from the validator's account.

Upon broadcasting the message, it goes through `HandleMsgStakeUpdate` handler which checks the basic sanity of the data in the transaction.

The `SideHandleMsgStakeUpdate` side-handler in all the existing (honest) validators then ensures the authenticity of the stake update transaction on L1. It fetches the transaction receipt from L1 contract and validates it with the data provided in the `MsgStakeUpdate` transaction. Upon successful validation, `YES` is voted.

The `PostHandleMsgStakeUpdate` post-handler then derives the new voting power for the validator and persists in the state via the keeper:

```
// set validator amount
p, err := helper.GetPowerFromAmount(msg.NewAmount.BigInt())
if err != nil {
	return hmCommon.ErrInvalidMsg(k.Codespace(), fmt.Sprintf("Invalid amount %v for validator %v", msg.NewAmount, msg.ID)).Result()
}

validator.VotingPower = p.Int64()

// save validator
err = k.AddValidator(ctx, validator)
if err != nil {
	k.Logger(ctx).Error("Unable to update signer", "ValidatorID", validator.ID, "error", err)
	return hmCommon.ErrSignerUpdateError(k.Codespace()).Result()
}
```

The `EndBlocker` hook then updates the changes in the validator set.

### How to propose a MsgStakeUpdate transaction

The `bridge` service in an existing validator's heimdall process polls for `StakeUpdate` event periodically and generates and broadcasts the transaction once it detects and parses the event. An existing validator on the network can also leverage the CLI to send the transaction:

```
heimdallcli tx staking stake-update --proposer <PROPOSER_ADDRESS> --id <VALIDATOR_ID> --tx-hash <ETH_TX_HASH> --staked-amount <STAKED_AMOUNT> --nonce <VALIDATOR_NONCE> --log-index <LOG_INDEX> --block-number <BLOCK_NUMBER>
```

Or the REST server :

```
curl -X POST "localhost:1317/staking/stake-update?proposer=<PROPOSER_ADDRESS>&id=<VALIDATOR_ID>&tx-hash=<ETH_TX_HASH>&staked-amount=<STAKED_AMOUNT>&nonce=<VALIDATOR_NONCE>&log-index=<LOG_INDEX>&block-number=<BLOCK_NUMBER>"
```

## How does a validator update its signer address

A validator can update its signer address in the network by invoking the the `updateSigner` function on the `StakeManager` contract on L1, which emits an `SignerChange` event:

```
/// @param validatorId unique integer to identify a validator.
/// @param nonce to synchronize the events in heimdall.
/// @param oldSigner old address of the validator.
/// @param newSigner new address of the validator.
/// @param signerPubkey public key of the validator.
event SignerChange(
    uint256 indexed validatorId,
    uint256 nonce,
    address indexed oldSigner,
    address indexed newSigner,
    bytes signerPubkey
);
```

On Heimdall, this event is parsed and a `MsgSignerUpdate` transaction is broadcasted, which is represented by the data structure:

```
type MsgSignerUpdate struct {
	From            hmTypes.HeimdallAddress `json:"from"`
	ID              hmTypes.ValidatorID     `json:"id"`
	NewSignerPubKey hmTypes.PubKey          `json:"pubKey"`
	TxHash          hmTypes.HeimdallHash    `json:"tx_hash"`
	LogIndex        uint64                  `json:"log_index"`
	BlockNumber     uint64                  `json:"block_number"`
	Nonce           uint64                  `json:"nonce"`
}
```
where

* `From` represents the address of the validator that initiated the `MsgSignerUpdate` transaction on heimdall.
* `ID` represents the id of the validator whose signer address is to be updated.
* `NewSignerPubKey` is new public key of the validator.
* `TxHash` is the hash of the signer update transaction on L1.
* `LogIndex` is the index of the `SignerChange` log in the signer update transaction receipt.
* `BlockNumber` is the L1 block number in which the signer update transaction was included.
* `Nonce` is the the count representing all the staking related transactions performed from the validator's account.

Upon broadcasting the message, it goes through `HandleMsgSignerUpdate` handler which checks the basic sanity of the data in the transaction.

The `SideHandleMsgSignerUpdate` side-handler in all the existing (honest) validators then ensures the authenticity of the signer update transaction on L1. It fetches the transaction receipt from L1 contract and validates it with the data provided in the `MsgSignerUpdate` transaction. Upon successful validation, `YES` is voted.

The `PostHandleMsgSignerUpdate` post-handler then updates the signer address for the validator, "unstakes" the validator instance with the old signer address and persists the changes in the state via the keeper:

```
oldValidator := validator.Copy()

// more code
...

// check if we are actually updating signer
if !bytes.Equal(newSigner.Bytes(), validator.Signer.Bytes()) {
	// Update signer in prev Validator
	validator.Signer = hmTypes.HeimdallAddress(newSigner)
	validator.PubKey = newPubKey

	k.Logger(ctx).Debug("Updating new signer", "newSigner", newSigner.String(), "oldSigner", oldValidator.Signer.String(), "validatorID", msg.ID)
} else {
	k.Logger(ctx).Error("No signer change", "newSigner", newSigner.String(), "oldSigner", oldValidator.Signer.String(), "validatorID", msg.ID)
	return hmCommon.ErrSignerUpdateError(k.Codespace()).Result()
}

k.Logger(ctx).Debug("Removing old validator", "validator", oldValidator.String())

// remove old validator from HM
oldValidator.EndEpoch = k.moduleCommunicator.GetACKCount(ctx)

// more code
...

// save old validator
if err := k.AddValidator(ctx, *oldValidator); err != nil {
	k.Logger(ctx).Error("Unable to update signer", "validatorId", validator.ID, "error", err)
	return hmCommon.ErrSignerUpdateError(k.Codespace()).Result()
}

// adding new validator
k.Logger(ctx).Debug("Adding new validator", "validator", validator.String())

// save validator
err := k.AddValidator(ctx, validator)
if err != nil {
	k.Logger(ctx).Error("Unable to update signer", "ValidatorID", validator.ID, "error", err)
	return hmCommon.ErrSignerUpdateError(k.Codespace()).Result()
}
```

The `EndBlocker` hook then updates the validator set upon completion of the epoch.

### How to propose a MsgSignerUpdate transaction

The `bridge` service in an existing validator's heimdall process polls for `SignerChange` event periodically and generates and broadcasts the transaction once it detects and parses the event. An existing validator on the network can also leverage the CLI to send the transaction:

```
heimdallcli tx staking signer-update --proposer <PROPOSER_ADDRESS> --id <VALIDATOR_ID> --new-pubkey <NEW_PUBKEY> --tx-hash <ETH_TX_HASH> --nonce <VALIDATOR_NONCE> --log-index <LOG_INDEX> --block-number <BLOCK_NUMBER>
```

Or the REST server : 

```
curl -X POST "localhost:1317/staking/signer-update?proposer=<PROPOSER_ADDRESS>&id=<VALIDATOR_ID>&new-pubkey=<NEW_PUBKEY>&tx-hash=<ETH_TX_HASH>&nonce=<VALIDATOR_NONCE>&log-index=<LOG_INDEX>&block-number=<BLOCK_NUMBER>"
```

## Query commands

One can run the following query commands from the staking module :

* `validator-info` - Query validator information via validator id or validator address.
* `current-validator-set` - Query the current validator set.
* `staking-power` - Query the current staking power.
* `validator-status` - Query the validator status by validator address.
* `proposer` - Fetch the first `<TIMES>` validators from the validator set, sorted by priority as a checkpoint proposer.
* `current-proposer` - Fetch the validator info selected as proposer of the current checkpoint.
* `is-old-tx` - Check whether the staking transaction is old.

### CLI commands

```
heimdallcli query staking validator-info --id=<VALIDATOR_ID>

OR

heimdallcli query staking validator-info --validator=<VALIDATOR_ADDRESS>
```

```
heimdallcli query staking current-validator-set
```

```
heimdallcli query staking staking-power
```

```
heimdallcli query staking validator-status --validator=<VALIDATOR_ADDRESS>
```

```
heimdallcli query staking proposer --times=<TIMES>
```

```
heimdallcli query staking current-proposer 
```

```
heimdallcli query staking is-old-tx --tx-hash=<ETH_TX_HASH> --log-index=<LOG_INDEX>
```


### REST endpoints

```
curl localhost:1317/staking/validator/<VALIDATOR_ID>

OR

curl localhost:1317/staking/validator/<VALIDATOR_ADDRESS>
```

```
curl localhost:1317/staking/validator-set
```

```
curl localhost:1317/staking/totalpower
```

```
curl localhost:1317/staking/validator-status/<VALIDATOR_ADDRESS>
```

```
curl localhost:1317/staking/proposer/<TIMES>
```

```
curl "localhost:1317/staking/current-proposer
```

```
curl localhost:1317/staking/isoldtx?tx-hash=<ETH_TX_HASH>&log-index=<LOG_INDEX>
```

Some other utility REST URLs:

* To query validator by signer address:

```
curl "localhost:1317/staking/signer/<SIGNER_ADDRESS>
```

* To fetch the first `<TIMES>` validators from the validator set, sorted by priority as a milestone proposer:

```
curl "localhost:1317/staking/milestoneProposer/<TIMES>
```

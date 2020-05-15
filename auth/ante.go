package auth

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/chainmanager"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
)

var (
	// simulation signature values used to estimate gas consumption
	simSecp256k1Pubkey secp256k1.PubKeySecp256k1

	// DefaultFeeInMatic represents default fee in matic
	DefaultFeeInMatic = big.NewInt(10).Exp(big.NewInt(10), big.NewInt(15), nil)

	// DefaultFeeWantedPerTx fee wanted per tx
	DefaultFeeWantedPerTx = sdk.Coins{sdk.Coin{Denom: authTypes.FeeToken, Amount: sdk.NewIntFromBigInt(DefaultFeeInMatic)}}
)

func init() {
	// This decodes a valid hex string into a sepc256k1Pubkey for use in transaction simulation
	bz, _ := hex.DecodeString("035AD6810A47F073553FF30D2FCC7E0D3B1C0B74B61A1AAA2582344037151E143A")
	copy(simSecp256k1Pubkey[:], bz)
}

// SignatureVerificationGasConsumer is the type of function that is used to both consume gas when verifying signatures
// and also to accept or reject different types of PubKey's. This is where apps can define their own PubKey
type SignatureVerificationGasConsumer = func(meter sdk.GasMeter, sig authTypes.StdSignature, params authTypes.Params) sdk.Result

//
// Collect fees interface
//

// FeeCollector interface for fees collector
type FeeCollector interface {
	GetModuleAddress(string) types.HeimdallAddress
	SendCoinsFromAccountToModule(
		sdk.Context,
		types.HeimdallAddress,
		string,
		sdk.Coins,
	) sdk.Error
}

//
// MainTxMsg tx hash
//
type MainTxMsg interface {
	GetTxHash() types.HeimdallHash
	GetLogIndex() uint64
}

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(
	ak AccountKeeper,
	chainKeeper chainmanager.Keeper,
	feeCollector FeeCollector,
	contractCaller helper.IContractCaller,
	sigGasConsumer SignatureVerificationGasConsumer,
) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, res sdk.Result, abort bool) {
		// get module address
		if addr := feeCollector.GetModuleAddress(authTypes.FeeCollectorName); addr.Empty() {
			return newCtx, sdk.ErrInternal(fmt.Sprintf("%s module account has not been set", authTypes.FeeCollectorName)).Result(), true
		}

		// all transactions must be of type auth.StdTx
		stdTx, ok := tx.(authTypes.StdTx)
		if !ok {
			// Set a gas meter with limit 0 as to prevent an infinite gas meter attack
			// during runTx.
			newCtx = SetGasMeter(simulate, ctx, 0)
			return newCtx, sdk.ErrInternal("tx must be StdTx").Result(), true
		}

		// get account params
		params := ak.GetParams(ctx)

		// gas for tx
		gasForTx := params.MaxTxGas // stdTx.Fee.Gas

		amount, ok := sdk.NewIntFromString(params.TxFees)
		if !ok {
			return newCtx, sdk.ErrInternal("Invalid param tx fees").Result(), true
		}
		feeForTx := sdk.Coins{sdk.Coin{Denom: authTypes.FeeToken, Amount: amount}} // stdTx.Fee.Amount

		// new gas meter
		newCtx = SetGasMeter(simulate, ctx, gasForTx)

		// AnteHandlers must have their own defer/recover in order for the BaseApp
		// to know how much gas was used! This is because the GasMeter is created in
		// the AnteHandler, but if it panics the context won't be set properly in
		// runTx's recover call.
		defer func() {
			if r := recover(); r != nil {
				switch rType := r.(type) {
				case sdk.ErrorOutOfGas:
					log := fmt.Sprintf(
						"out of gas in location: %v; gasWanted: %d, gasUsed: %d",
						rType.Descriptor, gasForTx, newCtx.GasMeter().GasConsumed(),
					)
					res = sdk.ErrOutOfGas(log).Result()

					res.GasWanted = gasForTx
					res.GasUsed = newCtx.GasMeter().GasConsumed()
					abort = true
				default:
					panic(r)
				}
			}
		}()

		// validate tx
		if err := tx.ValidateBasic(); err != nil {
			return newCtx, err.Result(), true
		}

		if res := ValidateMemo(stdTx, params); !res.IsOK() {
			return newCtx, res, true
		}

		// stdSigs contains the sequence number, account number, and signatures.
		// When simulating, this would just be a 0-length slice.
		signerAddrs := stdTx.GetSigners()

		if len(signerAddrs) == 0 {
			return newCtx, sdk.ErrNoSignatures("no signers").Result(), true
		}

		if len(signerAddrs) > 1 {
			return newCtx, sdk.ErrUnauthorized("wrong number of signers").Result(), true
		}

		isGenesis := ctx.BlockHeight() == 0

		// fetch first signer, who's going to pay the fees
		signerAcc, res := GetSignerAcc(newCtx, ak, types.AccAddressToHeimdallAddress(signerAddrs[0]))
		if !res.IsOK() {
			return newCtx, res, true
		}

		// deduct the fees
		if !feeForTx.IsZero() {
			res = DeductFees(feeCollector, newCtx, signerAcc, feeForTx)
			if !res.IsOK() {
				return newCtx, res, true
			}

			// reload the account as fees have been deducted
			signerAcc = ak.GetAccount(newCtx, signerAcc.GetAddress())
		}

		// stdSigs contains the sequence number, account number, and signatures.
		// When simulating, this would just be a 0-length slice.
		stdSigs := stdTx.GetSignatures()

		// check signature, return account with incremented nonce
		signBytes := GetSignBytes(newCtx.ChainID(), stdTx, signerAcc, isGenesis)
		signerAcc, res = processSig(newCtx, signerAcc, stdSigs[0], signBytes, simulate, params, sigGasConsumer)
		if !res.IsOK() {
			return newCtx, res, true
		}

		ak.SetAccount(newCtx, signerAcc)

		// TODO: tx tags (?)
		return newCtx, sdk.Result{GasWanted: gasForTx}, false // continue...
	}
}

// GetSignerAcc returns an account for a given address that is expected to sign
// a transaction.
func GetSignerAcc(
	ctx sdk.Context,
	ak AccountKeeper,
	addr types.HeimdallAddress,
) (authTypes.Account, sdk.Result) {
	if acc := ak.GetAccount(ctx, addr); acc != nil {
		return acc, sdk.Result{}
	}
	return nil, sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", addr)).Result()
}

// ValidateMemo validates the memo size.
func ValidateMemo(stdTx authTypes.StdTx, params authTypes.Params) sdk.Result {
	memoLength := len(stdTx.GetMemo())
	if uint64(memoLength) > params.MaxMemoCharacters {
		return sdk.ErrMemoTooLarge(
			fmt.Sprintf(
				"maximum number of characters is %d but received %d characters",
				params.MaxMemoCharacters, memoLength,
			),
		).Result()
	}

	return sdk.Result{}
}

// verify the signature and increment the sequence. If the account doesn't have
// a pubkey, set it.
func processSig(
	ctx sdk.Context,
	acc authTypes.Account,
	sig authTypes.StdSignature,
	signBytes []byte,
	simulate bool,
	params authTypes.Params,
	sigGasConsumer SignatureVerificationGasConsumer,
) (updatedAcc authTypes.Account, res sdk.Result) {

	if res := sigGasConsumer(ctx.GasMeter(), sig, params); !res.IsOK() {
		return nil, res
	}

	if !simulate {
		var pk secp256k1.PubKeySecp256k1
		p, err := authTypes.RecoverPubkey(signBytes, sig.Bytes())
		copy(pk[:], p[:])

		if err != nil || !bytes.Equal(acc.GetAddress().Bytes(), pk.Address().Bytes()) {
			return nil, sdk.ErrUnauthorized("signature verification failed; verify correct account sequence and chain-id").Result()
		}

		if acc.GetPubKey() == nil {
			var cryptoPk crypto.PubKey = pk
			if err := acc.SetPubKey(cryptoPk); err != nil {
				return nil, sdk.ErrUnauthorized("error while updating account pubkey").Result()
			}
		}
	}

	if err := acc.SetSequence(acc.GetSequence() + 1); err != nil {
		return nil, sdk.ErrUnauthorized("error while updating account sequence").Result()
	}

	return acc, res
}

// DefaultSigVerificationGasConsumer is the default implementation of SignatureVerificationGasConsumer. It consumes gas
// for signature verification based upon the public key type. The cost is fetched from the given params and is matched
// by the concrete type.
func DefaultSigVerificationGasConsumer(
	meter sdk.GasMeter, sig authTypes.StdSignature, params authTypes.Params,
) sdk.Result {
	meter.ConsumeGas(params.SigVerifyCostSecp256k1, "ante verify: secp256k1")
	return sdk.Result{}
}

// DeductFees deducts fees from the given account.
//
// NOTE: We could use the CoinKeeper (in addition to the AccountKeeper, because
// the CoinKeeper doesn't give us accounts), but it seems easier to do this.
func DeductFees(feeCollector FeeCollector, ctx sdk.Context, acc authTypes.Account, fees sdk.Coins) sdk.Result {
	blockTime := ctx.BlockHeader().Time
	coins := acc.GetCoins()

	if !fees.IsValid() {
		return sdk.ErrInsufficientFee(fmt.Sprintf("invalid fee amount: %s", fees)).Result()
	}

	// verify the account has enough funds to pay for fees
	_, hasNeg := coins.SafeSub(fees)
	if hasNeg {
		return sdk.ErrInsufficientFunds(
			fmt.Sprintf("insufficient funds to pay for fees; %s < %s", coins, fees),
		).Result()
	}

	// Validate the account has enough "spendable" coins
	spendableCoins := acc.SpendableCoins(blockTime)
	if _, hasNeg := spendableCoins.SafeSub(fees); hasNeg {
		return sdk.ErrInsufficientFunds(
			fmt.Sprintf("insufficient funds to pay for fees; %s < %s", spendableCoins, fees),
		).Result()
	}

	err := feeCollector.SendCoinsFromAccountToModule(ctx, acc.GetAddress(), authTypes.FeeCollectorName, fees)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

// SetGasMeter returns a new context with a gas meter set from a given context.
func SetGasMeter(simulate bool, ctx sdk.Context, gasLimit uint64) sdk.Context {
	// In various cases such as simulation and during the genesis block, we do not
	// meter any gas utilization.
	if simulate || ctx.BlockHeight() == 0 {
		return ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
	}

	return ctx.WithGasMeter(sdk.NewGasMeter(gasLimit))
}

// GetSignBytes returns a slice of bytes to sign over for a given transaction
// and an account.
func GetSignBytes(chainID string, stdTx authTypes.StdTx, acc authTypes.Account, genesis bool) []byte {
	var accNum uint64
	if !genesis {
		accNum = acc.GetAccountNumber()
	}

	return authTypes.StdSignBytes(chainID, accNum, acc.GetSequence(), stdTx.Msg, stdTx.Memo)
}

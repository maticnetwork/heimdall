package app

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	slashingTypes "github.com/maticnetwork/heimdall/slashing/types"
	"github.com/stretchr/testify/require"

	hmTypes "github.com/maticnetwork/heimdall/types"
)

func TestSlashingInfoRLPEncoding(t *testing.T) {
	var slashingInfoList []*hmTypes.ValidatorSlashingInfo

	// Input data
	slashingInfo1 := hmTypes.NewValidatorSlashingInfo(1, "1000", false)
	slashingInfo2 := hmTypes.NewValidatorSlashingInfo(2, "234", true)
	slashingInfoList = append(slashingInfoList, &slashingInfo1)
	slashingInfoList = append(slashingInfoList, &slashingInfo2)

	// Encoding
	encodedSlashInfos, err := slashingTypes.SortAndRLPEncodeSlashInfos(slashingInfoList)
	t.Log("RLP encoded", "encodedSlashInfos", hex.EncodeToString(encodedSlashInfos), "error", err)
	require.Empty(t, err)

	// Decoding
	decodedSlashInfoList, err := slashingTypes.RLPDecodeSlashInfos(encodedSlashInfos)
	require.Empty(t, err)
	t.Log("RLP Decoded data", "valID", decodedSlashInfoList[0].ID, "amount", decodedSlashInfoList[0].SlashedAmount, "isJailed", decodedSlashInfoList[0].IsJailed)
	t.Log("RLP Decoded data", "valID", decodedSlashInfoList[1].ID, "amount", decodedSlashInfoList[1].SlashedAmount, "isJailed", decodedSlashInfoList[1].IsJailed)

	// Assertions
	for i := 0; i < len(slashingInfoList); i++ {
		require.Equal(t, slashingInfoList[i].ID, decodedSlashInfoList[i].ID, "ID mismatch between slashInfoList and decodedSlashInfoList")
		require.Equal(t, slashingInfoList[i].SlashedAmount, decodedSlashInfoList[i].SlashedAmount, "Amount mismatch between slashInfoList and decodedSlashInfoList")
		require.Equal(t, slashingInfoList[i].IsJailed, decodedSlashInfoList[i].IsJailed, "JailStatus mismatch between slashInfoList and decodedSlashInfoList")
	}
}

func TestSlashingInfoRLPDecoding(t *testing.T) {
	// input data
	slashInfoEncodedBytesStr := "d9d8019532323735303030303030303030303030303030303080"
	slashInfoEncodedBytes, err := hex.DecodeString(slashInfoEncodedBytesStr)
	require.Empty(t, err)

	// decoding input
	slashInfos, err := slashingTypes.RLPDecodeSlashInfos(slashInfoEncodedBytes)
	require.Empty(t, err)
	t.Log("RLP decoded data", "slashInfos", slashInfos)

	// hash of slashInfos
	slashInfoHash, err := slashingTypes.GenerateInfoHash(slashInfos)
	require.Empty(t, err)
	t.Log("hashing", "slashInfoHash", hex.EncodeToString(slashInfoHash))

	// calculate hash manually of encoded slash info
	h := sha256.New()
	_, err = h.Write(slashInfoEncodedBytes)
	require.Empty(t, err)
	expectedSlashInfoHash := h.Sum(nil)
	t.Log("calculated slashinfo hash", "expectedSlashInfoHash", hex.EncodeToString(expectedSlashInfoHash))

	require.Equal(t, expectedSlashInfoHash, slashInfoHash)
}

func TestTickMsgVoteBytes(t *testing.T) {

	// input data
	vote := "f68f6865696d64616c6c2d71443073524a0282070380a03dd35fdab044f1c11e4a7a6e524bd5892b37be39e6fae14a70d7d16292d5e49c"
	sigs := "8e61a47b0481e1974b90d14518d27aaeece80cd181145446d749fc3871af4ecb431085d554f708fd0a7a3abd7046b0a008a548b164b2d5ffe42c6e00133b17dc01"
	slashInfoList := "d9d8019531343737303030303030303030303030303030303080"
	txData := "52FA7B1AF87BF6942CB71687CEB1FD18646A34BB16A6A1B1AC9FEF72A058B3FCB2AAC675F1F238FEB9EDD76933BE38FDF4E5386D408CD77BE19A807A77B8416C4A74484B5CAB6C92FAA8638825A30F244F1DB64F260369FA2EC41AE63491FB1C28DF6FA9D8BC2E981A26FDC1931340DDDD4C9FBBA483B0CE4DC20E4B8DFF3F0080"
	proposer := "2cb71687ceb1fd18646a34bb16a6a1b1ac9fef72"

	voteBytes, _ := hex.DecodeString(vote)
	sigsBytes, _ := hex.DecodeString(sigs)
	slashInfoBytes, _ := hex.DecodeString(slashInfoList)
	proposerBytes, _ := hex.DecodeString(proposer)
	txDataBytes, _ := hex.DecodeString(txData)

	proposerAddr := hmTypes.BytesToHeimdallAddress(proposerBytes)
	slashInfoHash := hmTypes.BytesToHeimdallHash(Sha256Hash(slashInfoBytes))
	// msg := NewMsgTick(proposerAddr, slashInfoHash)
	// name := fmt.Sprintf("%s::%s", msg.Route(), msg.Type())
	// txDataPulpBytes := append(GetPulpHash(name), txDataBytes[:]...)

	// 1. check voteData[4] == sha256(txData)
	tickVote, err := slashingTypes.RLPDeocdeTickVoteBytes(voteBytes)
	require.Empty(t, err)
	t.Log("RLP decoded tick vote data", "chainId", tickVote.ChainID, "voteData", hex.EncodeToString(tickVote.Data), "proposerAddr", proposerAddr, "slashInfoHash", slashInfoHash)

	txDataHash := Sha256Hash(txDataBytes[4:])
	t.Log("Hash of tx data", "txData", hex.EncodeToString(txDataHash), "tickVote.Data", hex.EncodeToString(tickVote.Data))
	require.Equal(t, hex.EncodeToString(txDataHash), hex.EncodeToString(tickVote.Data))

	t.Log("FYI", "sigsBytes", hex.EncodeToString(sigsBytes), "proposerBytes", hex.EncodeToString(proposerBytes), "slashInfoBytes", hex.EncodeToString(slashInfoBytes))

	// 2. check proposer address in txData
	cdc := MakeCodec()
	pulp := MakePulp()
	decoder := authTypes.RLPTxDecoder(cdc, pulp)
	tx, err := decoder(txDataBytes)
	require.Empty(t, err)
	t.Log("RLP decoded tx", "tx", tx)
	t.Log("RLP decoded tx Msgs", "tx", tx.GetMsgs())
	t.Log("RLP decoded tx 1st Msg ", "tx", tx.GetMsgs()[0])
	msg := tx.GetMsgs()[0]
	tickMsg := msg.(slashingTypes.MsgTick)
	t.Log("RLP decoded tx proposer in 1st Msg ", "tx", tickMsg.Proposer)
	require.Equal(t, proposerAddr, tickMsg.Proposer)
	t.Log("RLP decoded tx slashInfoHash in 1st Msg ", "tx", tickMsg.SlashingInfoHash)
	require.Equal(t, slashInfoHash, tickMsg.SlashingInfoHash)
	// stdTx, err := RLPDeocdeStdTxBytes(txDataBytes)
	// t.Log("RLP decoded tx data", "stdTx", stdTx, "error", err)

	// pulp := app.MakePulp()
	// msgsEncodedBytes, err := rlp.EncodeToBytes(msgs)

	// t.Log("RLP encoded msgs", "msgsEncodedBytes", hex.EncodeToString(msgsEncodedBytes))

	/*
		// calculate hash manually of encoded msgs
		h := sha256.New()
		_, err = h.Write(msgsEncodedBytes)
		require.Empty(t, err)
		expectedMsgHash := h.Sum(nil)
		t.Log("calculated msg hash", "expectedMsgHash", hex.EncodeToString(expectedMsgHash))

		expectedMsgTMHash := tmhash.Sum(msgsEncodedBytes)
		t.Log("calculated msg hash", "expectedMsgTMHash", hex.EncodeToString(expectedMsgTMHash)) */
}

func Sha256Hash(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	dataHash := h.Sum(nil)
	return dataHash
}

func TestMakePulp(t *testing.T) {
	pulp := MakePulp()
	require.NotNil(t, pulp, "Pulp should be nil")
}

func TestGetMaccPerms(t *testing.T) {
	dup := GetMaccPerms()
	require.Equal(t, maccPerms, dup, "duplicated module account permissions differed from actual module account permissions")
}

/*
 1. check voteData[4] == sha256(txData)
 2. check sha256(txData[0]) == sha256(encode(proposer, sha256(slashInfoList)))
*/

/*
Dummy Test Data - Sending new tick
vote=f48f6865696d64616c6c2d71443073524a023f80a0d989c8ce6edbd12eee8ff7b48cdd988a9c348275fecde1a8b52be9bcd4bc3728
sigs=b1439a46c62c7a7505489ba2cb0c5fd000d8b1f829ea59e7719808d975d547a1337975a8f81980a33443f7de3d00043204f04c783c174a4f2a06491c9724035c00
slashInfoList=d9d8019532323735303030303030303030303030303030303080
txData=f87bf6942cb71687ceb1fd18646a34bb16a6a1b1ac9fef72a01a1c9bf83d2124b3af3c8b8b72b9cd0b1b0b02cface6e3f58c1b006261922515b841f2f5732e6fb94a132faee17b37111970a7f54499e720e6f456583d5993b51ca820c186c0220d29a52b2faa15b62f258e665a660c71664babb1c8242640ee28ff0180
proposer=0x2cb71687ceb1fd18646a34bb16a6a1b1ac9fef72

Sending new tick
vote=f68f6865696d64616c6c2d71443073524a0282070380a03dd35fdab044f1c11e4a7a6e524bd5892b37be39e6fae14a70d7d16292d5e49c
sigs=8e61a47b0481e1974b90d14518d27aaeece80cd181145446d749fc3871af4ecb431085d554f708fd0a7a3abd7046b0a008a548b164b2d5ffe42c6e00133b17dc01
slashInfoList=d9d8019531343737303030303030303030303030303030303080
txData=f87bf6942cb71687ceb1fd18646a34bb16a6a1b1ac9fef72a058b3fcb2aac675f1f238feb9edd76933be38fdf4e5386d408cd77be19a807a77b8416c4a74484b5cab6c92faa8638825a30f244f1db64f260369fa2ec41ae63491fb1c28df6fa9d8bc2e981a26fdc1931340dddd4c9fbba483b0ce4dc20e4b8dff3f0080 proposer=0x2cb71687ceb1fd18646a34bb16a6a1b1ac9fef72

*/

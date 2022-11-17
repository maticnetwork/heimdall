package milestone

/*
Milestone module is responsible for validating milestone in heimdall.

Sending milestone is a 2 phase process.
1. Send `MsgMilestone`: Here the transaction sender proposes the new milestone by sending the start block, end block and the roothash of the new checkpoint
						which is basically the Merkle Root of all the blocks from start to end.
2. Validate this by `handleMsgMilestone`: Here the transaction is validated by fetching all the headers and validating if the incoming milestone proposal is valid.

*/

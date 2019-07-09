// consumes all events from respective queues
// Deposit Event --> Mint transaction on BOR on the basis of validator set% deposit index
// Withdraw Event --> Burn transaction on BOR
// Validator Join/Exit/Power-change --> Validator set changes on BOR
// Checkpoint Propose --> MsgCheckpoint on Heimdall
// Checkpoint ACK --> MsgCheckpointACK on Heimdall
// Checkpoint NO-ACK --> Sends MsgCheckpointNoACK after x interval on Heimdall
// Validator Join/Exit/Power-change --> Validator set changes on Heimdall
package pier

#!/bin/bash

EthAmount='10'
machine0='localhost'

echo "📍Transferring funds from ganache account[0] to others..."

src="${machine0}:~/matic-cli/devnet/devnet/signer-dump.json"

signerDump=$(cat $src | jq -c '.')

rootChainWeb3="http://${machine0}:9545"

for ((i = 1; i < ${#signerDump[@]}; i++)); do
  to_address=$(echo $signerDump | jq -r ".[$i].address")
  from_address=$(echo $signerDump | jq -r ".[0].address")

  txReceipt=$(curl -X POST --data '{"jsonrpc":"2.0","method":"eth_sendTransaction","params":[{"to":"'$to_address'","from":"'$from_address'","value":"'$EthAmount'"}],"id":1}' -H "Content-Type: application/json" $rootChainWeb3)

  txHash=$(echo $txReceipt | jq -r '.result')

  echo "📍Funds transferred from $from_address to $to_address with txHash $txHash"
done

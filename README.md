# TendermintExamples
It contains the sample code for tendermint

## Download CustomExample to your local Machine to test this Dapp

## Description
### Prerequisites
1. GoLang
2. Tendermint
3. MongoDb
4. Ubuntu

### Description
Here we are trying to perform simple addition and subtraction operations on tendermint and trying to explore its functionality. 
There are two files 
##### 1. myApp.go  : used as a Server. Creating a connection to the MongoDB
##### 2. Operations.go : Implement all the methods like Info, Check_tx, Deliver_tx, Commit and Query of Tendermint and some custom methods like Addition and Subtraction.
##### 3. Format for sending the Input:
  ###### 3.a: deliver_tx "username,int,int,operation" like deliver_tx "Prayag,1,2,add"
  ###### 3.b: check_tx "username,int,int,operation"
  ###### 3.c  Info
  ###### 3.d: commit
  ###### 3.e: query "operation" like query "add/sub"

#### Steps to start the App
##### 1. Open a terminal and start myApp.go using command "go run myApp.go" 
##### 2. There are two ways to send the transaction. 
 ######  2.a: Abci-cli console : https://tendermint.com/docs/app-dev/abci-cli.html#using-abci-cli 
 ######  2.b: Rpc client : https://tendermint.com/docs/tendermint-core/using-tendermint.html#transactions 
    
### Using Abci-cli console
##### 1. Run "abci-cli console" command on console
##### 2. Steps to send, check and query the transactions
  ###### 2.a: deliver_tx "Prayag,1,2,add" 
  ###### 2.b: check_tx "Prayag,1,2,add"
  ###### 2.c: info
  ###### 2.d: commit
  ###### 2.e: query "add/sub"
  
###### Using RPC client
1. To send, check and commit the txn -->  curl http://localhost:26657/broadcast_tx_commit?tx=\"Prayag,1,2,add\"
2. To query the txn ---> curl 'localhost:26657/abci_query?data="add"'
# NOTE:
1.  To send the txns using RPC first we need to start the tendermint node using command "tendermint node" on a console
2. Always run "tendermint unsafe_reset_all" if you don't want previous blockchain data (never run this command on production because it will wipe out all the data) and "tendermint init" before "tendermint node" command.

    

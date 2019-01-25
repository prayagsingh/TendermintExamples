package main

import (
	"fmt"
	"os"

	"CustomExample/OurApp"

	"github.com/tendermint/tendermint/abci/server"
	"github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	mgo "gopkg.in/mgo.v2"
)

func main() {
	initExample()
}

func initExample() error {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	fmt.Print("\n Value of Logger is: ", logger)
	// Create the application
	// application is an interface that enables any finite, deterministic state machine to be driven
	//   by the blockchain-based replication engine via the ABCI
	// All the methods under this interface take a Request args and returns Response args
	// except CheckTx/DeliverTx that takes `tx []byte` and Commit which takes nothing
	var app types.Application

	// Obtaining a session
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	// Creating a database name. if its empty then it takes same name provided in dial fucn
	// if its also empty then by default it takes test as a db name
	db := session.DB("MyFirstApp")

	//clean the db on each reboot. no need to clear it now
	//collections := [2]string{"add", "sub"}
	//for _, collection := range collections {
	//	db.C(collection).RemoveAll(nil)
	//}

	// passing mongoDb reference to our custom app
	app = OurApp.NewOperationApplication(db)

	// Start the listner
	srv, err := server.NewServer("tcp://0.0.0.0:26658", "socket", app)
	//fmt.Print("\n before panic at line 51 in myapp.go")
	if err != nil {
		panic(err)
	}

	//fmt.Print("\n after panic at line 56 in myapp.go")
	srv.SetLogger(logger.With("module", "abci-server"))
	//fmt.Print("\n value of checkError is: ", checkError)

	if err := srv.Start(); err != nil {
		return err
	}

	//Wait forever
	cmn.TrapSignal(func() {
		//Cleanup
		srv.Stop()
	})
	return nil

}

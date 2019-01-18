package OurApp

import (
	"fmt"
	//"reflect"
	"encoding/binary"
	"mint/code"
	"strconv"
	"strings"

	"github.com/tendermint/tendermint/abci/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ types.Application = (*OperationApplication)(nil)
var db *mgo.Database

// Structure to store the data to the mongoDb
type StoreResult struct {
	TxnID    bson.ObjectId `bson:"_id" json:"_id"`
	User     string        `bson:"username" json:"username"`
	ParamOne int           `bson:"paramone" json:"paramone"`
	ParamTwo int           `bson:"paramtwo" json:"paramtwo"`
	Result   string        `bson:"result" json:"result"`
}

// Our custom methods/structure
func add(a int, b int) string {
	c := 0
	c = a + b
	//fmt.Print("\n Inside add function and value of c is: ",c)
	return (fmt.Sprintf("%02x", c)) // '0' force using zero, '2' set the output size as two char, 'x' convert to hex
}

func sub(a int, b int) string {
	c := 0
	c = a - b
	//fmt.Print("\n Inside sub function and value of c is: ",c)
	return (fmt.Sprintf("%02x", c))
}

// our application which will use Baseapplication under application.go
type OperationApplication struct {
	types.BaseApplication
}

//
func NewOperationApplication(dbCopy *mgo.Database) *OperationApplication {
	db = dbCopy
	return &OperationApplication{}
}

// tendermint methods like Info, DeliverTx, CheckTx, Query, Commit, SetOption
func (app *OperationApplication) Info(req types.RequestInfo) (resInfo types.ResponseInfo) {
	return types.ResponseInfo{Data: fmt.Sprintf("{\"size\":%v}", 0)}
}

func (app *OperationApplication) DeliverTx(tx []byte) types.ResponseDeliverTx {
	//fmt.Print("\n Inside DeliverTx Module and value of tx is: ", tx)
	//fmt.Print("\n  Length of tx is: ", len(tx))

	// convert from ascii to respective values
	txStr := string(tx)
	//fmt.Printf("\n Value of txStr is %s ", txStr)

	//Splitting the above string using strings package
	txSplit := strings.Split(txStr, ",")

	//fmt.Print("\n Value of txSplit is: ", txSplit)

	// Convert bytes into int
	userNameParam := txSplit[0]
	firstParam, err := strconv.Atoi(txSplit[1])
	secondParam, err := strconv.Atoi(txSplit[2])

	if err != nil {
		panic(err)
	}

	//fmt.Println("\n Checking the type of tx[2] is ", reflect.TypeOf(txSplit[2]))
	//fmt.Printf("\n Values of tx[3] is %s, Value of tx[2] is %d, Value of tx[1] is %d, Value of username is %s", txSplit[3], firstParam, secondParam, userNameParam)

	var result string
	var storeresult StoreResult
	dataHash := make([]byte, 5)

	switch txSplit[3] {
	case "add":
		result = add(firstParam, secondParam)
	case "sub":
		result = sub(firstParam, secondParam)
	default:
		fmt.Println("\n None of the option is selected")
	}

	copy(dataHash[:], result)
	fmt.Println("\n Result is: ", result, " value of data is: ", dataHash)

	// assigning values to the storeresult struct
	storeresult.TxnID = bson.NewObjectId()
	storeresult.User = userNameParam
	storeresult.ParamOne = firstParam
	storeresult.ParamTwo = secondParam
	storeresult.Result = result

	mongoError := db.C(txSplit[3]).Insert(storeresult)

	if mongoError != nil {
		panic(mongoError)
	}

	return types.ResponseDeliverTx{Code: code.CodeTypeOK, Data: dataHash, Log: "Txn executed check data for result", Tags: nil}
}

func (app *OperationApplication) CheckTx(tx []byte) types.ResponseCheckTx {
	// TODO
	// 1. need to unmarshel the data since mongodb stores data in json format
	// 2. fetch the stored data to check if its not empty
	// 3. will fetch the data by username. it might return all the txns correspondig to the username

	return types.ResponseCheckTx{Code: code.CodeTypeOK}
}

func (app *OperationApplication) Commit() types.ResponseCommit {
	appHash := make([]byte, 8)
	var count int64 = 100
	binary.PutVarint(appHash, count)
	return types.ResponseCommit{Data: appHash}
}

func (app *OperationApplication) Query(reqQuery types.RequestQuery) (resQuery types.ResponseQuery) {
	return
}

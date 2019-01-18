package OurApp

import (
	"fmt"
	//"reflect"
	"strings"
	"strconv"
	"github.com/tendermint/tendermint/abci/types"
	mgo "gopkg.in/mgo.v2"
	"encoding/binary"
	"mint/code"
)


var _ types.Application = (*OperationApplication)(nil)
var db *mgo.Database

// Our custom methods/structure 
func add(a int, b int) string{
	c:= 0
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
	return  &OperationApplication{}
}

// tendermint methods like Info, DeliverTx, CheckTx, Query, Commit, SetOption
func(app *OperationApplication) Info(req types.RequestInfo) (resInfo types.ResponseInfo) {
	return types.ResponseInfo{Data: fmt.Sprintf("{\"size\":%v}",0)}
}

func(app *OperationApplication) DeliverTx(tx []byte) types.ResponseDeliverTx {
	//fmt.Print("\n Inside DeliverTx Module and value of tx is: ", tx)
	//fmt.Print("\n  Length of tx is: ", len(tx))

	// convert from ascii to respective values
	txStr := string(tx)
	//fmt.Printf("\n Value of txStr is %s ", txStr)

	//Splitting the above string using strings package
	txSplit := strings.Split(txStr, ",")

	//fmt.Print("\n Value of txSplit is: ", txSplit)

	// Convert bytes into int
	firstParam, err := strconv.Atoi(txSplit[0])
	secondParam, err := strconv.Atoi(txSplit[1])

	if err != nil {
		panic(err)
	}

	//fmt.Println("\n Checking the type of tx[2] is ", reflect.TypeOf(txSplit[2]))
	//fmt.Printf("\n Values of tx[2] is %s, Value of tx[1] is %d, Value of tx[0] is %d",txSplit[2],firstParam,secondParam)

	var result string 
	dataHash := make([]byte, 5)

	switch  txSplit[2] {
	case "add":
			result = add(firstParam, secondParam)
	case "sub":
			result = sub(firstParam, secondParam)
	default:
			fmt.Println("\n None of the option is selected")
	}

	copy(dataHash[:], result)
	fmt.Println("\n Result is: ",result, " value of data is: ", dataHash)
	return types.ResponseDeliverTx{Code: code.CodeTypeOK, Data: dataHash, Log: "Txn executed check data for result", Tags: nil}
}

func(app *OperationApplication) CheckTx(tx []byte) types.ResponseCheckTx {

	return types.ResponseCheckTx{Code: code.CodeTypeOK}
}  

func(app *OperationApplication) Commit() types.ResponseCommit {
	appHash := make([]byte, 8)
	var count int64 = 100
	binary.PutVarint(appHash, count)
	return types.ResponseCommit{Data: appHash}
}

func(app *OperationApplication) Query(reqQuery types.RequestQuery) (resQuery types.ResponseQuery){
	return
}

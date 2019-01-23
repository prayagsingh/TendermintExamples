package OurApp

import (
	"fmt"
	"time"

	//"reflect"
	"encoding/binary"
	"strconv"
	"strings"

	"github.com/tendermint/tendermint/abci/example/code"

	"github.com/tendermint/tendermint/abci/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ types.Application = (*OperationApplication)(nil)
var db *mgo.Database

// Structure to store the data to the mongoDb
type StoreResult struct {
	TxnID     bson.ObjectId `bson:"_id"       json:"_id"`
	User      string        `bson:"username"  json:"username"`
	ParamOne  int           `bson:"paramone"  json:"paramone"`
	ParamTwo  int           `bson:"paramtwo"  json:"paramtwo"`
	Result    string        `bson:"result"    json:"result"`
	Timestamp time.Time     `bson:"timestamp" json:"timestamp"`
}

// State of the block
type State struct {
	AppHash []byte `bson:"app_hash"`
	Height  int64  `bson:"height"`
}

// function to load the previous state
func loadState(loaddb *mgo.Database) State {
	var state State
	return state
}

// Our custom methods/structure
func add(a int, b int) string {
	c := 0
	c = a + b
	fmt.Print("\n Inside add function and value of c is: ", c)
	return (fmt.Sprintf("%05x", c)) // '0' force using zero, '2' set the output size as two char, 'x' convert to hex
}

func sub(a int, b int) string {
	c := 0
	c = a - b
	//fmt.Print("\n Inside sub function and value of c is: ",c)
	return (fmt.Sprintf("%05x", c))
}

// OurApplication which will use Baseapplication which is under application.go
type OperationApplication struct {
	types.BaseApplication
	txnCount int // stores the count of all the total deliver_tx
	hashCout int // stores the count of total commit
	//state    State
}

// NewOperationApplication which returns OperationApplication struct
func NewOperationApplication(dbCopy *mgo.Database) *OperationApplication {
	db = dbCopy
	fmt.Print("\n Inside NewOperation function")
	return &OperationApplication{}
}

// Tendermint methods

// Info method
func (app *OperationApplication) Info(req types.RequestInfo) (resInfo types.ResponseInfo) {
	return types.ResponseInfo{Data: fmt.Sprintf("{\"hashes\":%v,\"txns\":%v}", app.hashCout, app.txnCount)}
}

// DeliverTx method
func (app *OperationApplication) DeliverTx(tx []byte) types.ResponseDeliverTx {
	fmt.Print("\n Inside DeliverTx Module and value of tx is: ", tx)
	time.Sleep(5)
	fmt.Print("\n  Length of tx is: ", len(tx))

	// convert from ascii to respective values
	txStr := string(tx)
	fmt.Printf("\n Value of txStr is %s ", txStr)

	//Splitting the above string using strings package
	txSplit := strings.Split(txStr, ",")

	fmt.Print("\n Value of txSplit is: ", txSplit)

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

	switch txSplit[3] {
	case "add":
		result = add(firstParam, secondParam)
	case "sub":
		result = sub(firstParam, secondParam)
	default:
		fmt.Println("\n None of the option is selected")
	}

	dataHash := make([]byte, 8)
	copy(dataHash[:], result)
	//resultParse, err := strconv.Atoi(result)
	//binary.PutUvarint(dataHash, uint64(resultParse))
	fmt.Println("\n Result is: ", result, " value of data is: ", dataHash)

	// assigning values to the storeresult struct
	storeresult.TxnID = bson.NewObjectId()
	storeresult.User = userNameParam
	storeresult.ParamOne = firstParam
	storeresult.ParamTwo = secondParam
	storeresult.Result = result
	storeresult.Timestamp = time.Now()

	mongoError := db.C(txSplit[3]).Insert(storeresult)

	if mongoError != nil {
		panic(mongoError)
	}

	// converting time.now to string. Value inside Format is fixed by GoLang devs else get random val
	TxnTime := storeresult.Timestamp.Format("2006-01-02 15:04:05.000000000")
	fmt.Println("\n Val of TxnTime is: ", TxnTime)

	// increasing the txn Count
	app.txnCount++

	return types.ResponseDeliverTx{
		Code: code.CodeTypeOK,
		Data: dataHash,
		Log:  fmt.Sprintf("Txn executed with timestamp : %s and TxnCount is %d", TxnTime, app.txnCount)}
}

// CheckTx method
func (app *OperationApplication) CheckTx(tx []byte) types.ResponseCheckTx {
	// pattern for sending checktx = check_tx "username,val1,val2,method"
	storeresult := StoreResult{}

	// convert from ascii to respective values
	txStr := string(tx)
	fmt.Print("\n Inside the Check_tx module nd value of tx is: ", txStr)
	txSplit := strings.Split(txStr, ",")
	fmt.Print("\n Value of txSplit is: ", txSplit)

	//fetching the data from the mongodb for verification
	//errMongo := db.C(txSplit[1]).Find(bson.M{"username": txSplit[0]}).Select(bson.M{"username": txSplit[0]}).One(&temp)
	errMongo := db.C(txSplit[3]).Find(bson.M{"username": txSplit[0]}).Sort("-timestamp").One(&storeresult)

	if errMongo != nil {
		panic(errMongo)
	}

	//fmt.Printf("\n Value of query to mongo is: %+v", storeresult)
	fmt.Print("\n Txn Id is : ", storeresult.TxnID,
		" UserName is: ", storeresult.User,
		" Result is : ", storeresult.Result)

	// Data takes byte type in return statement
	dataHash := make([]byte, 5)
	copy(dataHash[:], storeresult.Result)

	// checking if the data is empty. If yes then return badCode
	fmt.Print("\n value of txnCount is: ", app.txnCount)
	if app.txnCount == 0 {
		return types.ResponseCheckTx{
			Code: code.CodeTypeUnknownError,
			Log:  fmt.Sprintf("No data found and TxnCount is %d ", app.txnCount)}

	}
	return types.ResponseCheckTx{
		Code: code.CodeTypeOK,
		Data: dataHash,
		Log:  fmt.Sprintf("Txn is not empty.Txn stored at time: %s", storeresult.Timestamp)}

	//return types.ResponseCheckTx{Code: code.CodeTypeOK}
}

// Commit method
func (app *OperationApplication) Commit() types.ResponseCommit {
	app.hashCout++

	if app.txnCount == 0 {
		return types.ResponseCommit{}
	}
	appHash := make([]byte, 8)
	binary.PutUvarint(appHash, uint64(app.txnCount))
	// new lines added
	//app.state.AppHash = appHash
	//app.state.Height++
	//saveState(app.state)
	return types.ResponseCommit{Data: appHash}
}

// Query method
func (app *OperationApplication) Query(reqQuery types.RequestQuery) (resQuery types.ResponseQuery) {
	//query returns last txn with operation applied i.e add/sub with timestamp
	// query will accept only single input, can't send like this query "a,b" or query "a b"
	fmt.Println("\n value of req Query Data is: ", string(reqQuery.Data))
	queryRequest := string(reqQuery.Data)
	var docsCount int
	var err error
	switch queryRequest {
	case "add":
		docsCount, err = db.C("add").Count()
	case "sub":
		docsCount, err = db.C("sub").Count()
	case "txn":
		docsCount = app.txnCount
	case "hash":
		docsCount = app.hashCout
	}
	if err != nil {
		panic(err)
	}
	appHash := make([]byte, 8)
	binary.PutUvarint(appHash, uint64(docsCount))
	//fmt.Print("\n Collection list is: ", temp)
	fmt.Printf("\n Total number of transaction of operation %s is %d  ", queryRequest, appHash)
	//fmt.Println("\n after formatting: ", []byte(fmt.Sprintf("%02x", docsCount)))
	return types.ResponseQuery{Value: appHash}
}

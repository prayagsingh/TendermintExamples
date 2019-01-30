package OurApp

import (
	"fmt"
	"time"

	//"reflect"
	"encoding/binary"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/tendermint/tendermint/abci/example/code"
	"github.com/tendermint/tendermint/abci/types"
	sm "github.com/tendermint/tendermint/state"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ types.Application = (*OperationApplication)(nil)
var db *mgo.Database

// StoreResult structure to store the data to the mongoDb
type StoreResult struct {
	TxnID     bson.ObjectId `bson:"_id"       json:"_id"`
	User      string        `bson:"username"  json:"username"`
	ParamOne  int           `bson:"paramone"  json:"paramone"`
	ParamTwo  int           `bson:"paramtwo"  json:"paramtwo"`
	Result    string        `bson:"result"    json:"result"`
	Timestamp time.Time     `bson:"timestamp" json:"timestamp"`
}

// MyState using State structure
type MyState struct {
	state   sm.State
	TotTxn  int64  `bson:"tottxn"`
	AppHash []byte `bson:"app_hash"`
	Height  int64  `bson:"height"`
}

// function to save the current state so that it can be loaded from the same state
func saveState(state MyState) {
	stateBytes, err := json.Marshal(state)
	fmt.Print("\n Val of stateBytes is: ", string(stateBytes))
	if err != nil {
		panic(err)
	}
	errMongo := db.C("CurrentState").Insert(state)
	if errMongo != nil {
		panic(errMongo)
	}
}

// function to load the previous state
func loadState(db *mgo.Database) MyState {
	var state MyState
	//err := db.C("CurrentState").Find(bson.M{"_id": ""}).Sort("-timestamp").One(&state)
	dbSize, errCount := db.C("CurrentState").Count()
	fmt.Print("\n Value of dbSize is: ", dbSize)
	if errCount != nil {
		panic(errCount)
	}
	if dbSize > 0 {
		err := db.C("CurrentState").Find(bson.M{}).Sort("-_id").One(&state)
		if err != nil {
			panic(err)
		}
	}
	//var state State
	fmt.Print("\n Inside loadState Func and Value of state is: txnCount: ", state.TotTxn,
		"\n Block height: ", state.Height,
		"\n Hash value: ", fmt.Sprintf("%x", state.AppHash))

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

// OperationApplication which will use Baseapplication which is under application.go
type OperationApplication struct {
	types.BaseApplication
	//txnCount int // stores the count of all the total deliver_tx
	//hashCout int // stores the count of total commit
	state MyState
}

// NewOperationApplication which returns OperationApplication struct
func NewOperationApplication(dbCopy *mgo.Database) *OperationApplication {
	db = dbCopy
	fmt.Print("\n Inside NewOperation function")
	state := loadState(db)
	return &OperationApplication{state: state}
}

// Tendermint methods

// Info method
func (app *OperationApplication) Info(req types.RequestInfo) (resInfo types.ResponseInfo) {
	fmt.Print("\n\n Inside INFO method and Values are: ",
		"\n Version: ", req.Version,
		"\n Block Version: ", req.BlockVersion,
		"\n P2P version: ", req.P2PVersion, "\n")

	resInfo.LastBlockHeight = app.state.Height
	resInfo.LastBlockAppHash = app.state.AppHash
	fmt.Print("\n Value of Last Block Height: ", resInfo.LastBlockHeight,
		"\n Last block App hash: ", fmt.Sprintf("%x", resInfo.LastBlockAppHash))
	//return resInfo
	return types.ResponseInfo{
		Data:             fmt.Sprintf("{\"Total txns so far\":%v}", app.state.TotTxn),
		Version:          req.Version,
		LastBlockHeight:  app.state.Height,
		LastBlockAppHash: app.state.AppHash}
}

// BeginBlock   ---> track the block hash and header Info
func (app *OperationApplication) BeginBlock(params types.RequestBeginBlock) types.ResponseBeginBlock {
	fmt.Print("\n Inside Begin Block function \n")
	return types.ResponseBeginBlock{}

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
	app.state.TotTxn++
	fmt.Print("\n value of state TotTxn is: ", app.state.TotTxn)

	return types.ResponseDeliverTx{
		Code: code.CodeTypeOK,
		Data: dataHash,
		Log:  fmt.Sprintf("Txn executed with timestamp : %s and TxnCount is %d", TxnTime, app.state.TotTxn)}
}

// CheckTx method
func (app *OperationApplication) CheckTx(tx []byte) types.ResponseCheckTx {
	// pattern for sending checktx = check_tx "username,val1,val2,method"
	//storeresult := StoreResult{}

	// convert from ascii to respective values
	txStr := string(tx)
	fmt.Print("\n Inside the Check_tx module nd value of tx is: ", txStr)

	// checking if the data is empty. If yes then return badCode
	fmt.Print("\n value of txnCount is: ", app.state.TotTxn)
	if app.state.TotTxn == 0 {
		return types.ResponseCheckTx{
			Code: code.CodeTypeUnknownError,
			Log:  fmt.Sprintf("No data found and TxnCount is %d ", app.state.TotTxn)}

	}
	return types.ResponseCheckTx{
		Code: code.CodeTypeOK,
		//		Data: dataHash,
		//		Log:  fmt.Sprintf("Txn is not empty.Txn stored at time: %s", storeresult.Timestamp)
		Log: fmt.Sprintf("Txn is not empty.Txn count is: %d", app.state.TotTxn)}

	//return types.ResponseCheckTx{Code: code.CodeTypeOK}
}

// Commit method
func (app *OperationApplication) Commit() types.ResponseCommit {

	fmt.Print("\n\n Inside Commmit method and txnCount is: ", app.state.TotTxn)

	if app.state.TotTxn == 0 {
		return types.ResponseCommit{}
	}
	appHash := make([]byte, 8)
	binary.PutUvarint(appHash, uint64(app.state.TotTxn))

	app.state.AppHash = appHash
	// this is the one of the those hash which is compared to ResponseInfo LastBlockHash when we starting
	// tendermint node. If this hash is not as same as resinfo hash then we won't be able to start
	// the node. will get Replay error which is under checkAppHash function in replay.go
	app.state.state.AppHash = appHash

	fmt.Printf("\n Value of app.state.state.AppHash in HEX is: %X", app.state.state.AppHash)
	app.state.Height++
	app.state.state.LastBlockHeight = app.state.Height
	fmt.Printf("\n Value of app.state.state.LastBlockHeight is: %d", app.state.state.LastBlockHeight)
	saveState(app.state)
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
		docsCount = int(app.state.TotTxn)
		//case "hash":
		//	docsCount = string(app.state.AppHash)
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

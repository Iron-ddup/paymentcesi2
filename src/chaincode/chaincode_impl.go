package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type PaymentChaincode struct {
}
type NOSTROID struct {
	ACTID string
}

const (
	//汇款交易

	TdNoStroBal = "TdNoStroBal"
)

type InitTableData struct {
	TdNoStroBal []*TdNoStroBalgit
}

//NOSTRO余额表
type TdNoStroBalgit struct {
	ACTID     string // 往来账户ID
	BKCODE    string // 银行号
	CLRBKCDE  string //资金清算银行号
	CURCDE    string //货币码
	NOSTROBAL string //Nostro余额
}

//初始化创建表
func (t *PaymentChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	log.Println("Init method.........")
	return nil, CreateTable(stub)
}

//调用方法
func (t *PaymentChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "InitData" { //初始化数据

		return InitData(stub, args)
	}
	return nil, errors.New("调用invoke失败")
}

//查询
func (t *PaymentChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "QueryTdNoStroBalRecordByKey" { //用户根据主键查询账户表记录(NOSTRO表)

		return QueryTdNoStroBalRecordByKey(stub, args)
	}

	return nil, errors.New("调用query失败")
}

func main() {

	err := shim.Start(new(PaymentChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}

}

//创建表
func CreateTable(stub shim.ChaincodeStubInterface) error {

	//创建NOSTRO余额表
	err := stub.CreateTable(TdNoStroBal, []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "ACTID", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "BKCODE", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "CLRBKCDE", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "CURCDE", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "NOSTROBAL", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		stub.DeleteTable(TdNoStroBal)
		return errors.New("create table TdNoStroBal is fail")
	}

	return nil
}

//初始化数据
func InitData(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	data := new(InitTableData)
	err := ParseJsonAndDecode(data, args)
	if err != nil {
		log.Println("Error occurred when parsing json")
		return nil, errors.New("Error occurred when parsing json.")
	}
	for j := 0; j < len(data.TdNoStroBal); j++ {
		InsertTdNoStroBal(stub, *data.TdNoStroBal[j])
	}
	log.Println("InitData success.")
	return nil, nil
}

type SuperObject interface{}

//处理传过来的json数据
func ParseJsonAndDecode(data SuperObject, args []string) error {
	fmt.Println(data)
	//base64解码
	arg, err := base64.StdEncoding.DecodeString(args[0])
	if err != nil {
		log.Println("ParseJson base64 decode error.")
		return err
	}
	log.Println("data after decode:" + string(arg[:]))
	fmt.Println("data after decode:" + string(arg[:]))

	//解析数据
	err = json.Unmarshal(arg, data)
	if err != nil {
		log.Println("ParseJson json Unmarshal error.")
		return err
	}
	fmt.Println("----------------------------------")
	fmt.Println(data)
	fmt.Println("Parse json is ok.")

	return nil
}

//插入数据
//NOSTRO余额表录入(初始化)
func InsertTdNoStroBal(stub shim.ChaincodeStubInterface, data TdNoStroBalgit) ([]byte, error) {
	//往账户表插入数据
	ok, err := stub.InsertRow(TdNoStroBal, shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: data.ACTID}},      // 往来账户ID
			&shim.Column{Value: &shim.Column_String_{String_: data.BKCODE}},     // 银行号
			&shim.Column{Value: &shim.Column_String_{String_: data.CLRBKCDE}},   //资金清算银行号
			&shim.Column{Value: &shim.Column_String_{String_: data.CURCDE}},     // 货币码
			&shim.Column{Value: &shim.Column_String_{String_: data.NOSTROBAL}}}, //Nostro余额
	})

	if !ok && err == nil {
		return nil, errors.New("Table TdNoStroBal insert failed.")
	}

	return nil, nil
}

//查询NOSTRO余额表的记录（只查一条数据）
func QueryTdNoStroBalRecordByKey(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	data := new(NOSTROID) //
	err := ParseJsonAndDecode(data, args)
	if err != nil {
		log.Println("Error occurred when parsing json")
		return nil, errors.New("Error occurred when parsing json.")
	}
	return QueryTdNoStroBalRecordByKeygit(stub, data.ACTID)
}

//根据主键查询nos记录
func QueryTdNoStroBalRecordByKeygit(stub shim.ChaincodeStubInterface, ACTID string) ([]byte, error) {

	var columns []shim.Column
	col := shim.Column{Value: &shim.Column_String_{String_: ACTID}}
	columns = append(columns, col)
	row, _ := stub.GetRow(TdNoStroBal, columns)
	if len(row.Columns) == 0 { //row是否为空
		var errorMsg = "Table TdNoStroBal: specified record doesn't exist,TdNoStroBal id = " + ACTID
		log.Println(errorMsg)
		return nil, errors.New(errorMsg)
	} else {

		jsonResp := `{"ACTID":"` + row.Columns[0].GetString_() + `","BKCODE":"` + row.Columns[1].GetString_() +
			`","CLRBKCDE":"` + row.Columns[2].GetString_() + `","CURCDE":"` + row.Columns[3].GetString_() +
			`","NOSTROBAL":"` + row.Columns[4].GetString_() + `"}`

		log.Println("jsonResp:" + jsonResp)
		return []byte(base64.StdEncoding.EncodeToString([]byte(`{"status":"OK","errMsg":"查询成功","data":` + jsonResp + `}`))), nil
	}
}

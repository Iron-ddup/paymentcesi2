package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type PaymentChaincode struct {
}

const (
	//汇款交易

	TdNoStroBal = "TdNoStroBal"
)

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

		//return InitData(stub, args)
	}
	return nil, errors.New("调用invoke失败")
}

//查询
func (t *PaymentChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "QueryTdNoStroBalRecordByKey" { //用户根据主键查询账户表记录(NOSTRO表)

		//return QueryTdNoStroBalRecordByKey(stub, args)
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

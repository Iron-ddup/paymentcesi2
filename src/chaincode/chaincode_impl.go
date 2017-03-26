package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type PaymentChaincode struct {
}

//初始化创建表
func (t *PaymentChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	defer End(Begin("Init"))

	log.Println("Init method.........")
	return nil, CreateTable(stub)
}

//调用方法
func (t *PaymentChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	defer End(Begin("Invoke"))

	if function == "InitData" { //初始化数据
		log.Println("InitData,args = " + args[0])
		return InitData(stub, args)
	}
	return nil, errors.New("调用invoke失败")
}

//查询
func (t *PaymentChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	defer End(Begin("Query"))

	if function == "QueryTdNoStroBalRecordByKey" { //用户根据主键查询账户表记录(NOSTRO表)
		log.Println("QueryTdNoStroBalRecordByKey,args = " + args[0])
		fmt.Println("QueryTdNoStroBalRecordByKey,args = " + args[0])
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

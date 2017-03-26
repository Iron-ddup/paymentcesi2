package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//初始化两张表NOSTRO余额表和表客户账户余额表
type InitTableData struct {
	TdNoStroBal []*TdNoStroBalgit
}

type NOSTROID struct {
	ACTID string
}

//汇款初始化数据(NOSTRO余额表和客户账户余额表)------------------start-----------------------------
//初始化数据
func InitData(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	defer End(Begin("InitData"))
	defer func() {
		e := recover()
		if e != nil {
			fmt.Println("初始化数据抛出异常:", e)
			log.Println("初始化数据抛出异常:", e)
		}
	}()

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

//*****************************查询功能start************************************

//查询NOSTRO余额表的记录（只查一条数据）
func QueryTdNoStroBalRecordByKey(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	defer End(Begin("QueryTdNoStroBalRecordByKey"))
	data := new(NOSTROID) //
	err := ParseJsonAndDecode(data, args)
	if err != nil {
		log.Println("Error occurred when parsing json")
		return nil, errors.New("Error occurred when parsing json.")
	}
	return QueryTdNoStroBalRecordByKeygit(stub, data.ACTID)
}

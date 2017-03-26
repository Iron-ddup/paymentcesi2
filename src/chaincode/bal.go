package main

import (
	"encoding/base64"
	"errors"

	"log"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//NOSTRO余额表
type TdNoStroBalgit struct {
	ACTID     string // 往来账户ID
	BKCODE    string // 银行号
	CLRBKCDE  string //资金清算银行号
	CURCDE    string //货币码
	NOSTROBAL string //Nostro余额
}

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
	//更新或者插入table_count表
	totalNumber, err := UpdateTableCount(stub, TdNoStroBal)
	if totalNumber == 0 || err != nil {
		stub.DeleteRow(TdNoStroBal, []shim.Column{shim.Column{Value: &shim.Column_String_{String_: data.ACTID}}})
		return nil, errors.New("InsertTdNoStroBal Table_Count insert failed")
	}
	//把获取的总数插入行号表中当主键
	err = UpdateRowNoTable(stub, TdNoStroBal_Rownum, data.ACTID, totalNumber)
	if err != nil {
		stub.DeleteRow(TdNoStroBal, []shim.Column{shim.Column{Value: &shim.Column_String_{String_: data.ACTID}}})
		var columns []shim.Column
		col := shim.Column{Value: &shim.Column_String_{String_: TdNoStroBal}}
		columns = append(columns, col)
		row, _ := stub.GetRow(Table_Count, columns) //row是否为空
		if len(row.Columns) == 1 {
			stub.DeleteRow(Table_Count, []shim.Column{shim.Column{Value: &shim.Column_String_{String_: TdNoStroBal}}})
		} else {
			stub.ReplaceRow(Table_Count, shim.Row{
				Columns: []*shim.Column{
					&shim.Column{Value: &shim.Column_String_{String_: TdNoStroBal}},  //表名
					&shim.Column{Value: &shim.Column_Int64{Int64: totalNumber - 1}}}, //总数
			})
		}
		return nil, errors.New("InsertTdNoStroBal insert TdNoStroBal_Rownum failed")
	}
	log.Println("InsertTdNoStroBal success.")
	return nil, nil
}

//根据主键查询nos记录
func QueryTdNoStroBalRecordByKeygit(stub shim.ChaincodeStubInterface, ACTID string) ([]byte, error) {
	defer End(Begin("QueryTdCstActBalRecordByKey"))
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

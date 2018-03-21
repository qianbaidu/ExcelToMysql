package main

import (
	"github.com/qianbaidu/ExcelToMysql/excel"
	"fmt"
)

func main() {
	filepath := "./testFile/2.xlsx"
	excelSheet := 0
	tableName := ""

	res, _ := excel.ExcelToMysql(filepath, excelSheet, tableName)
	fmt.Print(res)
}



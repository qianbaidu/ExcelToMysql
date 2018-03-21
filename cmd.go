/**
* 
* @author Alex
* @create 2018-03-21 11:37
* @package ExcelToMysql
**/


package main

import (
	"fmt"
	"flag"
	"os"
	"github.com/qianbaidu/ExcelToMysql/excel"
)

func main() {
	filepath := flag.String("F", " ", "Excel File path")
	excelSheet := flag.Int("S", 0, "Excel sheet index")
	tableName := flag.String("T", "", "create mysql table name (if empty auto date_time)")
	flag.Parse()

	if *filepath == "" {
		fmt.Println("文件路径不能为空")
		os.Exit(0)
	}
	res, _ := excel.ExcelToMysql(*filepath, *excelSheet, *tableName)
	fmt.Print(res)

}



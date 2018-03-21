### Excel导入Mysql

- 支持Excel导入到Mysql，解决因Excel列增多减少而需要不断改代码的需求问题。
- 原理:第一行标题作为Mysql字段名，内容自动插入转换Mysql
- 支持:xlsx文件转Mydql
- 不支持带合并单元格Excel文件
### 导入展示
- Excel文件

![](./_image/2018-03-21-13-39-16.jpg)
```
➜  ExcelToMysql git:(master) ✗ go run cmd.go -F ./testFile/2.xlsx
{"FilePath":"./testFile/2.xlsx","SheetIndex":0,"TableName":"2018_03_21_13_3844","SkipFirstLine":0,"LineNum":14,"RowNum":4}
```
- Mysql导入数据

![](./_image/2018-03-21-13-39-40.jpg)

###  使用
- 获取安装包
```
go get github.com/qianbaidu/ExcelToMysql/excel
```
- 导入示例
```
package main

import (
	"github.com/qianbaidu/ExcelToMysql/excel"
	"fmt"
)

func main() {
	filepath := "./testFile/1.xlsx"
	excelSheet := 0
	tableName := ""

	res, _ := excel.ExcelToMysql(filepath, excelSheet, tableName)
	fmt.Print(res)
}



```
- 命令行
```
go run cmd.go -F ./testFile/1.xlsx
```
- 实例代码
```
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
```




package excel

import (
	"github.com/tealeg/xlsx"
	"fmt"
	"strings"
	"unicode/utf8"
	"path"
	"os"
	"time"
	"github.com/qianbaidu/ExcelToMysql/db"
	"encoding/json"
)

type ImportInfo struct {
	FilePath      string
	SheetIndex    int
	TableName     string
	SkipFirstLine int
	LineNum       int
	RowNum        int
}

const sql_type_varchar = "varchar"
const sql_type_txt = "TEXT"
const debug = false

type ExcelData struct {
	Data        [][]string
	SqlFileType []int
}

type readExcelInterface interface {
	ReadExcel()
	checkFile()
	checkSeetNum(xlFile *xlsx.File)
	InsertData(excelData ExcelData)
}

func ( ImportInfo ImportInfo)checkSeetNum(xlFile *xlsx.File) {
	sheetLen := len(xlFile.Sheets)
	if (ImportInfo.SheetIndex > sheetLen || ImportInfo.SheetIndex < 0) {
		fmt.Println("sheet index out of range")
		os.Exit(0)
	}
}

func (ImportInfo *ImportInfo)checkFile() {

	fileSuff := path.Ext(ImportInfo.FilePath)
	if fileSuff != ".xlsx" {
		fmt.Println("不支持的文件类型，当前仅支持`.xlsx`文件格式,请修改文件格式类型")
		os.Exit(0)
	}
}

func exitError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func debugPrint(str string) {
	if debug == true {
		fmt.Println(str)
	}
}

//读取Excel内容
func (ImportInfo *ImportInfo)ReadExcel() (data ExcelData) {
	ImportInfo.checkFile()
	excelFileName := ImportInfo.FilePath
	xlFile, err := xlsx.OpenFile(excelFileName)
	exitError(err)
	ImportInfo.checkSeetNum(xlFile)

	sheet := xlFile.Sheets[ImportInfo.SheetIndex]
	//行数
	lineNum := len(sheet.Rows) + 1
	//列数
	rowNum := len(sheet.Rows[ImportInfo.SheetIndex].Cells) + 1 //这里包含该空列，后面根据内容过滤会会调整，需要判断处理
	ImportInfo.LineNum = lineNum
	ImportInfo.RowNum = rowNum

	sqlFieldsTypeList := make([]int, rowNum)

	excelData := make([][]string, lineNum)
	for rowIndex, row := range sheet.Rows {
		if row == nil {
			continue
		}

		rowValue := make([]string, rowNum)
		isEmptyLine := false
		for cellIndex, cell := range row.Cells {
			str, _ := cell.String()
			str = strings.Trim(str, " ")
			if str != "" {
				isEmptyLine = true
			}
			if (rowIndex == 0) {
				debugPrint(fmt.Sprintf("%d-%d-%s--%d \n", rowIndex, cellIndex, str, rowNum))
			}
			//mysql 字段长度类型判断
			fieldLen := utf8.RuneCountInString(str)
			if (fieldLen > sqlFieldsTypeList[cellIndex]) {
				sqlFieldsTypeList[cellIndex] = fieldLen
			}
			//空列过滤
			if len(str) == 0 && rowIndex == 0 {
				rowNum = cellIndex + 1
				ImportInfo.RowNum = rowNum
			}
			if cellIndex >= (rowNum - 1) {
				break
			} else {
				rowValue[cellIndex] = str
			}
		}
		//遇到一行全是空 认为是内容结束
		if isEmptyLine == false {
			lineNum = rowIndex + 1
			ImportInfo.LineNum = lineNum
			debugPrint(fmt.Sprintf("%d 空行\n", lineNum))
			break
		}
		excelData[rowIndex] = rowValue
	}
	data = ExcelData{excelData, sqlFieldsTypeList}

	return
}

func isInColumn(field *string, ColumnList []string) {
	existsNum := 0;
	for _, v := range ColumnList {
		if *field == v {
			existsNum ++
		}
	}
	if existsNum > 0 {
		*field = fmt.Sprintf("%s_%d", *field, existsNum + 1)
	}

}
func (importInfo *ImportInfo)createTable(excelData *ExcelData) {
	fieldSql := ""
	ColumnList := make([]string, importInfo.RowNum)
	for k, v := range excelData.Data[0] {

		isInColumn(&v, ColumnList)
		excelData.Data[0][k] = v
		ColumnList = append(ColumnList, v)
		tmpSql := ""
		if excelData.SqlFileType[k] < 255 {
			tmpSql = fmt.Sprintf(" `%s` varchar(%d) NOT NULL DEFAULT '' COMMENT '%s',\n\t", v, excelData.SqlFileType[k] + 10, v)
		} else {
			tmpSql = fmt.Sprintf(" `%s` text COMMENT '%s',\n\t", v, v)
		}
		if k >= (importInfo.RowNum - 1) || v == "" {
			break
		}
		fieldSql = fmt.Sprintf("%s %s ", fieldSql, tmpSql)

	}
	var sql = `
	CREATE TABLE %s (
	  _auto_id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
	  %s
	  PRIMARY KEY (_auto_id)
	) ENGINE=InnoDB  DEFAULT CHARSET=utf8;`

	if importInfo.TableName == "" {
		t := time.Now()
		importInfo.TableName = fmt.Sprintf("%d_%02d_%02d_%02d_%02d%02d",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second())
	}
	sql = fmt.Sprintf(sql, importInfo.TableName, fieldSql)
	//debugPrint(fmt.Sprintf("create table sql \n%s", sql))
	db.DbExecSql(sql)

}

func filterReplace(value string) (str string) {
	value = strings.Replace(value, ",", "", -1)
	value = strings.Replace(value, "'", "\\'", -1)
	value = strings.Replace(value, "\\", " ", -1)
	str = value
	return
}

func DataToMysql(importInfo ImportInfo, sqlFields string, values string) {
	insertSql := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s; ", importInfo.TableName, strings.Trim(sqlFields, ","), strings.Trim(values, ","))
	db.DbExecSql(insertSql)
}

func (importInfo ImportInfo)InsertData(excelData ExcelData) {
	sqlFields := ""
	sqlValues := ""
	for lineIndex, lineValue := range excelData.Data {
		values := ""
		for rowIndex, rowValue := range lineValue {
			if rowIndex >= (importInfo.RowNum - 1 ) {
				break
			}
			if lineIndex == 0 {
				sqlFields += fmt.Sprintf("`%s`,", rowValue)
			}
			rowValue = filterReplace(rowValue)
			values += fmt.Sprintf("'%s',", rowValue)
		}
		if importInfo.SkipFirstLine == 0 && lineIndex == 0 {
			continue
		}
		if values != "" {
			sqlValues += fmt.Sprintf("(%s),", strings.Trim(values, ","))
		}

		if ((lineIndex + 1) % 10 ) == 0 {
			DataToMysql(importInfo, sqlFields, sqlValues)
			sqlValues = ""
		}
	}
	DataToMysql(importInfo, sqlFields, sqlValues)
}

func ExcelToMysql(filePath string, sheetIndex int, tableName string) (res string, err error) {
	importInfo := ImportInfo{FilePath:filePath, SheetIndex:sheetIndex, TableName:tableName}
	excelData := importInfo.ReadExcel()
	importInfo.createTable(&excelData)
	importInfo.InsertData(excelData)

	result, err := json.Marshal(importInfo)
	res = string(result)
	return
}
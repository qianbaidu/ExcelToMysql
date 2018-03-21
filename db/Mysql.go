package db

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/qianbaidu/ExcelToMysql/conf"
)

type MysqlDb struct {
	Conn *sql.DB
}

type Mysql interface {
	DbExec(execSql string)
}

func Connect() (conn *sql.DB) {
	config := conf.InitConfig()
	connectUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		config.Mysql.Username,
		config.Mysql.Password,
		config.Mysql.Host,
		config.Mysql.Port,
		config.Mysql.Database)
	conn, err := sql.Open("mysql", connectUrl)

	if err != nil {
		panic(err.Error())
	}

	return conn
}

func Query() {

}

func (db *MysqlDb)DbExec(execSql string) {
	stmt, err := db.Conn.Prepare(execSql)
	if err != nil {
		fmt.Println(execSql)
		panic(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(execSql)
		panic(err)
	}
}

func DbExecSql(execSql string) {
	db := Connect()
	defer db.Close()
	mysql := MysqlDb{db}
	mysql.DbExec(execSql)
}
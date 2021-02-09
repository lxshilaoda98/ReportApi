package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"

)
import _ "github.com/denisenkom/go-mssqldb"

var SqlDB *sql.DB

var DBDriver ="mysql" //<mysql/mssql>

func init() {
	connString := ""
	if DBDriver=="mssql"{
		connString="server=127.0.0.1;port=1433;database=TF_LFT_CMS;user id=sa;password=sa@01;encrypt=disable"
	}else if DBDriver =="mysql"{
		connString="root:root@tcp(127.0.0.1:3307)/freeswitch?parseTime=true"
	}
	var err error
	SqlDB, err = sql.Open(DBDriver, connString)
	if err != nil {
		fmt.Println("数据库连接异常:err..>", err.Error())
	}
	err = SqlDB.Ping()
	if err != nil {
		fmt.Println("[ping]数据库连接异常:err..>", err.Error())
	}
}

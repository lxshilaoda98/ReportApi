package main

import (
	"database/sql"
	"fmt"
	//_ "github.com/go-sql-driver/mysql"

	"sync"
)
import _ "github.com/denisenkom/go-mssqldb"

var SqlDB *sql.DB
var wg sync.WaitGroup
func main() {

	defer SqlDB.Close()
	for i:=0;i<100;i++ {
		wg.Add(1)
		//开启10个进程添加测试数据。。。
		go InsterGateway(500000)
	}
	wg.Wait() //阻塞直到所有任务完成
	fmt.Println("添加完成咯")
}
func InsterGateway(Number int){

	for i:=0;i<=Number;i++  {
		rs,err:=SqlDB.Exec("INSERT INTO inboundservicerequest(Ani) VALUES (?)","6001")
		if err != nil {
			fmt.Println("插入库异常..>",err)
		}else {
			fmt.Println("正常循环插入数据..>",rs,".插入成功!")
		}
	}
	wg.Done()

}
func init() {

	//connString :="root:root@tcp(127.0.0.1:3307)/icrm_adms_core?parseTime=true"
	connString:="server=192.168.101.171;port=1433;database=cs;user id=sa;password=1qaz@WSX;encrypt=disable"
	var err error
	SqlDB, err = sql.Open("mssql", connString)
	if err != nil {
		fmt.Println("数据库连接异常:err..>", err.Error())
	}
	err = SqlDB.Ping()
	if err != nil {
		fmt.Println("[ping]数据库连接异常:err..>", err.Error())
	}
}

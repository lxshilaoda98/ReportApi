package main

import (
	db "github.com/n1n1n1_owner/ReportApi/database"
	//EslHelper "github.com/n1n1n1_owner/ReportApi/models/Helper"
)

func main() {
	defer db.SqlDB.Close()
	router := InitRouter()
	//开始esl 事件
	//go EslHelper.ConnectionEsl()
	//for i:=0;i<10;i++ {
	//	//开启10个进程添加测试数据。。。
	//	go EslHelper.InsterGateway(500000)
	//}

	router.Run(":8000")
}

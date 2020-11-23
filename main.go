package main

import (

	db "github.com/n1n1n1_owner/ReportApi/database"
  	EslHelper "github.com/n1n1n1_owner/ReportApi/models/Helper"
)


func main() {
	defer db.SqlDB.Close()
	router := InitRouter()
	go EslHelper.ConnectionEsl()
	router.Run(":8000")
}

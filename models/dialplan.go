package models

import (
	"fmt"
	db "github.com/n1n1n1_owner/ReportApi/database"
)

type Dialplan struct {
	Oid         int    `json:"Oid" form:"Oid"`
	Description string `json:"Description" form:"Description"`
	Context     string `json:"Context" form:"Context"`
}

/**
获取拨号计划数据
start :当前页
end：每页显示
*/
func (d *Dialplan) GetDialplan(start, end int) (dialplans []Dialplan, err error) {
	//如果是mssql版本的话 分页可能需要改变
	//Demo SQL 语句
	//select top 10 * from (select row_number() over(order by Oid asc) as rownumber,* from dialplan)
	// temp_row where rownumber>10;
	// top10 每页10条数据 rownumber>0 第一页数据 rownumber>10 第二页数据 20 -3 ....

	dialplans = make([]Dialplan, 0)
	fmt.Printf("开始%d,结束%d \n",start,end)
	rows, err := db.SqlDB.Query("select Oid, Description, Context from dialplan limit ?,?", start, end)
	if err != nil {
		fmt.Print("select SQL Err..>", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var dialplan Dialplan
		rows.Scan(&dialplan.Oid, &dialplan.Description, &dialplan.Context)
		dialplans = append(dialplans, dialplan)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return

	//c.JSON(http.StatusOK, gin.H{
	//	"result": dialplans,
	//	"count":  len(dialplans),
	//})
}

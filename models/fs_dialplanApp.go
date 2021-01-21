package models

import (
	"fmt"
	db "github.com/n1n1n1_owner/ReportApi/database"
)

type DialplanApp struct {
	Oid         string `json:"oid"`
	Application string `json:"application"`
	Data        string `json:"data"`
	Dialplan    string `json:"dialplan"`
	Sort        string `json:"sort"`
}

func (d *DialplanApp) GetDialplanApp(dialplanOid, start, end int) (dialplanApps []DialplanApp, err error) {
	//如果是mssql版本的话 分页可能需要改变
	//Demo SQL 语句
	//select top 10 * from (select row_number() over(order by Oid asc) as rownumber,* from dialplan)
	// temp_row where rownumber>10;
	// top10 每页10条数据 rownumber>0 第一页数据 rownumber>10 第二页数据 20 -3 ....

	dialplanApps = make([]DialplanApp, 0)
	fmt.Printf("开始%d,结束%d \n", start, end)
	rows, err := db.SqlDB.Query("select Oid, Application, Data,Dialplan,Sort from diaplan_app where dialplan=? limit ?,?", dialplanOid, start, end)
	if err != nil {
		fmt.Print("[GetDialplanApp]select SQL Err..>", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var dialplanApp DialplanApp
		rows.Scan(&dialplanApp.Oid, &dialplanApp.Application, &dialplanApp.Data, &dialplanApp.Dialplan, &dialplanApp.Sort)
		dialplanApps = append(dialplanApps, dialplanApp)
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
func (d *DialplanApp) GetDialplanAppCount(dialplan int) (count Count) {
	fmt.Println("查询dialplan的ID。。》", dialplan)
	err := db.SqlDB.QueryRow("select count(*) as Number from diaplan_app where dialplan=? ", dialplan).Scan(
		&count.Number,
	)
	if err != nil {
		fmt.Println("select Count Err..>", err)
	}
	return
}
func (d *DialplanApp) EditDialplanApp(dialplanapp DialplanApp) (id int64, err error) {
	rs, err := db.SqlDB.Exec("UPDATE diaplan_app SET application=?,data=?,sort=? where Oid=?",
		dialplanapp.Application, dialplanapp.Data, dialplanapp.Sort, dialplanapp.Oid)
	if err != nil {
		return
	}
	id, err = rs.LastInsertId()
	if err != nil {
		return
	} else {
		fmt.Println("id==============", id)
	}
	return
}
func (d *DialplanApp) AddDialplanApp(dialplanapp DialplanApp) (id int64, err error) {
	fmt.Println("新增拨号子模块..》", dialplanapp.Dialplan)
	rs, err := db.SqlDB.Exec("INSERT INTO diaplan_app (application,data,sort,dialplan) VALUES (?,?,?,?)",
		dialplanapp.Application, dialplanapp.Data, dialplanapp.Sort, dialplanapp.Dialplan)
	if err != nil {
		return
	}
	id, err = rs.LastInsertId()
	if err != nil {
		return
	} else {
		fmt.Println("id==============", id)
	}
	return
}
func (d *DialplanApp) DelDialplanApp() (id int64, err error) {
	rs, err := db.SqlDB.Exec("DELETE FROM diaplan_app WHERE Oid=?", d.Oid)
	if err != nil {
		fmt.Println("[diaplan_app].err ", err)
	}
	id, err = rs.RowsAffected()
	return
}

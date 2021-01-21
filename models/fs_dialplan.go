package models

import (
	"fmt"
	db "github.com/n1n1n1_owner/ReportApi/database"
)

type Dialplan struct {
	Oid         string `json:"oid"`
	Description string `json:"description"`
	Context     string `json:"context"`
	Name        string `json:"name"`
	Condition   string `json:"condition"`
	Expression  string `json:"expression"`
	Domain      string `json:"domain"`
}
type Count struct {
	Number int `json:"number" form:"number"`
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
	fmt.Printf("开始%d,结束%d \n", start, end)
	rows, err := db.SqlDB.Query("select Oid, Description, Context,Name,`Condition`,Expression,Domain from dialplan limit ?,?", start, end)
	if err != nil {
		fmt.Print("[GetDialplan]select SQL Err..>", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var dialplan Dialplan
		rows.Scan(&dialplan.Oid, &dialplan.Description, &dialplan.Context, &dialplan.Name, &dialplan.Condition, &dialplan.Expression, &dialplan.Domain)
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
func (d *Dialplan) GetDialplanCount() (count Count) {
	_ = db.SqlDB.QueryRow("select count(*) as Number from dialplan").Scan(
		&count.Number,
	)
	return
}
func (d *Dialplan) GetDialplanByOid(oid int) (dialplans Dialplan, err error) {

	//dialplans = make([]Dialplan, 0)

	rows, err := db.SqlDB.Query("select Oid, Description, Context,Name,`Condition`,Expression,Domain from dialplan where oid =?", oid)
	if err != nil {
		fmt.Print("[GetDialplanByOid]select SQL Err..>", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var dialplan Dialplan
		rows.Scan(&dialplan.Oid, &dialplan.Description, &dialplan.Context, &dialplan.Name, &dialplan.Condition, &dialplan.Expression, &dialplan.Domain)
		dialplans = dialplan
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
func (d *Dialplan) AddDialplan(dialplan Dialplan) (id int64, err error) {
	rs, err := db.SqlDB.Exec("INSERT INTO dialplan (Description,Context,`Name`,`Condition`,Expression,domain) VALUES(?,?,?,?,?,?)",
		dialplan.Description, dialplan.Context, dialplan.Name, dialplan.Condition, dialplan.Expression, dialplan.Domain)
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
func (d *Dialplan) EditDialplan(dialplan Dialplan) (id int64, err error) {
	rs, err := db.SqlDB.Exec("UPDATE dialplan SET Description=?,Context=?,`Name`=?,`Condition`=?,Expression=?,domain=? where Oid=?",
		dialplan.Description, dialplan.Context, dialplan.Name, dialplan.Condition, dialplan.Expression, dialplan.Domain, dialplan.Oid)
	if err != nil {
		return
	}
	id, err = rs.LastInsertId()
	if err != nil {
		return
	} else {
		fmt.Println("id==============", dialplan.Oid)
	}
	return
}
func (d *Dialplan) DelDialplan() (id int64, err error) {
	rs, err := db.SqlDB.Exec("DELETE FROM dialplan WHERE Oid=?", d.Oid)
	if err != nil {
		fmt.Println("[DelDialplan].err ", err)
	}
	id, err = rs.RowsAffected()
	return
}

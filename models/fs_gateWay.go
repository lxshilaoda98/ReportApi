package models

import (
	"fmt"
	db "github.com/n1n1n1_owner/ReportApi/database"
)

type GateWay struct {
	Oid      string `json:"oid"`
	Name     string `json:"name"`
	Realm    string `json:"realm"`
	Register string `json:"register"`
	Fromuser string `json:"fromuser"`
	Username string `json:"username"`
	Password string `json:"password"`
	Memo     string `json:"memo"`
}

func (g *GateWay) GetGateWay(start, end int) (gateWays []GateWay, err error) {
	gateWays = make([]GateWay, 0)
	fmt.Printf("开始%d,结束%d \n", start, end)
	rows, err := db.SqlDB.Query("select Oid,Name,Realm,Register,Fromuser,Username,Memo,`password` as Password from gateway limit ?,?", start, end)
	//fmt.Println("sql..>"+"select Oid,SIPUser,Password,Vm_password,Toll_allow,Accountcode,User_context,Effective_caller_id_name, Effective_caller_id_number,Outbound_caller_id_name,Outbound_caller_id_number,Callgroup,Callgroupid,GroupName from sipuser limit ?,?", start, end)
	if err != nil {
		fmt.Print("[GetSipUser]select SQL Err..>", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		//fmt.Println("rows....>",rows)
		var gateWay GateWay
		rows.Scan(&gateWay.Oid, &gateWay.Name, &gateWay.Realm, &gateWay.Register, &gateWay.Fromuser, &gateWay.Username, &gateWay.Memo, &gateWay.Password)
		gateWays = append(gateWays, gateWay)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}
func (g *GateWay) GetGateWayCount() (count Count) {
	_ = db.SqlDB.QueryRow("select count(*) as Number from gateway").Scan(
		&count.Number,
	)
	return
}
func (g *GateWay) AddGateWay(gateWay GateWay) (id int64, err error) {

	Us := gateWay.Username

	rs, err := db.SqlDB.Exec("INSERT INTO gateway (name,realm,register,fromuser,username,password,memo) VALUES (?,?,?,?,?,?,?)",
		gateWay.Name, gateWay.Realm, gateWay.Register, Us, Us, gateWay.Password, gateWay.Memo)

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
func (g *GateWay) DelGateWay() (id int64, err error) {
	rs, err := db.SqlDB.Exec("DELETE FROM gateway WHERE Oid=?", g.Oid)
	if err != nil {
		fmt.Println("[DelGateWay].err ", err)
	}
	id, err = rs.RowsAffected()
	return
}
func (g *GateWay) EditGateWay(gateWay GateWay) (id int64, err error) {

	Us := gateWay.Username
	rs, err := db.SqlDB.Exec("UPDATE gateway SET name=?,realm=?,register=?,fromuser=?,username=?,password=?,memo=? where Oid=?",
		gateWay.Name, gateWay.Realm, gateWay.Register, Us, Us, gateWay.Password, gateWay.Memo, gateWay.Oid)
	if err != nil {
		return
	}
	id, err = rs.LastInsertId()
	if err != nil {
		return
	} else {
		fmt.Println("id==============", gateWay.Oid)
	}
	return
}

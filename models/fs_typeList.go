package models

import (
	"fmt"
	db "github.com/n1n1n1_owner/ReportApi/database"
)

type TypeList struct {
	Oid   string `json:"oid" form:"oid"`
	Title string `json:"title" form:"title"`
	Name  string `json:"name" form:"name"`
	Val   string `json:"val" form:"val"`
	TypeL string `json:"typeL" form:"typeL"`
	State string `json:"state" form:"state"`
	Sort  string `json:"store" form:"store"`
	Label string `json:"label" form:"label"`
	Key   string `json:"key" form:"key"`
	Value string `json:"value" form:"value"`
}

func (t *TypeList) GetTypeLists(start, end int) (typeLists []TypeList, err error) {
	typeLists = make([]TypeList, 0)
	fmt.Printf("开始%d,结束%d \n", start, end)
	rows, err := db.SqlDB.Query("select Oid,Title,Name,Val,TypeL,State,Sort from typelist limit ?,?", start, end)
	if err != nil {
		fmt.Print("[GetTypeLists]select SQL Err..>", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var typeList TypeList
		rows.Scan(&typeList.Oid, &typeList.Title, &typeList.Name, &typeList.Val, &typeList.TypeL, &typeList.State, &typeList.Sort)
		typeLists = append(typeLists, typeList)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}
func (t *TypeList) GetTypeListCount() (count Count) {
	_ = db.SqlDB.QueryRow("select count(*) as Number from typelist").Scan(
		&count.Number,
	)
	return
}
func (t *TypeList) AddTypeList(Typelist TypeList) (id int64, err error) {

	fmt.Println("============", Typelist.Name)
	var StateVal = ""
	if Typelist.State == "true" {
		StateVal = "1"
	} else if Typelist.State == "false" {
		StateVal = "0"
	} else {
		StateVal = Typelist.State
	}
	rs, err := db.SqlDB.Exec("INSERT INTO typelist(title,Name,Val,typeL,state,sort) VALUES (?,?,?,?,?,?)",
		Typelist.Title, Typelist.Name, Typelist.Val, Typelist.TypeL, StateVal, Typelist.Sort)
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
func (t *TypeList) EditTypeList(Typelist TypeList) (id int64, err error) {

	fmt.Println("============", Typelist.Name)
	var StateVal = ""
	if Typelist.State == "true" {
		StateVal = "1"
	} else if Typelist.State == "false" {
		StateVal = "0"
	} else {
		StateVal = Typelist.State
	}
	rs, err := db.SqlDB.Exec("UPDATE typelist SET title=?,Name=?,Val=?,typeL=?,state=?,sort=? where Oid=?",
		Typelist.Title, Typelist.Name, Typelist.Val, Typelist.TypeL, StateVal, Typelist.Sort, Typelist.Oid)

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
func (t *TypeList) DelTypeList() (id int64, err error) {
	rs, err := db.SqlDB.Exec("DELETE FROM typelist WHERE Oid=?", t.Oid)
	if err != nil {
		fmt.Println("[DelTypeList].err ", err)
	}
	id, err = rs.RowsAffected()
	return
}
func (t *TypeList) GetTypeListByType(types string) (typeLists []TypeList, err error) {
	typeLists = make([]TypeList, 0)

	rows, err := db.SqlDB.Query("select Oid,Name,Val from typelist where TypeL=? order by sort asc", types)
	if err != nil {
		fmt.Print("[GetTypeListByType]select SQL Err..>", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var typeList TypeList
		rows.Scan(&typeList.Key, &typeList.Label, &typeList.Value)
		typeLists = append(typeLists, typeList)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}
func (t *TypeList) GetTypeListByVal(Val string) (typeList TypeList, err error) {
	rows, err := db.SqlDB.Query("select Name from typelist where val=? ", Val)
	if err != nil {
		fmt.Print("[GetTypeListByVal]select SQL Err..>", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&typeList.Name)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

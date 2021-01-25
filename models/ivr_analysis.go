package models

import (
	"encoding/json"
	"fmt"
	db "github.com/n1n1n1_owner/ReportApi/database"
)

type NodeLists struct {
	Id    string
	Name  string
	Type  string
	Left  string
	Top   string
	Ico   string
	State string
}

type LineList struct {
	From  string
	To    string
	Label string
}

type IvrModel struct {
	Name     string
	JsonStr  string
	NodeList []NodeLists
	LineList []LineList
}

func (i *IvrModel) ResJson(Oid, jsonStr string) (err error) {
	b := []byte(jsonStr)
	if err := json.Unmarshal(b, &i); err != nil {
		fmt.Println("err json..>", err)
	} else {
		/**
		直接删除foid 的数据,然后重新新增..
		*/
		if r, err := db.SqlDB.Exec("delete from ivr_com where foid = ?", Oid); err != nil {
			fmt.Println("delete ivr_com err..", err)
		} else {
			if count, err := r.RowsAffected(); err != nil {
				fmt.Println("aff err..>", err)
			} else {
				fmt.Println("delete count ", count)
				//继续新增...
				rows, err := db.SqlDB.Exec("insert into ivr_com(name,foid,JsonStr)values(?,?,?)", i.Name, Oid, jsonStr)
				if err != nil {
					fmt.Println("insert Err..>", err)
				} else {
					if id, err := rows.LastInsertId(); err != nil {
						fmt.Println("insert Err.last.>", err)
					} else {
						fmt.Printf("insert into id ..>%d \n", id)
					}
				}
			}
		}
	}
	return

	//fmt.Println("i==>",i.Name)
	//for k,v := range i.NodeList  {
	//	fmt.Println(k,v)
	//}
	//for k,v :=range i.LineList  {
	//	fmt.Println(k,v)
	//}

}

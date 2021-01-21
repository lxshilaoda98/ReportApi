package models

import "fmt"
import db "github.com/n1n1n1_owner/ReportApi/database"

type Tiers struct {
	Oid      string `json:"oid"`
	Queue    string `json:"queue"`
	Agent    string `json:"agent"`
	State    string `json:"state"`
	Level    string `json:"level"`
	Position string `json:"position"`
}

func (t *Tiers) GetTiers(start, end int) (tierss []Tiers, err error) {
	tierss = make([]Tiers, 0)
	fmt.Printf("开始%d,结束%d \n", start, end)
	rows, err := db.SqlDB.Query("select Oid,Queue,Agent,State,Level,Position from tiers limit ?,?", start, end)
	if err != nil {
		fmt.Print("[GetTiers]select SQL Err..>", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		//fmt.Println("rows....>",rows)
		var tiers Tiers
		rows.Scan(&tiers.Oid, &tiers.Queue, &tiers.Agent, &tiers.State, &tiers.Level, &tiers.Position)
		tierss = append(tierss, tiers)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func (t *Tiers) GetTiersCount() (count Count) {
	_ = db.SqlDB.QueryRow("select count(*) as Number from tiers").Scan(
		&count.Number,
	)
	return
}

func (t *Tiers) AddTiers(tiers Tiers) (id int64, err error) {
	rs, err := db.SqlDB.Exec("insert into tiers (queue,agent,state,level,position)values(?,?,?,?,?)",
		tiers.Queue, tiers.Agent, "Ready", tiers.Level, tiers.Position)
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

func (t *Tiers) EditTiers(tiers Tiers) (id int64, err error) {
	rs, err := db.SqlDB.Exec("UPDATE tiers SET queue=?,agent=?,level=?,position=? where Oid = ?",
		tiers.Queue, tiers.Agent, tiers.Level, tiers.Position, tiers.Oid)
	if err != nil {
		return
	}
	id, err = rs.LastInsertId()
	if err != nil {
		return
	} else {
		fmt.Println("id==============", tiers.Oid)
	}
	return
}

func (t *Tiers) DelTiers() (id int64, err error) {
	rs, err := db.SqlDB.Exec("DELETE FROM tiers WHERE Oid=?", t.Oid)
	if err != nil {
		fmt.Println("[DelTiers].err ", err)
	}
	id, err = rs.RowsAffected()
	return
}

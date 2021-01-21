package models

import "fmt"
import db "github.com/n1n1n1_owner/ReportApi/database"

type Member struct {
	Queue           string `json:"queue"`
	Uuid            string `json:"uuid"`
	Session_uuid    string `json:"session_uuid"`
	Cid_number      string `json:"cid_number"`
	Cid_name        string `json:"cid_name"`
	System_epoch    string `json:"system_epoch"`
	Joined_epoch    string `json:"joined_epoch"`
	Rejoined_epoch  string `json:"rejoined_epoch"`
	Bridge_epoch    string `json:"bridge_epoch"`
	Abandoned_epoch string `json:"abandoned_epoch"`
	Base_score      string `json:"base_score"`
	Skill_score     string `json:"skill_score"`
	Serving_agent   string `json:"serving_agent"`
	Serving_system  string `json:"serving_system"`
	State           string `json:"state"`
}

func (m *Member) GetMembers(start, end int) (Members []Member, err error) {
	Members = make([]Member, 0)

	rows, err := db.SqlDB.Query("select Queue,Uuid,Session_uuid,Cid_number,Cid_name,System_epoch"+
		",Joined_epoch,Rejoined_epoch,Bridge_epoch,Abandoned_epoch,Base_score,Abandoned_epoch,Base_score,"+
		"Skill_score,Serving_agent,Serving_system,State from members limit ?,?", start, end)
	if err != nil {
		fmt.Print("[GetTypeLists]select SQL Err..>", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var member Member
		rows.Scan(&member.Queue, &member.Uuid, &member.Session_uuid, &member.Cid_number, &member.Cid_name,
			&member.System_epoch, &member.Joined_epoch, &member.Rejoined_epoch, &member.Bridge_epoch,
			&member.Abandoned_epoch, &member.Base_score, &member.Abandoned_epoch, &member.Base_score,
			&member.Skill_score, &member.Serving_agent, &member.Serving_system, &member.State)
		Members = append(Members, member)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}
func (m *Member) GetMemberCount() (count Count) {
	_ = db.SqlDB.QueryRow("select count(*) as Number from members").Scan(
		&count.Number,
	)
	return
}
func (m *Member) DelMember() (id int64, err error) {
	rs, err := db.SqlDB.Exec("DELETE FROM members WHERE uuid=?", m.Uuid)
	if err != nil {
		fmt.Println("[DelMember].err ", err)
	}
	id, err = rs.RowsAffected()
	return
}

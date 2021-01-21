package models

import (
	"fmt"
	db "github.com/n1n1n1_owner/ReportApi/database"
)

type SipUser struct {
	Oid                        string `json:"oid"`
	SIPUser                    string `json:"sipuser"`
	Password                   string `json:"password"`
	Vm_password                string `json:"vm_password"`
	Toll_allow                 string `json:"toll_allow"`
	Accountcode                string `json:"accountcode"`
	User_context               string `json:"user_context"`
	Effective_caller_id_name   string `json:"effective_caller_id_name"`
	Effective_caller_id_number string `json:"effective_caller_id_number"`
	Outbound_caller_id_name    string `json:"outbound_caller_id_name"`
	Outbound_caller_id_number  string `json:"outbound_caller_id_number"`
	Callgroup                  string `json:"callgroup"`
	Callgroupid                string `json:"callgroupid"`
	GroupName                  string `json:"groupName"`
}

func (s *SipUser) GetSipUser(start, end int) (sipUsers []SipUser, err error) {
	sipUsers = make([]SipUser, 0)
	fmt.Printf("开始%d,结束%d \n", start, end)
	rows, err := db.SqlDB.Query("select Oid,SIPUser,Password,Vm_password,Toll_allow,Accountcode,User_context,Effective_caller_id_name, Effective_caller_id_number,Outbound_caller_id_name,Outbound_caller_id_number,Callgroup,Callgroupid,GroupName from sipuser limit ?,?", start, end)
	//fmt.Println("sql..>"+"select Oid,SIPUser,Password,Vm_password,Toll_allow,Accountcode,User_context,Effective_caller_id_name, Effective_caller_id_number,Outbound_caller_id_name,Outbound_caller_id_number,Callgroup,Callgroupid,GroupName from sipuser limit ?,?", start, end)
	if err != nil {
		fmt.Print("[GetSipUser]select SQL Err..>", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		//fmt.Println("rows....>",rows)
		var sipUser SipUser
		rows.Scan(&sipUser.Oid, &sipUser.SIPUser, &sipUser.Password, &sipUser.Vm_password, &sipUser.Toll_allow, &sipUser.Accountcode,
			&sipUser.User_context, &sipUser.Effective_caller_id_name, &sipUser.Effective_caller_id_number, &sipUser.Outbound_caller_id_name, &sipUser.Outbound_caller_id_number,
			&sipUser.Callgroup, &sipUser.Callgroupid, &sipUser.GroupName)
		sipUsers = append(sipUsers, sipUser)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func (s *SipUser) GetSipUserCount() (count Count) {
	_ = db.SqlDB.QueryRow("select count(*) as Number from sipuser").Scan(
		&count.Number,
	)
	return
}

func (s *SipUser) AddSipUser(sipuser SipUser) (id int64, err error) {

	Pw := sipuser.Password
	Us := sipuser.SIPUser
	fmt.Println("callgroup,..>", sipuser.Callgroup)
	fmt.Println("u.>", Us)
	fmt.Println("p.>", Pw)
	rs, err := db.SqlDB.Exec("INSERT INTO sipuser(SIPUser,Password,Vm_password,Toll_allow,Accountcode,User_context,Effective_caller_id_name, Effective_caller_id_number,"+
		"Outbound_caller_id_name,Outbound_caller_id_number,Callgroup,callgroupid,GroupName) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)",
		Us, Pw, Pw, "domestic,international,local", Us, "default", Us, Us, Us, Us, sipuser.Callgroup, 1, "测试")

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

func (s *SipUser) EditSipUser(sipuser SipUser) (id int64, err error) {
	Pw := sipuser.Password
	Us := sipuser.SIPUser

	rs, err := db.SqlDB.Exec("UPDATE sipuser SET SIPUser=?,password=?,vm_password=?,toll_allow=?,accountcode=?, "+
		"user_context=?,effective_caller_id_name=?,effective_caller_id_number=?,outbound_caller_id_name=?,outbound_caller_id_number=?,"+
		"callgroup=?,callgroupid=?,groupName=? where Oid=?",
		Us, Pw, Pw, "domestic,international,local", Us, "default", Us, Us, Us, Us, sipuser.Callgroup, 1, "测试", sipuser.Oid)

	if err != nil {
		return
	}
	id, err = rs.LastInsertId()
	if err != nil {
		return
	} else {
		fmt.Println("id==============", sipuser.Oid)
	}
	return
}

func (s *SipUser) DelSipUser() (id int64, err error) {
	rs, err := db.SqlDB.Exec("DELETE FROM sipuser WHERE Oid=?", s.Oid)
	if err != nil {
		fmt.Println("[DelSipUser].err ", err)
	}
	id, err = rs.RowsAffected()
	return
}

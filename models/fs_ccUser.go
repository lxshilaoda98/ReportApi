package models

import "fmt"
import db "github.com/n1n1n1_owner/ReportApi/database"

type CCUser struct {
	Name                 string `json:"name"`
	Contact              string `json:"contact"`
	Status               string `json:"status"`
	State                string `json:"state"`
	Max_no_answer        string `json:"max_no_answer"`
	Wrap_up_time         string `json:"wrap_up_time"`
	Reject_delay_time    string `json:"reject_delay_time"`
	Busy_delay_time      string `json:"busy_delay_time"`
	No_answer_delay_time string `json:"no_answer_delay_time"`
	Last_bridge_end      string `json:"last_bridge_end"`
	Last_status_change   string `json:"last_status_change"`
	Key                  string `json:"key"`
	Label                string `json:"label"`
	Value                string `json:"value"`
}

func (c *CCUser) GetCCUserByName() (CCUsers []CCUser, err error) {
	CCUsers = make([]CCUser, 0)

	rows, err := db.SqlDB.Query("select Name,Name,Name from agents ")
	if err != nil {
		fmt.Print("[GetCCUserByName]select SQL Err..>", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var CCUser CCUser
		rows.Scan(&CCUser.Key, &CCUser.Label, &CCUser.Value)
		CCUsers = append(CCUsers, CCUser)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}
func (c *CCUser) GetAllCCUser(start, end int) (ccUsers []CCUser, err error) {

	ccUsers = make([]CCUser, 0)
	fmt.Printf("开始%d,结束%d \n", start, end)
	rows, err := db.SqlDB.Query("SELECT Name,`Contact`,`Status`,State,Max_no_answer,Wrap_up_time,Reject_delay_time,Busy_delay_time,No_answer_delay_time,Last_bridge_end,Last_status_change  from agents limit ?,?", start, end)
	if err != nil {
		fmt.Print("[GetAllCCUser]select SQL Err..>", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var ccUser CCUser
		rows.Scan(&ccUser.Name, &ccUser.Contact, &ccUser.Status, &ccUser.State, &ccUser.Max_no_answer, &ccUser.Wrap_up_time, &ccUser.Reject_delay_time, &ccUser.Busy_delay_time, &ccUser.No_answer_delay_time, &ccUser.Last_bridge_end, &ccUser.Last_status_change)
		ccUsers = append(ccUsers, ccUser)
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
func (c *CCUser) GetCCUserCount() (count Count) {
	_ = db.SqlDB.QueryRow("select count(*) as Number from agents").Scan(
		&count.Number,
	)
	return
}
func (c *CCUser) AddCCUser(ccUser CCUser) (id int64, err error) {
	rs, err := db.SqlDB.Exec("INSERT INTO agents (NAME,instance_id,type,contact,STATUS,state,max_no_answer,wrap_up_time,reject_delay_time,busy_delay_time,no_answer_delay_time) "+
		"VALUES(?,'single_box','callback',?,'Logged Out','Waiting',?,?,?,?,?)",
		ccUser.Name, ccUser.Contact, ccUser.Max_no_answer, ccUser.Wrap_up_time, ccUser.Reject_delay_time, ccUser.Busy_delay_time, ccUser.No_answer_delay_time)
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
func (c *CCUser) EditCCUser(ccUser CCUser) (id int64, err error) {
	fmt.Println("UPDATE agents SET `name`=?,contact=?,max_no_answer=?,wrap_up_time=?,reject_delay_time=?,busy_delay_time=?,no_answer_delay_time=? where `name`=?",
		ccUser.Name, ccUser.Contact, ccUser.Max_no_answer, ccUser.Wrap_up_time, ccUser.Reject_delay_time, ccUser.Busy_delay_time, ccUser.No_answer_delay_time, ccUser.Name)

	rs, err := db.SqlDB.Exec("UPDATE agents SET `name`=?,contact=?,max_no_answer=?,wrap_up_time=?,reject_delay_time=?,busy_delay_time=?,no_answer_delay_time=? where `name`=?",
		ccUser.Name, ccUser.Contact, ccUser.Max_no_answer, ccUser.Wrap_up_time, ccUser.Reject_delay_time, ccUser.Busy_delay_time, ccUser.No_answer_delay_time, ccUser.Name)
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
func (c *CCUser) DelCCUser() (id int64, err error) {
	rs, err := db.SqlDB.Exec("DELETE FROM agents WHERE name =?", c.Name)
	if err != nil {
		fmt.Println("[DelCCUser].err ", err)
	}
	id, err = rs.RowsAffected()
	return
}

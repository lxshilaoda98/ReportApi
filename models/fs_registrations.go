package models

import (
	"fmt"
	db "github.com/n1n1n1_owner/ReportApi/database"
)

type Registrations struct {
	Reg_user      string `json:"reg_user" form:"reg_user"`
	Realm         string `json:"realm" form:"realm"`
	Token         string `json:"token" form:"token"`
	Url           string `json:"url" form:"url"`
	Expires       string `json:"expires" form:"expires"`
	Network_ip    string `json:"network_ip" form:"network_ip"`
	Network_port  string `json:"network_port" form:"network_port"`
	Network_proto string `json:"network_proto" form:"network_proto"`
	Hostname      string `json:"hostname" form:"hostname"`
	Metadata      string `json:"metadata" form:"metadata"`
}

func (r *Registrations) GetRegistrations(start, end int) (registrations []Registrations, err error) {
	registrations = make([]Registrations, 0)
	fmt.Printf("开始%d,结束%d \n", start, end)
	rows, err := db.SqlDB.Query("select Reg_user,Realm,Token,Url,Expires,Network_ip,Network_port"+
		",Network_proto,Hostname,Metadata from registrations limit ?,?", start, end)
	if err != nil {
		fmt.Print("[GetRegistrations]select SQL Err..>", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var registration Registrations
		rows.Scan(&registration.Reg_user, &registration.Realm, &registration.Token, &registration.Url, &registration.Expires,
			&registration.Network_ip, &registration.Network_port, &registration.Network_proto, &registration.Hostname, &registration.Metadata)
		registrations = append(registrations, registration)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func (r *Registrations) GetRegistrationCount() (count Count) {
	_ = db.SqlDB.QueryRow("select count(*) as Number from registrations").Scan(
		&count.Number,
	)
	return
}

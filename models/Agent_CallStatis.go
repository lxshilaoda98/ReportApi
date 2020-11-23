package models

import (
	"fmt"
	db "github.com/n1n1n1_owner/ReportApi/database"
)
//坐席呼叫量统计
type Agent_CallStatis struct {
	AgentName string `json:"坐席名称" from:"AgentName"`
	QueueNum string `json:"排队数" from:"QueueNum"`
	AnswerNum string `json:"应答数" from:"AnswerNum"`
	HangUpNumber string `json:"放弃数" from:"HangUpNumber"`
}

func (a Agent_CallStatis) GetAgent_CallStatis(startTime, EndTime string, start, end int)(agent_CallStatis []Agent_CallStatis, err error) {

	fmt.Println("Stime.>",startTime)
	fmt.Println("Etime.>",EndTime)
	sql := "select agents.name as AgentName," +
		"sum(org is not null) as QueueNum," +
		"sum(CCAgentAnsweredTime is not null) as AnswerNum," +
		"sum(CCancelReason = 'BREAK_OUT') as HangUpNumber " +
		"from agents left JOIN callstart on agents.name = callstart.CallerANI " +
		"where callstart.IvrStartTime BETWEEN ? and ? " +
		"group by agents.name LIMIT ?,?"

	agent_CallStatis = make([]Agent_CallStatis, 0)
	rows,err :=db.SqlDB.Query(sql,startTime, EndTime, start, end)
	if err !=nil{
		fmt.Println("query 【Agent_CallStatis】。Err.",err)
	}
	defer rows.Close()
	for rows.Next() {
		var agent_CallStatisNew Agent_CallStatis
		rows.Scan(&agent_CallStatisNew.AgentName, &agent_CallStatisNew.AnswerNum, &agent_CallStatisNew.HangUpNumber,
			&agent_CallStatisNew.QueueNum)
		agent_CallStatis = append(agent_CallStatis, agent_CallStatisNew)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return

}

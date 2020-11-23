package models

import (
	"fmt"
	db "github.com/n1n1n1_owner/ReportApi/database"
	"strings"
)
type CallCountStatis struct {
	Rq string `json:"日期" from:"rq"`
	Hrl string `json:"呼入量" from:"hrl"`
	Pdl string `json:"进入排队量" from:"pdl"`
	Wfq string `json:"未转坐席放弃量" from:"wfq"`
	Yd string `json:"坐席应答量" from:"yd"`
	Wyd string `json:"转坐席未答量" from:"wyd"`
	Fqs string `json:"三秒放弃数" from:"fqs"`
	Fql string `json:"三秒放弃率" from:"fql"`
	Jtl string `json:"排队接通率" from:"jtl"`
	Thsc string `json:"通话时长" from:"thsc"`
	Pjthsc string `json:"平均通话时长" from:"pjthsc"`
}
func convertStr(SelectType,name string)(str string) {
	str = "CONVERT("+name+", DATE)"
	switch SelectType {
	case "日":
		str = "DATE_FORMAT("+name+", '%m/%d')"
		break;
	case "月":
		str = "DATE_FORMAT("+name+", '%m' )"
		break;
	case "年":
		str = "DATE_FORMAT("+name+", '%Y' )"
		break;
	case "时":
		str = "DATE_FORMAT("+name+", '%H')"
		break
	}
	return
}
//综合呼叫统计
func (c CallCountStatis)GetCallCountStatis(startTime, EndTime string, start, end int,SelectType string)(callCountStatis []CallCountStatis, err error)  {
	var conver = convertStr(SelectType,"IvrStartTime")
	sql := "select :convert as rq,count(*) as hrl," +
		"sum(QueueStartTime is not null) as pdl," +
		"sum(QueueStartTime is null) as wfq," +
		"sum(CCAgentAnsweredTime is not NULL) as yd," +
		"sum(CCAgentAnsweredTime is NULL and CCAgent is not null) as wyd," +
		"sum(CCancelReason ='BREAK_OUT' and TIMESTAMPDIFF(SECOND,QueueStartTime,QueueEndTime)>3) AS fqs," +
		"concat(round((sum(CCancelReason ='BREAK_OUT' and TIMESTAMPDIFF(SECOND,QueueStartTime,QueueEndTime)>3) / sum(QueueStartTime is not null))*100,2),'%') AS fql," +
		"concat(round((sum(CCAgentAnsweredTime is not NULL) / sum(QueueStartTime is not null))*100,2),'%') AS jtl," +
		"sum(TIMESTAMPDIFF(SECOND,CCAgentCalledTime,CCAgentAnsweredTime)) as thsc," +
		"round(sum(TIMESTAMPDIFF(SECOND,CCAgentCalledTime,CCAgentAnsweredTime)) / sum(CCAgentAnsweredTime is not NULL),0) as pjthsc" +
		" from callstart where IvrStartTime BETWEEN ? and ? " +
		"GROUP BY :convert LIMIT ?,?"
	callCountStatis = make([]CallCountStatis, 0)
	sql = strings.Replace(sql,":convert",conver,-1)

	rows,err := db.SqlDB.Query(sql,startTime, EndTime, start, end)
	if err !=nil{
		fmt.Println("[综合呼叫统计] sql Err.",err)
	}
	defer rows.Close()

	for rows.Next() {
		var callCountStatisNew CallCountStatis
		rows.Scan(&callCountStatisNew.Rq,&callCountStatisNew.Hrl,&callCountStatisNew.Pdl,&callCountStatisNew.Wfq,
			&callCountStatisNew.Yd,
			&callCountStatisNew.Wyd,&callCountStatisNew.Fqs,&callCountStatisNew.Fql,
			&callCountStatisNew.Jtl,&callCountStatisNew.Thsc,&callCountStatisNew.Pjthsc,)
		callCountStatis = append(callCountStatis, callCountStatisNew)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return

}

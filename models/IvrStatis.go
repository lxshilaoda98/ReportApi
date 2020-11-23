package models

import (
	"fmt"
	db "github.com/n1n1n1_owner/ReportApi/database"
)

type IvrStatis struct {
	Time    string `json:"时间" from:"time"`
	Num     string `json:"个数" from:"num"`
	SumTime string `json:"总时长" from:"sumTime"`
}

func (i *IvrStatis) GetIvrStatis(types, startTime, EndTime string, start, end int) (ivrStatis []IvrStatis, err error) {
	sql := "select DATE_FORMAT(IvrStartTime, '%m/%d') as Time,count(*) as num,sum(TIMESTAMPDIFF(SECOND,IvrStartTime,CallEndTime)) as sumTime from callstart " +
		"where IvrStartTime BETWEEN ? and ? group by CONVERT(IvrStartTime, DATE) LIMIT ?,?"
	switch types {
	case "月":
		sql = "select DATE_FORMAT( IvrStartTime, '%m' ) as Time,count(*) as num,sum(TIMESTAMPDIFF(SECOND,IvrStartTime,CallEndTime)) as sumTime " +
			" from callstart where IvrStartTime BETWEEN ? and ? " +
			"group by DATE_FORMAT( IvrStartTime, '%m' ) LIMIT ?,?"
	case "年":
		sql = "select DATE_FORMAT( IvrStartTime, '%Y' ) as Time,count(*) as num,sum(TIMESTAMPDIFF(SECOND,IvrStartTime,CallEndTime)) as sumTime " +
			" from callstart where IvrStartTime BETWEEN ? and ? " +
			"group by DATE_FORMAT( IvrStartTime, '%Y' ) LIMIT ?,?"
	case "时":
		sql = "select DATE_FORMAT(IvrStartTime, '%H') as Time,count(*) as num,sum(TIMESTAMPDIFF(SECOND,IvrStartTime,CallEndTime)) as sumTime " +
			" from callstart where IvrStartTime BETWEEN ? and ? " +
			"group by DATE_FORMAT(IvrStartTime, '%H') LIMIT ?,?"
	default:
		sql = "select DATE_FORMAT(IvrStartTime, '%m/%d') as Time,count(*) as num,sum(TIMESTAMPDIFF(SECOND,IvrStartTime,CallEndTime)) as sumTime from callstart " +
			"where IvrStartTime BETWEEN ? and ? group by CONVERT(IvrStartTime, DATE) LIMIT ?,?"
	}
	fmt.Printf("IVR呼叫量统计呼入参数：[type] %s..> stime:%s .>etime:%s .> start:%d .> end:%d .>", types,startTime, EndTime, start, end)
	ivrStatis = make([]IvrStatis, 0)

	rows, err := db.SqlDB.Query(sql, startTime, EndTime, start, end)
	if err != nil {
		fmt.Print("select SQL Err..>", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var ivrStatisNew IvrStatis
		rows.Scan(&ivrStatisNew.Time, &ivrStatisNew.Num, &ivrStatisNew.SumTime)
		ivrStatis = append(ivrStatis, ivrStatisNew)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

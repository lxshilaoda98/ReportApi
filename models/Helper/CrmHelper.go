package Helper

import (
	"fmt"
	db "github.com/n1n1n1_owner/ReportApi/database"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func InsterGateway(Number int) {
	name := ""
	for i := 0; i <= Number; i++ {
		name = "测试" + strconv.Itoa(i)
		rs, err := db.SqlDB.Exec("INSERT INTO typelist(title, Name,typeL,sort) VALUES (?, ?,?,?)", name, name, "", i)
		if err != nil {
			fmt.Println("插入库异常..>", err)
		} else {
			fmt.Println("正常循环插入数据..>", rs, ".插入成功!")
		}
	}

}

//查询归属地，通过传入的号码
func getGsdForSql(Ani string) (CityName string) {
	//CityName 默认为宁波,防止代码异常导致电话无法接入
	CityName = "宁波"
	fmt.Println("截取后传入的值为:", Ani)
	rows, err := db.SqlDB.Query("select City from phoneterritory where Number = ? ", Ani)
	if err != nil {
		fmt.Println("select sql Err..>", err)
		return
	} else {
		if rows.Next() {
			fmt.Println("找到对应的归属地城市")
			rows.Scan(&CityName)
			fmt.Println("城市：", CityName)

		} else {
			fmt.Println("没有找到对应归属地的城市.")
		}
	}
	defer rows.Close()
	return
}

//周几 转成 中文
func tranWeekday(week int) (weekName string) {
	//防止报错，weekName 默认为周一
	switch week {
	case 1:
		weekName = "周一"
		break
	case 2:
		weekName = "周二"
		break
	case 3:
		weekName = "周三"
		break
	case 4:
		weekName = "周四"
		break
	case 5:
		weekName = "周五"
		break
	case 6:
		weekName = "周六"
		break
	case 0:
		weekName = "周日"
		break
	default:
		weekName = "周一"
		break
	}
	return
}

//判断dqtime时间是否在stime和etime之间
func tranTimeSub(dqTime, sTime, eTime string) bool {
	dqNew, _ := strconv.Atoi(strings.Split(dqTime, ":")[0])
	sbNew, _ := strconv.Atoi(strings.Split(sTime, ":")[0])
	xbNew, _ := strconv.Atoi(strings.Split(eTime, ":")[0])
	if dqNew >= sbNew && dqNew < xbNew {
		return true
	}
	return false
}

//根据号码查询归属地
func GetGsdForAni(Ani string) (City string) {
	fmt.Printf("要查询的电话号码为:%v \n", Ani)
	aniLen := len(Ani)
	if aniLen == 11 { //号码正常11位的
		fmt.Println("号码正常11位,继续去判断是否为正常的号码！")
		match, _ := regexp.MatchString(`(13[0-9]{9}$|14[0-9]{9}|15[0-9]|16[0-9]{9}|17[0-9]{9}${9}$|18[0-9]{9})$|(^(0\d{10})|^(0\d{2}-\d{8}))`, Ani)
		//fmt.Println(match)    //输出true
		if match {
			fmt.Println("号码正确，开始查找号码的归属地")
			//截取号码前7位
			AniNew := string([]rune(Ani)[:7])
			City = getGsdForSql(AniNew) //通过

		} else {
			fmt.Println("号码错误.号码为：", Ani)
		}
	} else {
		fmt.Println("号码不是11位.号码为：", Ani)
	}

	return
}

//查看上下班时间
func GetZgssxbsjForCity(city string) (boNew bool) {
	boNew = false
	sql := "select Oid from WorkTimeSetting where title like '%" + city + "%'"
	rows, err := db.SqlDB.Query(sql)
	if err != nil {
		fmt.Println("查询公司上下班时间Err.>", err)
	} else {
		if rows.Next() {
			t := time.Now()
			//dd, _ := time.ParseDuration("144h")
			//fmt.Println("dd..>",dd)
			//ts:=time.Now().Add(dd)
			//fmt.Println("周",tranWeekday(int(t.Weekday())))
			weekName := tranWeekday(int(t.Weekday())) //算出今天是周几
			var CityOid int
			rows.Scan(&CityOid)
			fmt.Println("找到对应的城市配置上下班时间的信息...如下:")
			fmt.Println("城市ID：", CityOid)
			fmt.Println("继续查询时间关联表[WorkTimeNormal]")
			rowsNew, errNew := db.SqlDB.Query("select DayOfWeek,WorkShift,ClosingTime from WorkTimeNormal where SetID = ? "+
				"and Checklist = 1 and DayOfWeek = ?", CityOid, weekName)
			if errNew != nil {
				fmt.Println("查询子表【WorkTimeNormal】，err..>", err)
			} else {
				if rowsNew.Next() {
					fmt.Println("找对对应的上下班时间数据.Begin")
					var DayOfWeek string
					var WorkShift string
					var ClosingTime string
					rowsNew.Scan(&DayOfWeek, &WorkShift, &ClosingTime)
					fmt.Println("DayOfWeek:", DayOfWeek)

					WorkShiftNew := strings.Split(strings.Split(WorkShift, "T")[1], "Z")[0]
					fmt.Println("上班转换成时间：", WorkShiftNew)

					ClosingTimeNew := strings.Split(strings.Split(ClosingTime, "T")[1], "Z")[0]
					fmt.Println("下班转换成时间：", ClosingTimeNew)
					fmt.Println("找对对应的上下班时间数据.End")

					fmt.Println("判断当前时间是否在上班..>")
					now1 := time.Now().Format("15:04:05")
					fmt.Println("当前时间：", now1)
					//截取时间
					boNew = tranTimeSub(now1, WorkShiftNew, ClosingTimeNew)
					fmt.Println("是否上班.>", boNew) //true 上班，false 下班
				} else {
					fmt.Println("没有数据了！")
				}
			}

			defer rowsNew.Close()
		} else {
			fmt.Println("没有找到对应的城市配置上下班时间的信息...")
		}
	}
	defer rows.Close()
	return boNew
}

//根据传入的秒 返回对应的天，或者小时 或者分钟

func ResolveTime(seconds int) (day int, hour int, minute int) {
	var (
		//定义每分钟的秒数
		SecondsPerMinute = 60
		//定义每小时的秒数
		SecondsPerHour = SecondsPerMinute * 60
		//定义每天的秒数
		SecondsPerDay = SecondsPerHour * 24
	)
	//每分钟秒数
	minute = seconds / SecondsPerMinute
	//每小时秒数
	hour = seconds / SecondsPerHour
	//每天秒数
	day = seconds / SecondsPerDay
	return
}

//转换计算机kb g tb ..等
func FormatFileSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

//返回G
func GetGbFileSize(fileSize uint64) (size string) {

	if fileSize < 1024 {
		return fmt.Sprintf("%.0f", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.0f", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.0f", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.0f", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.0f", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.0f", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

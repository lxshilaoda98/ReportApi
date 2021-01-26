package models

import (
	"encoding/json"
	"fmt"
	db "github.com/n1n1n1_owner/ReportApi/database"
	helper "github.com/n1n1n1_owner/ReportApi/models/Helper"
	"strings"
)

//工作时间结构体
type IvrModel struct {
	Name     string
	JsonStr  string
	NodeList []NodeLists
	LineList []LineList
}
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

//语音结构体
type MusicMode struct {
	Name     string
	FileList []fileList
	IsCheck  string
}
type fileList struct {
	Status     string
	Name       string
	Size       int64
	Percentage int64
	Uid        int64
	Raw        raw
	Response   response
}
type raw struct {
	Uid int64
}
type response struct {
	Code int8
	Msg  string
	Data data
}
type data struct {
	FileName string
	FilePath string
}

type IvrType struct {
	Oid           string
	TimeStart     []string
	TimeEnd       []string
	CheckedCities []string
	Time4         []string
	Time5         []string
}
type workTime struct {
	Sktime        string
	Setime        string
	Xktime        string
	Xetime        string
	Sweek         string
	Tdworktime    string
	Tdoffworktime string
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
						fmt.Println("添加成功!.继续添加子模块数据!")
						if rows, err := db.SqlDB.Exec("delete from ivr_views where fcomoid = ?", Oid); err != nil {
							fmt.Println("delete views err", err)
						} else {
							if count, err := rows.RowsAffected(); err != nil {
								fmt.Println("delete views aff err..>", err)
							} else {
								fmt.Println("delete views success ..>", count)
								//继续新增views
								for _, v := range i.NodeList {
									res, err := db.SqlDB.Exec("insert into ivr_views (mid,`name`,`type`,`left`,top,ico,state,fcomoid)values(?,?,?,?,?,?,?,?)", v.Id, v.Name, v.Type, v.Left, v.Top, v.Ico, v.State, Oid)
									if err != nil {
										fmt.Println("insert ivr_views err ..>", err)
									} else {
										count, err := res.RowsAffected()
										if err != nil {
											fmt.Println("insert vr_views aff err..>", err)
										} else {
											fmt.Println("insert success count .>", count)
										}
									}
								}
								fmt.Println("node insert success! next insert line..>")
								_, err := db.SqlDB.Exec("delete from ivr_viewsline where fcomoid = ?", Oid)
								if err != nil {
									fmt.Println("delete ivr_viewsline err..>", err)
								} else {
									for _, v := range i.LineList {
										res, err := db.SqlDB.Exec("insert into ivr_viewsline (`from`,`to`,`label`,fcomoid)values(?,?,?,?)", v.From, v.To, v.Label, Oid)
										if err != nil {
											fmt.Println("insert ivr_viewsline err ..>", err)
										} else {
											count, err := res.RowsAffected()
											if err != nil {
												fmt.Println("insert ivr_viewsline aff err..>", err)
											} else {
												fmt.Println("insert ivr_viewsline count .>", count)
											}
										}
									}
								}
							}
						}
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

func (i *IvrType) ResViewsType(id, JsonStr string) (err error) {
	wt := workTime{}
	fmt.Printf("id..>%s ,viewsTypeJson.>>%s \n", id, JsonStr)
	b := []byte(JsonStr)
	if err := json.Unmarshal(b, &i); err != nil {
		fmt.Println("err json..>", err)
	} else {
		sktime, err := helper.TranVueTime(i.TimeStart[0])
		if err != nil {
			fmt.Println("【sktime】err..>", err)
		} else {
			wt.Sktime = sktime
		}
		setime, err := helper.TranVueTime(i.TimeStart[1])
		if err != nil {
			fmt.Println("【setime】err..>", err)
		} else {
			wt.Setime = setime
		}
		xktime, err := helper.TranVueTime(i.TimeEnd[0])
		if err != nil {
			fmt.Println("【xktime】err..>", err)
		} else {
			wt.Xktime = xktime
		}
		xetime, err := helper.TranVueTime(i.TimeEnd[1])
		if err != nil {
			fmt.Println("【xetime】err..>", err)
		} else {
			wt.Xetime = xetime
		}

		weekStr, err := helper.TranWeekSplic(i.CheckedCities)
		if err != nil {
			fmt.Println("[tranWeek]..Err..>", err)
		} else {
			wt.Sweek = weekStr
		}
		if len(i.Time4) > 0 {
			tim4Str := ""
			for _, v := range i.Time4 {
				s, _ := helper.TranVueTimeForYY(v)
				tim4Str += s + ","
			}
			if len(tim4Str) > 0 {
				wt.Tdworktime = strings.TrimRight(tim4Str, ",")
			}
		}
		if len(i.Time5) > 0 {
			time5Str := ""
			for _, v := range i.Time5 {
				s, _ := helper.TranVueTimeForYY(v)
				time5Str += s + ","
			}
			if len(time5Str) > 0 {
				wt.Tdoffworktime = strings.TrimRight(time5Str, ",")
			}
		}
		count, err := insertWt(id, wt)
		if err != nil {
			fmt.Println("insertWt Err..>", err)
		} else {
			fmt.Println("insertWt success ..>", count)
		}

	}
	return
}

func (m *MusicMode) ResViewsTypeForMusic(id, JsonStr string) (err error) {
	fmt.Printf("id..>%s . views TypeJson.>>%s \n", id, JsonStr)
	b := []byte(JsonStr)
	if err := json.Unmarshal(b, &m); err != nil {
		fmt.Println("json Err..>", err)
	}
	err = insertFile(id, m)

	return
}
func insertFile(id string, m *MusicMode) (err error) {
	_, err = db.SqlDB.Exec("delete from ivr_viewstype where viewsOid = ?", id)
	if err != nil {
		fmt.Println("delete ivr_viewstype Err..>", err)
	} else {
		_, err = db.SqlDB.Exec("insert into ivr_viewstype (viewsOid,bf_relNumber,bf_interrupt)values(?,?,?)", id, m.Name, m.IsCheck)
		if err != nil {
			fmt.Println("insert ivr_viewstype Err..>", err)
		} else {
			_, err = db.SqlDB.Exec("delete from ivr_viewsfile where viewsOid =?", id)
			if err != nil {
				fmt.Println("clear ivr_viewFile Err..>", err)
			} else {
				fmt.Println("clear ivr_viewFile Success ")
				for _, v := range m.FileList {
					//if v.Status == "success" {
					//添加数据到数据库
					fmt.Println("insert ivr_viewFile begin")
					_, err = db.SqlDB.Exec("insert into ivr_viewsfile(fileName,filePath,viewsOid)values(?,?,?)", v.Response.Data.FileName, v.Response.Data.FilePath, id)
					if err != nil {
						fmt.Println("insert ivr_viewFile Err..>", err)
					} else {
						fmt.Println("insert ivr_viewFile Success.")
					}
					//} else {
					//	fmt.Println("Response Err.>", v.Status)
					//}
				}
			}
		}
	}
	return
}

//插入关于工作时间的属性
func insertWt(id string, wt workTime) (count int64, err error) {
	result, err := json.Marshal(wt)
	if err != nil {
		fmt.Println("toMarshal err..>", err)
	} else {
		fmt.Println(string(result))
	}
	_, err = db.SqlDB.Exec("delete from ivr_viewstype where viewsOid = ?", id)
	if err != nil {
		fmt.Println("delete ivr_viewstype err..>", err)
	} else {
		result, err := db.SqlDB.Exec("insert into ivr_viewstype(sktime,setime,xktime,xetime,sweek,tdworktime,tdoffworktime,viewsOid,`type`)values(?,?,?,?,?,?,?,?,?)",
			wt.Sktime, wt.Setime, wt.Xktime, wt.Xetime, wt.Sweek, wt.Tdworktime, wt.Tdoffworktime, id, "offTime")
		if err != nil {
			fmt.Println("insert into ivr_viewstype Err..>", err)
		} else {
			if count, err = result.RowsAffected(); err != nil {
				fmt.Println("err.>", err)
			} else {
				fmt.Println("insert into ivrviewsType count ..>", count)
			}

		}
	}

	return
}

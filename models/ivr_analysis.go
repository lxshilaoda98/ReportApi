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

//语音结构体begin
type MusicMode struct {
	Oid         string
	Title 		string
	Min         string
	Max         string
	Timeout     string
	Terminators string
	KxCheck     []string
	KxCheckStr  string
	Name        string
	FileList    []fileList
	IsCheck     string
	DataList    []resData
}
type fileList struct {
	Status     string
	Name       string
	Size       int64
	Percentage int64
	Uid        string
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
type resData struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Uid  string `json:"uid"`
}
type InFileData struct {
	FileName string
	FilePath string
}

//语音结构体end

//人员信息结构体
type AgentModel struct {
	IsCheck string `json:"isCheck"`
	Name    string `json:"name"`
	Err     string `json:"err"`
}

//技能组信息结构体
type GroupModel struct {
	Name string `json:"name"`
}

type GetIvrType struct {
	Id   string
	Type string
}

type IvrType struct {
	Oid           string
	TimeStart     []string
	TimeEnd       []string
	CheckedCities []string
	Time4         []string
	Time5         []string
}
type WorkTime struct {
	Sktime        string
	Setime        string
	Xktime        string
	Xetime        string
	Sweek         string
	Tdworktime    string
	Tdoffworktime string
}

func (i *IvrModel) GetJsonStr(Oid string) (JsonStr string, err error) {
	fmt.Println("Oid===>", Oid)
	row := db.SqlDB.QueryRow("select JsonStr from ivr_com where foid =?", Oid)
	row.Scan(&JsonStr)
	fmt.Println(JsonStr)
	return
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

//工作时间
func (i *IvrType) ResViewsType(id, JsonStr string) (err error) {
	wt := WorkTime{}
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

//插入关于工作时间的属性
func insertWt(id string, wt WorkTime) (count int64, err error) {
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

//语音模块
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
		fmt.Println("ll...>",m.KxCheck)
		var cs =""

		for _,v:=range m.KxCheck {
			cs += v + ","
		}
		var kxChen=""
		if len(cs) > 0 {
			kxChen =strings.TrimRight(cs, ",")
		}
		_, err = db.SqlDB.Exec("insert into ivr_viewstype (title,`type`,viewsOid,bf_relNumber,bf_interrupt,bf_Min,bf_Max,bf_Timeout,bf_Terminators,bf_kxCheck)values(?,?,?,?,?,?,?,?,?,?)", m.Title,"music", id, m.Name, m.IsCheck,m.Min,m.Max,m.Timeout,m.Terminators,kxChen)
		if err != nil {
			fmt.Println("insert ivr_viewstype Err..>", err)
		}
	}
	return
}

//语音添加成功后方法
func (i *InFileData) InsSuccess(id, JsonStr, fid string) (err error) {
	b := []byte(JsonStr)
	if err := json.Unmarshal(b, &i); err != nil {
		fmt.Println("json Err..>", err)
	}
	_, err = db.SqlDB.Exec("insert into ivr_viewsfile(fileName,filePath,viewsOid,fileOid)values(?,?,?,?)", i.FileName, i.FilePath, id, fid)
	if err != nil {
		fmt.Println("insert ivr_viewFile Err..>", err)
	} else {
		fmt.Println("insert ivr_viewFile Success.")
	}
	return
}

//语音删除后方法
func (i *InFileData) RmSuccess(fid string) (err error) {
	_, err = db.SqlDB.Exec("delete from ivr_viewsfile where fileOid =?", fid)
	if err != nil {
		fmt.Println("clear ivr_viewFile Err..>", err)
	} else {
		fmt.Println("clear ivr_viewFile Success ")
	}
	return
}

//坐席模块
func (a *AgentModel) ResViewsTypeForAgent(id, JsonStr string) (err error) {
	fmt.Printf("id..>%s . views TypeJson.>>%s \n", id, JsonStr)
	b := []byte(JsonStr)
	if err := json.Unmarshal(b, &a); err != nil {
		fmt.Println("json Err..>", err)
	}
	err = insertAgent(id, a)
	return
}
func insertAgent(id string, a *AgentModel) (err error) {
	//直接插入库数据
	_, err = db.SqlDB.Exec("delete from ivr_viewstype where viewsOid = ?", id)
	if err != nil {
		fmt.Println("delete ivr_viewstype Err..>", err)
	} else {
		_, err = db.SqlDB.Exec("insert into ivr_viewstype (`type`,viewsOid,ag_interrupt,ag_name,ag_errTran)values(?,?,?,?,?)", "agent", id, a.IsCheck, a.Name, a.Err)
		if err != nil {
			fmt.Println("insert ivr_viewstype Err..>", err)
		}
	}
	return
}

//技能组模块
func (g *GroupModel) ResViewsTypeForGroup(id, JsonStr string) (err error) {
	fmt.Printf("id..>%s . views TypeJson.>>%s \n", id, JsonStr)
	b := []byte(JsonStr)
	if err := json.Unmarshal(b, &g); err != nil {
		fmt.Println("json Err..>", err)
	}
	err = insertGroup(id, g)
	return
}
func insertGroup(id string, g *GroupModel) (err error) {
	//直接插入库数据
	_, err = db.SqlDB.Exec("delete from ivr_viewstype where viewsOid = ?", id)
	if err != nil {
		fmt.Println("delete ivr_viewstype Err..>", err)
	} else {
		_, err = db.SqlDB.Exec("insert into ivr_viewstype (`type`,viewsOid,gp_name)values(?,?,?)", "group", id, g.Name)
		if err != nil {
			fmt.Println("insert ivr_viewstype Err..>", err)
		}
	}
	return
}

//查询模块属性信息
//上下班时间属性
func (g *GetIvrType) OffTimeM(id, types string) (ivrType IvrType, err error) {
	rows := db.SqlDB.QueryRow("select Oid,Sktime,Setime,Xktime,Xetime,Sweek,Tdworktime,Tdoffworktime from ivr_viewstype where viewsOid = ? and `type` = ?", id, types)
	if err != nil {
		fmt.Println("select ivr_viewsType Err..>", err)
	} else {
		typeN := WorkTime{}
		rows.Scan(&ivrType.Oid,&typeN.Sktime, &typeN.Setime, &typeN.Xktime, &typeN.Xetime, &typeN.Sweek, &typeN.Tdworktime, &typeN.Tdoffworktime)
		ivrType.TimeStart = append(ivrType.TimeStart, "2016-10-10 "+typeN.Sktime, "2016-10-10 "+typeN.Setime)
		ivrType.TimeEnd = append(ivrType.TimeEnd, "2016-10-10 "+typeN.Xktime, "2016-10-10 "+typeN.Xetime)
		strs := strings.Split(typeN.Sweek, ",")
		for _, v := range strs {
			ivrType.CheckedCities = append(ivrType.CheckedCities, v)
		}
		t4 := strings.Split(typeN.Tdworktime, ",")
		for _, v := range t4 {
			ivrType.Time4 = append(ivrType.Time4, v)
		}
		t5 := strings.Split(typeN.Tdoffworktime, ",")
		for _, v := range t5 {
			ivrType.Time5 = append(ivrType.Time5, v)
		}
		fmt.Println(ivrType.Time4)
	}
	return
}

//播放音乐属性
func (g *GetIvrType) MusicM(id, types string) (musicMode MusicMode, err error) {
	datas := make([]resData, 0)
	data := resData{}
	row := db.SqlDB.QueryRow("select Oid,Title,bf_relNumber,bf_interrupt,bf_Min,bf_Max,bf_Timeout,bf_Terminators,bf_KxCheck from ivr_viewsType where type = ? and viewsOid=?", types, id)
fmt.Println("select Title,bf_relNumber,bf_interrupt,bf_Min,bf_Max,bf_Timeout,bf_Terminators,bf_KxCheck from ivr_viewsType where type = ? and viewsOid=?", types, id)
	row.Scan(&musicMode.Oid,&musicMode.Title,&musicMode.Name, &musicMode.IsCheck,&musicMode.Min,&musicMode.Max,&musicMode.Timeout,&musicMode.Terminators,&musicMode.KxCheckStr)
	rows, err := db.SqlDB.Query("select fileName,filePath,fileOid from ivr_viewsfile where viewsOid = ?", id)
	for rows.Next() {
		rows.Scan(&data.Name, &data.Url, &data.Uid)
		datas = append(datas, data)
	}
	if len(musicMode.KxCheckStr) > 0 {
		fs := strings.Split(musicMode.KxCheckStr,",")
		for _,v:=range fs  {
			musicMode.KxCheck = append(musicMode.KxCheck, v)
		}
	}
	fmt.Println(musicMode.KxCheck)
	musicMode.DataList = datas
	return
}

//坐席属性
func (g *GetIvrType) AgentM(id, types string) (agnetModel AgentModel, err error) {
	row := db.SqlDB.QueryRow("select ag_interrupt,ag_name,ag_errTran from ivr_viewsType where type = ? and viewsOid=?", types, id)

	row.Scan(&agnetModel.IsCheck, &agnetModel.Name, &agnetModel.Err)
	return
}

//技能组属性
func (g *GetIvrType) GroupM(id, types string) (groupModel GroupModel, err error) {
	row := db.SqlDB.QueryRow("select gp_name from ivr_viewsType where type = ? and viewsOid=?", types, id)
	row.Scan(&groupModel.Name)
	return
}

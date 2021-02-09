package models

import (
	"errors"
	"fmt"
	db "github.com/n1n1n1_owner/ReportApi/database"
	"github.com/n1n1n1_owner/ReportApi/models/Helper"
	"github.com/spf13/viper"
	"io/ioutil"
	"strings"
)
//ivr的配置文件
type IVRConfig struct {
	dir string
}
//公共的结构体
type ggModel struct {
	Label string
}

type MusicFile struct {
	FileName string
	FileFile string
	ViewsOid string
	FileOid  string
}
func GetIVRConfig()(config *viper.Viper) {
	config = Helper.GetIVRConfig()
	//直接反序列化为Struct
	var configjson IVRConfig
	if err := config.Unmarshal(&configjson); err != nil {
		fmt.Println(err)
	}
	return
}

//json转换成lua文件
func JsonAsLua(id string) (err error) {
	var Mid string
	var Mtype string
	g := ggModel{}
	//第一步 先通过id查找是否有开始节点type = begin
	row := db.SqlDB.QueryRow("select Mid from ivr_views where `type`='begin' and fcomoid = ?", id)
	row.Scan(&Mid)
	if Mid != "" {
		LuaPZ("begin", Mid,id,g) //先初始化一个lua文件..
		//查看下一个节点是什么？
		row = db.SqlDB.QueryRow("select `type`,mid from ivr_views where mid = (select `to` from ivr_viewsline where `from` =? and fcomoid =?)", Mid, id)
		row.Scan(&Mtype, &Mid)
		if Mid != "" && Mtype=="offTime" {

			LuaPZ(Mtype, Mid,id,g)
		}
	} else {
		err = errors.New("Error, IVR Model")
	}

	return
}
func LuaPZ(t, id,cid string,g ggModel) {
	//初始化格式
	switch t {
	case "begin":
		cshLua()
	case "offTime":
		readLua(t, id,cid)
		//时间节点下面控制 一个是 一个否
		//如果是时间offtime节点，就去查找是否有“是”和“否”的连接线
		//继续查看下一级目录
		readLua("workTimeMusic", id,cid)
		//rowx:= db.SqlDB.QueryRow("select bf_relNumber,bf_interrupt,bf_Min,bf_Max,bf_Timeout,bf_Terminators,bf_kxCheck from ivr_viewstype where mid = (select `to` from ivr_viewsline where `from` =? and fcomoid =? and label ='否')", Mid, id)
		readLua("offTimeMusic",id,cid)
		//继续判断下个节点是什么？？一个有几个节点？是否还是继续播放音乐 ，还是？
	case "agent":
		readLua("agent",id,cid)
	default:
		fmt.Println("NO. format file",t,id,cid)
	}
}

//读取文件
func readLua(s, id ,cid string) {
	config := GetIVRConfig()
	var FsDir = config.GetString("IVRConfig.dir")
	var Str =""
	fmt.Println("=======================fs地址====>",FsDir)
	ret, err := ioutil.ReadFile(FsDir+`/scripts/GOWelcome.lua`)
	if err != nil {
		fmt.Println("read file err.>", err)
		return
	}
	Str = string(ret)
	switch s {
	case "offTime":
		newStr := IFWorkTime(id)
		Str = strings.Replace(Str, "--funcMenu", newStr, -1)
		Str1 := "function ErrKey(errKey)\n"
		Str1 += "session:streamFile(\"IVRWav/ErrKey.wav\");\n"
		Str1 += " if errKey ==\"WorkTime\" then \n"
		Str1 += " WorkTime()\n"
		Str1 += " end\n"
		Str1 += "end\n"
		Str1 += "function WorkTime()\n"
		Str1 += "--WorkTimeMod\n"
		Str1 += "end\n"
		Str2 := "function OffTime()\n"
		Str2 += "--OffTimeMod\n"
		Str2 += "end\n"
		Str = strings.Replace(Str, "--funcWorkTime", Str1, -1)
		Str = strings.Replace(Str, "--funcOffTime", Str2, -1)
	case "workTimeMusic":
		sqlStr:="select bf_relNumber,bf_interrupt,bf_Min,bf_Max,bf_Timeout,bf_Terminators,bf_kxCheck,viewsOid " +
			"from ivr_viewstype where viewsOid = (select `to` from ivr_viewsline where `from` =? and fcomoid =? and label ='是')"
		musicStr := toMusic(s,id,cid,sqlStr)
		Str = strings.Replace(Str, "--WorkTimeMod", musicStr, -1)
	case "offTimeMusic":
		sqlStr:="select bf_relNumber,bf_interrupt,bf_Min,bf_Max,bf_Timeout,bf_Terminators,bf_KxCheck,viewsOid " +
			"from ivr_viewstype where viewsOid = (select `to` from ivr_viewsline where `from` =? and fcomoid =? and label ='否')"
		musicStr := toMusic(s,id,cid,sqlStr)
		Str = strings.Replace(Str, "--OffTimeMod", musicStr, -1)
	}
	resWrite(Str)
}
//初始化写入..
func cshLua() {
	var wireteString strings.Builder
	wireteString.WriteString("--funcWorkTime\n")
	wireteString.WriteString("--funcOffTime\n")
	wireteString.WriteString("--funcMusic\n")
	wireteString.WriteString("--funcAgent\n")
	wireteString.WriteString("--funcGroup\n")
	wireteString.WriteString("--funcMenu\n")
	wireteString.WriteString("session:answer()\n")
	wireteString.WriteString("session:setAutoHangup(false)\n")
	wireteString.WriteString("Menu()\n")

	resWrite(wireteString.String())

}
//写入工作时间
func IFWorkTime(id string) (newStr string) {
	row := db.SqlDB.QueryRow("select Sweek,Sktime,Setime,Xktime,Xetime,Tdworktime,Tdoffworktime from ivr_viewstype where viewsOid = ? ", id)
	wk := WorkTime{}
	row.Scan(&wk.Sweek, &wk.Sktime, &wk.Setime, &wk.Xktime, &wk.Xetime, &wk.Tdworktime, &wk.Tdoffworktime)
	wkSweek, _ := Helper.TranWeekNumber(wk.Sweek)
	localskSp := strings.Split(wk.Sktime, ":")
	localseSp := strings.Split(wk.Setime, ":")
	localxkSp := strings.Split(wk.Xktime, ":")
	localxeSp := strings.Split(wk.Xetime, ":")

	//早上上班时间
	zsbBfTime := localskSp[0] + localskSp[1] + localskSp[2]
	//早上下班时间
	zxbBfTime := localseSp[0] + localseSp[1] + localseSp[2]
	//下午上班时间
	xsbBfTime := localxkSp[0] + localxkSp[1] + localxkSp[2]
	//下午下班时间
	xxbBfTilem := localxeSp[0] + localxeSp[1] + localxeSp[2]

	newStr = "function Menu() \n"
	newStr += "	if not session:ready() then return end\n"
	newStr += "		local yy = os.date(\"%Y%m%d\");\n"
	newStr += "		local getWeek = os.date(\"%w\");\n"
	newStr += "		local getHour = os.date(\"%H\");\n"
	newStr += "		local getMinute = os.date(\"%M\");\n"
	newStr += "		local getSecond = os.date(\"%S\");\n"
	newStr += fmt.Sprintf("		local week= string.find(\"%v\", getWeek,1)\n", wkSweek)
	newStr += "		local ifWorkTime=0;\n"
	newStr += "		if week ~=nil\n"
	newStr += "		then\n"
	newStr += "		session:consoleLog(\"info\",\"今天为工作日,继续查找验证.>\\n\");\n"
	newStr += "		ifWorkTime =1\n"
	newStr += "		-- body\n"
	newStr += "		--上班--看下今天是否是特殊非工作日\n"
	newStr += fmt.Sprintf("		local xx= string.find(\"%v\", yy,1)\n", strings.Replace(wk.Tdoffworktime, "-", "", -1))
	newStr += "		if xx ~=nil then\n"
	newStr += "		--找到休息日了\n"
	newStr += "		ifWorkTime=0\n"
	newStr += "		session:consoleLog(\"info\",\"找到特殊非工作日，今天休息.>\\n\");\n"
	newStr += "		else\n"
	newStr += "		session:consoleLog(\"info\",\"没有找到特殊非工作日！！！今天上班.>\\n\");\n"
	newStr += "		end\n"
	newStr += "		else\n"
	newStr += "		--不上班，看下是否今天是上班日\n"
	newStr += fmt.Sprintf("		local td= string.find(\"%v\", yy,1)\n", strings.Replace(wk.Tdworktime, "-", "", -1))
	newStr += "		if td ~=nil then\n"
	newStr += "		--今天上班\n"
	newStr += "		ifWorkTime=1\n"
	newStr += "		session:consoleLog(\"info\",\"找到特殊工作日，今天上班.>\\n\");\n"
	newStr += "		else\n"
	newStr += "		session:consoleLog(\"info\",\"没有找到特殊工作日！！！今天休息.>\\n\");\n"
	newStr += "		end\n"
	newStr += "		end\n"
	newStr += "		if ifWorkTime ==1 then\n"
	newStr += "		local ttTime = getHour..getMinute..getSecond\n"
	newStr += "		session:consoleLog(\"info\",\"今天上班.>查找是否在上班时间内\"..ttTime..\"\\n\");\n"
	newStr += "		local intTTime = tonumber(ttTime)\n"
	newStr += fmt.Sprintf("		if intTTime >=%v and intTTime < %v then\n", zsbBfTime, zxbBfTime)
	newStr += "		-- body\n"
	newStr += "		session:consoleLog(\"info\",\"早上上班中\\n\");\n"
	newStr += fmt.Sprintf("		else if intTTime >=%v and intTTime <%v then\n", xsbBfTime, xxbBfTilem)
	newStr += "		session:consoleLog(\"info\",\"下午上班中\\n\");\n"
	newStr += "		else\n"
	newStr += "		session:consoleLog(\"info\",\"下班了\\n\");\n"
	newStr += "		ifWorkTime = 0\n"
	newStr += "		end\n"
	newStr += "		end\n"
	newStr += "		end\n"
	newStr += "		if ifWorkTime ==1\n"
	newStr += "		then\n"
	newStr += "		WorkTime()\n"
	newStr += "		else\n"
	newStr += "		OffTime()\n"
	newStr += "		end\n"
	newStr += "		session:consoleLog(\"info\",ifWorkTime..\"\\n\");\n"
	newStr += "		session:consoleLog(\"info\",\"week.>\"..getWeek..\"..>hour..>\"..getHour..\"..>getMinute..>\"..getMinute..\"..getSecond..>\"..getSecond);\n"
	newStr += "		end\n"
	return
}
//上下班的语音拼接
func toMusic(s,id,cid,sqlStr string)(Str string){

	fmt.Println("节点类型..>",s,id,sqlStr)
	mm := MusicMode{}
	mf := MusicFile{}
	fmt.Println("传过来的sql..>",sqlStr)
	roww := db.SqlDB.QueryRow(sqlStr, id, cid)
	roww.Scan(&mm.Name, &mm.IsCheck, &mm.Min, &mm.Max, &mm.Timeout, &mm.Terminators, &mm.KxCheckStr,&id)
	fmt.Println("mmid....>",mm.Name,id)
	newStr := strings.Builder{}
	fmt.Println("查找文件viewsOid....>",id,mm.Max,mm.Terminators)
	rows,err:= db.SqlDB.Query("select fileName,filePath,fileOid from ivr_viewsfile where viewsOid = ?",id)
	if err != nil {
		fmt.Println("查找录音文件失败..>",err)
		return
	}else{
		for rows.Next() {
			rows.Scan(&mf.FileName,&mf.FileFile,&mf.FileOid)
			spCheck:= strings.Split(mm.KxCheckStr,",")
			keyString:=""
			//val:=[]int{}
			for _,v:=range spCheck {
				switch v {
				case "*":
					keyString +="|\\\\*"
				case "#":
					keyString +="|#"
				default:
					//iv,_:=strconv.Atoi(v)
					//val=append(val, iv)
				}
			}
			//zz:="["+strconv.Itoa(Helper.Min(val...))+"-"+strconv.Itoa(Helper.Max(val...))+"]"+keyString
			zz:="["+mm.KxCheckStr+"]"+keyString
			kName:= "digit_"+mf.FileOid
			newStr.WriteString(fmt.Sprintf(kName+" = session:playAndGetDigits(%v,%v,%v,%v,\"\",\"IVRWav/%v\",\"IVRWav/ErrKey.wav\",\"^%v$\")\n",mm.Min,mm.Max,mm.Name,mm.Timeout,mf.FileName,zz))

			//拼接完成后，继续查找下一个节点
			rowsView,err:=db.SqlDB.Query("select `type`,mid,label from ivr_views left JOIN ivr_viewsline ON ivr_views.mid = ivr_viewsline.`to` where mid in (select `to` from ivr_viewsline where `from` =? and fcomoid =?)",id,cid)
			check(err)
			g :=ggModel{}
			for rowsView.Next() {
				rowsView.Scan(&s,&id,&g.Label)
				if g.Label!="" {
					newStr.WriteString(fmt.Sprintf("if(%v==\"%v\") then \n",kName,g.Label))
					if s == "agent" {
						serStr:=toAgent(id)
						newStr.WriteString(fmt.Sprintf("%v",serStr))
					}else if s == "group" {
						serStr:=toGroup(id)
						newStr.WriteString(fmt.Sprintf("%v",serStr))
					}else if s == "music" {
						 musicById(id,&newStr)
						//sql :="select bf_relNumber,bf_interrupt,bf_Min,bf_Max,bf_Timeout,bf_Terminators,bf_kxCheck,viewsOid from ivr_viewstype where viewsOid = ?"
						//如果找到了语音节点就继续生成
						//musp := toMusic(s,id,cid,sql)

						//fmt.Println("音乐模块生成单：",musp)

						//newStr.WriteString(fmt.Sprintf("%v \n",nextStr))
						//fmt.Println("音乐节点...>",id)
					}else{
						newStr.WriteString(fmt.Sprintf("--ErrKey%v(); \n",g.Label))
					}
					//newStr.WriteString(fmt.Sprintf("else \n"))
					//
					//newStr.WriteString(fmt.Sprintf("ErrKey(\"%v\"); \n","WorkTime"))

					newStr.WriteString(fmt.Sprintf("end \n"))
				}
			}
		}
	}
	fmt.Println("本次字符串===>",newStr.String())
	fmt.Println("本次字符串===id>",id)
	Str = newStr.String()

	return
}
//单独的语音模块拼接
func musicById(id string,n *strings.Builder)(nextStr string){
	fmt.Println("================begin进入语音节点模块处理!!!",id)
	mm := MusicMode{}
	mf := MusicFile{}
	row := db.SqlDB.QueryRow("select bf_relNumber,bf_interrupt,bf_Min,bf_Max,bf_Timeout,bf_Terminators,bf_KxCheck,viewsOid " +
		"from ivr_viewstype where viewsOid = ?", id)
	row.Scan(&mm.Name, &mm.IsCheck, &mm.Min, &mm.Max, &mm.Timeout, &mm.Terminators, &mm.KxCheckStr,&id)
	fmt.Println(fmt.Sprintf("本次查询的id为:%v, 本节点超时时间:%v",id,mm.Timeout))
	newStr := strings.Builder{}
	fmt.Println("查找文件viewsOid....>",id,mm.Max,mm.Terminators)
	rows,err:= db.SqlDB.Query("select fileName,filePath,fileOid from ivr_viewsfile where viewsOid = ?",id)
	if err != nil {
		fmt.Println("查找录音文件失败..>",err)
		return
	}else{
		for rows.Next() {
			rows.Scan(&mf.FileName,&mf.FileFile,&mf.FileOid)
			spCheck:= strings.Split(mm.KxCheckStr,",")
			keyString:=""
			//val:=[]int{}
			for _,v:=range spCheck {
				switch v {
				case "*":
					keyString +="|\\\\*"
				case "#":
					keyString +="|#"
				default:
					//iv,_:=strconv.Atoi(v)
					//val=append(val, iv)
				}
			}
			//strconv.Itoa(Helper.Min(val...))+"-"+strconv.Itoa(Helper.Max(val...))
			zz:="["+mm.KxCheckStr+"]"+keyString
			kName:= "digit_"+mf.FileOid
			newStr.WriteString(fmt.Sprintf(kName+" = session:playAndGetDigits(%v,%v,%v,%v,\"\",\"IVRWav/%v\",\"IVRWav/ErrKey.wav\",\"^%v$\")\n",mm.Min,mm.Max,mm.Name,mm.Timeout,mf.FileName,zz))

			fmt.Println("继续查询数据..>查看自己下面是否还有子节点..>",id)
			//拼接完成后，继续查找下一个节点
			rowsView,err:=db.SqlDB.Query("select `type`,mid,label from ivr_views left JOIN ivr_viewsline ON ivr_views.mid = ivr_viewsline.`to` where mid in (select `to` from ivr_viewsline where `from` =?)",id)
			check(err)
			g :=ggModel{}
			s:=""
			for rowsView.Next() {
				rowsView.Scan(&s,&id,&g.Label)
				if g.Label!="" {
					newStr.WriteString(fmt.Sprintf("if(%v==\"%v\") then \n",kName,g.Label))
					if s == "agent" {
						serStr:=toAgent(id)
						newStr.WriteString(fmt.Sprintf("%v",serStr))
					}else if s == "group" {
						serStr:=toGroup(id)
						newStr.WriteString(fmt.Sprintf("%v",serStr))
					}else if s == "music" {
						musicById(id,&newStr)
					}else{
						newStr.WriteString(fmt.Sprintf("--ErrKey%v(); \n",g.Label))
					}
					newStr.WriteString(fmt.Sprintf("end \n"))
				}
			}
		}
	}
	fmt.Println("本次字符串===>",n.String())

	//newStr.WriteString(fmt.Sprintf("%v \n",newStr.String()))
	fmt.Println("================end进入语音节点模块处理!!!",newStr.String())
	n.WriteString(fmt.Sprintf("%v \n",newStr.String()))
	nextStr = n.String()
	return
}
//组拼接
func toGroup(id string)(str string){
	groupName:=""
	res:=strings.Builder{}
	row:= db.SqlDB.QueryRow("select gp_name from ivr_viewstype where viewsOid = ?",id)
	row.Scan(&groupName)
	fmt.Println("查询组信息..>",groupName,id)
	if groupName!="" {
		res.WriteString("if not session:ready() then return end \n")
		res.WriteString("session:streamFile(\"/usr/local/freeswitch/sounds/ivrwav/key0.wav\"); \n")
		res.WriteString(fmt.Sprintf("session:execute(\"callcenter\",%v) \n",groupName))
		res.WriteString("return \"break\";\n")
	}else{
		res.WriteString("session:hangup();\n")
		res.WriteString("session:destroy();\n")
	}
	str = res.String()
	return
}
//转人拼接
func toAgent(id string)(str string){
	agentName:=""
	res:=strings.Builder{}
	row:= db.SqlDB.QueryRow("select ag_name from ivr_viewstype where viewsOid = ?",id)
	row.Scan(&agentName)
	if agentName!="" {
		res.WriteString("if not session:ready() then return end \n")
		res.WriteString("session:streamFile(\"/usr/local/freeswitch/sounds/ivrwav/key0.wav\"); \n")
		res.WriteString("prefix = \"{ignore_early_media=true}user/\".."+agentName+";\n")
		res.WriteString("session = freeswitch.Session(prefix);\n")
		res.WriteString("return \"break\";\n")
	}else{
		res.WriteString("session:hangup();\n")
		res.WriteString("session:destroy();\n")
	}
	str = res.String()
	return
}
//重新写入
func resWrite(wireteString string) {
	var d1 = []byte(wireteString)
	config := GetIVRConfig()
	var FsDir = config.GetString("IVRConfig.dir")
	err2 := ioutil.WriteFile(FsDir+`/scripts/GOWelcome.lua`, d1, 0666) //写入文件(字节数组)
	check(err2)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

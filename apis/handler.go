package apis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/n1n1n1_owner/ReportApi/models"
	Crmhelper "github.com/n1n1n1_owner/ReportApi/models/Helper"
	"github.com/n1n1n1_owner/ReportApi/models/ivrConter/jys"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"net/http"
	"strconv"
	"time"
	//jys "github.com/n1n1n1_owner/ReportApi/models/ivrConter/jys"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE ,PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

//type TypeList struct {
//	Name string
//	state string
//	store string
//	title string
//	typeL string
//}
//保存ivr返回的模块json
func IvrSave(c *gin.Context) {
	id := c.Request.FormValue("id")
	sJson := c.Request.FormValue("sJson")
	//fmt.Println(id, sJson)
	fmt.Println("step 1 Start")
	m := models.IvrModel{}
	err := m.ResJson(id, sJson)
	if err != nil {
		fmt.Println("step 1 Err ..", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"msg":  "错误:",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"msg":  "成功",
		})
	}
	fmt.Println("step 1 end")
}
func IvrSaveForType(c *gin.Context) {
	fmt.Println("ivrSaveForType  Start")
	id := c.Request.FormValue("id")
	sJson := c.Request.FormValue("sJson")
	sType := c.Request.FormValue("type")
	code := 20000
	msg := "成功"
	//fmt.Println(id)
	//fmt.Println(sJson)
	//fmt.Println(sType)
	switch sType {
	case "offTime":
		t := models.IvrType{}
		err := t.ResViewsType(id, sJson)
		if err != nil {
			code = 50000
			msg = "保存失败"
		} else {
			code = 20000
			msg = "成功"
		}
	case "music":
		t := models.MusicMode{}
		err := t.ResViewsTypeForMusic(id, sJson)
		if err != nil {
			code = 50000
			msg = "保存失败"
		} else {
			code = 20000
			msg = "成功"
		}
	case "agent":
		t := models.AgentModel{}
		err := t.ResViewsTypeForAgent(id, sJson)
		if err != nil {
			code = 50000
			msg = "保存失败"
		} else {
			code = 20000
			msg = "成功"
		}
	case "group":
		t := models.GroupModel{}
		err := t.ResViewsTypeForGroup(id, sJson)
		if err != nil {
			code = 50000
			msg = "保存失败"
		} else {
			code = 20000
			msg = "成功"
		}
	default:
		fmt.Println(sType)
		code = 50000
		msg = "失败，未找到方法"
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})

}
func GetViewsType(c *gin.Context) {
	id := c.Request.FormValue("id")
	typeStr := c.Request.FormValue("type")
	m := models.GetIvrType{}
	switch typeStr {
	case "offTime":
		Gtype, err := m.OffTimeM(id, typeStr)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 20000,
				"msg":  "错误:",
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":     20000,
				"msg":      "成功",
				"nodeList": Gtype,
			})
		}
	case "music":
		MType, err := m.MusicM(id, typeStr)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 20000,
				"msg":  "错误:",
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":     20000,
				"msg":      "成功",
				"nodeList": MType,
			})
		}
	case "agent":
		AGType, err := m.AgentM(id, typeStr)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 20000,
				"msg":  "错误:",
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":     20000,
				"msg":      "成功",
				"nodeList": AGType,
			})
		}
	case "group":
		GPType, err := m.GroupM(id, typeStr)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 20000,
				"msg":  "错误:",
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":     20000,
				"msg":      "成功",
				"nodeList": GPType,
			})
		}
	}

}
func MusicSaveFile(c *gin.Context) {
	id := c.Request.FormValue("id")
	sJson := c.Request.FormValue("sJson")
	fid := c.Request.FormValue("fid")
	m := models.InFileData{}
	err := m.InsSuccess(id, sJson, fid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"msg":  "错误:",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"msg":  "成功",
		})
	}
}
func MusicRmFile(c *gin.Context) {
	fid := c.Request.FormValue("fid")
	m := models.InFileData{}
	err := m.RmSuccess(fid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 50000,
			"msg":  "错误:",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"msg":  "成功",
		})
	}
}
//通过id 获取ivr的json数据
func GetIvrModel(c *gin.Context) {
	fid := c.Request.FormValue("fid")
	m := models.IvrModel{}
	JSonStr, err := m.GetJsonStr(fid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 50000,
			"msg":  "错误:",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"msg":  "成功",
			"data": JSonStr,
		})
	}

}
//通过id 解析成lua文件，解析成ivr
func SenIvrModel(c *gin.Context){
	result := models.ResultClass{}
	id:=c.Request.FormValue("id")
	//写入流程，json解析
	err:= models.JsonAsLua(id)
	if err != nil {
		result.Code = 50000
		result.Msg = err.Error()
	}else{
		result.Code = 20000
		result.Msg = "写入成功！"
	}
	c.JSON(http.StatusOK, result)
}

//文件相关
func FileUpload(c *gin.Context) {
	result := models.ResultClass{}
	file, header, err := c.Request.FormFile("file")
	fmt.Println("upLoad.Name.File.>", header.Filename)
	if err == nil {
		result.Up(file, header)
	} else {
		result.Code = -1
		result.Msg = "接收文件出错"
	}
	c.JSON(http.StatusOK, result)
}

//排队相关 begin
func GetMembers(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pagesize := c.DefaultQuery("pagesize", "10")
	var data_page, _ = strconv.Atoi(page)
	var data_pagesize, _ = strconv.Atoi(pagesize)

	start := (data_page - 1) * data_pagesize
	end := data_pagesize * 1
	r := models.Member{}
	registrations, err := r.GetMembers(start, end)
	//查询一共多少条数据
	count := r.GetMemberCount()

	if err != nil {
		fmt.Printf("sql error .>>", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":   20000,
		"result": registrations,
		"count":  len(registrations),
		"limit":  count.Number,
	})

	fmt.Printf("查看当前页:%d ,每页显示个数:%d \n", data_page, end)
}
func DelMember(c *gin.Context) {
	r := models.Member{}
	r.Uuid = c.Param("Oid")
	if id, err := r.DelMember(); err != nil {
		fmt.Println("【DelMember】 Err.>", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 50000,
			"msg":  err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"oid":  id,
		})
	}
}

//排队相关 end

//值列表 相关 begin
func EditTypeList(c *gin.Context) {
	r := models.TypeList{}
	r.Oid = c.Request.FormValue("oid")
	r.Name = c.Request.FormValue("name")
	r.Val = c.Request.FormValue("val")
	r.State = c.Request.FormValue("state")
	r.Sort = c.Request.FormValue("store")
	r.Title = c.Request.FormValue("title")
	r.TypeL = c.Request.FormValue("typeL")

	if id, err := r.EditTypeList(r); err != nil {
		fmt.Println("【EditTypeList】接受到错误信息Err.>", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 50000,
			"msg":  err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"oid":  id,
		})
	}
}
func DelTypeList(c *gin.Context) {
	r := models.TypeList{}
	r.Oid = c.Param("Oid")
	if id, err := r.DelTypeList(); err != nil {
		fmt.Println("【DelTypeList】 Err.>", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 50000,
			"msg":  err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"oid":  id,
		})
	}
}
func AddTypeList(c *gin.Context) {
	r := models.TypeList{}
	r.Name = c.Request.FormValue("name")
	r.Val = c.Request.FormValue("val")
	r.State = c.Request.FormValue("state")
	r.Sort = c.Request.FormValue("store")
	r.Title = c.Request.FormValue("title")
	r.TypeL = c.Request.FormValue("typeL")

	if id, err := r.AddTypeList(r); err != nil {
		fmt.Println("【AddTypeList】接受到错误信息Err.>", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 50000,
			"msg":  err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"oid":  id,
		})
	}
}
func GetTypeList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pagesize := c.DefaultQuery("pagesize", "10")
	var data_page, _ = strconv.Atoi(page)
	var data_pagesize, _ = strconv.Atoi(pagesize)

	start := (data_page - 1) * data_pagesize
	end := data_pagesize * 1
	r := models.TypeList{}
	registrations, err := r.GetTypeLists(start, end)
	//查询一共多少条数据
	count := r.GetTypeListCount()

	if err != nil {
		fmt.Printf("sql error .>>", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":   20000,
		"result": registrations,
		"count":  len(registrations),
		"limit":  count.Number,
	})

	fmt.Printf("查看当前页:%d ,每页显示个数:%d \n", data_page, end)
}
func GetTypeListByType(c *gin.Context) {
	types := c.DefaultQuery("type", "1")
	r := models.TypeList{}
	registrations, err := r.GetTypeListByType(types)
	if err != nil {
		fmt.Printf("sql error .>>", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":   20000,
		"result": registrations,
	})
}
func GetTypeListByVal(c *gin.Context) {
	types := c.DefaultQuery("val", "1")
	r := models.TypeList{}
	registrations, err := r.GetTypeListByVal(types)
	if err != nil {
		fmt.Printf("sql error .>>", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":   20000,
		"result": registrations,
	})
}

//值列表 相关 end

//sip用户相关 begin
func GetAllSipUser(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pagesize := c.DefaultQuery("pagesize", "10")
	var data_page, _ = strconv.Atoi(page)
	var data_pagesize, _ = strconv.Atoi(pagesize)

	start := (data_page - 1) * data_pagesize
	end := data_pagesize * 1
	s := models.SipUser{}
	registrations, err := s.GetSipUser(start, end)
	//查询一共多少条数据
	count := s.GetSipUserCount()

	if err != nil {
		fmt.Printf("sql error .>>", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":   20000,
		"result": registrations,
		"count":  len(registrations),
		"limit":  count.Number,
	})

	fmt.Printf("查看当前页:%d ,每页显示个数:%d \n", data_page, end)
}
func EditSipUser(c *gin.Context) {
	s := models.SipUser{}
	s.Oid = c.Request.FormValue("oid")
	s.SIPUser = c.Request.FormValue("sipuser")
	s.Password = c.Request.FormValue("password")
	s.Callgroup = c.Request.FormValue("callgroup")
	s.GroupName = c.DefaultPostForm("groupName", "测试")

	if id, err := s.EditSipUser(s); err != nil {
		fmt.Println("【EditSipUser】错误信息Err.>", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 50000,
			"msg":  err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"oid":  id,
		})
	}
}
func DelSipUser(c *gin.Context) {
	s := models.SipUser{}
	s.Oid = c.Param("Oid")
	if id, err := s.DelSipUser(); err != nil {
		fmt.Println("【DelSipUser】 Err.>", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 50000,
			"msg":  err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"oid":  id,
		})
	}
}
func AddSipUser(c *gin.Context) {
	s := models.SipUser{}

	s.SIPUser = c.Request.FormValue("sipuser")
	s.Password = c.Request.FormValue("password")
	s.Callgroup = c.Request.FormValue("callgroup")
	s.GroupName = c.DefaultPostForm("groupName", "测试")

	if id, err := s.AddSipUser(s); err != nil {
		fmt.Println("【AddSipUser】接受到错误信息Err.>", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 50000,
			"msg":  err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"oid":  id,
		})
	}
}

//sip用户相关 end

// 网关相关 begin
func GetGateWay(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pagesize := c.DefaultQuery("pagesize", "10")
	var data_page, _ = strconv.Atoi(page)
	var data_pagesize, _ = strconv.Atoi(pagesize)

	start := (data_page - 1) * data_pagesize
	end := data_pagesize * 1
	g := models.GateWay{}
	gateWays, err := g.GetGateWay(start, end)
	//查询一共多少条数据
	count := g.GetGateWayCount()

	if err != nil {
		fmt.Printf("sql error .>>", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":   20000,
		"result": gateWays,
		"count":  len(gateWays),
		"limit":  count.Number,
	})

	fmt.Printf("查看当前页:%d ,每页显示个数:%d \n", data_page, end)
}
func AddGateWay(c *gin.Context) {
	g := models.GateWay{}
	g.Name = c.Request.FormValue("name")
	g.Realm = c.Request.FormValue("realm")
	g.Username = c.Request.FormValue("username")
	g.Register = c.Request.FormValue("register")
	g.Memo = c.Request.FormValue("memo")
	g.Password = c.Request.FormValue("password")

	if id, err := g.AddGateWay(g); err != nil {
		fmt.Println("【AddGateWay】错误信息Err.>", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 50000,
			"msg":  err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"oid":  id,
		})
	}
}
func DelGateWay(c *gin.Context) {
	g := models.GateWay{}
	g.Oid = c.Param("Oid")
	if id, err := g.DelGateWay(); err != nil {
		fmt.Println("【DelGateWay】 Err.>", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 50000,
			"msg":  err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"oid":  id,
		})
	}
}
func EditGateWay(c *gin.Context) {
	g := models.GateWay{}
	g.Oid = c.Request.FormValue("oid")
	g.Name = c.Request.FormValue("name")
	g.Realm = c.Request.FormValue("realm")
	g.Username = c.Request.FormValue("username")
	g.Register = c.Request.FormValue("register")
	g.Memo = c.Request.FormValue("memo")
	g.Password = c.Request.FormValue("password")

	if id, err := g.EditGateWay(g); err != nil {
		fmt.Println("【EditGateWay】错误信息Err.>", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 50000,
			"msg":  err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"oid":  id,
		})
	}
}

// 网关相关 end

//查询注册用户
func GetRegistrations(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pagesize := c.DefaultQuery("pagesize", "10")
	var data_page, _ = strconv.Atoi(page)
	var data_pagesize, _ = strconv.Atoi(pagesize)

	start := (data_page - 1) * data_pagesize
	end := data_pagesize * 1
	r := models.Registrations{}
	registrations, err := r.GetRegistrations(start, end)
	//查询一共多少条数据
	count := r.GetRegistrationCount()

	if err != nil {
		fmt.Printf("sql error .>>", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":   20000,
		"result": registrations,
		"count":  len(registrations),
		"limit":  count.Number,
	})

	fmt.Printf("查看当前页:%d ,每页显示个数:%d \n", data_page, end)

}

//查询拨号计划 begin
func GetDialplan(c *gin.Context) {

	page := c.DefaultQuery("page", "1")
	pagesize := c.DefaultQuery("pagesize", "10")

	var data_page, _ = strconv.Atoi(page)
	var data_pagesize, _ = strconv.Atoi(pagesize)

	start := (data_page - 1) * data_pagesize
	end := data_pagesize * 1
	p := models.Dialplan{}
	dialplans, err := p.GetDialplan(start, end)
	//查询一共多少条数据
	count := p.GetDialplanCount()

	if err != nil {
		fmt.Printf("sql error .>>", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":   20000,
		"result": dialplans,
		"count":  len(dialplans),
		"limit":  count.Number,
	})

	fmt.Printf("查看当前页:%d ,每页显示个数:%d \n", data_page, end)
}
func GetDialplanByOid(c *gin.Context) {
	oid := c.DefaultQuery("oid", "1")
	var roid, _ = strconv.Atoi(oid)

	p := models.Dialplan{}
	dialplans, err := p.GetDialplanByOid(roid)

	if err != nil {
		fmt.Printf("sql error .>>", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":   20000,
		"result": dialplans,
	})
}
func EditDialplan(c *gin.Context) {
	d := models.Dialplan{}
	d.Oid = c.Request.FormValue("oid")
	d.Condition = c.Request.FormValue("condition")
	d.Context = c.Request.FormValue("context")
	d.Domain = c.Request.FormValue("domain")
	d.Expression = c.Request.FormValue("expression")
	d.Description = c.Request.FormValue("description")
	d.Name = c.Request.FormValue("name")

	fmt.Println("要修改的dialplan..Oid...>", d.Oid)
	if id, err := d.EditDialplan(d); err != nil {
		fmt.Println("【EditDialplan】错误信息Err.>", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 50000,
			"msg":  err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"oid":  id,
		})
	}
}
func AddDialplan(c *gin.Context) {
	d := models.Dialplan{}
	d.Oid = c.Request.FormValue("oid")
	d.Condition = c.Request.FormValue("condition")
	d.Context = c.Request.FormValue("context")
	d.Domain = c.Request.FormValue("domain")
	d.Expression = c.Request.FormValue("expression")
	d.Description = c.Request.FormValue("description")
	d.Name = c.Request.FormValue("name")
	if id, err := d.AddDialplan(d); err != nil {
		fmt.Println("【AddDialplan】错误信息Err.>", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 50000,
			"msg":  err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"oid":  id,
		})
	}
}
func DelDialplan(c *gin.Context) {
	d := models.Dialplan{}
	d.Oid = c.Param("Oid")
	if id, err := d.DelDialplan(); err != nil {
		fmt.Println("【DelDialplan】 Err.>", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 50000,
			"msg":  err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"oid":  id,
		})
	}
}

//查询拨号计划 end

//针对Dialplan_APP dialplanOid 查询
func GetDialplanAppByBOid(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pagesize := c.DefaultQuery("pagesize", "10")
	dialplan := c.DefaultQuery("dialplan", "1")

	var data_page, _ = strconv.Atoi(page)
	var data_pagesize, _ = strconv.Atoi(pagesize)

	dialplanOid, _ := strconv.Atoi(dialplan)
	start := (data_page - 1) * data_pagesize
	end := data_pagesize * 1
	d := models.DialplanApp{}
	dialplans, err := d.GetDialplanApp(dialplanOid, start, end)
	//查询一共多少条数据
	count := d.GetDialplanAppCount(dialplanOid)

	if err != nil {
		fmt.Printf("sql error .>>", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":   20000,
		"result": dialplans,
		"count":  len(dialplans),
		"limit":  count.Number,
	})

	fmt.Printf("查看当前页:%d ,每页显示个数:%d \n", data_page, end)
}
func EditDialplanApp(c *gin.Context) {
	d := models.DialplanApp{}
	d.Oid = c.Request.FormValue("oid")
	d.Application = c.Request.FormValue("application")
	d.Data = c.Request.FormValue("data")
	d.Sort = c.Request.FormValue("sort")

	if id, err := d.EditDialplanApp(d); err != nil {
		fmt.Println("【EditDialplanApp】错误信息Err.>", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 50000,
			"msg":  err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"oid":  id,
		})
	}
}
func AddDialplanApp(c *gin.Context) {
	d := models.DialplanApp{}
	d.Application = c.Request.FormValue("application")
	d.Data = c.Request.FormValue("data")
	d.Sort = c.Request.FormValue("sort")
	d.Dialplan = c.Request.FormValue("dialplanOid")

	if id, err := d.AddDialplanApp(d); err != nil {
		fmt.Println("【AddDialplanApp】错误信息Err.>", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 50000,
			"msg":  err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"oid":  id,
		})
	}
}
func DelDialplanApp(c *gin.Context) {
	d := models.DialplanApp{}
	d.Oid = c.Param("Oid")
	if id, err := d.DelDialplanApp(); err != nil {
		fmt.Println("【DelDialplanApp】 Err.>", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 50000,
			"msg":  err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"oid":  id,
		})
	}
}

//end

//呼叫中心相关 begin
func GetCCUserByName(c *gin.Context) {
	//name := c.DefaultQuery("name", "1")
	cc := models.CCUser{}
	dialplans, err := cc.GetCCUserByName()
	//查询一共多少条数据
	//count := cc.GetCCUserCount()
	if err != nil {
		fmt.Printf("sql error .>>", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":   20000,
		"result": dialplans,
		"count":  len(dialplans),
		//"limit":  count.Number,
	})
}
func GetCCUserAll(c *gin.Context) {

	page := c.DefaultQuery("page", "1")
	pagesize := c.DefaultQuery("pagesize", "10")

	var data_page, _ = strconv.Atoi(page)
	var data_pagesize, _ = strconv.Atoi(pagesize)

	start := (data_page - 1) * data_pagesize
	end := data_pagesize * 1
	cc := models.CCUser{}
	dialplans, err := cc.GetAllCCUser(start, end)
	//查询一共多少条数据
	count := cc.GetCCUserCount()

	if err != nil {
		fmt.Printf("sql Error .>>", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":   20000,
		"result": dialplans,
		"count":  len(dialplans),
		"limit":  count.Number,
	})

	fmt.Printf("查看当前页:%d ,每页显示个数:%d \n", data_page, end)
}
func AddCCUser(c *gin.Context) {
	cc := models.CCUser{}
	cc.Name = c.Request.FormValue("name")
	cc.Contact = c.Request.FormValue("contact")
	cc.Max_no_answer = c.Request.FormValue("max_no_answer")
	cc.Wrap_up_time = c.Request.FormValue("wrap_up_time")
	cc.Reject_delay_time = c.Request.FormValue("reject_delay_time")
	cc.Busy_delay_time = c.Request.FormValue("busy_delay_time")
	cc.No_answer_delay_time = c.Request.FormValue("no_answer_delay_time")

	if id, err := cc.AddCCUser(cc); err != nil {
		fmt.Println("【AddCCUser】错误信息Err.>", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 50000,
			"msg":  err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"oid":  id,
		})
	}
}
func EditCCUser(c *gin.Context) {
	cc := models.CCUser{}
	cc.Name = c.Request.FormValue("name")
	cc.Contact = c.Request.FormValue("contact")
	cc.Max_no_answer = c.Request.FormValue("max_no_answer")
	cc.Wrap_up_time = c.Request.FormValue("wrap_up_time")
	cc.Reject_delay_time = c.Request.FormValue("reject_delay_time")
	cc.Busy_delay_time = c.Request.FormValue("busy_delay_time")
	cc.No_answer_delay_time = c.Request.FormValue("no_answer_delay_time")

	if id, err := cc.EditCCUser(cc); err != nil {
		fmt.Println("【AddCCUser】错误信息Err.>", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 50000,
			"msg":  err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"oid":  id,
		})
	}
}
func DelCCUser(c *gin.Context) {
	d := models.CCUser{}
	d.Name = c.Param("Oid")
	if id, err := d.DelCCUser(); err != nil {
		fmt.Println("【DelCCUser】 Err.>", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 50000,
			"msg":  err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"oid":  id,
		})
	}
}

//呼叫中心相关 end

//获取呼入策略 begin
func GetTiers(c *gin.Context) {

	page := c.DefaultQuery("page", "1")
	pagesize := c.DefaultQuery("pagesize", "10")

	var data_page, _ = strconv.Atoi(page)
	var data_pagesize, _ = strconv.Atoi(pagesize)

	start := (data_page - 1) * data_pagesize
	end := data_pagesize * 1
	cc := models.Tiers{}
	dialplans, err := cc.GetTiers(start, end)
	//查询一共多少条数据
	count := cc.GetTiersCount()

	if err != nil {
		fmt.Printf("sql error .>>", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":   20000,
		"result": dialplans,
		"count":  len(dialplans),
		"limit":  count.Number,
	})

	fmt.Printf("查看当前页:%d ,每页显示个数:%d \n", data_page, end)
}
func AddTiers(c *gin.Context) {
	cc := models.Tiers{}
	cc.Queue = c.Request.FormValue("queue")
	cc.Agent = c.Request.FormValue("agent")
	cc.Level = c.Request.FormValue("level")
	cc.Position = c.Request.FormValue("position")

	if id, err := cc.AddTiers(cc); err != nil {
		fmt.Println("【AddTiers】错误信息Err.>", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 50000,
			"msg":  err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"oid":  id,
		})
	}
}
func EditTiers(c *gin.Context) {
	cc := models.Tiers{}
	cc.Oid = c.Request.FormValue("oid")
	cc.Queue = c.Request.FormValue("queue")
	cc.Agent = c.Request.FormValue("agent")
	cc.Level = c.Request.FormValue("level")
	cc.Position = c.Request.FormValue("position")

	if id, err := cc.EditTiers(cc); err != nil {
		fmt.Println("【EditTiers】错误信息Err.>", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 50000,
			"msg":  err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"oid":  id,
		})
	}
}
func DelTiers(c *gin.Context) {
	d := models.Tiers{}
	d.Oid = c.Param("Oid")
	if id, err := d.DelTiers(); err != nil {
		fmt.Println("【DelTiers】 Err.>", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 50000,
			"msg":  err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"oid":  id,
		})
	}
}

//获取呼入策略 end

//IVR呼叫量统计呼入参数
func GetIvrStatis(c *gin.Context) {

	page := c.DefaultQuery("page", "1")
	pagesize := c.DefaultQuery("pagesize", "10")
	types := c.DefaultQuery("type", "日")
	sTime_epoch := c.DefaultQuery("sTime_epoch", string(time.Now().Unix()))
	eTime_epoch := c.DefaultQuery("eTime_epoch", string(time.Now().Unix()))

	timeLayout := "2006-01-02 15:04:05"
	//sTime_epoch string 转换成int64
	sTime, err := strconv.ParseInt(sTime_epoch, 10, 64)
	if err != nil {
		fmt.Println("ParseInt. sTime_epoch .>Err.>", err)
	}
	var data_sTime = time.Unix(sTime, 0).Format(timeLayout)

	//eTime_epoch string 转换成int64
	eTime, err := strconv.ParseInt(eTime_epoch, 10, 64)
	if err != nil {
		fmt.Println("ParseInt. eTime_epoch .>Err.>", err)
	}
	var data_eTime = time.Unix(eTime, 0).Format(timeLayout)

	//分页相关 转换
	var data_page, _ = strconv.Atoi(page)
	var data_pagesize, _ = strconv.Atoi(pagesize)

	start := (data_page - 1) * data_pagesize
	end := data_pagesize * 1
	p := models.IvrStatis{}
	ivrStatis, err := p.GetIvrStatis(types, data_sTime, data_eTime, start, end)
	if err != nil {
		fmt.Printf("sql error .>>", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"result": ivrStatis,
		"count":  len(ivrStatis),
	})

	fmt.Printf("查看当前页:%d ,每页显示个数:%d \n", data_page, end)
}

//坐席呼叫量统计
func GetAgent_CallStatis(c *gin.Context) {

	page := c.DefaultQuery("page", "1")
	pagesize := c.DefaultQuery("pagesize", "10")
	sTime_epoch := c.DefaultQuery("sTime_epoch", string(time.Now().Unix()))
	eTime_epoch := c.DefaultQuery("eTime_epoch", string(time.Now().Unix()))

	timeLayout := "2006-01-02 15:04:05"
	//sTime_epoch string 转换成int64
	sTime, err := strconv.ParseInt(sTime_epoch, 10, 64)
	if err != nil {
		fmt.Println("ParseInt. sTime_epoch .>Err.>", err)
	}
	var data_sTime = time.Unix(sTime, 0).Format(timeLayout)

	//eTime_epoch string 转换成int64
	eTime, err := strconv.ParseInt(eTime_epoch, 10, 64)
	if err != nil {
		fmt.Println("ParseInt. eTime_epoch .>Err.>", err)
	}
	var data_eTime = time.Unix(eTime, 0).Format(timeLayout)

	//分页相关 转换
	var data_page, _ = strconv.Atoi(page)
	var data_pagesize, _ = strconv.Atoi(pagesize)

	start := (data_page - 1) * data_pagesize
	end := data_pagesize * 1
	p := models.Agent_CallStatis{}
	ivrStatis, err := p.GetAgent_CallStatis(data_sTime, data_eTime, start, end)
	if err != nil {
		fmt.Printf("[GetAgent_CallStatis] sql error .>>", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"result": ivrStatis,
		"count":  len(ivrStatis),
	})

	fmt.Printf("查看当前页:%d ,每页显示个数:%d \n", data_page, end)
}

//综合呼叫统计
func GetCallCountStatis(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pagesize := c.DefaultQuery("pagesize", "10")
	sTime_epoch := c.DefaultQuery("sTime_epoch", string(time.Now().Unix()))
	eTime_epoch := c.DefaultQuery("eTime_epoch", string(time.Now().Unix()))

	selectType := c.DefaultQuery("selectType", "日")

	timeLayout := "2006-01-02 15:04:05"
	//sTime_epoch string 转换成int64
	sTime, err := strconv.ParseInt(sTime_epoch, 10, 64)
	if err != nil {
		fmt.Println("ParseInt. sTime_epoch .>Err.>", err)
	}
	var data_sTime = time.Unix(sTime, 0).Format(timeLayout)

	//eTime_epoch string 转换成int64
	eTime, err := strconv.ParseInt(eTime_epoch, 10, 64)
	if err != nil {
		fmt.Println("ParseInt. eTime_epoch .>Err.>", err)
	}
	var data_eTime = time.Unix(eTime, 0).Format(timeLayout)

	//分页相关 转换
	var data_page, _ = strconv.Atoi(page)
	var data_pagesize, _ = strconv.Atoi(pagesize)

	start := (data_page - 1) * data_pagesize
	end := data_pagesize * 1
	p := models.CallCountStatis{}
	ivrStatis, err := p.GetCallCountStatis(data_sTime, data_eTime, start, end, selectType)
	if err != nil {
		fmt.Printf("[GetCallCountStatis] sql error .>>", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"result": ivrStatis,
		"count":  len(ivrStatis),
	})

	fmt.Printf("查看当前页:%d ,每页显示个数:%d \n", data_page, end)

}

//弃用 （给美康检验所写的ivr流程）
func GetJysAgent(c *gin.Context) {
	Ani := c.DefaultQuery("Ani", "17600082595")
	a := jys.Agent{}
	agent, err := a.GetAgentIDForJys(Ani)
	if err != nil {
		fmt.Printf("sql error .>>", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"result": agent,
	})
}

//#####系统监控方法

//系统信息 返回json
func GetComputeInfoJson(c *gin.Context) {

	info, _ := host.Info() //获取计算机信息
	fmt.Println(info)

	c.JSON(http.StatusOK, gin.H{
		"result": info,
	})
}

//系统信息，映射到页面
func GetComputeInfoMonitor(c *gin.Context) {
	info, _ := host.Info() //获取计算机信息
	day, hour, minute := Crmhelper.ResolveTime(int(info.Uptime))
	upTimeDay := strconv.Itoa(day) + "天"
	bootTime := time.Unix(int64(info.BootTime), 0).Format("2006-01-02 15:04:05")
	fmt.Println(day, hour, minute, upTimeDay, bootTime)

	//stats, e := host.Users()
	//if e!=nil{
	//	fmt.Println("get host users failed .>",e)
	//}else{
	//	fmt.Println(stats)
	//}

	c.HTML(200, "index.html", gin.H{
		"hostname":             info.Hostname,
		"uptime":               upTimeDay,
		"bootTime":             bootTime,
		"procs":                info.Procs,
		"os":                   info.OS,
		"platform":             info.Platform,
		"platformFamily":       info.PlatformFamily,
		"platformVersion":      info.PlatformVersion,
		"kernelVersion":        info.KernelVersion,
		"kernelArch":           info.KernelArch,
		"virtualizationSystem": info.VirtualizationSystem,
		"virtualizationRole":   info.VirtualizationRole,
		"hostID":               info.HostID,
	})
}

//获取cpu信息，映射到页面

func GetCpuInfoMonitor(c *gin.Context) {
	InfoStat, err := cpu.Info()
	if err != nil {
		fmt.Println("Get cpu info Err.>", err)
		return
	}
	iCpu, _ := cpu.Counts(true) //cpu逻辑数量
	fmt.Println("cpu逻辑核心数量", iCpu)
	wCpu, _ := cpu.Counts(false)
	fmt.Println("cpu物理核心数量", wCpu)
	fmt.Println(InfoStat)
	fmt.Println(len(InfoStat))

	//执行所消耗的时间
	//TimesStat包含CPU执行各种工作所花费的时间。 时间单位以秒为单位。 它基于linux / proc / stat文件。
	//用户时间（User  Time）即us所对应的列，表示CPU执行用户进程所占用的时间，通常情况下希望us的占比越高越好
	//系统时间（System Time）即sy所对应该的列，表示CPU自内核态所花费的时间，sy占比比较高通常意味着系统在某些方面设计得不合理，比如频繁的系统调用导致的用户态和内核态的频繁切换
	//Nice时间（Nice Time）即ni所对应的列，表示系统在调整进程优先级的时候所花费的时间
	//空闲时间（Idle Time）即id所对应的列，表示系统处于空闲期，等待进程运行，这个过程所占用的时间。当然，我们希望id的占比越低越好
	//等待时间（Waiting Time）即wa所对应的列，表示CPU在等待I/O操作所花费的时间，系统不应该花费大量的时间来进行等待，否则便表示可能有某个地方设计不合理
	//硬件中断处理时间（Hard Irq Time）即hi对应的列，表示系统处理硬件中断所占用的时间
	//软件中断处理时间（Soft Irq Time）即si对应的列，表示系统处理软件中断所占用的时间
	//丢失时间（Steal Time）即st对应的列，实在硬件虚拟化开始流行后操作系统新增的一列，表示被强制等待虚拟CPU的时间
	tsInfo, _ := cpu.Times(false)
	fmt.Println(tsInfo)
	//user User表示：CPU一共花了多少比例的时间运行在用户态空间或者说是用户进程(running user space processes)。
	// 典型的用户态空间程序有：Shells、数据库、web服务器……

	//Nice表示：可理解为，用户空间进程的CPU的调度优先级，范围为[-20,19]

	//System System的含义与User相似。System表示：CPU花了多少比例的时间在内核空间运行。分配内存、IO操作、创建子进程……都是内核操作。这也表明，当IO操作频繁时，System参数会很高。

	//Wait  在计算机中，读写磁盘的操作远比CPU运行的速度要慢，CPU负载处理数据，而数据一般在磁盘上需要读到内存中才能处理。
	//当CPU发起读写操作后，需要等着磁盘驱动器将数据读入内存(可参考：JAVA IO 以及 NIO 理解)，
	// 从而导致CPU 在等待的这一段时间内无事可做。CPU处于这种等待状态的时间由Wait参数来衡量。

	//Idle Idel表示：CPU处于空闲状态时间比例。一般而言，idel + user + nice 约等于100%

	//iowait iowait其实是一种特殊形式的CPU空闲

	//irq 中断请求？

	//softirq 软中断？

	//steal 你的虚拟机（VM）会与虚拟环境的宿主机上的多个虚拟机实例共享物理资源。
	// 其中之一共享的就是CPU时间切片。如果你的VM的物理机虚拟比是1/4，
	// 那么它的CPU使用率不会限制于25%的CPU时间切片－它能够使用超过它设置的虚拟比。（有别于内存的使用，内存大小是严格控制的）。

	//guest 当客户CPU执行hlt指令，它自愿产生CPU时间给其他客户。

	//百分比计算每个CPU使用或组合使用的cpu的百分比。 如果给定的间隔为0，它将把当前的CPU时间与上次调用进行比较。
	// 每个cpu返回一个值，如果percpu设置为false，则返回一个值。
	syinfo, _ := cpu.Percent(time.Duration(time.Second), false)
	fmt.Println(syinfo)

}

//获取物理内存和交换区内存信息
func GetMemInfoMonitor(c *gin.Context) {
	memInfo, _ := mem.SwapMemory()
	fmt.Println(memInfo)
	fmt.Println(Crmhelper.FormatFileSize(int64(memInfo.Total)))
	fmt.Println(Crmhelper.FormatFileSize(int64(memInfo.Used)))
	fmt.Println(Crmhelper.FormatFileSize(int64(memInfo.Free)))
	info2, _ := mem.SwapMemory()
	fmt.Println(info2)
}

//可以通过psutil获取磁盘分区、磁盘使用率和磁盘IO信息
func GetDiskInfoMonitor(c *gin.Context) {
	query := c.DefaultQuery("id", "0")
	info, _ := disk.Partitions(true) //所有分区
	//fmt.Println(info)
	numbers := []string{}          //磁盘信息
	SumNumber := []string{}        //磁盘总容量
	UseNumber := []string{}        //磁盘使用量
	DisBL := []string{}            //磁盘使用情况，用了多少+剩余多少
	infoStr := []*disk.UsageStat{} //磁盘列表详细信息
	infoStrForDiskName := []disk.IOCountersStat{}

	RWIoNumber := make(map[string]uint64)
	for k, _ := range info {
		numbers = append(numbers, info[k].Device)
		info2, _ := disk.Usage(info[k].Device) //指定某路径的硬盘使用情况

		i, e := strconv.Atoi(Crmhelper.GetGbFileSize(info2.Total))
		if e != nil {
			fmt.Println("convert String to Int Err..>", e)
		}
		fi, e := strconv.Atoi(Crmhelper.GetGbFileSize(info2.Free))
		if e != nil {
			fmt.Println("convert String to Int Err..>", e)
		}
		ui, _ := strconv.Atoi(Crmhelper.GetGbFileSize(info2.Used))

		infoStr = append(infoStr, info2)

		infoStr[k].Total = uint64(i)
		infoStr[k].Free = uint64(fi)
		infoStr[k].Used = uint64(ui)
		upc, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", info2.UsedPercent), 64)
		infoStr[k].UsedPercent = upc
		ipc, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", info2.InodesUsedPercent), 64)
		infoStr[k].InodesUsedPercent = ipc
		SumStr := Crmhelper.GetGbFileSize(info2.Total)
		SumNumber = append(SumNumber, SumStr)

		UseStr := Crmhelper.GetGbFileSize(info2.Used)
		UseNumber = append(UseNumber, UseStr)

		FeeStr := Crmhelper.GetGbFileSize(info2.Free)
		//添加指定的磁盘使用情况
		if query == info[k].Device {
			DisBL = append(DisBL, UseStr)
			DisBL = append(DisBL, FeeStr)
		}
	}
	info3, _ := disk.IOCounters("D:") //所有硬盘的io信息
	for m, v := range info3 {
		if query != "0" {
			if query == m {
				infoStrForDiskName = append(infoStrForDiskName, v)

				RWIoNumber["ReadCount"] = v.ReadCount
				RWIoNumber["WriteCount"] = v.WriteCount
				RWIoNumber["IopsInProgress"] = v.IopsInProgress
				//RWIoNumber=append(RWIoNumber,v.ReadCount,v.WriteCount,v.IopsInProgress )
			}
			//fmt.Println("找到值..>",query)
		}

		//fmt.Println("v..>",v,"m..>",m)
	}
	c.HTML(200, "DiskInfo.html", gin.H{
		"RWIoNumber":    RWIoNumber,
		"DisBL":         DisBL,
		"IoDiskInfo":    infoStrForDiskName, //渲染磁盘io
		"cMapDiskForIo": info3,
		"DiskInfo":      infoStr,
		"DiskLen":       numbers,
		"SumNumber":     SumNumber,
		"UseNumber":     UseNumber,
	})
}

//获取当前网络连接信息
func GetNetInfoMonitor(c *gin.Context) {
	info, _ := net.Connections("all") //可填入tcp、udp、tcp4、udp4等等 all 查询全部的连接
	//fmt.Println(info)

	//获取网络读写字节／包的个数
	info2, _ := net.IOCounters(false)
	//fmt.Println(info2)
	//c.HTML(200,"index.html",gin.H{
	//	"data":info2,
	//})
	c.HTML(http.StatusOK, "index.html", gin.H{
		"code": 0,
		"res":  info2,
		"nres": info,
	})
}

func GetProcessByIdInfoMonitor(c *gin.Context) {
	pid := c.DefaultQuery("id", "")
	pidInt, _ := strconv.Atoi(pid)

	fmt.Println("进入byID..>", c.DefaultQuery("id", ""))
	newProcess, _ := process.NewProcess(int32(pidInt))

	fmt.Println(newProcess)

	//info2,_ := process.GetWin32Proc(1120) //对应pid的进程信息
	//fmt.Println(info2)
	//
	//fmt.Println(info2[0].ParentProcessID) //获取父进程的pid

}

func GetFreeSwitchInfoMonitor(c *gin.Context) {
	c.HTML(http.StatusOK, "FreeswitchMonitor.html", gin.H{})
}

/**
转接到 软电话测试页面
*/
func GetSoftPhoneHtml(c *gin.Context) {
	c.HTML(http.StatusOK, "SoftPhone.html", gin.H{
		"message": "ssss",
	})
}
func GetSoftPhoneRTCHtml(c *gin.Context) {
	c.HTML(http.StatusOK, "SoftPhoneRTC.html", gin.H{})
}

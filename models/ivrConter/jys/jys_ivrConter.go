package jys

import (
	"fmt"
	helper "github.com/n1n1n1_owner/ReportApi/models/Helper"

)

type Agent struct {
	Id string
	Type string
}

/**
检验所流程 逻辑
先判断来电归属地--》根据归属地找子公司--》再判断当前时间上下班时间-》找坐席
接收参数 Ani = 来电号码
暂停
 */
func (a Agent)GetAgentIDForJys(Ani string) (agent Agent, err error){
	fmt.Println("检验所流程开始..>")
	CityStr:=helper.GetGsdForAni(Ani) //查询归属地
	fmt.Println("归属地：",CityStr)
	fmt.Println("根据归属地然后再去查询上下班时间.")
	sf := helper.GetZgssxbsjForCity(CityStr)//查询上下班时间
	if sf {
		fmt.Println("上班了，准备查询转入的坐席")
		agent.Type ="agent"
		agent.Id="6001"
	}else{
		agent.Type ="group"
		agent.Id="ningbo"
		fmt.Println("下班了，转到大组")
	}
	return
	//根据来电号码去查找归属地号码
}

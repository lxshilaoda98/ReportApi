package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/n1n1n1_owner/ReportApi/apis"
	"net/http"
	"time"
)

//定义一个中间件
func MiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		fmt.Println("开始中间件.")
		c.Set("request", "中间件")
		//执行
		c.Next()
		status := c.Writer.Status()
		fmt.Println("中间件执行完毕.>", status)
		t2 := time.Since(t)
		fmt.Println("执行时间：", t2)
	}
}
func InitRouter() *gin.Engine {
	router := gin.Default()
	//加载模板文件
	router.LoadHTMLGlob("./static/html/*")

	//router.LoadHTMLGlob("templates/monitor/*")

	router.StaticFS("/static", http.Dir("./static"))

	//router.Static("js","js")

	//加载中间件
	router.Use(MiddleWare())
	//router.GET("/", IndexApi)
	//
	//router.POST("/person", AddPersonApi)

	router.GET("/dialplans", GetDialplan)

	//router.GET("/person/:id", GetPersonApi)
	//
	//router.PUT("/person/:id", ModPersonApi)
	//
	//router.DELETE("/person/:id", DelPersonApi)

	router.GET("/IvrStatis", GetIvrStatis)

	router.GET("/Agent_CallStatis", GetAgent_CallStatis)

	router.GET("/CallCountStatis", GetCallCountStatis)

	router.GET("/GetJysAgent", GetJysAgent)

	//########测试路由

	//1.异步
	router.GET("/go_async", MiddleWare(), func(c *gin.Context) {
		cpStr := c.Copy()
		waremStr, _ := c.Get("request")
		fmt.Println("获取", waremStr)
		go func() {
			time.Sleep(3 * time.Second)
			fmt.Println("异步执行...", cpStr.Request.URL.Path)
		}()
	})

	groupName := router.Group("/monitor")
	{
		groupName.GET("/info", GetComputeInfoMonitor)
		groupName.GET("/cpu", GetCpuInfoMonitor)
		groupName.GET("/mem", GetMemInfoMonitor)
		groupName.GET("/disk", GetDiskInfoMonitor)
		groupName.GET("/net", GetNetInfoMonitor)

		//http://localhost:8000/monitor/ProcessInfo.html?id=5696
		groupName.GET("/ProcessInfo:id", GetProcessByIdInfoMonitor)

		groupName.GET("/freeswitch", GetFreeSwitchInfoMonitor)
	}
	router.GET("/softphone", GetSoftPhoneHtml)

	router.GET("/cookie", func(c *gin.Context) {
		s, e := c.Cookie("namekey")
		if e != nil {
			s = "未知cookie."
			//name, value string,
			//maxAge 过期时间,
			//path 路径
			//domain  作用域
			//secure  是否 【只能】 通过https访问，false是http就可以访问
			//httpOnly 是否运行别人获取cookie
			c.SetCookie("namekey", "lixin", 60, "/", "localhost", false, true)
		}
		fmt.Println(s)
	})

	//########系统信息路由
	router.GET("/GetComputeInfo", GetComputeInfoJson)

	return router
}

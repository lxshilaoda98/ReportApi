package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	//_ "github.com/mkevac/debugcharts"
	"github.com/shirou/gopsutil/host"
	"net/http"
	//_ "net/http/pprof"
)

func CInfo(w http.ResponseWriter, r *http.Request) {
	info, _ := host.Info()
	fmt.Println(info.String())
	w.Write([]byte(info.String()))
}

func main() {

	//创建路由
	r:=gin.Default()
	//绑定路由规则，执行函数
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK,"Hello Word!")
	})
	r.Run(":6633")
	//提供给负载均衡探活以及pprof调试


	//http.HandleFunc("/CpuInfo", CInfo)

	//http.Handle("/metrics", promhttp.Handler())

	//http.ListenAndServe(":10108", nil)
}

package main

import (
	"fmt"
	db "github.com/n1n1n1_owner/ReportApi/database"
	"github.com/n1n1n1_owner/ReportApi/models/Helper"
	"github.com/spf13/viper"
	//EslHelper "github.com/n1n1n1_owner/ReportApi/models/Helper"
)

type AppConfig struct {
	port string
}

func GetIVRConfig() (config *viper.Viper) {
	config = Helper.GetIVRConfig()
	//直接反序列化为Struct
	var configjson AppConfig
	if err := config.Unmarshal(&configjson); err != nil {
		fmt.Println(err)
	}
	return
}
func main() {
	defer db.SqlDB.Close()
	router := InitRouter()
	//开始esl 事件
	//go EslHelper.ConnectionEsl()
	//for i:=0;i<10;i++ {
	//	//开启10个进程添加测试数据。。。
	//	go EslHelper.InsterGateway(500000)
	//}
	config := GetIVRConfig()
	var appPort = config.GetString("AppConfig.port")

	router.Run(appPort)
}

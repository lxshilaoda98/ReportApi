package Helper

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"strings"
	. "github.com/0x19/goesl"
)
type EslConfig struct {
	fshost string
	fsport uint
	password string
	timeout int
}

//连接到FS，并监听数据
func ConnectionEsl()(config *viper.Viper){

	config = viper.New()
	config.AddConfigPath("./")
	config.SetConfigName("config")
	config.SetConfigType("json")
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}
	config.WatchConfig()
	config.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := config.ReadInConfig(); err != nil {
			panic(err)
		}
	})

	//直接反序列化为Struct
	var configjson EslConfig
	if err :=config.Unmarshal(&configjson);err !=nil{
		fmt.Println(err)
	}
	fmt.Printf("connection User:%s, Port:%d , PWd:%s , Tiemout:%d \n",
		config.GetString("EslConfig.fshost"),config.GetUint("EslConfig.fsport"),
		config.GetString("EslConfig.password"),config.GetInt("EslConfig.timeout"))

	client, err := NewClient(config.GetString("EslConfig.fshost"), config.GetUint("EslConfig.fsport"),
		config.GetString("EslConfig.password"), config.GetInt("EslConfig.timeout"))
	if err!=nil{
		fmt.Println("connect Go Esl Failed Err.>",err)
		return
	}else{
		fmt.Println("Connection Success ")
		go client.Handle()
		client.Send("events json ALL")

		for {
			msg, err := client.ReadMessage()
			if err != nil {
				// If it contains EOF, we really dont care...
				if !strings.Contains(err.Error(), "EOF") && err.Error() != "unexpected end of JSON input" {
					Error("Error while reading Freeswitch message: %s", err)
				}
				break
			}
			Debug("Got new message: %s", msg)
		}
	}
	return
}

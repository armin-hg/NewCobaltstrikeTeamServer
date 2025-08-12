package main

import (
	"NewCsTeamServer/client"
	"NewCsTeamServer/profile"
	"NewCsTeamServer/server/http"
	"NewCsTeamServer/server/manager"
	"NewCsTeamServer/utils"
	"flag"
	"fmt"
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
	profile.ProfileConfig = infos
	return nil
}
func main() {
	client.GlobalClientManager = client.NewClientManager()
	go http.HttpServer("8085", false)    //TODO 后续通过监听器，监听指定端口
	manager.RunTeamServer("8088", false) //管理端服务 TODO 后续使用websocket或其他方式进行通信
}

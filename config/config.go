package config

var PrivateKey []byte

var PublicKey []byte

var Iv = []byte("abcdefghijklmnop")

//var ProfileConfig Profile
//
//type Profile struct { //TODO 读取profile适配，现已有方案
//	CookieName  string
//	GetUrl      string
//	GetRetBody  string //空任务返回的字符串
//	PostUrl     string
//	PostQuery   string
//	PostRetBody string //Post后返回的字符串
//}

var Config Cfg

var ConfigFilePath = "./config.yaml"

type Cfg struct {
	Debug struct {
		Gin bool `yaml:"gin"`
		Sql bool `yaml:"sql"`
	} `yaml:"debug,omitempty"`
	Mysql struct {
		Addr     string `yaml:"addr"`
		UserName string `yaml:"user_name"`
		PassWord string `yaml:"pass_word"`
		DbName   string `yaml:"db_name"`
	} `yaml:"mysql,omitempty"`
	WebSocket struct {
		SSL       bool   `yaml:"ssl"` //是否开启ssl
		Port      string `yaml:"port"`
		SSLPort   string `yaml:"ssl_port"`
		CertFile  string `yaml:"cert_file"`
		KeyFile   string `yaml:"key_file"`
		LoginPass string `yaml:"login_pass"` //连接密码
	} `yaml:"websocket_server,omitempty"`
}

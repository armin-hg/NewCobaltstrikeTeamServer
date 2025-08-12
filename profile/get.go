package profile

import (
	"fmt"
	malleable "github.com/D00Movenok/goMalleable"
	"os"
	"strings"
)

func GetHttpPost(postdata malleable.HTTPPost) *ProConfig {
	var idcrypt []string
	p := new(ProConfig)
	p.HttpPost.Url = postdata.URI[0]

	ID := postdata.Client.ID
	for i := 0; i < len(ID); i++ {
		d := ID[i].String()
		d = strings.ReplaceAll(d, "\n", "")
		index := strings.Index(d, " ")
		if index != -1 {
			// 分割字符串成两部分
			firstPart := d[:index]
			secondPart := d[index+1:]
			strinfo := secondPart[1:]          //去掉第一个
			strinfo = strinfo[:len(strinfo)-1] //去掉最后一个
			strinfo = strinfo[:len(strinfo)-1] //去掉最后一个
			strinfo = strings.Replace(strinfo, "\\\"", "\"", -1)
			strinfo = strings.Replace(strinfo, "\\\\", "\\", -1)
			switch firstPart {
			case "header":
				p.HttpPost.IdType = "header"
				p.HttpPost.IdTypeValue = strinfo
			case "parameter":
				p.HttpPost.IdType = "parameter"
				p.HttpPost.IdTypeValue = strinfo
			case "append":
				p.HttpPost.IdAppend = strinfo
			case "prepend":
				p.HttpPost.IdPrepend = strinfo
			}
		}
		if d == "mask;" || d == "base64;" || d == "base64url;" || d == "netbios;" || d == "netbiosu;" {
			d = strings.Replace(d, ";", "", 1)
			idcrypt = append(idcrypt, d) //Http_post_id_crypt
		}
	}
	OutPut := postdata.Client.Output
	var outcrypt []string
	ServerOutPut := postdata.Server.Output
	ServerHeader := postdata.Server.Headers
	p.HttpPost.ServerHeader = ServerHeader
	var serverout string
	//var appendsss string
	for i := 0; i < len(ID); i++ {
		d := OutPut[i].String()
		d = strings.ReplaceAll(d, "\n", "")
		index := strings.Index(d, " ")
		if index != -1 {
			// 分割字符串成两部分
			firstPart := d[:index]

			secondPart := d[index+1:]
			strinfo := secondPart[1:]          //去掉第一个
			strinfo = strinfo[:len(strinfo)-1] //去掉最后一个
			strinfo = strinfo[:len(strinfo)-1] //去掉最后一个
			strinfo = strings.Replace(strinfo, "\\\"", "\"", -1)
			strinfo = strings.Replace(strinfo, "\\\\", "\\", -1)
			switch firstPart {
			case "header":
				p.HttpPost.ClientOutputType = "header"
				p.HttpPost.ClientOutputTypeValue = strinfo
			case "parameter":
				p.HttpPost.ClientOutputType = "parameter"
				p.HttpPost.ClientOutputTypeValue = strinfo
			case "append":
				p.HttpPost.ClientOutputAppend = strinfo
			case "prepend":
				p.HttpPost.ClientOutputPrepend = strinfo
			}
		}
		if d == "mask;" || d == "base64;" || d == "base64url;" || d == "netbios;" || d == "netbiosu;" {
			d = strings.Replace(d, ";", "", 1)
			outcrypt = append(outcrypt, d) //Http_post_id_crypt
		}

		if d == "print;" {
			p.HttpPost.ClientOutputType = "print"
			p.HttpPost.ClientOutputTypeValue = ""
		}
	}
	for i := 0; i < len(ServerOutPut); i++ { //Server Out put 解析
		d := ServerOutPut[i].String()
		d = strings.ReplaceAll(d, "\n", "")
		index := strings.Index(d, " ")
		if index != -1 {
			// 分割字符串成两部分
			firstPart := d[:index]
			secondPart := d[index+1:]
			strinfo := secondPart[1:] //去掉第一个
			lastIndex := strings.LastIndex(strinfo, "\";")
			if lastIndex != -1 {
				strinfo = strinfo[:lastIndex] + strings.Replace(strinfo[lastIndex:], "\";", "", 1)
			}
			strinfo = strings.Replace(strinfo, "\\\"", "\"", -1)
			strinfo = strings.Replace(strinfo, "\\\\", "\\", -1)
			strinfo = strings.Replace(strinfo, "\\r", "\r", -1)
			switch firstPart {
			case "prepend":
				serverout = strinfo + serverout
			case "append":
				serverout = serverout + strinfo
			}
		}

	}
	p.HttpPost.ServerOutput = serverout
	p.HttpPost.IdCrypt = idcrypt
	p.HttpPost.ClientOutputCrypt = outcrypt
	return p
}
func GetHttpGet(getdata malleable.HTTPGet) *ProConfig {
	var metadatacrypt []string
	p := new(ProConfig)
	p.HttpGet.Url = getdata.URI[0] //Get url
	Metadata := getdata.Client.Metadata
	for i := 0; i < len(Metadata); i++ {
		d := Metadata[i].String()
		d = strings.ReplaceAll(d, "\n", "")
		index := strings.Index(d, " ")
		if index != -1 {
			// 分割字符串成两部分
			firstPart := d[:index]
			secondPart := d[index+1:]
			strinfo := secondPart[1:]          //去掉第一个
			strinfo = strinfo[:len(strinfo)-1] //去掉最后一个
			strinfo = strinfo[:len(strinfo)-1]
			strinfo = strings.Replace(strinfo, "\\\"", "\"", -1)
			strinfo = strings.Replace(strinfo, "\\\\", "\\", -1)
			switch firstPart {
			case "append":
				p.HttpGet.MetadataAppend = strinfo
			case "header":
				p.HttpGet.MetadataType = "header"
				p.HttpGet.MetadataTypeValue = strinfo
			case "prepend":
				p.HttpGet.MetadataPrepend = strinfo
			}
		}
		if d == "mask;" || d == "base64;" || d == "base64url;" || d == "netbios;" || d == "netbiosu;" {
			d = strings.Replace(d, ";", "", 1)
			metadatacrypt = append(metadatacrypt, d) //获取get_metadata_crypt
		}
	}
	p.HttpGet.MetadataCrypt = metadatacrypt
	ServerHeader := getdata.Server.Headers
	p.HttpGet.ServerHeader = ServerHeader
	MetadataOut := getdata.Server.Output //http —— get  server 的配置
	//fmt.Println("MetadataOut:", MetadataOut)
	var outputcrypt []string
	var prependsss string
	var appendsss string
	for i := 0; i < len(MetadataOut); i++ {
		d := MetadataOut[i].String()
		d = strings.ReplaceAll(d, "\n", "")
		index := strings.Index(d, " ")
		if index != -1 {
			// 分割字符串成两部分
			firstPart := d[:index]
			secondPart := d[index+1:]
			strinfo := secondPart[1:] //去掉第一个
			lastIndex := strings.LastIndex(strinfo, "\";")
			if lastIndex != -1 {
				strinfo = strinfo[:lastIndex] + strings.Replace(strinfo[lastIndex:], "\";", "", 1)
			}
			strinfo = strings.Replace(strinfo, "\\\"", "\"", -1)
			strinfo = strings.Replace(strinfo, "\\\\", "\\", -1)
			strinfo = strings.Replace(strinfo, "\\r", "\r", -1)
			//fmt.Println("secondPart:", len(secondPart))
			switch firstPart {
			case "prepend":
				//prependlen = prependlen + len(strinfo)
				prependsss = strinfo + prependsss
			case "append":
				appendsss = appendsss + strinfo
			}
		}
		if d == "mask;" || d == "base64;" || d == "base64url;" || d == "netbios;" || d == "netbiosu;" {
			d = strings.Replace(d, ";", "", 1)
			outputcrypt = append(outputcrypt, d) //获取get_output_crypt
		}
	}
	//prependsss = strings.ReplaceAll(prependsss, "\n", "")
	//appendsss = strings.ReplaceAll(appendsss, "\n", "")
	p.HttpGet.OutPutPrepend = prependsss
	p.HttpGet.OutPutAppend = appendsss
	p.HttpGet.OutPutCrypt = outputcrypt
	return p
}
func GetProfile(data *os.File) *ProConfig {
	parsed, err := malleable.Parse(data)
	if err != nil {
		fmt.Println("[info]", "Please Use C2init Check Your Profile!")
		os.Exit(0)
		return nil
	}
	p := new(ProConfig)
	p.Header.UserAgent = strings.Replace(parsed.UserAgent, "\n", "", -1)
	p.Setting.Sleep = parsed.SleepTime
	p.Setting.Jitter = parsed.Jitter
	get := GetHttpGet(parsed.HTTPGet[0]).HttpGet
	if get.Url == "" {
		fmt.Println("[Info]", "Build Err.")
		return nil
	}
	post := GetHttpPost(parsed.HTTPPost[0]).HttpPost
	if post.Url == "" {
		fmt.Println("[Info]", "Build Err.")
		return nil
	}
	p.HttpGet = get
	p.HttpPost = post
	return p
	//fmt.Println(parsed.HTTPGet[0].Client.Metadata)
}

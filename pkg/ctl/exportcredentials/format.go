package exportcredentials

import (
	"encoding/json"
	"strings"
	"text/template"
)

type AliyunCLIConfig struct {
	Mode            string `json:"mode"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	StsToken        string `json:"sts_token,omitempty"`
}

type AliyunCLIURIBody struct {
	Code            string `json:"Code"`
	AccessKeyId     string `json:"AccessKeyId"`
	AccessKeySecret string `json:"AccessKeySecret"`
	SecurityToken   string `json:"SecurityToken"`
	Expiration      string `json:"Expiration,omitempty"`
}

var credentialFileFormat = `[default]
enable = true
{{- if eq .Mode "AK" }}
type = access_key
access_key_id = {{.AccessKeyId}}
access_key_secret = {{.AccessKeySecret}}
{{- else }}
type = sts
access_key_id = {{.AccessKeyId}}
access_key_secret = {{.AccessKeySecret}}
security_token = {{.StsToken}}
{{- end }}
`

func toAliyunCLIConfig(cred Credentials) AliyunCLIConfig {
	config := AliyunCLIConfig{
		Mode:            "AK",
		AccessKeyId:     cred.AccessKeyId,
		AccessKeySecret: cred.AccessKeySecret,
		StsToken:        cred.SecurityToken,
	}
	if config.StsToken != "" {
		config.Mode = "StsToken"
	}
	return config
}

func toAliyunCLIConfigJSON(cred Credentials) string {
	config := toAliyunCLIConfig(cred)
	data, _ := json.MarshalIndent(config, " ", " ")
	return string(data)
}

func toAliyunCLIURIBody(cred Credentials) string {
	config := AliyunCLIURIBody{
		Code:            "Success",
		AccessKeyId:     cred.AccessKeyId,
		AccessKeySecret: cred.AccessKeySecret,
		SecurityToken:   cred.SecurityToken,
		Expiration:      cred.Expiration,
	}
	data, _ := json.MarshalIndent(config, " ", " ")
	return string(data)
}

func toCredentialFileIni(cred Credentials) string {
	config := toAliyunCLIConfig(cred)
	t, _ := template.New("ini").Parse(credentialFileFormat)
	var buf strings.Builder
	_ = t.Execute(&buf, config)
	return buf.String()
}

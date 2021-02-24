package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	log "github.com/sirupsen/logrus"
	"os"
	tmpl "text/template"
)

// GetToken 根据用户id、用户密码、项目id获取 token, 返回的 token 用于调用 CallRawAPI 方法
//   @tokenURL: "http://192.168.1.100:9090/identity/v3/auth/tokens"
func GetToken(tokenURL string, userID string, password string, projectID string) (string, error) {
	js :=
`{
    "auth": {
        "identity": {
            "methods": [
                "password"
            ],
            "password": {
                "user": {
                    "id": "{{.UserID}}",
                    "password": "{{.Password}}"
                }
            }
        },
        "scope": {
            "project": {
                "id": "{{.ProjectID}}"
            }
        }
    }
}`

	mp := map[string]string{
		"UserID": userID,
		"Password": password,
		"ProjectID": projectID,
	}

	jsbody, err := CreateJsonByTmpl(js, mp)
	if err != nil {
		log.Error("create json with template failed: ", err)
		return "", err
	}

	header := make(map[string]string)
	header["Content-Type"] = "application/json"

	res, err := CallRestAPI(tokenURL, "POST", header, jsbody)
	if err != nil {
		log.Error("Failed in call rest api:", err)
	}

	token := res.Header.Get("X-Subject-Token")
	return token, err
}

func CreateJsonByTmpl(jstmpl string, mp map[string]string) ([]byte, error) {
	t, err := tmpl.New("tmp").Parse(jstmpl)
	if err != nil {
		log.Error("template parse error: ", err)
		return []byte{}, err
	}

	buf := &bytes.Buffer{}
	err = t.Execute(buf, mp)
	if err != nil {
		log.Error("template execute error: ", err)
		return []byte{}, err
	}

	m := make(map[string]interface{})
	_ = json.Unmarshal(buf.Bytes(), &m)
	return json.Marshal(m)
}

// AuthAndGetToken 先利用 apikey, tenantID 等进行认证，然后返回一个 token
// TODO 读取数据库, 用户认证
func AuthAndGetToken(apikey string, tenantID string, platformUserID string) (string, error) {
	resultTenantInfo,err:=QueryTenantInfoByTenantIdAndApikey(tenantID,apikey)
	if err!=nil {
		return "", err
	}
	if  resultTenantInfo.IsEmpty(){
		return "", errors.New("apikey无效，没有权限获取token")
	}
	tokenUrl:=os.Getenv("TOKEN_URL")
	token,err1:=GetToken(tokenUrl,resultTenantInfo.OpenstackUserid,resultTenantInfo.OpenstackPassword,resultTenantInfo.OpenstackProjectid)
	if err1!=nil {
		return "", err1
	}
	return token, nil
}

// CallRawAPI 返回 http response 的 body 部分
//   @url:     "http://192.168.66.131/compute/v2.1/servers"
//   @method:  "POST"
func CallRawAPI(url string, method string, bodyJson []byte, token string)([]byte, error) {
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["X-Auth-Token"] = token

	res, err := CallRestAPI(url, method, header, bodyJson)
	if err != nil {
		log.Error("Failed in call rest api:", err)
		return []byte{}, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("Failed in read the body data:", err)
		return []byte{}, err
	}

	return body, nil
}

func CallRestAPI(url string, method string, header map[string]string, bodyJson []byte) (*http.Response, error ){
	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyJson))
	if err != nil {
		log.Error("create new http request failed:", err)
		return nil, err
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Error("http client do failed:", err)
		return nil, err
	}

	return res, nil
}

package common

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

func MakeURLWithParams(u string, p map[string]string) string {
	url := u
	count := len(p)
	if count > 0 {
		url += "?"
	}
	for k, v := range p {
		url += k
		url += "="
		url += v
		if count > 1 {
			url += "&"
		}
		count--
	}
	return url
}

type CephMgrRestful struct {
	Endpoint string
	User     string
	Passwd   string
	CephfsID string
}

type CephMgrRESTRequest interface {
	DoRequest(string, string) error
}

func (o *CephMgrRestful) Process(r CephMgrRESTRequest) error {
	login := CephMgrRESTLogin{Url: "/api/auth"}
	token, err := login.getToken(o.Endpoint, o.User, o.Passwd)
	if err != EOK {
		return err
	}
	err = r.DoRequest(o.Endpoint, token)
	if err != nil && err != EOK {
		// LOG TRACE ERROR
		return err
	}
	return EOK
}

type CephMgrRESTListCephFSDirectory struct {
	Url    string            // input
	Params map[string]string // input
	Dirs   []string          // output
}

func (o *CephMgrRESTListCephFSDirectory) DoRequest(endpoint string, token string) error {
	url := MakeURLWithParams(o.Url, o.Params)
	header := make(map[string]string)
	header["Context-Type"] = "application/json"
	header["Authorization"] = "Bearer " + token
	res, err := CallRestAPI(endpoint+url, "GET", header, nil)

	defer res.Body.Close()
	resbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("Failed in read the body data:", err)
		return EIO
	}

	response := []interface{}{}
	json.Unmarshal(resbody, &response)
	for _, item := range response {
		//mitem := map[string]interface{}{}
		mitem, _ := item.(map[string]interface{})
		name, _ := mitem["name"].(string)
		o.Dirs = append(o.Dirs, name)
	}
	return nil
}

func (o *CephMgrRestful) ListCephFSDirectory(path string) ([]string, error) {
	params := make(map[string]string)
	params["path"] = path
	req := &CephMgrRESTListCephFSDirectory{Url: "/api/cephfs/" + o.CephfsID + "/ls_dir", Params: params}

	err := o.Process(req)
	if err != EOK {
		return nil, err
	}
	return req.Dirs, err
}

type CephMgrRESTMakeCephFSDirectory struct {
	Url  string // input
	Path string // input
}

func (o *CephMgrRESTMakeCephFSDirectory) DoRequest(endpoint string, token string) error {
	url := o.Url
	header := make(map[string]string)
	header["Context-Type"] = "application/json"
	header["Authorization"] = "Bearer " + token

	bodyTmpl := `{"path": "{{.path}}"}`
	mp := map[string]string{
		"path": o.Path,
	}
	body, err := CreateJsonByTmpl(bodyTmpl, mp)

	_, err = CallRestAPI(endpoint+url, "POST", header, body)
	if err != nil {
		return ECEPHMGRMKDIR
	}
	return EOK
}

func (o *CephMgrRestful) MakeCephFSDirectory(path string) error {
	req := &CephMgrRESTMakeCephFSDirectory{Url: "/api/cephfs/" + o.CephfsID + "/mk_dirs", Path: path}
	return o.Process(req)
}

type CephMgrRESTRemoveCephFSDirectory struct {
	Url  string // input
	Path string // input
}

func (o *CephMgrRESTRemoveCephFSDirectory) DoRequest(endpoint string, token string) error {
	url := o.Url
	header := make(map[string]string)
	header["Context-Type"] = "application/json"
	header["Authorization"] = "Bearer " + token

	bodyTmpl := `{"path": "{{.path}}"}`
	mp := map[string]string{
		"path": o.Path,
	}
	body, err := CreateJsonByTmpl(bodyTmpl, mp)

	_, err = CallRestAPI(endpoint+url, "POST", header, body)
	if err != nil {
		return ECEPHMGRRMDIR
	}
	return EOK
}

func (o *CephMgrRestful) RemoveCephFSDirectory(path string) error {
	req := &CephMgrRESTRemoveCephFSDirectory{Url: "/api/cephfs/" + o.CephfsID + "/rm_dir", Path: path}
	return o.Process(req)
}

type CephMgrRESTLogin struct {
	Url string
}

func (o *CephMgrRESTLogin) getToken(endpoint string, user string, passwd string) (string, error) {
	token := ""
	if endpoint == "" || user == "" || passwd == "" {
		return token, EPARAM
	}

	header := make(map[string]string)
	header["Context-Type"] = "application/json"
	header["Accept"] = "application/vnd.ceph.api.v1.0+json"
	bodyTmpl := `{"username": "{{.UserName}}", "password": "{{.Password}}"}`
	mp := map[string]string{
		"UserName": user,
		"Password": passwd,
	}
	body, err := CreateJsonByTmpl(bodyTmpl, mp)
	res, err := CallRestAPI(endpoint+o.Url, "POST", header, body)

	defer res.Body.Close()
	resbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("Failed in read the body data:", err)
		return token, EIO
	}

	response := map[string]interface{}{}
	json.Unmarshal(resbody, &response)
	token, ok := response["token"].(string)
	if !ok {
		return token, EPARSE
	}
	return token, EOK
}

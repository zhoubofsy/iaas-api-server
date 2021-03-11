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
	header["Content-Type"] = "application/json"
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

func (o *CephMgrRestful) ListCephFSDirectory(cephfsID string, path string) ([]string, error) {
	params := make(map[string]string)
	params["path"] = path
	req := &CephMgrRESTListCephFSDirectory{Url: "/api/cephfs/" + cephfsID + "/ls_dir", Params: params}

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
	header["Content-Type"] = "application/json"
	header["Authorization"] = "Bearer " + token

	bodyTmpl := `{"path": "{{.path}}"}`
	mp := map[string]string{
		"path": o.Path,
	}
	body, err := CreateJsonByTmpl(bodyTmpl, mp)

	res, err := CallRestAPI(endpoint+url, "POST", header, body)
	if err != nil || res.StatusCode != 200 {
		return ECEPHMGRMKDIR
	}
	return EOK
}

func (o *CephMgrRestful) MakeCephFSDirectory(cephfsID string, path string) error {
	req := &CephMgrRESTMakeCephFSDirectory{Url: "/api/cephfs/" + cephfsID + "/mk_dirs", Path: path}
	return o.Process(req)
}

type CephMgrRESTRemoveCephFSDirectory struct {
	Url  string // input
	Path string // input
}

func (o *CephMgrRESTRemoveCephFSDirectory) DoRequest(endpoint string, token string) error {
	url := o.Url
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["Authorization"] = "Bearer " + token

	bodyTmpl := `{"path": "{{.path}}"}`
	mp := map[string]string{
		"path": o.Path,
	}
	body, err := CreateJsonByTmpl(bodyTmpl, mp)

	res, err := CallRestAPI(endpoint+url, "POST", header, body)
	if err != nil || res.StatusCode != 200 {
		return ECEPHMGRRMDIR
	}
	return EOK
}

func (o *CephMgrRestful) RemoveCephFSDirectory(cephfsID string, path string) error {
	req := &CephMgrRESTRemoveCephFSDirectory{Url: "/api/cephfs/" + cephfsID + "/rm_dir", Path: path}
	return o.Process(req)
}

type CephMgrRESTGetCephFSQuotas struct {
	Url      string            // input
	Params   map[string]string // input
	MaxSize  int               // output
	MaxFiles int               // output
}

func (o *CephMgrRESTGetCephFSQuotas) DoRequest(endpoint string, token string) error {
	url := MakeURLWithParams(o.Url, o.Params)
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["Authorization"] = "Bearer " + token

	res, err = CallRestAPI(endpoint+url, "GET", header, nil)

	defer res.Body.Close()
	resbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("Failed in read the body data:", err)
		return EIO
	}

	response := make(map[string]int)
	err = json.Unmarshal(resbody, &response)
	if err != nil {
		return EPARSE
	}
	o.MaxSize = response["max_bytes"]
	o.MaxFiles = response["max_files"]
	return EOK
}

func (o *CephMgrRestful) GetCephFSQuotas(cephfsID string, path string) (int, int, error) {
	params := make(map[string]string)
	params["path"] = path

	req := &CephMgrRESTGetCephFSQuotas{Url: "/api/cephfs/" + cephfsID + "/get_quotas", Params: params}
	err := o.Process(req)
	if err != common.EOK {
		return 0, 0, err
	}
	return req.MaxSize, req.MaxFiles, err
}

type CephMgrRESTSetCephFSQuotas struct {
	Url      string // input
	Path     string // input
	MaxSize  int    // input
	MaxFiles int    // input
}

func (o *CephMgrRESTSetCephFSQuotas) DoRequest(endpoint string, token string) error {
	url := o.Url
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["Authorization"] = "Bearer " + token

	bodyTmpl := `{"path": "{{.path}}", "max_bytes": "{{.max_bytes}}", "max_files": "{{.max_files}}"}`
	mp := map[string]string{
		"path":      o.Path,
		"max_bytes": o.MaxSize,
		"max_files": o.MaxFiles,
	}
	body, err := CreateJsonByTmpl(bodyTmpl, mp)
	res, err := CallRestAPI(endpoint+url, "POST", header, body)
	if err != nil || res.StatusCode != 200 {
		return ECEPHMGRSETQUOTA
	}
	return EOK
}

func (o *CephMgrRestful) SetCephFSQuotas(cephfsID string, path string, maxsize int, maxfiles int) error {
	req := &CephMgrRESTSetCephFSQuotas{Url: "/api/ceph/" + cephfsID + "/set_quotas", Path: path, MaxSize: strconv.Itoa(maxsize), MaxFiles: strconv.Itoa(maxfiles)}
	return o.Process(req)
}

type GaneshaDaemonInfo struct {
	ClusterID   string `json:"cluster_id"`
	DaemonID    string `json:"daemon_id"`
	ClusterType string `json:"cluster_type"`
	Status      int    `json:"status"`
	StatusDesc  string `json:"status_desc"`
}
type CephMgrRESTListGaneshaDaemons struct {
	Url     string              // input
	Daemons []GaneshaDaemonInfo // output
}

func (o *CephMgrRESTListGaneshaDaemons) DoRequest(endpoint string, token string) error {
	url := o.Url
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["Authorization"] = "Bearer " + token
	res, err := CallRestAPI(endpoint+url, "GET", header, nil)
	if err != nil || res.StatusCode != 200 {
		return ECEPHMGRLISTGANESHADAEMON
	}

	defer res.Body.Close()
	resbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("Failed in read the body data:", err)
		return EIO
	}
	err = json.Unmarshal(resbody, &(o.Daemons))
	if err != nil {
		return EPARSE
	}
	return EOK
}

func (o *CephMgrRestful) ListGaneshaDaemons() ([]GaneshaDaemonInfo, error) {
	req := &CephMgrRESTListGaneshaDaemons{Url: "/api/nfs-ganesha/daemon"}
	err := o.Process(req)
	if err != EOK {
		return nil, err
	}
	return req.Daemons, err
}

type CephMgrRESTListGaneshaExport struct {
	Url string // input
}

func (o *CephMgrRESTListGaneshaExport) DoRequest(endpoint string, token string) error {
}

func (o *CephMgrRestful) ListGaneshaExport() error {
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
	header["Content-Type"] = "application/json"
	header["Accept"] = "application/vnd.ceph.api.v1.0+json"

	bodyOrigin := "{\"username\": \"" + user + "\", \"password\": \"" + passwd + "\"}"
	body := []byte(bodyOrigin)
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

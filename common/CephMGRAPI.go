package common

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"strconv"
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
	if err == nil {
		switch res.StatusCode {
		case 200:
			return EOK
		default:
			return ECEPHMGRRMDIR
		}
	}
	return ECEPHMGRRMDIR
}

func (o *CephMgrRestful) RemoveCephFSDirectory(cephfsID string, path string) error {
	req := &CephMgrRESTRemoveCephFSDirectory{Url: "/api/cephfs/" + cephfsID + "/recursive_rm_dir", Path: path}
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

	res, err := CallRestAPI(endpoint+url, "GET", header, nil)

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
	if err != EOK {
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
		"max_bytes": strconv.Itoa(o.MaxSize),
		"max_files": strconv.Itoa(o.MaxFiles),
	}
	body, err := CreateJsonByTmpl(bodyTmpl, mp)
	res, err := CallRestAPI(endpoint+url, "POST", header, body)
	if err != nil || res.StatusCode != 200 {
		return ECEPHMGRSETQUOTA
	}
	return EOK
}

func (o *CephMgrRestful) SetCephFSQuotas(cephfsID string, path string, maxsize int, maxfiles int) error {
	req := &CephMgrRESTSetCephFSQuotas{Url: "/api/cephfs/" + cephfsID + "/set_quotas", Path: path, MaxSize: maxsize, MaxFiles: maxfiles}
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

type FSALCephfs struct {
	FSName        string  `json:"fs_name"`
	Name          string  `json:"name"`
	SecLabelXattr *string `json:"sec_label_xattr"`
	UserID        string  `json:"user_id"`
}

type GaneshaClients struct {
	AccessType string   `json:"access_type"`
	Addresses  []string `json:"addresses"`
	Squash     string   `json:"squash"`
}

type GaneshaExportInfo struct {
	AccessType    string           `json:"access_type"`
	Clients       []GaneshaClients `json:"clients"`
	ClusterID     string           `json:"cluster_id",omitempty`
	Daemons       []string         `json:"daemons"`
	Path          string           `json:"path"`
	Protocols     []int            `json:"protocols"`
	Pseudo        string           `json:"pseudo"`
	ReloadDaemons bool             `json:"reload_daemons",omitempty`
	SecurityLabel bool             `json:"security_label"`
	Squash        string           `json:"squash"`
	Transports    []string         `json:"transports"`
	ExportID      int              `json:"export_id",omitempty`
	Tag           *string          `json:"tag"`
	FSAL          FSALCephfs       `json:"fsal"`
}

type CephMgrRESTListGaneshaExport struct {
	Url     string              // input
	Exports []GaneshaExportInfo // output
}

func (o *CephMgrRESTListGaneshaExport) DoRequest(endpoint string, token string) error {
	url := o.Url
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["Authorization"] = "Bearer " + token
	res, err := CallRestAPI(endpoint+url, "GET", header, nil)
	if err != nil || res.StatusCode != 200 {
		return ECEPHMGRLISTGANESHAEXPORT
	}
	defer res.Body.Close()
	resbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("Failed in read the body data:", err)
		return EIO
	}
	err = json.Unmarshal(resbody, &(o.Exports))
	if err != nil {
		return EPARSE
	}
	return EOK
}

func (o *CephMgrRestful) ListGaneshaExport() ([]GaneshaExportInfo, error) {
	req := &CephMgrRESTListGaneshaExport{Url: "/api/nfs-ganesha/export"}
	err := o.Process(req)
	if err != EOK {
		return nil, err
	}
	return req.Exports, err
}

type CephMgrRESTCreateGaneshaExport struct {
	Url  string // input
	Body []byte // input
}

func (o *CephMgrRESTCreateGaneshaExport) DoRequest(endpoint string, token string) error {
	url := o.Url
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["Authorization"] = "Bearer " + token

	res, err := CallRestAPI(endpoint+url, "POST", header, o.Body)
	if err != nil || res.StatusCode != 201 {
		return ECEPHMGRCREATEGANESHAEXPORT
	}
	return EOK
}

type CreateGaneshaSExportWithCephFS struct {
	AccessType    string           `json:"access_type"`
	Clients       []GaneshaClients `json:"clients"`
	ClusterID     string           `json:"cluster_id"`
	Daemons       []string         `json:"daemons"`
	FSAL          FSALCephfs       `json:"fsal"`
	Path          string           `json:"path"`
	Protocols     []int            `json:"protocols"`
	Pseudo        string           `json:"pseudo"`
	ReloadDaemons bool             `json:"reload_daemons"`
	SecurityLable bool             `json:"security_label"`
	Squash        string           `json:"squash"`
	Tag           *string          `json:"tag"`
	Transports    []string         `json:"transports"`
}

func (o *CephMgrRestful) CreateGaneshaExport(clusterID string, userID string, path string, pseudo string, dispatch []string, ips []string) error {
	exportConfig := &CreateGaneshaSExportWithCephFS{
		AccessType: "NONE",
		Clients: []GaneshaClients{
			GaneshaClients{
				AccessType: "RW",
				Squash:     "no_root_squash",
				Addresses:  ips,
			},
		},
		ClusterID: clusterID,
		FSAL: FSALCephfs{
			FSName: "cephfs",
			Name:   "CEPH",
			UserID: userID,
		},
		Path:          path,
		Protocols:     []int{4},
		Daemons:       dispatch,
		Pseudo:        pseudo,
		ReloadDaemons: true,
		SecurityLable: true,
		Squash:        "no_root_squash",
		Transports:    []string{"TCP"},
	}
	body, err := json.Marshal(exportConfig)
	if err != nil {
		return EPARSE
	}
	req := &CephMgrRESTCreateGaneshaExport{Url: "/api/nfs-ganesha/export", Body: body}
	return o.Process(req)
}

type CephMgrRESTPutGaneshaExport struct {
	Url  string // input
	Body []byte // input
}

func (o *CephMgrRESTPutGaneshaExport) DoRequest(endpoint string, token string) error {
	url := o.Url
	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["Authorization"] = "Bearer " + token

	res, err := CallRestAPI(endpoint+url, "PUT", header, o.Body)
	if err == nil {
		switch res.StatusCode {
		case 200:
			return EOK
		case 204:
			return EOK
		default:
			return ECEPHMGRPUTGANESHAEXPORT
		}
	}
	return ECEPHMGRPUTGANESHAEXPORT
}

func (o *CephMgrRestful) PutGaneshaExport(clusterID string, exportID string, context GaneshaExportInfo) error {
	context.ReloadDaemons = true
	body, err := json.Marshal(context)
	if err != nil {
		return EPARSE
	}
	req := &CephMgrRESTPutGaneshaExport{Url: "/api/nfs-ganesha/export/" + clusterID + "/" + exportID, Body: body}
	return o.Process(req)
}

type CephMgrRESTDeleteGaneshaExport struct {
	Url string //input
}

func (o *CephMgrRESTDeleteGaneshaExport) DoRequest(endpoint string, token string) error {
	url := o.Url
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + token

	res, err := CallRestAPI(endpoint+url, "DELETE", header, nil)
	if err == nil {
		switch res.StatusCode {
		case 204:
			return EOK
		case 200:
			return EOK
		default:
			return ECEPHMGRDELETEGANESHAEXPORT
		}
	}
	return ECEPHMGRDELETEGANESHAEXPORT
}

func (o *CephMgrRestful) DeleteGaneshaExport(clusterID string, exportID string) error {
	req := &CephMgrRESTDeleteGaneshaExport{Url: "/api/nfs-ganesha/export/" + clusterID + "/" + exportID}
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

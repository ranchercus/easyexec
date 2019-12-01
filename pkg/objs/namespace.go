package objs

import (
	"crypto/tls"
	"easyexec/pkg/common"
	"easyexec/pkg/login"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type V3Namespaces struct {
	Id        string `json:"Id"`
	Name      string `json:"name"`
	ProjectId string `json:"projectId"`
}
type V3NamespacesList struct {
	Data []V3Namespaces `json:"data"`
}

func GetNsByCookie() ([]string, error) {
	c, err := login.GetCookiePath()
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadFile(c)
	if err != nil {
		return nil, err
	}
	var cookie http.Cookie
	err = json.Unmarshal(b, &cookie)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   30 * time.Second,
			ExpectContinueTimeout: 10 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	req, err := http.NewRequest("GET", common.RANCHER_NAMESPACES, nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(&cookie)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("无法获取命名空间信息")
	}
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var list V3NamespacesList
	err = json.Unmarshal(r, &list)
	if err != nil {
		return nil, err
	}
	namespaces := make([]string, 0)
	for _, ns := range list.Data {
		namespaces = append(namespaces, ns.Name)
	}
	return namespaces, nil
}

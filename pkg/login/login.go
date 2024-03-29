package login

import (
	"crypto/tls"
	"easyexec/pkg/common"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var kubeconfig = ""

type Login struct {
	Username     string       `json:"username"`
	Password     string       `json:"password"`
	ResponseType string       `json:"responseType"`
	Ttl          int64        `json:"ttl"`
	client       *http.Client `json: "-"`
}

type UserInfo struct {
	Username string `json:"username"`
	Enabled  bool   `json:"Enabled"`
	Id       string `json:"Id"`
}

type UserCollection struct {
	Data []UserInfo `json:"data"`
}

type KubeConfig struct {
	BaseType string `json:"baseType"`
	Config   string `json:"config"`
}

func NewLogin(username, password string) *Login {
	l := &Login{
		Username:     username,
		Password:     password,
		ResponseType: "cookie",
		Ttl:          57600000,
		client: &http.Client{
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
		},
	}
	return l
}

func (l *Login) Login() error {
	if l.Username == "" || l.Password == "" {
		return errors.New("用户名或密码不能为空")
	}
	err := l.request()
	if err != nil {
		return err
	}

	return l.store()
}

func (l *Login) request() error {
	j, err := json.Marshal(l)
	if err != nil {
		return err
	}
	js := string(j)
	req, err := http.NewRequest("POST", common.RANCHER_URL, strings.NewReader(js))
	if err != nil {
		return err
	}
	resp, err := l.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		var cookie *http.Cookie
		for _, c := range resp.Cookies() {
			if c != nil && strings.HasPrefix(c.Value, "token") {
				cookie = c
				break
			}
		}
		if cookie == nil {
			return errors.New("认证失败，无法获取Cookie信息")
		}
		req, err = http.NewRequest("POST", common.RANCHER_kUBECONFIG_URL, nil)
		if err != nil {
			return err
		}
		req.AddCookie(cookie)

		defer l.storeCookie(cookie)

		resp, err = l.client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode == 200 {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			var jconfig KubeConfig

			err = json.Unmarshal(b, &jconfig)
			if err != nil {
				return err
			}
			kubeconfig = jconfig.Config

			return nil
		} else {
			return errors.New("认证失败，无法获取用户信息, " + resp.Status)
		}
	} else {
		return errors.New("认证失败, " + resp.Status)
	}
}

func (l *Login) store() error {
	path, err := common.GetConfigPath()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, []byte(kubeconfig), 0666)
	return err
}

func (l *Login) storeCookie(cookie *http.Cookie) error {
	path, err := common.GetCookiePath()
	if err != nil {
		return err
	}
	c, err := json.Marshal(cookie)
	if err != nil {
		return nil
	}
	err = ioutil.WriteFile(path, []byte(c), 0666)
	return err
}
package login

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var kubeconfig = `apiVersion: v1
kind: Config
clusters:
- name: "${RANCHER_CLUSTER_NAME}"
  cluster:
    server: "${RANCHER_KUBECONFIG_SERVER}"
    certificate-authority-data: "${RANCHER_CLUSTER_CA}"

users:
- name: "${RANCHER_CLUSTER_NAME}"
  user:
    token: "${TOKEN}"

contexts:
- name: "${RANCHER_CLUSTER_NAME}"
  context:
    user: "${RANCHER_CLUSTER_NAME}"
    cluster: "${RANCHER_CLUSTER_NAME}"

current-context: "${RANCHER_CLUSTER_NAME}"`

var (
	RANCHER_URL               = "https://rancher.i.fbank.com/v3-public/localProviders/local?action=login"
	RANCHER_USER_INFO_URL     = "https://rancher.i.fbank.com/v3/users?me=true&limit=-1&sort=name"
	RANCHER_KUBECONFIG_SERVER = "https://rancher.i.fbank.com/k8s/clusters/local"
	RANCHER_CLUSTER_NAME      = "pre-xpf"
	RANCHER_CLUSTER_CA        = `LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM3akNDQ\
      WRhZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFvTVJJd0VBWURWUVFLRXdsMGFHVXQKY\
      21GdVkyZ3hFakFRQmdOVkJBTVRDV05oZEhSc1pTMWpZVEFlRncweE9URXdNamt3TnpJNE5EUmFGd\
      zB5T1RFdwpNall3TnpJNE5EUmFNQ2d4RWpBUUJnTlZCQW9UQ1hSb1pTMXlZVzVqYURFU01CQUdBM\
      VVFQXhNSlkyRjBkR3hsCkxXTmhNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ\
      0tDQVFFQXhTU3RNbzZDcDBzdXhIS1IKL3hDNjVTQnZmSlQvUWd6NGRaaGpITlVsNndoK3ZaeXljU\
      lFJY2ZudExmWUlTYlhFY0JrM2swVDFoMnBXN0hXbwpNanVUQ0lxWFVhSjIxY0UvSWY1OEtpNHZoR\
      ldaakduOStjS3dkSFpHOWNTVnVTSFZsTTNTSFdOUWlhL2IvZVNPCmQxbVdWTUgyaXN0T3Q4Wjhma\
      WN5MGdHYnZGVnFWSHQ0a3ozUEhlVzJNU3REMDRDMDVZajFjb0ttMmtDLytwcVYKRkUrT21uYTBIa\
      kVnOWNiN3VMMURlMHJkL3ovK3VsZ20xQmRDd2R0WFhUWjM5bGpBbHdpVjRKMHFFTnBkY2hvawoyM\
      DZWNC9GQnNoVzZ1bURFMlVmd2k3OW92N2dDMS8vbUp3YXFvWVVtZWw1VXl4dWsyY2dvd1BVU0NlR\
      GtIMEJaCno4V0lpUUlEQVFBQm95TXdJVEFPQmdOVkhROEJBZjhFQkFNQ0FxUXdEd1lEVlIwVEFRS\
      C9CQVV3QXdFQi96QU4KQmdrcWhraUc5dzBCQVFzRkFBT0NBUUVBZ1FOZU93MnhDMURLQUhPd2o0b\
      WNmSlZUcnB4Y0hsUUxSNE5WbDJ1cwp1Qmx0TUordnl6amxKVG1FemR2MDV5YjdjZlJPY2srOCtXK\
      2R5L042N1poTkZSQm1XYUxGSDQrWEdqZjVNOWdQCjYyYUVpVWo5RnNtRWRoTlN1MTBIYkRIODJMc\
      2hpT00zNWNHL3JYckYxcVV4cXB4NkVIM1kzN0tObmxraW1hWkkKNlExeGx3VnZ2Nm1HZDNoRHRJb\
      2RqaTF5dEJhNmdIUC9oMk5XK1Q1dmpTY0k4THhObWxRcGNrbGNvZENwZUtJRApCMEo4RnJvaTBpa\
      EtPUXJKTkpZMVVOaEw5QzliRUxGcGZQYS9idWZ2T0ZjQU9odG56RDFMenk2Y1oyazNOL3ZYCmFQU\
      WRBeEJYVER4TmNHRk1vbG1aRC9PdWtCNGgrOWR6bmxaL0thMHdHT1BpREE9PQotLS0tLUVORCBDR\
      VJUSUZJQ0FURS0tLS0t`
)

type Login struct {
	Username     string       `json:"username"`
	Password     string       `json:"password"`
	ResponseType string       `json:"responseType"`
	Ttl          int64        `json:"ttl"`
	client       *http.Client `json: "-"`
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
	l.request()
	return nil
}

func (l *Login) request() error {
	j, err := json.Marshal(l)
	if err != nil {
		return err
	}
	js := string(j)
	fmt.Println(js)
	req, err := http.NewRequest("POST", RANCHER_URL, strings.NewReader(js))
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

		return nil
	} else {
		return errors.New("认证失败")
	}
}

func (l *Login) store() error {
	dir, err := os.UserHomeDir()
	if err != nil {
		dir = os.TempDir()
	}
	path, err := GetStoreFilePath(dir)
	if err != nil {
		return err
	}
	fmt.Println(path)
	return nil
}

package web

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"time"
)

var client *Client
var clientSingleton sync.Once

type Client struct {
	api    *http.Client
	config *Config
}

func GetClient() *Client {
	clientSingleton.Do(func() {
		config := GetConfig()
		jar, _ := cookiejar.New(nil)
		user := []*http.Cookie{
			{Name: "COOKIE_LOGIN_USER", Value: config.Auth},
		}
		jar.SetCookies(&url.URL{Scheme: "https", Host: "cloud.189.cn"}, user)
		jar.SetCookies(&url.URL{Scheme: "https", Host: "m.cloud.189.cn"}, user)
		jar.SetCookies(&url.URL{Scheme: "https", Host: "api.cloud.189.cn"}, user)
		jar.SetCookies(&url.URL{Scheme: "https", Host: "open.e.189.cn"}, []*http.Cookie{
			{Name: "SSON", Value: config.SSON},
		})
		client = &Client{
			config: config,
			api:    &http.Client{Jar: jar},
		}
	})
	return client
}

func (client *Client) refresh() {
	req, _ := http.NewRequest(http.MethodGet, "https://cloud.189.cn/api/portal/loginUrl.action", nil)

	resp, err := client.api.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	cookies := client.api.Jar.Cookies(resp.Request.URL)
	for _, cookie := range cookies {
		if cookie.Name == "COOKIE_LOGIN_USER" {
			config := client.config
			config.Auth = cookie.Value
			config.SessionKey = getUserBriefInfo(*config).SessionKey
			config.Save()
			return
		}
	}
	NewContentWithResp(resp).QrCode().Login()
}
func (client *Client) rsa() *rsa {
	config := client.config
	rsa := client.config.RSA
	now := time.Now().UnixMilli()
	if rsa.Expire > now {
		return &rsa
	}
	for rsa.Expire < now {
		req, _ := http.NewRequest(http.MethodGet, "https://cloud.189.cn/api/security/generateRsaKey.action", nil)
		req.Header.Add("accept", "application/json;charset=UTF-8")
		resp, _ := client.api.Do(req)
		json.NewDecoder(resp.Body).Decode(&rsa)
		resp.Body.Close()
		if resp.StatusCode != 200 || rsa.Expire == 0 {
			client.refresh()
		}
	}
	config.RSA = rsa
	config.Save()
	return &rsa
}
func (client *Client) initSesstion() {
	user := getUserBriefInfo(*client.config)
	if user.SessionKey == "" {
		client.refresh()
	} else {
		config.SessionKey = user.SessionKey
		config.Save()
	}
}
func (client *Client) sesstionKey() string {
	config := client.config
	key := config.SessionKey
	if key != "" {
		return key
	}
	user := getUserBriefInfo(*client.config)
	if user.SessionKey == "" {
		client.refresh()
	} else {
		config.SessionKey = user.SessionKey
		config.Save()
	}
	return config.SessionKey
}

type briefInfo struct {
	SessionKey  string `json:"sessionKey,omitempty"`
	UserAccount string `json:"userAccount,omitempty"`
}

func getUserBriefInfo(config Config) *briefInfo {
	u := fmt.Sprintf("https://cloud.189.cn/v2/getUserBriefInfo.action?noCache=%v", rand.Float64())
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	req.AddCookie(&http.Cookie{Name: "COOKIE_LOGIN_USER", Value: config.Auth})
	req.Header.Add("accept", "application/json;charset=UTF-8")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	var user briefInfo
	if strings.Index(resp.Header.Get("Content-Type"), "html") > 0 {
		return &user
	}
	json.NewDecoder(resp.Body).Decode(&user)
	return &user
}
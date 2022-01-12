// @File:  config.go
// @Time:  2021/12/27 2:08 PM
// @Author: ClassmateLin
// @Email: classmatelin.site@gmail.com
// @Site: https://www.classmatelin.top

package jd

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"scripts/config"
	"scripts/global"
	"strings"
)

var UserList []User

type User struct {
	Username   string            // 用户名
	CookieStr  string            // 英文分号分隔的cookie字符串
	CookieMap  map[string]string // cookie键值对
	PtPin      string            // jd手机网页版cookie中的pt_pin
	PtKey      string            // jd手机网页版cookie中的pt_key
	WsKey      string            // jd app 中的ws_key
	Sort       int               // 账号排序, 影响助力顺序
	TgBotToken string            // Tg 机器人配置
	TgUserId   string            // Tg 用户ID
	ServerJ    string            // server酱 sendKey 配置
	PushPlus   string            // push+ token 配置
	NotifyType string            // 通知类型: 默认: 1, 开启单独通知且推送到总通知, 2:开启单独通知但不推送到总通知
}

// GetPtKeyByWsKey 通过ws_key换取pt_key
//  @Description:
//  @param ptPin
//  @param wsKey
//  @return string
//
func GetPtKeyByWsKey(ptPin string, wsKey string) string {
	client := resty.New()

	headers := map[string]string{
		"user-agent":   global.GetJdUserAgent(),
		"content-type": "application/x-www-form-urlencoded; charset=UTF-8",
		"cookie":       fmt.Sprintf("pin=%s;wskey=%s", ptPin, wsKey),
	}
	getTokenUrl := "https://api.m.jd.com/client.action?functionId=genToken&clientVersion=10.1.2&build=89743&client=android&d_brand=&d_model=&osVersion=&screen=&partner=&oaid=&openudid=a27b83d3d1dba1cc&eid=&sdkVersion=30&lang=zh_CN&uuid=a27b83d3d1dba1cc&aid=a27b83d3d1dba1cc&area=19_1601_36953_50397&networkType=wifi&wifiBssid=&uts=&uemps=0-2&harmonyOs=0&st=1630413012009&sign=ca712dabc123eadd584ce93f63e00207&sv=121"
	body := "body=%7B%22to%22%3A%22https%253a%252f%252fplogin.m.jd.com%252fjd-mlogin%252fstatic%252fhtml%252fappjmp_blank.html%22%7D&"

	resp, err := client.R().
		SetHeaders(headers).
		SetBody(body).
		Post(getTokenUrl)

	if err != nil {
		return ""
	}

	token := gjson.Get(resp.String(), `tokenKey`).String()
	toUrl := gjson.Get(resp.String(), `url`).String()

	resp, err = client.
		SetRedirectPolicy(resty.RedirectPolicyFunc(func(request *http.Request, requests []*http.Request) error {
			return http.ErrUseLastResponse
		})).
		R().
		SetQueryParams(map[string]string{
			"tokenKey": token,
			"to":       "https://plogin.m.jd.com/jd-mlogin/static/html/appjmp_blank.html",
		}).
		Get(toUrl)

	if err != nil {
		return ""
	}

	for _, cookie := range resp.Cookies() {
		if strings.HasPrefix(cookie.String(), "pt_key") {
			ptKey := strings.Split(cookie.String(), ";")[0]
			data := strings.Split(strings.ReplaceAll(ptKey, " ", ""), "=")
			if len(data) == 2 {
				return data[1]
			} else {
				return ""
			}
		}
	}
	return ""
}

// init
// @description: 初始化JD账号列表
func init() {
	cookies := config.VP.GetStringSlice(`jd.cookies`)
	for index, cookie := range cookies { // sort 用于助力排序

		user := User{CookieMap: map[string]string{}}
		cookie = strings.ReplaceAll(cookie, " ", "")

		if strings.HasSuffix(cookie, ";") == false {
			log.Panicf("cookie: (%s), 字符串末尾缺少英文分号\n", cookie)
		}

		r, _ := regexp.Compile("([^=]+)=([^;]+);?\\s*")
		itemList := r.FindAllString(cookie, 10)

		for _, item := range itemList {

			temp := strings.Split(item, "=")

			switch temp[0] {
			case "pt_pin":
				str := strings.ReplaceAll(temp[1], ";", "")
				user.PtPin = str
				user.CookieStr += "pt_pin=" + str + ";"
				user.CookieMap[temp[0]] = str
			case "pt_key":
				str := strings.ReplaceAll(temp[1], ";", "")
				user.PtKey = str
				user.CookieStr += "pt_key=" + str + ";"
				user.CookieMap[temp[0]] = str
			case "remark":
				user.Username = strings.ReplaceAll(temp[1], ";", "")
			case "ws_key":
				user.WsKey = strings.ReplaceAll(temp[1], ";", "")
			case "tg_bot_token":
				user.TgBotToken = strings.ReplaceAll(temp[1], ";", "")
			case "tg_user_id":
				user.TgUserId = strings.ReplaceAll(temp[1], ";", "")
			case "server_j":
				user.ServerJ = strings.ReplaceAll(temp[1], ";", "")
			case "push_plus":
				user.PushPlus = strings.ReplaceAll(temp[1], ";", "")
			case "notify_type":
				user.NotifyType = strings.ReplaceAll(temp[1], ";", "")
			}
		}

		if user.PtKey == "" || user.WsKey == "" {
			if strings.HasSuffix(cookie, ";") == false {
				log.Panicf("cookie: (%s), 缺少pt_key或ws_key.\n", cookie)
			}
		}

		if user.WsKey != "" {
			user.PtKey = GetPtKeyByWsKey(user.PtPin, user.WsKey)
			user.CookieStr += "pt_key=" + user.PtKey + ";"
			user.CookieMap["pt_key"] = user.PtKey
		}

		if user.Username == "" {
			username, _ := url.QueryUnescape(user.PtPin)
			user.Username = username
		}

		if user.NotifyType == "" {
			user.NotifyType = "1"
		}

		user.Sort = index + 1
		UserList = append(UserList, user)
	}
}

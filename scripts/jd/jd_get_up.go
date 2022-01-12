// @File:  {jd_get_up}.go
// @Time:  {2022/1/12} 20:24
// @Author: ClassmateLin
// @Email: classmatelin.site@gmail.com
// @Site: https://www.classmatelin.top
// @Description:
// @Cron: 0 7 * * *

package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"scripts/config/jd"
	"scripts/constracts"
	"scripts/global"
	"scripts/structs"
)

// JdGetUp
// @Description:
type JdGetUp struct {
	structs.JdBase
	client *resty.Request
}

// New
// @description: 初始化
// @receiver : JdGetUp
// @param:  user
// @return: JdGetUp
func (JdGetUp) New(user jd.User) constracts.Jd {
	obj := JdGetUp{}
	obj.User = user
	obj.client = resty.New().R().SetHeaders(map[string]string{
		"Cookie":          obj.User.CookieStr,
		"Accept":          "*/*",
		"Connection":      "keep-alive",
		"Accept-Encoding": "gzip, deflate, br",
		"User-Agent":      global.GetJdUserAgent(),
		"Accept-Language": "zh-Hans-CN;q=1",
		"Host":            "api.m.jd.com",
	})
	return obj
}

// GetTitle
// @description: 返回脚本名称
// @return: interface{}
func (j JdGetUp) GetTitle() interface{} {
	return "早起福利"
}

// Exec
// @description: 执行任务
// @receiver : j
func (j JdGetUp) Exec() {
	url := "https://api.m.jd.com/client.action?functionId=morningGetBean&area=22_1930_50948_52157&body=%7B%22rnVersion%22%3A%224.7%22%2C%22fp%22%3A%22-1%22%2C%22eid%22%3A%22%22%2C%22shshshfp%22%3A%22-1%22%2C%22userAgent%22%3A%22-1%22%2C%22shshshfpa%22%3A%22-1%22%2C%22referUrl%22%3A%22-1%22%2C%22jda%22%3A%22-1%22%7D&build=167724&client=apple&clientVersion=10.0.6&d_brand=apple&d_model=iPhone12%2C8&eid=eidI1aaf8122bas5nupxDQcTRriWjt7Slv2RSJ7qcn6zrB99mPt31yO9nye2dnwJ/OW%2BUUpYt6I0VSTk7xGpxEHp6sM62VYWXroGATSgQLrUZ4QHLjQw&isBackground=N&joycious=60&lang=zh_CN&networkType=wifi&networklibtype=JDNetworkBaseAF&openudid=32280b23f8a48084816d8a6c577c6573c162c174&osVersion=14.4&partner=apple&rfs=0000&scope=01&screen=750%2A1334&sign=0c19e5962cea97520c1ef9a2e67dda60&st=1625354180413&sv=112&uemps=0-0&uts=0f31TVRjBSsqndu4/jgUPz6uymy50MQJSPYvHJMKdY9TUw/AQc1o/DLA/rOTDwEjG4Ar9s7IY4H6IPf3pAz7rkIVtEeW7XkXSOXGvEtHspPvqFlAueK%2B9dfB7ZbI91M9YYXBBk66bejZnH/W/xDy/aPsq2X3k4dUMOkS4j5GHKOGQO3o2U1rhx5O70ZrLaRm7Jy/DxCjm%2BdyfXX8v8rwKw%3D%3D&uuid=hjudwgohxzVu96krv/T6Hg%3D%3D&wifiBssid=c99b216a4acd3bce759e369eaeeafd7"

	resp, err := j.client.Get(url)

	if err != nil {
		j.Println("无法完成任务...")
	}

	code := gjson.Get(resp.String(), `data.awardResultFlag`).String()

	if code == "0" {
		j.Println(fmt.Sprintf("%s, %s", gjson.Get(resp.String(), `data.bizMsg`), gjson.Get(resp.String(), `data.beanNum`)))
	} else {
		j.Println(fmt.Sprintf("%s", gjson.Get(resp.String(), `data.bizMsg`)))
	}

}

func main() {
	structs.RunJd(JdGetUp{}, jd.UserList)
}

// @File:  jd_special.go
// @Time:  2022/1/19 20:28
// @Author: ClassmateLin
// @Email: classmatelin.site@gmail.com
// @Site: https://www.classmatelin.top
// @Description: 京东app下拉 特物
// @Cron: * */1 * * *

package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	url "net/url"
	"scripts/config/jd"
	"scripts/global"
	"scripts/structs"
	"time"
)

type JdSpecial struct {
	structs.JdBase
	client           *resty.Request
	activityId       int
	encryptProjectId string
}

func (j JdSpecial) New(user jd.User) JdSpecial {
	obj := JdSpecial{}
	obj.User = user
	obj.client = resty.New().R().SetHeaders(map[string]string{
		"Cookie":       user.CookieStr,
		"Host":         "api.m.jd.com",
		"Origin":       "https://pro.m.jd.com",
		"Accept":       "application/json, text/plain, */*",
		"User-Agent":   global.GetJdUserAgent(),
		"Referer":      "https://pro.m.jd.com/",
		"Content-Type": "application/json",
	})
	return obj
}

func (j JdSpecial) Request(fn string, body map[string]interface{}) string {
	content, err := json.Marshal(body)

	link := fmt.Sprintf("https://api.m.jd.com/api?functionId=%s&appid=ProductZ4Brand&client=wh5&t=%d&body=%s",
		fn, time.Now().UnixNano()/1e6, url.QueryEscape(string(content)))
	fmt.Println(link)
	resp, err := j.client.SetBody(body).Post(link)
	if err != nil {
		return ""
	}
	return resp.String()
}

// GetActivityInfo
// @description: 获取活动信息
// @receiver : j
// @return: bool
func (j JdSpecial) GetActivityInfo() bool {

	data := j.Request("superBrandSecondFloorMainPage", map[string]interface{}{"source": "secondfloor"})

	if code := gjson.Get(data, `code`); code.Int() != 0 {
		j.Println("获取任务列表失败...")
		return false
	}

	j.activityId = int(gjson.Get(data, `data.result.activityBaseInfo.activityId`).Int())
	j.encryptProjectId = gjson.Get(data, `data.result.activityBaseInfo.encryptProjectId`).String()
	fmt.Println(gjson.Get(data, `data.result.activityBaseInfo.activityId`))
	if j.activityId == 0 || j.encryptProjectId == "" {
		j.Println("获取活动ID失败...")
		return false
	}
	return true
}

// DoTask
// @description: 做任务
// @receiver : j
func (j JdSpecial) DoTask() {

	data := j.Request("superBrandTaskList", map[string]interface{}{
		"source":         "secondfloor",
		"activityId":     j.activityId,
		"assistInfoFlag": 1,
	})

	if code := gjson.Get(data, `code`); code.Int() != 0 {
		j.Println("获取任务列表失败...")
	}

	taskList := gjson.Get(data, `task.taskList`).Array()

	//if len(taskList) == 0 {
	//	j.Println("今日暂无任务...")
	//	return
	//}
	for _, task := range taskList {
		j.Println(task)
	}

}

func (j JdSpecial) Exec() {
	success := j.GetActivityInfo()
	if !success {
		return
	}
	j.DoTask()
}

func main() {
	j := JdSpecial{}.New(jd.UserList[0])
	j.Exec()
}

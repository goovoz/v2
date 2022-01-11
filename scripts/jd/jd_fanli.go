// @File:  jd_city_cash.go
// @Time:  2022/1/11 22:34
// @Author: ClassmateLin
// @Email: classmatelin.site@gmail.com
// @Site: https://www.classmatelin.top
// @Description: 京东饭粒
// @Cron: 0 4 * * *

package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"scripts/config/jd"
	"scripts/constracts"
	"scripts/global"
	"scripts/structs"
	"time"
)

type JdFanLi struct {
	structs.JdBase
	client *resty.Request
}

// New
// @description: 初始化
// @receiver : JdFanLi
// @param:  user
// @return: construct.Jd
func (JdFanLi) New(user jd.User) constracts.Jd {
	obj := JdFanLi{}
	obj.User = user
	obj.client = resty.New().R()
	return obj
}

// Get
// @description: get 请求
// @receiver : j
// @param:  fn
// @param:  body
// @return: string
func (j JdFanLi) Get(fn string) string {
	url := fmt.Sprintf("https://ifanli.m.jd.com/rebateapi/task/%s", fn)

	resp, err := j.client.SetHeaders(map[string]string{
		"cookie":       j.User.CookieStr,
		"user-agent":   global.GetJdUserAgent(),
		"host":         "ifanli.m.jd.com",
		"content-type": "application/json;charset=UTF-8",
		"origin":       "https://ifanli.m.jd.com",
		"referer":      "https://ifanli.m.jd.com/rebate/earnBean.html?paltform=null",
		"authority":    "ifanli.m.jd.com",
	}).Get(url)

	if err != nil {
		return ""
	}
	return resp.String()
}

// Post
// @description: post请求
// @receiver : j
// @param:  fn
// @param:  body
// @return: string
func (j JdFanLi) Post(fn string, body map[string]interface{}) string {

	url := fmt.Sprintf("https://ifanli.m.jd.com/rebateapi/task/%s", fn)

	resp, err := j.client.SetHeaders(map[string]string{
		"authority":       "ifanli.m.jd.com",
		"user-agent":      global.GetJdUserAgent(),
		"content-type":    "application/json;charset=UTF-8",
		"accept":          "application/json, text/plain, */*",
		"origin":          "https://ifanli.m.jd.com",
		"referer":         "https://ifanli.m.jd.com/rebate/earnBean.html?paltform=null",
		"accept-language": "zh-CN,zh;q=0.9",
		"cookie":          j.User.CookieStr,
		"Content-Type":    "application/json; charset=UTF-8",
	}).SetBody(body).Post(url)

	if err != nil {
		return ""
	}

	return resp.String()
}

// GetTitle
// @description: 返回脚本名称
// @receiver : j
// @return: interface{}
func (j JdFanLi) GetTitle() interface{} {
	return "京东饭粒"
}

// DoTasks
// @description: 做任务
// @receiver : j
func (j JdFanLi) DoTasks() {

	resp := j.Get("getTaskList")

	taskList := gjson.Get(resp, `content`).Array()

	for _, task := range taskList {
		statusName := gjson.Get(task.String(), `statusName`).String()
		taskName := gjson.Get(task.String(), `taskName`)
		status := gjson.Get(task.String(), `status`).Int()

		if statusName == "活动结束" || status == 2 {
			j.Println(fmt.Sprintf("跳过任务: 《%s》, 任务已结束或已完成...", taskName))
			continue
		}

		res := j.Post("saveTaskRecord", map[string]interface{}{
			"taskId":     gjson.Get(task.String(), `taskId`).Int(),
			"taskType":   gjson.Get(task.String(), `taskType`).Int(),
			"businessId": gjson.Get(task.String(), `businessId`).String(),
		})

		watchTime := gjson.Get(task.String(), `watchTime`).Int()

		j.Println(fmt.Sprintf("正在做任务: 《%s》, 等等%d秒...", taskName, watchTime))
		time.Sleep(time.Second * time.Duration(watchTime))
		time.Sleep(time.Millisecond * 500)

		res = j.Post("saveTaskRecord", map[string]interface{}{
			"taskId":     gjson.Get(task.String(), `taskId`).Int(),
			"taskType":   gjson.Get(task.String(), `taskType`).Int(),
			"businessId": gjson.Get(task.String(), `businessId`).String(),
			"uid":        gjson.Get(res, `content.uid`).String(),
			"tt":         gjson.Get(res, `content.tt`).Int(),
		})
		j.Println(fmt.Sprintf("完成任务:《%s》, 结果:%s", taskName, res))
		time.Sleep(time.Second * 2)
	}
}

// Exec
// @description: 任务入口
// @receiver : j
func (j JdFanLi) Exec() {
	resp := j.Get("getTaskFinishCount")

	finishCount := gjson.Get(resp, `content.finishCount`).Int()   // 已完成次数
	maxTaskCount := gjson.Get(resp, `content.maxTaskCount`).Int() // 最多完成次数

	if finishCount >= maxTaskCount {
		j.Println("今日任务已完成...")
		return
	}

	for i := finishCount; i < maxTaskCount; i++ {
		j.DoTasks()
	}
}

func main() {
	structs.RunJd(JdFanLi{}, jd.UserList)
}

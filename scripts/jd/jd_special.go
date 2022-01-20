// @File:  jd_special.go
// @Time:  2022/1/19 20:28
// @Author: ClassmateLin
// @Email: classmatelin.site@gmail.com
// @Site: https://www.classmatelin.top
// @Description: 京东app下拉 特物
// @Cron: 10 10 * * *

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

// JdSpecial
// @Description:
type JdSpecial struct {
	structs.JdBase
	client           *resty.Request
	activityId       int
	encryptProjectId string
}

// New
// @Description:
// @receiver j
// @param user
// @return JdSpecial
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

// Request
// @Description:
// @receiver j
// @param fn
// @param body
// @return string
func (j JdSpecial) Request(fn string, body map[string]interface{}) string {
	content, err := json.Marshal(body)

	link := fmt.Sprintf("https://api.m.jd.com/api?functionId=%s&appid=ProductZ4Brand&client=wh5&t=%d&body=%s",
		fn, time.Now().UnixNano()/1e6, url.QueryEscape(string(content)))
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

	if len(taskList) == 0 {
		j.Println("今日暂无任务...")
		return
	}
	for _, task := range taskList {

		completionCnt := gjson.Get(task.String(), `completionCnt`).Int()
		assignmentTimesLimit := gjson.Get(task.String(), `assignmentTimesLimit`).Int()
		data := gjson.Get(task.String(), `ext.shoppingActivity`).Array()

		if len(data) < 1 {
			continue
		}

		ext := data[0].String()
		title := gjson.Get(ext, `title`)
		encryptAssignmentId := gjson.Get(ext, `encryptAssignmentId`)
		advId := gjson.Get(ext, `advId`)

		if completionCnt >= assignmentTimesLimit {
			j.Println(fmt.Sprintf("%s, 已完成...", title))
			continue
		}

		res := j.Request("superBrandDoTask", map[string]interface{}{
			"source":              "secondfloor",
			"activityId":          j.activityId,
			"encryptProjectId":    j.encryptProjectId,
			"encryptAssignmentId": encryptAssignmentId,
			"assignmentType":      1,
			"itemId":              advId,
			"actionType":          0,
		})

		j.Println(fmt.Sprintf("%s任务结果:%s", title, gjson.Get(res, `data.bizMsg`)))

		time.Sleep(time.Second * 2)

	}

}

// Lottery
// @Description: 抽奖
// @receiver j
func (j JdSpecial) Lottery() {
	data := j.Request("superBrandSecondFloorMainPage", map[string]interface{}{"source": "secondfloor"})
	userStarNum := int(gjson.Get(data, `data.result.activityUserInfo.userStarNum`).Int())
	j.Println(fmt.Sprintf("可抽奖次数:%d", userStarNum))

	for i := 0; i < userStarNum; i++ {
		res := j.Request("superBrandTaskLottery", map[string]interface{}{"source": "secondfloor", "activityId": j.activityId})
		j.Println(fmt.Sprintf("第%d次抽奖结果:%s", i+1, res))
		time.Sleep(time.Second * 2)
	}
}

// Exec
// @Description:
// @receiver j
func (j JdSpecial) Exec() {
	success := j.GetActivityInfo()
	if !success {
		return
	}
	j.DoTask()
	j.Lottery()
}

func main() {
	j := JdSpecial{}.New(jd.UserList[0])
	j.Exec()
}

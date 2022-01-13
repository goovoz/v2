// @File:  jd_live.go
// @Time:  2022/1/13 11:57 AM
// @Author: ClassmateLin
// @Email: classmatelin.site@gmail.com
// @Site: https://www.classmatelin.top
// @Description: 京东直播
// @Cron: 10-20/5 12 * * *
package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"net/url"
	"scripts/config/jd"
	"scripts/constracts"
	"scripts/global"
	"scripts/structs"
	"strconv"
	"strings"
	"time"
)

// JdLive
// @Description: 京东直播
type JdLive struct {
	structs.JdBase
	client *resty.Request
}

// New
// @Description: 实例化
// @receiver JdLive
// @param user
// @return JdLive
func (JdLive) New(user jd.User) constracts.Jd {
	obj := JdLive{}
	obj.client = resty.New().R().SetHeaders(map[string]string{
		"Cookie":       user.CookieStr,
		"origin":       "https://h5.m.jd.com",
		"referer":      "https://h5.m.jd.com",
		"Content-Type": "application/x-www-form-urlencoded",
		"User-Agent":   global.GetJdUserAgent(),
	})
	obj.User = user
	return obj
}

// GetTitle
// @description: 返回脚本名称
// @return: interface{}
func (j JdLive) GetTitle() interface{} {
	return "京东直播"
}

// Get
// @Description: Get请求
// @receiver j
// @param functionId
// @param body
// @return string
func (j JdLive) Get(functionId string, body map[string]interface{}) string {

	bodyString, err := json.Marshal(body)
	if err != nil {
		return ""
	}

	link := fmt.Sprintf("https://api.m.jd.com/client.action?functionId=%s&appid=h5-live&body=%s&v=%d", functionId, url.QueryEscape(string(bodyString)), time.Now().UnixNano()/1e6)

	resp, err := j.client.Get(link)

	if err != nil {
		return ""
	}

	return resp.String()
}

// Post
// @description: Post请求
// @receiver : j
// @param:  functionId
// @param:  body
// @return: string
func (j JdLive) Post(functionId string, body string) string {

	link := fmt.Sprintf("https://api.m.jd.com/client.action?functionId=%s", functionId)

	resp, err := j.client.SetBody(body).Post(link)

	if err != nil {
		return ""
	}
	return resp.String()
}

// sign
// @description: 签到
// @receiver : j
// @param:  signDay
func (j JdLive) sign(signDay int64) {
	resp := j.Get("getChannelTaskRewardToM", map[string]interface{}{"type": "signTask", "itemId": "1"})
	if code := gjson.Get(resp, `subCode`).String(); code == "0" {
		j.Println(fmt.Sprintf("第%d天签到成功, 获得%d京豆...", signDay, gjson.Get(resp, `sum`).Int()))
	} else {
		j.Println(fmt.Sprintf("第%d天签到失败!", signDay))
	}
}

// doViewListTask
// @description: 观看直播
// @receiver : j
// @param:  seconds
// @param:  liveTitle
func (j JdLive) doViewListTask(seconds int, liveTitle string) {
	j.Post("liveChannelReportDataV912", "body=%7B%22liveId%22%3A%223008300%22%2C%22type%22%3A%22viewTask%22%2C%22authorId%22%3A%22644894%22%2C%22extra%22%3A%7B%22time%22%3A120%7D%7D&build=167408&client=apple&clientVersion=9.2.0&eid=eidIF3CF0112RTIyQTVGQTEtRDVCQy00Qg%3D%3D6HAJa9%2B/4Vedgo62xKQRoAb47%2Bpyu1EQs/6971aUvk0BQAsZLyQAYeid%2BPgbJ9BQoY1RFtkLCLP5OMqU&isBackground=N&joycious=194&openudid=53f4d9c70c1c81f1c8769d2fe2fef0190a3f60d2&osVersion=14.2&partner=TF&rfs=0000&scope=01&sign=90e14adc21c4bf31232a1ded5f4ba40e&st=1607561111999&sv=111&uts=0f31TVRjBSsxGLJHVBkddxFxBqY/8qFkrfEYLL0gkhB/JVGyEYIoD8r5rLvootZziQYAUyvIPogdJpesEuOMmvlisDx6AR2SEsfp381xPoggwvq8XaMYlOnHUV66TZiSfC%2BSgcLpB2v9cy/0Z41tT%2BuLheoEwBwDDYzANkZjncUI9PDCWpCg5/i0A14XfnsUTfQHbMqa3vwsY6QtsbNsgA%3D%3D&uuid=hjudwgohxzVu96krv/T6Hg%3D%3D")
}

// doSuperLiveTask
// @description: 超级推荐任务
// @receiver : j
func (j JdLive) doSuperStarTask(liveTitle string) {
	j.Post("liveChannelReportDataV912", "body=%7B%22liveId%22%3A%223019486%22%2C%22type%22%3A%22viewTask%22%2C%22authorId%22%3A%22681523%22%2C%22extra%22%3A%7B%22time%22%3A200%7D%7D&build=167454&client=apple&clientVersion=9.3.0&d_brand=apple&d_model=iPhone10%2C2&eid=eidIF3CF0112RTIyQTVGQTEtRDVCQy00Qg%3D%3D6HAJa9%2B/4Vedgo62xKQRoAb47%2Bpyu1EQs/6971aUvk0BQAsZLyQAYeid%2BPgbJ9BQoY1RFtkLCLP5OMqU&isBackground=N&joycious=194&lang=zh_CN&networkType=wifi&networklibtype=JDNetworkBaseAF&openudid=53f4d9c70c1c81f1c8769d2fe2fef0190a3f60d2&osVersion=14.2&partner=apple&rfs=0000&scope=01&screen=1242%2A2208&sign=68c0d87a5de62711bab9e6e796c08170&st=1607778652966&sv=100&uts=0f31TVRjBSsySvX9aqk89gHBMqz5E28EYCqc3cRu/4%2Bv0EzRuStHwMI1R5P9RqeizLow/pAquaX1v5IJQGVxUzSfExCFmfO0L7BEMvXnkeCZhKEsmSkbQm54W7ig8aRsmHiXp7YT/SOV7sEKxXauv59O/SAAFkr1egGgKev7Uj81nJRFDnNRSomlrOj2jQzH6iddCTSpydcSYRnDyDcodA%3D%3D&uuid=hjudwgohxzVu96krv/T6Hg%3D%3D")
}

// getTaskAward
// @description: 领取任务奖励
// @receiver : j
// @param:  taskType
// @param:  liveId
// @param:  taskTitle
func (j JdLive) getTaskAward(taskType string, taskTitle string) {

	body := "functionId=getChannelTaskRewardToM&appid=h5-live&body=%7B%22type%22%3A%22" + taskType + "%22%2C%22liveId%22%3A%222942545%22%7D&v=" + strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	link := "https://api.m.jd.com/api"

	resp, err := j.client.
		SetBody(body).Post(link)

	if err != nil {
		fmt.Println(err)
		j.Println(fmt.Sprintf("无法领取<%s>任务奖励", taskTitle))
		return
	}

	if code := gjson.Get(resp.String(), `subCode`).String(); code == "0" {
		j.Println(fmt.Sprintf("成功领取《%s》任务奖励, 获得%d京豆...", taskTitle, gjson.Get(resp.String(), `sum`).Int()))
	} else {
		j.Println(fmt.Sprintf("无法领取《%s》任务奖励, 原因:", taskTitle), gjson.Get(resp.String(), `msg`))
	}
}

// doShareTask
// @description: 分享任务
// @receiver : j
func (j JdLive) doShareTask() {
	j.Post("liveChannelReportDataV912", "body=%7B%22liveId%22%3A%222995233%22%2C%22type%22%3A%22shareTask%22%2C%22authorId%22%3A%22682780%22%2C%22extra%22%3A%7B%22num%22%3A1%7D%7D&build=167408&client=apple&clientVersion=9.2.0&eid=eidIF3CF0112RTIyQTVGQTEtRDVCQy00Qg%3D%3D6HAJa9%2B/4Vedgo62xKQRoAb47%2Bpyu1EQs/6971aUvk0BQAsZLyQAYeid%2BPgbJ9BQoY1RFtkLCLP5OMqU&isBackground=Y&joycious=194&lang=zh_CN&networkType=wifi&networklibtype=JDNetworkBaseAF&openudid=53f4d9c70c1c81f1c8769d2fe2fef0190a3f60d2&osVersion=14.2&partner=TF&rfs=0000&scope=01&screen=1242%2A2208&sign=457d557a0902f43cbdf9fb735d2bcd64&st=1607559819969&sv=110&uts=0f31TVRjBSsxGLJHVBkddxFxBqY/8qFkrfEYLL0gkhB/JVGyEYIoD8r5rLvootZziQYAUyvIPogdJpesEuOMmvlisDx6AR2SEsfp381xPoggwvq8XaMYlOnHUV66TZiSfC%2BSgcLpB2v9cy/0Z41tT%2BuLheoEwBwDDYzANkZjncUI9PDCWpCg5/i0A14XfnsUTfQHbMqa3vwsY6QtsbNsgA%3D%3D&uuid=hjudwgohxzVu96krv/T6Hg%3D%3D")
}

// Exec
// @description: 执行脚本
// @receiver : j
func (j JdLive) Exec() {
	resp := j.Get("liveChannelTaskListToM", map[string]interface{}{
		"timestamp": time.Now().UnixNano() / 1e6,
	})
	signDay := gjson.Get(resp, `data.sign.today`).Int()
	for _, item := range gjson.Get(resp, `data.sign.list`).Array() {
		m := item.String()
		day := gjson.Get(m, `day`)
		state := gjson.Get(m, `state`)
		if day.Int() != signDay {
			continue
		}
		if state.Int() == 1 {
			j.sign(signDay)
		} else {
			j.Println(fmt.Sprintf("第%d天已签到...", signDay))
		}
	}

	starLiveList := gjson.Get(resp, `data.starLiveList`).Array()

	for _, item := range starLiveList {
		task := item.String()
		title := gjson.Get(task, `liveTitle`).String()
		state := gjson.Get(task, `state`).Int()
		taskType := gjson.Get(task, `type`).String()

		if state == 3 {
			j.Println(fmt.Sprintf("今日已完成观看超级推荐:《%s》任务...", title))
			continue
		}
		if state == 2 {
			j.getTaskAward(taskType, title)
		}

		j.doSuperStarTask(title)
		time.Sleep(time.Second * 1)
		j.getTaskAward(taskType, title)
	}

	moreTaskList := gjson.Get(resp, `data.task`).Array()

	for _, item := range moreTaskList {
		task := item.String()
		taskType := gjson.Get(task, `type`).String()
		taskTitle := gjson.Get(task, `title`).String()
		state := gjson.Get(task, `state`).Int()

		if state == 3 {
			j.Println(fmt.Sprintf("今日已完成:《%s》任务...", taskTitle))
			continue
		}

		if state == 2 {
			j.getTaskAward(taskType, taskTitle)
			continue
		}

		switch taskType {
		case "commonViewTask":
			if exist := strings.Contains(taskTitle, "30"); exist {
				j.doViewListTask(30, taskTitle)
			} else {
				j.doViewListTask(60, taskTitle)
			}
		case "shareTask":
			j.doShareTask()
		}
		j.getTaskAward(taskType, taskTitle)
	}
}

func main() {
	structs.RunJd(JdLive{}, jd.UserList)
}

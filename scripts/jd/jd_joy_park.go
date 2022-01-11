// @File:  jd_joy_park.go
// @Time:  2022/1/7 3:17 PM
// @Author: ClassmateLin
// @Email: classmatelin.site@gmail.com
// @Site: https://www.classmatelin.top
// @Description: 汪汪乐园
// @Cron: 0 5 * * *
package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"scripts/config/jd"
	"scripts/constracts"
	"scripts/global"
	"scripts/structs"
	"time"
)

// JdJoyPark
// @Description: 汪汪公园
type JdJoyPark struct {
	structs.JdBase
	client *resty.Request
	linkId string
}

// New
// @description: 初始化JdJoyPark
// @param:  user
// @return: JdJoyPark
func (JdJoyPark) New(user jd.User) constracts.Jd {
	obj := JdJoyPark{}
	obj.User = user
	obj.client = resty.New().R().
		SetHeaders(map[string]string{
			"cookie":       obj.User.CookieStr,
			"user-agent":   global.GetJdUserAgent(),
			"referer":      "https://joypark.jd.com/",
			"origin":       "https://joypark.jd.com/",
			"host":         "api.m.jd.com",
			"content-Type": "application/x-www-form-urlencoded",
		})
	obj.linkId = "LsQNxL7iWDlXUs6cFl-AAg"
	return obj
}

func (j JdJoyPark) GetTitle() interface{} {
	return "汪汪公园"
}

// request
// @description: 请求数据
// @receiver : j
// @param:  fn
// @param:  body
// @return: string
func (j JdJoyPark) request(fn string, body map[string]interface{}) string {
	url := "https://api.m.jd.com/"
	body["linkId"] = j.linkId
	temp, _ := json.Marshal(body)
	data := fmt.Sprintf("functionId=%s&body=%s&_t=%d&appid=activities_platform", fn, string(temp), 1555600)
	resp, err := j.client.SetBody(data).Post(url)
	if err != nil {
		return ""
	}
	return resp.String()
}

// doTasks
// @description: 做任务
// @receiver : j
func (j JdJoyPark) doTasks() {
	taskData := j.request("apTaskList", map[string]interface{}{})
	if code := gjson.Get(taskData, "code").Int(); code != 0 {
		fmt.Println("请求失败")
		return
	}

	taskList := gjson.Get(taskData, `data`).Array()

	for _, task := range taskList {
		taskType := gjson.Get(task.String(), `taskType`).String()
		taskTitle := gjson.Get(task.String(), `taskShowTitle`)

		switch taskType {
		case "ORDER_MARK":
			j.Println(fmt.Sprintf("无法执行任务:《%s》, 请手动完成...", taskTitle))
		case "SIGN":
			j.sign(task.String())
		case "BROWSE_CHANNEL":
			j.browse(task.String())
		case "BROWSE_PRODUCT":
			j.browse(task.String())
		default:
			fmt.Printf("任务:%s 不支持...\n", taskTitle)
		}
	}
}

// sign
// @description: 签到任务
// @receiver : j
// @param:  task
func (j JdJoyPark) sign(task string) {

	if taskFinished := gjson.Get(task, `taskFinished`).Bool(); taskFinished {
		j.Println(fmt.Sprintf("今日已完成任务:《%s》...", gjson.Get(task, `taskTitle`)))
		return
	}
	resp := j.request("apDoTask", map[string]interface{}{
		"taskType": gjson.Get(task, `taskType`).String(),
		"taskId":   gjson.Get(task, `id`).String(),
	})

	if finished := gjson.Get(resp, `data.finished`).Bool(); !finished {
		j.Println(fmt.Sprintf("签到失败, %s...", gjson.Get(resp, `errMsg`)))
		return
	}

	time.Sleep(time.Second * 1)

	resp = j.request("apTaskDrawAward", map[string]interface{}{
		"taskType": gjson.Get(task, `taskType`).String(),
		"taskId":   gjson.Get(task, `id`).String(),
	})

	if code := gjson.Get(task, `code`).Int(); code == 0 {
		j.Println("成功领取签到任务奖励...")
	} else {
		j.Println("无法领取签到任务奖励, " + gjson.Get(resp, `errMsg`).String())
	}

}

// browseProduct
// @description: 浏览商品任务
// @receiver : j
// @param:  task
func (j JdJoyPark) browse(task string) {

	taskItemData := j.request("apTaskDetail", map[string]interface{}{
		"taskType": gjson.Get(task, `taskType`).String(),
		"taskId":   gjson.Get(task, `id`).String(),
		"channel":  4,
	})

	if code := gjson.Get(taskItemData, `code`).Int(); code != 0 {
		j.Println(fmt.Sprintf("无法获取任务:《%s》详情, %s...",
			gjson.Get(task, `taskTitle`), gjson.Get(taskItemData, `errMsg`)))
		return
	}

	if finished := gjson.Get(taskItemData, `data.status.finished`).Bool(); finished {
		j.Println(fmt.Sprintf("今日已完成任务: 《%s》...", gjson.Get(task, `taskTitle`)))
		return
	}

	taskItemList := gjson.Get(taskItemData, `data.taskItemList`).Array()

	times := gjson.Get(taskItemData, `data.status.finishNeed`).Int()

	if times > int64(len(taskItemList)) {
		times = int64(len(taskItemList))
	}

	for i := 0; int64(i) < times; i++ {

		item := taskItemList[i].String()

		resp := j.request("apDoTask", map[string]interface{}{
			"taskType": gjson.Get(task, `taskType`).String(),
			"taskId":   gjson.Get(task, `id`).String(),
			"channel":  4,
			"itemId":   gjson.Get(item, `itemId`).String(),
		})

		if code := gjson.Get(resp, `code`).Int(); code == 0 {
			j.Println(fmt.Sprintf("成功浏览ID:%s", gjson.Get(item, `itemName`)))
		}
		time.Sleep(time.Second * 2)

		resp = j.request("apTaskDrawAward", map[string]interface{}{
			"taskType": gjson.Get(task, `taskType`).String(),
			"taskId":   gjson.Get(task, `id`).String(),
		})

		if code := gjson.Get(resp, `code`).Int(); code == 0 {
			j.Println(fmt.Sprintf("成功领取浏览:%s的奖励...", gjson.Get(item, `itemName`)))
		}
		time.Sleep(time.Second * 2)
	}
}

func (j JdJoyPark) Exec() {
	j.doTasks()
}

func main() {
	structs.RunJd(JdJoyPark{}, jd.UserList)
}

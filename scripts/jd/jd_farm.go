// @File:  jd_farm.go
// @Time:  2022/2/14 3:24 PM
// @Author: ClassmateLin
// @Email: classmatelin.site@gmail.com
// @Site: https://www.classmatelin.top
// @Description:
// @Cron: 30 6,12,18 * * *
package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"scripts/config/jd"
	"scripts/constracts"
	"scripts/structs"
	"time"
)

type JdFarm struct {
	structs.JdBase
	client *resty.Request
}

// New
// @Description: 初始化
// @receiver j
// @param user
// @return JdFarm
func (j JdFarm) New(user jd.User) constracts.Jd {
	obj := JdFarm{}
	obj.User = user
	obj.client = resty.New().R()
	return obj
}

// request
// @Description: 请求数据
// @receiver j
// @param functionId
// @param body
// @return string
func (j JdFarm) request(functionId string, body map[string]interface{}) string {
	body["version"] = 13
	body["channel"] = 1
	temp, _ := json.Marshal(body)
	url := fmt.Sprintf("https://api.m.jd.com/client.action?functionId=%s&body=%s&appid=wh5", functionId, string(temp))
	resp, err := j.client.SetHeaders(map[string]string{
		"cookie": j.User.CookieStr,
	}).Post(url)
	if err != nil {
		return ""
	}
	return resp.String()
}

// initFarm
// @Description: 初始化农场
// @receiver j
// @return string
func (j JdFarm) initFarm() string {

	resp := j.request("initForFarm", map[string]interface{}{})

	farmInfo := gjson.Get(resp, `farmUserPro`).String()

	return farmInfo
}

// sign
// @Description: 签到
// @receiver j
func (j JdFarm) sign() {
	resp := j.request("signForFarm", map[string]interface{}{})

	if code := gjson.Get(resp, `code`).Int(); code == 0 {
		totalSinged := gjson.Get(resp, `signDay`).Int()
		j.Println(fmt.Sprintf("签到成功, 已连续签到%d天...", totalSinged))
	}
}

// doBrowseTasks
// @Description: 做浏览广告任务
// @receiver j
// @param taskList
func (j JdFarm) doBrowseTasks(taskList []gjson.Result) {

	for _, item := range taskList {

		task := item.String()

		taskName := gjson.Get(task, `mainTitle`).String()

		advertId := gjson.Get(task, `advertId`).String()

		taskRes := j.request("browseAdTaskForFarm", map[string]interface{}{
			"advertId": advertId,
			"type":     0,
		})

		if code := gjson.Get(taskRes, `code`).Int(); code == 0 || code == 7 {
			awardRes := j.request("browseAdTaskForFarm", map[string]interface{}{
				"advertId": advertId,
				"type":     1,
			})

			if code := gjson.Get(awardRes, `code`).Int(); code == 0 {
				amount := gjson.Get(awardRes, `amount`).Int()
				j.Println(fmt.Sprintf("成功领取任务:《%s》奖励, 获得水滴:%dg...", taskName, amount))
			}
		}
	}
}

// doThreeMealTask
// @Description: 三餐定时领水
// @receiver j
func (j JdFarm) doThreeMealTask() {
	resp := j.request("gotThreeMealForFarm", map[string]interface{}{})
	if code := gjson.Get(resp, `code`).Int(); code == 0 {
		amount := gjson.Get(resp, `amount`).Int()
		j.Println(fmt.Sprintf("完成定时领水任务, 获得%dg水滴...", amount))
	}
}

// doWaterRainTask
// @Description: 水滴雨任务
// @receiver j
func (j JdFarm) doWaterRainTask(task string) {
	if flag := gjson.Get(task, `f`).Bool(); flag {
		j.Println("今日已完成全部水滴雨任务...")
		return
	}

	lastTime := gjson.Get(task, `lastTime`).Int()
	nowTime := time.Now().Unix()
	times := gjson.Get(task, `winTimes`).Int() + 1

	if nowTime*1000 < lastTime+3*60*60*1000 {
		j.Println(fmt.Sprintf("第%d次水滴雨未到时间...", times))
		return
	}

	taskRes := j.request("waterRainForFarm", map[string]interface{}{})

	if code := gjson.Get(taskRes, `code`).Int(); code == 0 {
		amount := gjson.Get(taskRes, `addEnergy`).Int()
		j.Println(fmt.Sprintf("第%d次水滴雨任务完成, 获得%dg水滴...", times, amount))
	}

}

// doClickDuckTask
// @Description: 点鸭子任务
// @receiver j
func (j JdFarm) doClickDuckTask() {
	for i := 0; i < 10; i++ {
		resp := j.request("getFullCollectionReward", map[string]interface{}{"type": 2, "version": 14, "channel": 1,
			"babelChannel": 0})
		if code := gjson.Get(resp, `code`).Int(); code == 0 {
			title := gjson.Get(resp, `title`)
			j.Println(fmt.Sprintf("第%d次点击鸭子, %s...", i+1, title))
			time.Sleep(time.Second * 1)
		} else {
			j.Println("今日点击鸭子次数已达上限...")
			break
		}
	}
}

// waterTenTimes
// @Description: 浇水十次
// @receiver j
func (j JdFarm) waterTenTimes() {
	farmInfo := j.initFarm()
	totalEnergy := gjson.Get(farmInfo, `totalEnergy`).Int()

	cardData := j.request("myCardInfoForFarm", map[string]interface{}{})

	beanCard := gjson.Get(cardData, `beanCard`).Int()
	signCard := gjson.Get(cardData, `signCard`).Int()
	doubleCard := gjson.Get(cardData, `doubleCard`).Int()
	fastCard := gjson.Get(cardData, `fastCard`).Int()

	j.Println(fmt.Sprintf("背包道具:\n\t快速浇水卡:%d, 水滴翻倍卡:%d, 水滴换豆卡:%d, 加签卡: %d...", fastCard,
		doubleCard, beanCard, signCard))

	if totalEnergy > 100 && doubleCard > 0 {
		r := j.request("userMyCardForFarm", map[string]interface{}{
			"cardType": doubleCard,
		})
		j.Println("使用水滴翻倍卡: ", r)
	}

	if signCard > 0 {
		r := j.request("userMyCardForFarm", map[string]interface{}{
			"cardType": "signCard",
		})
		j.Println("使用加签卡: ", r)
	}

	taskData := j.request("taskInitForFarm", map[string]interface{}{})
	tenWaterTaskFinishedFlag := gjson.Get(taskData, `totalWaterTaskInit.f`).Bool()

	farmInfo = j.initFarm()
	treeTotalEnergy := gjson.Get(farmInfo, `treeTotalEnergy`).Int()
	totalEnergy = gjson.Get(farmInfo, `totalEnergy`).Int()
	keepEnergy := 80

	availableEnergy := int(totalEnergy) - keepEnergy

	if availableEnergy < 0 {
		availableEnergy = 0
	}

	if tenWaterTaskFinishedFlag && availableEnergy < 10 { // 十次浇水任务完成且保留水滴不足则退出浇水
		j.Println(fmt.Sprintf("总水滴%dg, 保留%dg, 可用水滴:%dg...", totalEnergy, keepEnergy, availableEnergy))
		j.Println("当前可用水滴小于10g, 无法进行浇水...")
		return
	}

	maxWaterTimes := int(availableEnergy / 10)

	if !tenWaterTaskFinishedFlag && maxWaterTimes < 10 { // 十次浇水任务未完成
		maxWaterTimes = 10
	}

	for i := 1; i <= maxWaterTimes; i++ {
		waterRes := j.request("waterGoodForFarm", map[string]interface{}{})
		if code := gjson.Get(waterRes, `code`).Int(); code == 0 {
			treeEnergy := gjson.Get(waterRes, `treeEnergy`).Int()
			j.Println(fmt.Sprintf("成功浇水10g, 距离水果成熟还需浇水%dg...", treeTotalEnergy-treeEnergy))

			if finished := gjson.Get(waterRes, `finished`).Bool(); finished {
				j.Println("水果已成熟, 无需继续浇水...")
				break
			}
		} else {
			j.Println("浇水失败, 不再继续浇水...")
			break
		}
	}
}

// doWaterFriendTask
// @Description: 为好友浇水
// @receiver j
func (j JdFarm) doWaterFriendTask() {
	friendData := j.request("friendListInitForFarm", map[string]interface{}{})

	friendList := gjson.Get(friendData, `friends`).Array()

	if len(friendList) < 1 {
		j.Println("暂无好友...")
		return
	}

	count := 0

	for _, item := range friendList {

		if count >= 2 {
			break
		}

		friend := item.String()
		if friendState := gjson.Get(friend, `friendState`).Int(); friendState != 1 {
			continue
		}

		friendShareCode := gjson.Get(friend, `shareCode`).String()

		waterRes := j.request("waterFriendForFarm", map[string]interface{}{
			"shareCode": friendShareCode,
		})

		code := gjson.Get(waterRes, `code`).Int()

		friendName := gjson.Get(friend, `nickName`).String()

		if code == 0 {
			j.Println(fmt.Sprintf("成功为好友: 《%s》浇水...", friendName))
			count += 1
		} else {
			j.Println(fmt.Sprintf("无法为好友:《%s》浇水, %s", friendName, waterRes))
			if code == 11 { // 水滴不足
				break
			}
		}

	}
}

// doDailyTasks
// @Description: 做每日任务
// @receiver j
func (j JdFarm) doDailyTasks() {
	resp := j.request("taskInitForFarm", map[string]interface{}{})

	if code := gjson.Get(resp, `code`).Int(); code != 0 {
		j.Println("查询每日任务失败...")
		return
	}

	todaySigned := gjson.Get(resp, `signInit.todaySigned`).Bool()
	totalSigned := gjson.Get(resp, `signInit.totalSigned`).Int()

	if !todaySigned {
		j.sign()
	} else {
		j.Println(fmt.Sprintf("今日已签到, 已连续签到%d天...", totalSigned))
	}

	browseTaskFlag := gjson.Get(resp, `gotBrowseTaskAdInit.f`).Bool()
	if !browseTaskFlag {
		browseTaskList := gjson.Get(resp, `gotBrowseTaskAdInit.userBrowseTaskAds`).Array()
		j.doBrowseTasks(browseTaskList)
	} else {
		j.Println("今日已完成所有浏览广告任务...")
	}

	threeMealTaskFlag := gjson.Get(resp, `gotThreeMealInit.f`).Bool()
	if !threeMealTaskFlag {
		j.doThreeMealTask()
	} else {
		j.Println("今日已完成定时领水任务...")
	}

	waterRainTask := gjson.Get(resp, `waterRainInit`).String()
	j.doWaterRainTask(waterRainTask)

	waterFriendTaskFlag := gjson.Get(resp, `waterFriendTaskInit.f`).Bool()
	if !waterFriendTaskFlag {
		j.doWaterFriendTask()
	} else {
		j.Println("今日已完成为2位好友浇水任务...")
	}
}

// getFirstWaterAward
// @Description: 领取首次浇水奖励
// @receiver j
func (j JdFarm) getFirstWaterAward() {

	taskData := j.request("taskInitForFarm", map[string]interface{}{})

	firstWaterFlag := gjson.Get(taskData, `firstWaterInit.f`).Bool()

	totalWaterTimes := gjson.Get(taskData, `firstWaterInit.totalWaterTimes`).Int()

	if firstWaterFlag {
		j.Println("今日已领取首次浇水任务奖励...")
		return
	}

	if totalWaterTimes == 0 {
		j.Println("今日未浇水，无法领取首次浇水任务奖励...")
		return
	}

	awardRes := j.request("firstWaterTaskForFarm", map[string]interface{}{})

	if code := gjson.Get(awardRes, `code`).Int(); code == 0 {
		amount := gjson.Get(awardRes, `amount`).Int()
		j.Println(fmt.Sprintf("成功领取首次浇水任务奖励, 获得%dg水滴...", amount))
	}
}

// getTenWaterAward
// @Description: 领取十次浇水奖励
// @receiver j
func (j JdFarm) getTenWaterAward() {

	taskData := j.request("taskInitForFarm", map[string]interface{}{})

	taskLimitTimes := gjson.Get(taskData, `totalWaterTaskInit.totalWaterTaskLimit`).Int()
	curTimes := gjson.Get(taskData, `totalWaterTaskInit.totalWaterTaskTimes`).Int()
	taskFinishedFlag := gjson.Get(taskData, `totalWaterTaskInit.f`).Bool()

	if taskFinishedFlag {
		j.Println("今日已领取十次浇水奖励...")
		return
	}

	if curTimes >= taskLimitTimes {
		awardRes := j.request("totalWaterTaskForFarm", map[string]interface{}{})

		if code := gjson.Get(awardRes, `code`).Int(); code == 0 {
			amount := gjson.Get(awardRes, `totalWaterTaskEnergy`).Int()
			j.Println(fmt.Sprintf("成功领取十次浇水任务奖励， 获得%dg水滴...", amount))
		}
	}

}

// getWaterFriendTaskAward
// @Description: 领取为两位好友浇水奖励
// @receiver j
func (j JdFarm) getWaterFriendTaskAward() {
	taskData := j.request("taskInitForFarm", map[string]interface{}{})
	waterFriendTaskFlag := gjson.Get(taskData, `waterFriendTaskInit`).Bool()

	if waterFriendTaskFlag {
		j.Println("今日已领取为两位好友浇水任务奖励...")
		return
	}

	awardRes := j.request("waterFriendGotAwardForFarm", map[string]interface{}{})

	if code := gjson.Get(awardRes, `code`).Int(); code == 0 {
		amount := gjson.Get(awardRes, `addWater`).Int()
		j.Println(fmt.Sprintf("成功领取未两位好友浇水任务奖励, 获得水滴%dg...", amount))
	} else {
		j.Println("无法领取为两位好友浇水任务奖励, ", awardRes)
	}
}

// clockIn
// @Description: 打卡领水
// @receiver j
func (j JdFarm) clockIn() {
	data := j.request("clockInInitForFarm", map[string]interface{}{})

	if code := gjson.Get(data, `code`).Int(); code != 0 {
		j.Println("获取打卡领水数据失败...")
		return
	}

	todaySigned := gjson.Get(data, `todaySigned`).Bool()
	signedDay := gjson.Get(data, `totalSigned`).Int()

	if !todaySigned {
		signRes := j.request("clockInForFarm", map[string]interface{}{
			"type": 1,
		})
		if code := gjson.Get(signRes, `code`).Int(); code == 0 {
			j.Println("今日打卡成功...")
			signedDay = gjson.Get(signRes, `signDay`).Int()
		} else {
			j.Println("今日打卡失败, ", signRes)
		}
	} else {
		j.Println(fmt.Sprintf("今日已打卡, 已连续打卡%d天...", signedDay))
	}

	if signedDay == 7 {
		awardRes := j.request("clockInForFarm", map[string]interface{}{
			"type": 2,
		})
		if code := gjson.Get(awardRes, `code`).Int(); code == 0 {
			amount := gjson.Get(awardRes, `amount`).Int()
			j.Println(fmt.Sprintf("成功领取惊喜礼包, 获得水滴%dg...", amount))
		} else {
			j.Println("领取惊喜礼包失败, ", awardRes)
		}
	}

	themes := gjson.Get(data, `themes`).Array()

	if len(themes) > 0 { // 限时关注领水滴
		for _, item := range themes {
			theme := item.String()
			name := gjson.Get(theme, `name`).String()
			hadGot := gjson.Get(theme, `hadGot`).Bool()
			id := gjson.Get(theme, `id`).String()
			if hadGot {
				continue
			}
			followRes := j.request("clockInFollowForFarm", map[string]interface{}{
				"id":   id,
				"type": "theme",
				"step": 1,
			})

			if code := gjson.Get(followRes, `code`).Int(); code != 0 {
				continue
			}

			awardRes := j.request("clockInFollowForFarm", map[string]interface{}{
				"id":   id,
				"type": "theme",
				"step": 2,
			})
			if code := gjson.Get(awardRes, `code`).Int(); code == 0 {
				amount := gjson.Get(awardRes, `amount`).Int()
				j.Println(fmt.Sprintf("成功关注《%s》, 获得水滴%dg...", name, amount))
			}
		}
	}

}

// getWater
// @Description: 领取农场赠送水滴
// @receiver j
func (j JdFarm) getWater() {
	res := j.request("gotWaterGoalTaskForFarm", map[string]interface{}{
		"type": 3, "version": 14, "channel": 1, "babelChannel": 0,
	})
	if code := gjson.Get(res, `code`).Int(); code == 0 {
		j.Println("成功领取农场赠送水滴, ", res)
	}
}

// getHelpAward
// @Description: 领取助力奖励
// @receiver j
func (j JdFarm) getHelpAward() {
	for i := 0; i < 5; i++ {
		res := j.request("receiveStageEnergy", map[string]interface{}{})
		if code := gjson.Get(res, `code`).Int(); code == 0 {
			amount := gjson.Get(res, `amount`)
			j.Println(fmt.Sprintf("成功领取好友助力奖励, 获得水滴%dg...", amount))
			time.Sleep(time.Second * 1)
		} else {
			break
		}
	}
}

// getStageAward
// @Description: 领取阶段性浇水奖励
// @receiver j
func (j JdFarm) getStageAward() {
	waterRes := j.request("waterGoodForFarm", map[string]interface{}{})

	waterStatus := gjson.Get(waterRes, `waterStatus`).Int()
	treeEnergy := gjson.Get(waterRes, `treeEnergy`).Int()

	if waterStatus == 0 && treeEnergy == 10 {
		awardRes := j.request("gotStageAwardForFarm", map[string]interface{}{
			"type": 1,
		})
		j.Println("领取浇水第一阶段奖励: ", awardRes)
	} else if waterStatus == 1 {
		awardRes := j.request("gotStageAwardForFarm", map[string]interface{}{
			"type": 2,
		})
		j.Println("领取浇水第二阶段奖励: ", awardRes)
	} else if waterStatus == 2 {
		awardRes := j.request("gotStageAwardForFarm", map[string]interface{}{
			"type": 3,
		})
		j.Println("领取浇水第三阶段奖励: ", awardRes)
	}
}

// doDdParkTask
// @Description: 东东乐园任务
// @receiver j
func (j JdFarm) doDdParkTask() {
	data := j.request("ddnc_farmpark_Init", map[string]interface{}{
		"version": "1", "channel": 1,
	})

	if code := gjson.Get(data, `code`).Int(); code != 0 {
		j.Println("无法获取东东乐园任务列表...")
		return
	}

	itemList := gjson.Get(data, `buildings`).Array()

	if len(itemList) < 1 {
		return
	}

	idx := -1

	for _, item := range itemList {
		idx += 1
		task := gjson.Get(item.String(), `topResource.task`).String()
		if task == "" {
			continue
		}
		name := gjson.Get(item.String(), `name`).String()
		if status := gjson.Get(task, `status`).Int(); status != 1 {
			j.Println(fmt.Sprintf("今日已完成东东乐园浏览:《%s》任务...", name))
			continue
		}
		advertId := gjson.Get(task, `advertId`).String()

		taskRes := j.request("ddnc_farmpark_markBrowser", map[string]interface{}{
			"version":  "1",
			"channel":  1,
			"advertId": advertId,
		})
		if code := gjson.Get(taskRes, `code`).Int(); code != 0 {
			continue
		}
		time.Sleep(time.Second * 1)

		awardRes := j.request("ddnc_farmpark_browseAward", map[string]interface{}{
			"version":  "1",
			"channel":  1,
			"advertId": advertId,
			"index":    idx,
			"type":     1,
		})

		if code := gjson.Get(awardRes, `code`).Int(); code == 0 {
			amount := gjson.Get(awardRes, `result.waterEnergy`).Int()
			j.Println(fmt.Sprintf("完成《%s》任务, 获得水滴%dg...", name, amount))
		} else {
			message := gjson.Get(awardRes, `message`).String()
			j.Println(fmt.Sprintf("无法完成《%s》任务, 原因: %s", name, message))
		}
		time.Sleep(time.Second * 1)
	}
}

// turntable
// @Description: 天天抽奖
// @receiver j
func (j JdFarm) turntable() {
	data := j.request("initForTurntableFarm", map[string]interface{}{})

	if code := gjson.Get(data, `code`).Int(); code != 0 {
		j.Println("无法获取天天抽奖数据...")
		return
	}

	timingGotStatus := gjson.Get(data, `timingGotStatus`).Bool()
	sysTime := int(gjson.Get(data, `sysTime`).Int())
	timingLastSysTime := int(gjson.Get(data, `timingLastSysTime`).Int())
	timingIntervalHours := int(gjson.Get(data, `timingIntervalHours`).Int())

	if !timingGotStatus && sysTime > timingLastSysTime+60*60*timingIntervalHours*1000 {
		awardRes := j.request("timingAwardForTurntableFarm", map[string]interface{}{})
		j.Println("领取天天抽奖定时奖励, ", awardRes)
	}

	taskList := gjson.Get(data, `turntableBrowserAds`).Array()

	for _, item := range taskList {
		task := item.String()
		title := gjson.Get(task, `main`)
		if status := gjson.Get(task, `status`).Bool(); status {
			j.Println(fmt.Sprintf("天天抽奖任务: 《%s》已完成...", title))
			continue
		}
		adId := gjson.Get(task, `adId`).String()
		taskRes := j.request("browserForTurntableFarm", map[string]interface{}{
			"type": 1,
			"adId": adId,
		})
		if code := gjson.Get(taskRes, `code`).Int(); code != 0 {
			continue
		}
		time.Sleep(time.Second * 1)
		awardRes := j.request("browserForTurntableFarm", map[string]interface{}{
			"type": 2,
			"adId": adId,
		})
		j.Println(fmt.Sprintf("领取天天抽奖任务:《%s》奖励, ", title), awardRes)
		time.Sleep(time.Second * 1)
	}

	data = j.request("initForTurntableFarm", map[string]interface{}{})

	lotteryTimes := int(gjson.Get(data, `remainLotteryTimes`).Int())

	if lotteryTimes < 1 {
		j.Println("天天抽奖次数已用完, 无法抽奖...")
		return
	}

	for i := 1; i <= lotteryTimes; i++ {
		res := j.request("lotteryForTurntableFarm", map[string]interface{}{})
		j.Println(fmt.Sprintf("天天抽奖->第%d次抽奖结果: ", i), res)
		time.Sleep(time.Second * 1)
	}
}

// Exec
// @Description: 程序入口
// @receiver j
func (j JdFarm) Exec() {
	farmInfo := j.initFarm()
	if farmInfo == "" {
		j.Println("无法初始化农场, 退出程序...")
		return
	}

	awardName := gjson.Get(farmInfo, `name`)
	shareCode := gjson.Get(farmInfo, `shareCode`).String()
	j.Println(fmt.Sprintf("奖品名称-> 《%s》, \n\t\t\t助力码-> 《%s》", awardName, shareCode))

	j.doDailyTasks()

	j.clockIn()

	j.doClickDuckTask()

	j.waterTenTimes()

	j.getFirstWaterAward()

	j.getTenWaterAward()

	j.getWater()

	j.doDdParkTask()

	j.turntable()

	j.getHelpAward()

	j.getStageAward()

}

// GetTitle
// @Description: 返回脚本名称
// @receiver j
func (j JdFarm) GetTitle() interface{} {
	return "东东农场"
}

func main() {
	structs.RunJd(JdFarm{}, jd.UserList)
}

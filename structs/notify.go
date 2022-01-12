// @File:  notify.go
// @Time:  2022/1/9 14:13
// @Author: ClassmateLin
// @Email: classmatelin.site@gmail.com
// @Site: https://www.classmatelin.top
// @Description:

package structs

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/russross/blackfriday"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"scripts/config/jd"
	"scripts/constracts"
	"scripts/global"
	"sync"
)

// TgNotify
// @Description: Tg消息推送
type TgNotify struct {
	BotToken string
	UserId   string
}

// New
// @description: 初始化
// @receiver : tg
// @param:  botToken
// @param:  userId
// @return: TgNotify
func (tg TgNotify) New(botToken string, userId string) TgNotify {
	tg.UserId = userId
	tg.BotToken = botToken
	return tg
}

// Title
// @Description: 返回通知类型
// @receiver tg
// @return string
//
func (tg TgNotify) Title() string {
	return "telegram"
}

// Send
// @description: 发送消息
// @receiver : tg
// @param:  title
// @param:  message
func (tg TgNotify) Send(title string, message string) bool {

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", tg.BotToken)

	client := resty.New()

	resp, err := client.R().
		SetHeaders(map[string]string{
			"content-type": "application/json",
			"user-agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36",
		}).SetBody(map[string]string{
		"chat_id": tg.UserId,
		"text":    title + "\n\n" + message,
	}).Post(url)

	if err != nil {
		return false
	}

	if ok := gjson.Get(resp.String(), `ok`).Bool(); ok {
		return true
	} else {
		global.Log.Error(resp.String())
		return false
	}
}

// PushPlusNotify
// @Description: push plus 消息推送, url: https://pushplus.hxtrip.com
type PushPlusNotify struct {
	Token string
}

// New
// @description: 初始化
// @receiver : p
// @param:  token
// @return: PushPlusNotify
func (p PushPlusNotify) New(token string) PushPlusNotify {
	p.Token = token
	return p
}

// Title
// @Description: 返回通知类型
// @receiver tg
// @return string
//
func (p PushPlusNotify) Title() string {
	return "push+"
}

// Send
// @description: 推送消息
// @receiver : p
// @param:  title
// @param:  message
func (p PushPlusNotify) Send(title string, message string) bool {
	url := "https://pushplus.hxtrip.com/send"

	client := resty.New()

	resp, err := client.R().
		SetHeaders(map[string]string{
			"content-type": "application/json",
			"user-agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36",
		}).SetBody(map[string]string{
		"token":    p.Token,
		"title":    title,
		"content":  string(blackfriday.MarkdownCommon([]byte(message))),
		"template": "html",
	}).Post(url)

	if err != nil {
		return false
	}
	if code := gjson.Get(resp.String(), `code`).Int(); code == 200 {
		return true
	} else {
		global.Log.Error(resp.String())
		return false
	}
}

// ServerJNotify
// @Description: server酱消息推送
type ServerJNotify struct {
	SendKey string
}

// New
// @description: 初始化
// @receiver : s
// @param:  sendKey
// @return: ServerJNotify
func (s ServerJNotify) New(sendKey string) ServerJNotify {
	s.SendKey = sendKey
	return s
}

// Title
// @Description: 返回通知类型
// @receiver tg
// @return string
//
func (s ServerJNotify) Title() string {
	return "telegram"
}

// Send
// @description: 发送消息
// @receiver : s
// @param:  title
// @param:  message
func (s ServerJNotify) Send(title string, message string) bool {
	url := fmt.Sprintf("https://sc.ftqq.com/%s.send", s.SendKey)

	client := resty.New()

	resp, err := client.R().
		SetHeaders(map[string]string{
			"content-type": "application/json",
			"user-agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36",
		}).SetQueryParams(map[string]string{
		"text": title,
		"desp": message,
	}).Post(url)

	if err != nil {
		return false
	}

	if code := gjson.Get(resp.String(), `code`).Int(); code == 0 {
		return true
	} else {
		global.Log.Error(resp.String())
		return false
	}
}

// SingleNotify
//  @Description:
//  @param user
//  @param title
//  @param message
//
func SingleNotify(user jd.User, title string, message string) {

	if user.TgUserId != "" && user.TgBotToken != "" {
		notify := TgNotify{}.New(user.TgBotToken, user.TgUserId)
		isSuccess := notify.Send(title, message)
		if isSuccess {
			fmt.Printf("%s, %s消息推送成功...\n", user.Username, notify.Title())
		} else {
			fmt.Printf("%s, %s消息推送失败...\n", user.Username, notify.Title())
		}
	}

	if user.ServerJ != "" {
		notify := ServerJNotify{}.New(user.ServerJ)
		isSuccess := notify.Send(title, message)
		if isSuccess {
			fmt.Printf("%s, %s消息推送成功...\n", user.Username, notify.Title())
		} else {
			fmt.Printf("%s, %s消息推送失败...\n", user.Username, notify.Title())
		}
	}

	if user.PushPlus != "" {
		notify := PushPlusNotify{}.New(user.PushPlus)
		isSuccess := notify.Send(title, message)
		if isSuccess {
			fmt.Printf("%s, %s消息推送成功...\n", user.Username, notify.Title())
		} else {
			fmt.Printf("%s, %s消息推送失败...\n", user.Username, notify.Title())
		}
	}
}

// Notify
// @description: 消息推送
// @param:  vp
// @param:  title
// @param:  message
func Notify(vp *viper.Viper, title string, message string) {
	var notifyList []constracts.Notify

	notifyConfig := vp.GetStringMap(`notify`)

	for key, val := range notifyConfig {
		switch key {
		case "tg":
			temp := val.(map[string]interface{})
			botToken := temp["bot_token"].(string)
			userId := temp["user_id"].(string)
			if botToken != "" && userId != "" {
				notifyList = append(notifyList, TgNotify{}.New(botToken, userId))
			} else {
				fmt.Println("由于未配置notify.tg.bot_key或notify.tg.user_id, 无法推送消息至tg...")
			}
		case "push_plus":
			if val.(string) != "" {
				notifyList = append(notifyList, PushPlusNotify{}.New(val.(string)))
			} else {
				fmt.Println("由于未配置notify.push_plus, 无法推送消息至push+...")
			}
		case "server_j":
			if val.(string) != "" {
				notifyList = append(notifyList, ServerJNotify{}.New(val.(string)))
			} else {
				fmt.Println("由于未配置notify.server_j, 无法推送消息至server酱...")
			}
		}

	}

	if notifyList == nil {
		return
	}

	var wg sync.WaitGroup
	for _, notify := range notifyList {
		wg.Add(1)
		go func(notify constracts.Notify, title string, message string, wg *sync.WaitGroup) {
			defer wg.Done()
			isSuccess := notify.Send(title, message)
			if isSuccess {
				fmt.Printf("%s消息推送成功...\n", notify.Title())
			} else {
				fmt.Printf("%s消息推送失败...\n", notify.Title())
			}
		}(notify, title, message, &wg)
	}
	wg.Wait()
}

// @File:  jd.go
// @Time:  2022/1/6 4:51 PM
// @Author: ClassmateLin
// @Email: classmatelin.site@gmail.com
// @Site: https://www.classmatelin.top
// @Description:
// @Cron: * */1 * * *

package structs

import (
	"fmt"
	"reflect"
	"scripts/config/jd"
	"scripts/constracts"
	"strings"
	"sync"
)

// JdBase
// @Description: JD活动基类
type JdBase struct {
	User jd.User
}

// Exec
// @description: 执行任务
// @receiver : JdBase
// @param:  wg
func (JdBase) Exec() {
}

// Help
// @description: 助力好友
// @receiver : JdBase
// @param:  wg
func (JdBase) Help() {
}

func (jd JdBase) GetTitle() interface{} {
	return nil
}

func (jd JdBase) GetHelpKey() interface{} {
	return nil
}

// Println
// @description: 输出内容
// @receiver : jd
// @param:  args
func (jd JdBase) Println(args ...interface{}) {
	fmt.Printf("账号%d(%s):", jd.User.Sort, jd.User.Username)
	fmt.Println(args...)
}

// RunJd
// @description: 执行JD脚本
// @param:  cj
// @param:  userList
func RunJd(cj constracts.Jd, userList []jd.User) {
	var wg sync.WaitGroup

	scriptName := cj.GetTitle()
	if scriptName == nil {
		scriptName = strings.Split(reflect.TypeOf(cj).String(), ".")[1]
	}

	fmt.Println(fmt.Sprintf("*********************%s:开始任务*********************", scriptName))
	for _, user := range userList {
		wg.Add(1)
		go func(wg *sync.WaitGroup, user jd.User) {
			defer wg.Done()
			app := cj.New(user)
			app.Exec()
		}(&wg, user)
	}
	wg.Wait()
	fmt.Println(fmt.Sprintf("*********************%s:完成任务*********************", scriptName))

	helpKey := cj.GetHelpKey()
	if helpKey == nil { // 不需要助力
		return
	}

	fmt.Println(fmt.Sprintf("*********************%s:开始助力*********************", scriptName))
	for _, user := range userList {
		wg.Add(1)
		go func(wg *sync.WaitGroup, user jd.User) {
			defer wg.Done()
			app := cj.New(user)
			app.Help()
		}(&wg, user)
	}
	wg.Wait()
	fmt.Println(fmt.Sprintf("*********************%s:完成助力*********************", scriptName))
}

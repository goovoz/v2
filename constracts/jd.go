// @File:  jd.go
// @Time:  2022/1/6 4:44 PM
// @Author: ClassmateLin
// @Email: classmatelin.site@gmail.com
// @Site: https://www.classmatelin.top
// @Description:
// @Cron: * */1 * * *

package constracts

import (
	"scripts/config/jd"
)

// Jd
// @Description: JD活动脚本interface
type Jd interface {
	Exec()                   // 执行任务
	Help()                   // 助力好友
	New(user jd.User) Jd     // 初始化方法
	GetTitle() interface{}   // 脚本名称
	GetHelpKey() interface{} // 助力码关键字
}

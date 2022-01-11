// @File:  notify.go
// @Time:  2022/1/9 14:16
// @Author: ClassmateLin
// @Email: classmatelin.site@gmail.com
// @Site: https://www.classmatelin.top
// @Description:

package constracts

type Notify interface {
	Send(title string, message string)
}

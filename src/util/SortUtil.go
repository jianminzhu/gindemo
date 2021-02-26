package util

import (
	"sort"
	"time"
)

//通用排序
//结构体排序，必须重写数组Len() Swap() Less()函数
type body_wrapper struct {
	Bodys []interface{}
	by    func(p, q *interface{}) bool //内部Less()函数会用到
}
type SortBodyBy func(p, q *interface{}) bool //定义一个函数类型

//数组长度Len()
func (acw body_wrapper) Len() int {
	return len(acw.Bodys)
}

//元素交换
func (acw body_wrapper) Swap(i, j int) {
	acw.Bodys[i], acw.Bodys[j] = acw.Bodys[j], acw.Bodys[i]
}

func (acw body_wrapper) Less(i, j int) bool {
	return acw.by(&acw.Bodys[i], &acw.Bodys[j])
}

func SortBody(bodys []interface{}, by SortBodyBy) {
	sort.Sort(body_wrapper{bodys, by})
}

type ByDate []map[string]interface{}

func ParseYMDDate(dateStr string) time.Time {
	t, _ := time.Parse("2006-01-02", dateStr)
	return t
}
func (a ByDate) Len() int      { return len(a) }
func (a ByDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool {
	return ParseYMDDate(a[i]["date"].(string)).Before(ParseYMDDate(a[i]["date"].(string)))
}

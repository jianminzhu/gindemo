package models

import (
	"github.com/astaxie/beego/orm"
)

func GetArticle(artileTypeName string) ([]Article, error) {
	items := []Article{}
	_, err := orm.NewOrm().QueryTable(Article{}).Filter("type", artileTypeName).All(&items)
	if err != nil {
		return []Article{}, nil
	}
	return items, err
}

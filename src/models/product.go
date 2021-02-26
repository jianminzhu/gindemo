package models

import (
	"github.com/astaxie/beego/orm"
)

func GetProducts() ([]Goods, error) {
	items := []Goods{}
	_, err := orm.NewOrm().QueryTable(Goods{}).All(&items)
	if err != nil {
		return []Goods{}, nil
	}
	return items, err
}

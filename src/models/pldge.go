package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/prometheus/common/log"
)

type PladgeModel struct {
}

func (p *PladgeModel) PlageList(userIds ...int64) []UserPledge {
	table := orm.NewOrm().QueryTable(UserPledge{})
	if len(userIds) > 0 {
		table = table.Filter("user_id__in", userIds)
	}
	pledgeArr := []UserPledge{}
	if _, err := table.All(&pledgeArr); err == nil {
		return pledgeArr
	}
	log.Error("__queryUSerPledge__error_", ToJson(map[string]interface{}{"userIds": userIds}))
	return []UserPledge{}
}

func (p *PladgeModel) UserPlageList(userId int64) []UserPledge {
	return p.PlageList(userId)
}

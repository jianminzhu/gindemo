package service

import (
	"errors"
	"gin/src/models"
	"github.com/astaxie/beego/orm"
)

func GetUserForShow(id int) (models.User, error) {
	orm.NewOrm()
	t := models.T("user").Filter("id", id)
	count, _ := t.Count()
	var users models.User
	if count > 0 {
		t.One(&users, "id", "username", "mobile", "invitation_code", "references_user_id", "reg_time", "status", "per_t_config_special")
		return users, nil
	} else {
		return users, errors.New("未找到用户")
	}
}

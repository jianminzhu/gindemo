package system

import (
	c "gin/src/controller"
	"gin/src/models"
	"gin/src/service"
	"strconv"
)

func Health(c *c.Base) {
	count, _ := models.T("system_user").Count()
	c.Ok(map[string]interface{}{
		"db": count > 0,
	}, "")
}
func User(c *c.Base) {
	//1 参数校验
	//2 读取数据
	//3 读取数据失败
	//4 读取数据成功

	//1
	idStr := c.Gin.DefaultQuery("id", "")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.FailDefObj("需要指定正确的userId")
		return
	}
	//2
	var user models.User
	user, err = service.GetUserForShow(id)
	//3
	if err != nil {
		c.FailDefObj(err.Error())
		return
	}
	//4
	trace, _ := c.Gin.Get("trace")
	c.Ok(map[string]interface{}{
		"trace":             trace,
		"Id":                user.Id,
		"InvitationCode":    user.InvitationCode,
		"ReferencesUserId":  user.ReferencesUserId,
		"Mobile":            user.Mobile,
		"Status":            user.Status,
		"PerTConfigSpecial": user.PerTConfigSpecial,
		"RegTime":           user.RegTime,
		"Username":          user.Username,
	}, "")
}

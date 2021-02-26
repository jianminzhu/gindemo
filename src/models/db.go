package models

import (
	"fmt"
	"gin/src/util"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/ini.v1"
	"io"
	"time"
)

func init() {
	cfg, _ := ini.Load("app.conf")
	mysqlConfig := cfg.Section("mysql")
	mysqlHost := mysqlConfig.Key("host").String()
	mysqlPort := mysqlConfig.Key("port").String()
	mysqlUser := mysqlConfig.Key("user").String()
	mysqlEncodePassword := mysqlConfig.Key("encodePassword").String()
	mysqlDbname := mysqlConfig.Key("dbname").String()
	mysqlpassword := util.DecodePassword(mysqlEncodePassword, util.DefaultSalt)
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=Local", mysqlUser, mysqlpassword, mysqlHost, mysqlPort, mysqlDbname)
	logs.Debug("start with mysql: ", fmt.Sprintf("%s@tcp(%s:%s)/%s?charset=utf8&loc=Local", mysqlUser, mysqlHost, mysqlPort, mysqlDbname))
	orm.RegisterDataBase("default", "mysql", dataSource)
	// register model
	orm.RegisterModel(new(User))
	orm.RegisterModel(new(UserOrderCoinage))
	orm.RegisterModel(new(UserCoinage))
	orm.RegisterModel(new(Goods))
	orm.RegisterModel(new(Article))
	orm.RegisterModel(new(Orders))
	orm.RegisterModel(new(UserMultiLogin))
	orm.RegisterModel(new(SystemUser))
	orm.RegisterModel(new(SystemUserMultiLogin))
	orm.RegisterModel(new(SystemConfig))
	orm.RegisterModel(new(SystemPerTProducingCoinage))
	orm.RegisterModel(new(UserPerTProducingCoinage))
	orm.RegisterModel(new(UserWithdraw))
	orm.RegisterModel(new(UserWithdrawRecommend))
	orm.RegisterModel(new(UserPledge))
	orm.RegisterModel(new(UserBonus))
	orm.RegisterModel(new(Sms))
	orm.RegisterModel(new(UserWithdrawalAddress))
	// create table
	orm.RunSyncdb("default", false, true)
	o := orm.NewOrm()
	salt, _ := GenerateSalt()
	password, _ := GeneratePassHash("server123", salt)
	if created, id, err := o.ReadOrCreate(&SystemUser{
		Mobile:   "admin",
		Username: "admin",
		Password: password,
		Salt:     salt,
	}, "mobile"); err == nil {
		if created {
			logs.Debug("create admin user", id)
		} else {
			logs.Debug("admin has exists")
		}
	}
	date := time.Now()
	id := 1
	invitationCode := "VFLL" //util.InviteCodeDefault.IdToCode(uint64(id))
	if _, err := o.Raw(`
INSERT INTO  user  ( id, invitation_code, references_user_id, mobile, password, trade_password, salt, status, reg_time, username, create_at, update_at )
VALUES
  ( ?, ?, ?, '', '', '', '', '0', ?, ?, ?, ? );

`, id, invitationCode, id, date, "系统默认", date, date).Exec(); err == nil {
		logs.Debug("create default references default ")
	} else {
		logs.Debug("references default has exists")
	}
	var w io.Writer
	orm.DebugLog = orm.NewLog(w)
}
func T(v interface{}) orm.QuerySeter {
	return orm.NewOrm().QueryTable(v)
}

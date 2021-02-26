package models

import (
	"fmt"
	"gin/src/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/v2"
	"github.com/prometheus/common/log"
	"gopkg.in/ini.v1"
)

var err error
var engin *gorose.Engin

func init() {
	cfg, _ := ini.Load("app.conf")
	mysqlConfig := cfg.Section("mysql")
	mysqlHost := mysqlConfig.Key("host").String()
	mysqlPort := mysqlConfig.Key("port").String()
	mysqlUser := mysqlConfig.Key("user").String()
	mysqlEncodePassword := mysqlConfig.Key("encodePassword").String()
	mysqlDbname := mysqlConfig.Key("dbname").String()
	mysqlpassword := util.DecodePassword(mysqlEncodePassword, util.DefaultSalt)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=Local&parseTime=true", mysqlUser, mysqlpassword, mysqlHost, mysqlPort, mysqlDbname)
	log.Debug("__gorose_mysql__", fmt.Sprintf("%s@tcp(%s:%s)/%s?charset=utf8&loc=Local&parseTime=true", mysqlUser, mysqlHost, mysqlPort, mysqlDbname))
	engin, err = gorose.Open(&gorose.Config{Driver: "mysql", Dsn: dsn})
}

func DB() gorose.IOrm {
	return engin.NewOrm()
}
func TRose(tableName string) gorose.IOrm {
	return DB().Table(tableName)
}

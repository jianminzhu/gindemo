package models

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
)

func GetConfigArr(configName string) ([]map[string]interface{}, error) {
	sysConfig := SystemConfig{
		Name: configName,
	}

	if err = orm.NewOrm().Read(&sysConfig, "name"); err == nil {
		items := []map[string]interface{}{}
		if err = json.Unmarshal([]byte(sysConfig.Value), &items); err == nil {
			return items, err
		}
	}
	return []map[string]interface{}{}, nil
}

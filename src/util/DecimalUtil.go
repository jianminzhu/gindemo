package util

import (
	"encoding/json"
	"github.com/prometheus/common/log"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
	"strconv"
)

func DecimalSum(params ...interface{}) decimal.Decimal {
	var sum decimal.Decimal = decimal.NewFromInt(0)
	if len(params) > 0 {
		if len(params) == 1 {
			sum, _ = decimal.NewFromString(params[0].(string))
		} else {
			for i := 0; i < len(params); i++ {
				if data, err := decimal.NewFromString(params[i].(string)); err == nil {
					sum = sum.Add(data)
				} else {
					log.Error("不是有效数值", params[i])
				}
			}
		}
	}
	return sum
}

func DecimalSumFixed(fixedLen int32, params ...interface{}) string {
	return DecimalSum(params...).StringFixed(fixedLen)
}

func D6(data decimal.Decimal) string {
	return data.StringFixed(6)
}
func D(value interface{}) decimal.Decimal {
	if val, err := decimal.NewFromString(Strval(value)); err == nil {
		return val
	} else {
		return decimal.NewFromInt(0)
	}
}
func DFromJson(json string, key string) decimal.Decimal {
	if val, err := decimal.NewFromString(gjson.Get(json, key).String()); err == nil {
		return val
	} else {
		return decimal.NewFromInt(0)
	}
}

func Strval(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}
	return key
}

func ToJson(tempMap interface{}) string {
	data, _ := json.Marshal(tempMap)
	return string(data)
}

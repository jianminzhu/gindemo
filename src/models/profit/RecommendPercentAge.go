package profit

import (
	"fmt"
	"gin/src/models"
	d "gin/src/util/date"
	"github.com/shopspring/decimal"
	"sort"
	"time"
)

type RecommendPercentAge struct {
	Power      decimal.Decimal
	Percentage decimal.Decimal
}

// 按照 RecommendPercentAge.Power 从小到大排序
type RecommendPercentAgeSlice []RecommendPercentAge

func (a RecommendPercentAgeSlice) Len() int { // 重写 Len() 方法
	return len(a)
}
func (a RecommendPercentAgeSlice) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}
func (a RecommendPercentAgeSlice) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return a[j].Power.GreaterThan(a[i].Power)
}

/** 不同推荐服务费等级配置 */
func RecommendLevelConfigArr() []RecommendPercentAge {
	datas, err := models.GetConfigArr("RecommendationServiceFee")
	ages := []RecommendPercentAge{}
	if err == nil {
		for _, data := range datas {
			power, err := decimal.NewFromString(data["Power"].(string))
			if err == nil {
				percentNum, err := decimal.NewFromString(data["Percentage"].(string))
				percentage := percentNum.Div(decimal.NewFromInt(100))
				if err == nil {
					ages = append(ages, RecommendPercentAge{
						power,
						percentage,
					})
				}
			}
		}
	}
	sort.Sort(RecommendPercentAgeSlice(ages))
	return ages
}

func GetRecommendPercent(power decimal.Decimal, feeConfigArr []RecommendPercentAge) decimal.Decimal {
	per := DZero
	for _, a := range feeConfigArr {
		if power.GreaterThanOrEqual(a.Power) {
			per = a.Percentage
		}
	}
	return per
}

func RecommendYmdConfig(orderArr []models.Orders, computeDate time.Time) map[string]decimal.Decimal {
	recommendYmdConfig := map[string]decimal.Decimal{}
	orderLen := len(orderArr)
	if orderLen > 0 {
		feeConfigArr := RecommendLevelConfigArr()
		totalPower := DZero
		startDate := d.AddDays(d.ParseYmd(orderArr[0].OrderTime), 1)
		for _, orders := range orderArr {
			runDate := d.ParseYmd(orders.OrderTime)
			runYMD := d.Format2Ymd(runDate)
			totalPower = totalPower.Add(decimal.NewFromFloat(orders.ActualPower))
			recommendYmdConfig[runYMD] = GetRecommendPercent(totalPower, feeConfigArr)
		}
		preConfig := recommendYmdConfig[orderArr[0].OrderTime]
		for {
			if startDate.After(computeDate) {
				break
			}
			runYMD := d.Format2Ymd(startDate)
			now, ok := recommendYmdConfig[runYMD]
			if ok == false {
				recommendYmdConfig[runYMD] = preConfig
			} else {
				preConfig = now
			}
			startDate = d.AddDays(startDate, 1)
			fmt.Println("config", runYMD, preConfig)
		}
	}

	return recommendYmdConfig
}

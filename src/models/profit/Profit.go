package profit

import (
	"fmt"
	"gin/src/models"
	d "gin/src/util/date"
	"github.com/astaxie/beego/orm"
	"github.com/shopspring/decimal"
	"sort"
	"time"
)

var DZero decimal.Decimal = decimal.Zero
var DPer25 decimal.Decimal = decimal.NewFromFloat(0.25)
var DPer75 decimal.Decimal = decimal.NewFromFloat(0.75)
var D180 decimal.Decimal = decimal.NewFromInt(180)

type Profit struct{}

/** 获取指定日期真实算力 */
func PowerRealActual(orders models.Orders, computeDate time.Time) (decimal.Decimal, decimal.Decimal) {
	startDate := d.ParseYmd(orders.OrderTime)
	endDate := d.ParseYmd(orders.OrderEndTime)
	if computeDate.Before(startDate) {
		return DZero, DZero
	}
	if d.DiffDays(computeDate, endDate) > 180 {
		return DZero, DZero
	}
	total := decimal.NewFromFloat(orders.ActualPower)
	splitDaysD := decimal.NewFromInt(orders.SplitDays)
	perDaySplit := total.Div(splitDaysD)
	diffFromStart := d.DiffDays(computeDate, startDate)
	theNewAdd := perDaySplit
	realActualPower := perDaySplit.Mul(decimal.NewFromInt(diffFromStart + 1))
	if realActualPower.GreaterThan(total) || diffFromStart >= orders.SplitDays {
		theNewAdd = DZero
		realActualPower = total
	}
	return realActualPower, theNewAdd
}

func FIL2575180(power decimal.Decimal, perTFil decimal.Decimal) map[string]decimal.Decimal {
	if power.LessThanOrEqual(DZero) {
		return map[string]decimal.Decimal{
			"p100":    DZero,
			"p25":     DZero,
			"p75":     DZero,
			"p75_180": DZero,
		}
	}
	p100 := power.Mul(perTFil)
	p25 := p100.Mul(DPer25)
	p75 := p100.Mul(DPer75)
	p75_180 := p75.Div(D180)
	data := map[string]decimal.Decimal{
		"p100":    p100,
		"p25":     p25,
		"p75":     p75,
		"p75_180": p75_180,
	}
	return data
}
func PerTYmdConfigWithIsSpecial(computeEndDate time.Time, userId int64) map[string]interface{} {
	config, isSpecial := computePertConfig(computeEndDate, userId)
	ymd := d.Format2Ymd(computeEndDate)
	return map[string]interface{}{
		"perTConfigMapping": config,
		"isSpecial":         isSpecial,
		"lastYmd":           ymd,
		"lastPerTConfig":    config[ymd],
	}
}
func PerTYmdConfig(computeEndDate time.Time, userId int64) map[string]decimal.Decimal {
	pertMapping, _ := computePertConfig(computeEndDate, userId)
	return pertMapping
}

/**获取单T 产币配置*/
func computePertConfig(computeEndDate time.Time, userId int64) (map[string]decimal.Decimal, bool) {
	o := orm.NewOrm()
	configArr := []models.SystemPerTProducingCoinage{}
	user := models.User{Id: userId}
	err := o.Read(&user)
	isSpecial := false
	if err == nil && user.PerTConfigSpecial == 1 {
		isSpecial = true

		userPerTProducingCoinageArr := []models.UserPerTProducingCoinage{}
		o.QueryTable(models.UserPerTProducingCoinage{}).Filter("user_id", userId).OrderBy("configDateStr").All(&userPerTProducingCoinageArr)
		for _, config := range userPerTProducingCoinageArr {
			configArr = append(configArr, models.SystemPerTProducingCoinage{
				ConfigDateStr:    config.ConfigDateStr,
				ProducingCoinage: config.ProducingCoinage,
			})
		}
	} else {
		o.QueryTable(models.SystemPerTProducingCoinage{}).OrderBy("configDateStr").All(&configArr)
	}
	perConfigMap := map[string]decimal.Decimal{}
	for _, config := range configArr {
		perConfigMap[config.ConfigDateStr] = decimal.NewFromFloat(config.ProducingCoinage)
	}

	configLen := len(configArr)
	if configLen > 0 {
		it := configArr[0]
		runYMD := it.ConfigDateStr
		runTime := d.ParseYmd(runYMD)
		preConfig := perConfigMap[runYMD]
		if configLen > 1 {
			for {
				runYMD = d.Format2Ymd(runTime)
				if runTime.After(computeEndDate) {
					break
				}
				runConfig, ok := perConfigMap[runYMD]
				if ok {
					preConfig = runConfig
				} else {
					perConfigMap[runYMD] = preConfig
				}
				runTime = d.AddDays(runTime, 1)
			}
		}
	}
	return perConfigMap, isSpecial
}

var STATUS_PASSED int64 = -1
var STATUS_NEED_EXAMINE int64 = 0

/** 某用户，推荐服务费计算 */
func UserRecommend(userId int64, computeDate time.Time) (map[int64]decimal.Decimal, decimal.Decimal) {
	//1 找到所有该用户推荐的好友
	//2 找到所有推荐用户的单子，按时间排序
	//3 计算那天开始收计算推荐服务费
	//1
	o := orm.NewOrm()
	ordersArr := []models.Orders{}
	o.Raw(`
select
  *
from
  orders as o  
where  o.user_id in
  (select
    id
  from
    user as u
  where u.references_user_id = ?)
  and  o.status = ?
order by o.order_time asc
`, userId, STATUS_PASSED).QueryRows(&ordersArr)
	userRecommendMap := map[int64]decimal.Decimal{}
	sort.Sort(models.SortOrderByOrderTimeAsc(ordersArr))
	recommendYmdPercent := RecommendYmdConfig(ordersArr, computeDate)

	totalRecommend := DZero
	for _, orders := range ordersArr {
		dataArr, _ := OrderProfit(orders, computeDate, PerTYmdConfig(computeDate, orders.UserId))
		recommend := ComputeRecommend(dataArr, recommendYmdPercent)
		userTotalRecommend, ok := userRecommendMap[orders.UserId]
		if ok {
			userTotalRecommend = userTotalRecommend.Add(recommend)
		} else {
			userRecommendMap[orders.UserId] = recommend
		}
		totalRecommend = totalRecommend.Add(recommend)
	}
	fmt.Print(totalRecommend)
	return userRecommendMap, totalRecommend
}

/**计算出订单开始时间到计算日期为止的每一天数据*/
func OrderProfit(orders models.Orders, computeDateDate time.Time, pertYmdConfig map[string]decimal.Decimal) ([]map[string]interface{}, map[string]map[string]interface{}) {
	startDate := d.ParseYmd(orders.OrderTime)
	endDate := d.ParseYmd(orders.OrderEndTime)
	runingDate := d.AddDays(startDate, 0)
	computeDate := d.ParseYmd(d.Format2Ymd(computeDateDate))
	runingDays := 1
	runingYMD := runingDate.Format(d.YMD)
	dayRuningMap := map[string]map[string]interface{}{}
	dayRuningArr := []map[string]interface{}{}
	p100_total := DZero
	p25_total := DZero
	p75_total := DZero
	p75_180_total := DZero
	for {
		if d.DiffDays(runingDate, endDate) > 181 || computeDate.Before(startDate) {
			break
		}
		if runingDate.After(computeDate) {
			break
		}
		runingYMD = d.Format2Ymd(runingDate)
		realAllocatedPower, pTheDay_add := PowerRealActual(orders, runingDate)
		perT, ok := pertYmdConfig[runingYMD] //单T算力产币
		if ok == false {
			perT = DZero
		}
		m := FIL2575180(realAllocatedPower, perT)
		p100 := m["p100"]
		p25 := m["p25"]
		p75 := m["p75"]
		p75_180 := m["p75_180"]
		p75_180_theDay := DZero // 当天180累计释放
		if runingDays > 1 {
			for _, dayData := range dayRuningArr {
				p75_180_theDay = p75_180_theDay.Add(dayData["p75_180"].(decimal.Decimal))
			}

			p75_180_total = p75_180_total.Add(p75_180_theDay)
		}
		pTotal_theDay := p25.Add(p75_180_theDay)
		p25_total = p25_total.Add(p25)
		p75_total = p75_total.Add(p75)
		p100_total = p100_total.Add(p100)
		pTotal := p25_total.Add(p75_180_total)
		totalPower := decimal.NewFromFloat(orders.ActualPower)
		theDay := map[string]interface{}{
			"userId":               orders.UserId,
			"id":                   orders.Id,
			"runingDays":           runingDays,                         // 运行天数
			"runingYMD":            runingYMD,                          // 运行日期
			"totalPower":           totalPower,                         // 总算力
			"allocatedPower":       realAllocatedPower,                 // 有效算力
			"needReleased":         totalPower.Sub(realAllocatedPower), // 得释放算力
			"today_allocatedPower": pTheDay_add,                        // 新增算力
			"pTheDay_add":          pTheDay_add,                        // 新增算力
			"perT":                 pertYmdConfig[runingYMD],           // 单T算力产币
			"p100":                 p100,                               // 当天产币
			"p25":                  p25,                                // 当天有效算力产币*25%
			"p75":                  p75,                                // 当天有效算力产币*75%
			"p75_180":              p75_180,                            // 当天产币 75% 的 1/180(隔天开始释放的平均值 )
			"pTotal":               pTotal,                             // 累计产币（今天25%释放+ 历史75%180释放 ）
			"p100_total":           p100_total,                         // 累计产币
			"p25_total":            p25_total,                          // 25%累计释放
			"p75_total":            p75_total,                          // 75%累计释放
			"p75_180_theDay":       p75_180_theDay,                     // 当天180释放总数
			"p75_180_total":        p75_180_total,                      // 累计180释放
			"pTotal_theDay":        pTotal_theDay,                      // 当天实际释放 （当天25%释放+  【每天要释放的180】的总和）

		}
		fmt.Println(runingYMD, realAllocatedPower, totalPower)
		dayRuningArr = append(dayRuningArr, theDay)
		dayRuningMap[runingYMD] = theDay
		runingDate = d.AddDays(runingDate, 1)
		runingDays++
	}
	return dayRuningArr, dayRuningMap
}

func ComputeRecommend(dayRuningArr []map[string]interface{}, recommendYmdPercent map[string]decimal.Decimal) decimal.Decimal {
	theOrderToRecommendTotal := DZero
	for _, m := range dayRuningArr {
		runingYMD := m["runingYMD"].(string)
		//给推荐人的币
		percentToRecommend, ok := recommendYmdPercent[runingYMD]
		if ok == false {
			percentToRecommend = DZero
		}
		theOrderToRecommendTotal = theOrderToRecommendTotal.Add(m["pTotal_theDay"].(decimal.Decimal).Mul(percentToRecommend))
	}
	return theOrderToRecommendTotal
}

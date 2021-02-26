package models

import (
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type CreateUpdate struct {
	CreateAt time.Time `json:"createAt" orm:"auto_now_add;type(datetime)"`
	UpdateAt time.Time `json:"updateAt" orm:"auto_now;type(datetime)"`
}

type Orders struct {
	Id           int64   `json:"id" orm:"column(id);auto"`
	UserId       int64   `json:"userId"  `
	GoodsId      int64   `json:"goodsId" `
	GoodsNums    int64   `json:"goodsNums" `
	OrderTime    string  `json:"orderTime" orm:"size(10)"`
	OrderEndTime string  `json:"orderEndTime"  orm:"size(10)"`
	ActualPower  float64 `json:"actualPower" orm:"digits(12);decimals(6);description(实际分配算力)"`
	SplitDays    int64   `json:"splitDays" orm:"description(算力释放天数)"`
	Status       int     `json:"status" orm:"default(-1);description(算力释放天数-1)"`
	CreateUpdate
}

// 按照  订单时间 从小到大排序
type SortOrderByOrderTimeAsc []Orders

func (a SortOrderByOrderTimeAsc) Len() int { // 重写 Len() 方法
	return len(a)
}
func (a SortOrderByOrderTimeAsc) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}
func (a SortOrderByOrderTimeAsc) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return a[j].OrderTime > a[i].OrderTime
}

type UserOrderCoinage struct {
	Id              int64     `json:"id" orm:"column(id);auto"`
	UserId          int64     `json:"userId" orm:"column(userId);"`
	OrderId         int64     `json:"orderId" orm:"column(orderId);"`
	CoinageDateStr  string    `json:"coinageDateStr" orm:"size(10)"`
	PowerEffective  float64   `json:"powerEffective" orm:"digits(12);decimals(6);description(当前释放总算力)"`
	Coinage         float64   `json:"coinage" orm:"digits(12);decimals(6)"`
	Coinage25       float64   `json:"coinage25" orm:"digits(12);decimals(6)"`
	Coinage75       float64   `json:"coinage75" orm:"digits(12);decimals(6)"`
	Coinage180total float64   `json:"coinage180total" orm:"digits(12);decimals(6)"`
	CoinageDate     time.Time `json:"coinageDate" orm:"type(datetime)"`
	CreateUpdate
}

type UserCoinage struct {
	Id                   int64     `json:"id" orm:"column(id);auto"`
	UserId               int64     `json:"userId"  `
	StatisticsDateStr    string    `json:"statisticsDateStr" orm:"size(10)"`
	Coinage25total       float64   `json:"coinage25total" orm:"digits(12);decimals(6)"`
	Coinage75total       float64   `json:"coinage75total" orm:"column(coinage75total);digits(12);decimals(6)"`
	Coinage180total      float64   `json:"coinage180total" orm:"column(coinage180total);digits(12);decimals(6)"`
	CoinageTotal         float64   `json:"coinageTotal" orm:"column(coinageTotal);digits(12);decimals(6)"`
	CoinageCanWithdrawal float64   `json:"coinageCanWithdrawal" orm:"column(coinageCanWithdrawal);digits(12);decimals(6)"`
	PowerEffective       float64   `json:"powerEffective" orm:"digits(12);decimals(6);description(有效算力)"`
	PowerTodayNew        float64   `json:"powerTodayNew" orm:"digits(12);decimals(6);description(今日新增算力)"`
	StatisticsDate       time.Time `json:"statisticsDate" orm:"type(datetime)"`
	CreateUpdate
}

type Sms struct {
	Id           int64     `json:"id" orm:"column(id);auto"`
	Mobile       string    `json:"userIdOrMobile" orm:"size(32)"`
	BusinessType string    `json:"businessType" orm:"size(8)"` //reg withdraw forgetPassword
	SmsCode      string    `json:"verifyCode" orm:"size(64)"`
	GenDate      time.Time `json:"genDate" orm:"auto_now;type(datetime);"`
	GenTimes     int       `json:"genTimes" orm:"default(1);description(生成次数)"`
	CreateUpdate
}

type Goods struct {
	Id              int64   `json:"id"`
	Name            string  `json:"name" orm:"size(32)"`
	Price           float64 `json:"price" orm:"digits(12);decimals(6)"`
	PriceUnit       string  `json:"priceUnit" orm:"size(8)"`
	ComputerRoomFee float64 `json:"computerRoomFee" orm:"digits(12);decimals(6);"`
	Power           float64 `json:"power" orm:"digits(12);decimals(6);description(算力)"`
	SplitDays       int64   `json:"splitDays" orm:"description(算力释放天数)"`
	RunDays         int64   `json:"runDays" orm:"description(产品运行天数)"`
	ExtAttrJsonStr  string  `json:"extAttrJsonStr" orm:"type(text)"`
	Status          int     `json:"status" orm:"default(0)"`
	CreateUpdate
}

type Article struct {
	Id          int64     `json:"id" orm:"auto"`
	Title       string    `json:"title" orm:"size(32)"`
	Type        string    `json:"type" orm:"size(32)"`
	Author      string    `json:"author" orm:"size(32)"`
	Content     string    `json:"content" orm:"type(text)"`
	PublishDate time.Time `json:"publishDate" orm:"auto_now_add;type(datetime)"`
	Status      int       `json:"status" orm:"default(0)"`
	CreateUpdate
}

type SystemConfig struct {
	Id        int64  `json:"id" orm:"auto"`
	Name      string `json:"name,omitempty" orm:"size(32)"`
	Value     string `json:"value,omitempty" orm:"type(text)"`
	ValueType string `json:"valueType,omitempty" orm:"size(8)"`
	CreateUpdate
}

type SystemPerTProducingCoinage struct {
	Id               int64   `json:"id" orm:"auto"`
	ConfigDateStr    string  `json:"configDateStr" orm:"column(configDateStr);size(10);description(生效日期)"`
	ProducingCoinage float64 `json:"producingCoinage" orm:"column(producingCoinage);digits(12);decimals(6);description(单T产币)"`
	CreateUpdate
}

type UserPerTProducingCoinage struct {
	Id               int64   `json:"id" orm:"auto"`
	UserId           int64   `json:"userId" `
	ConfigDateStr    string  `json:"configDateStr" orm:"column(configDateStr);size(10);description(生效日期)"`
	ProducingCoinage float64 `json:"producingCoinage" orm:"column(producingCoinage);digits(12);decimals(6);description(单T产币)"`
	CreateUpdate
}

type User struct {
	Id                int64     `json:"id" orm:"auto"`
	InvitationCode    string    `json:"invitationCode" orm:"size(8)"`
	ReferencesUserId  int64     `json:"referencesUserId" `
	Mobile            string    `json:"mobile" orm:"size(32)"`
	Password          string    `json:"password" orm:"column(password);size(128)"`
	TradePassword     string    `json:"tradePassword" orm:" size(128)"`
	Salt              string    `json:"salt"  orm:"column(salt);size(128)"`
	Status            int       `json:"status" `
	PerTConfigSpecial int       `json:"perTConfigSpecial"  orm:"default(0);"`
	RegTime           time.Time `json:"regTime"  orm:"auto_now_add;type(datetime)"`
	Username          string    `json:"username" orm:"column(username);size(128)"`
	CreateUpdate
}

type UserWithdraw struct {
	Id               int64     `json:"id" orm:"auto"`
	UserId           int64     `json:"userId" `
	WithdrawQuantity float64   `json:"withdrawQuantity" orm:"digits(12);decimals(6);"`
	WithdrawDate     time.Time `json:"withdrawDate"  orm:"auto_now_add;type(datetime)"`
	Status           int       `json:"status" orm:"default(0);description(-1已删除，0待审核，1已审核，2已提取)"`
	finishDate       time.Time `json:"finishDate"  orm:"auto_now;type(datetime)"`
	CreateUpdate
}
type UserWithdrawRecommend struct {
	Id               int64     `json:"id" orm:"auto"`
	UserId           int64     `json:"userId" `
	WithdrawQuantity float64   `json:"withdrawQuantity" orm:"digits(12);decimals(6);"`
	WithdrawDate     time.Time `json:"withdrawDate"  orm:"auto_now_add;type(datetime)"`
	Status           int       `json:"status" orm:"default(0);description(-1已删除，0待审核，1已审核，2已提取)"`
	finishDate       time.Time `json:"finishDate"  orm:"auto_now;type(datetime)"`
	CreateUpdate
}

type UserWithdrawalAddress struct {
	Id                int64  `json:"id" orm:"auto"`
	UserId            int64  `json:"userId" `
	WithdrawalAddress string `json:"withdrawalAddress" orm:" size(128)"`
	CreateUpdate
}

type UserPledge struct {
	Id             int64     `json:"id" orm:"auto"`
	UserId         int64     `json:"userId" `
	PledgeQuantity float64   `json:"pledgeQuantity" orm:"digits(12);decimals(6);"`
	PledgeDate     time.Time `json:"pledgeDate"  orm:"auto_now_add;type(datetime)"`
	Status         int       `json:"status" orm:"default(0);description(-2已解质,-1已删除，0质押)"`
	CreateUpdate
}

type UserBonus struct {
	Id            int64     `json:"id" orm:"auto"`
	UserId        int64     `json:"userId" `
	BonusQuantity float64   `json:"pledgeQuantity" orm:"digits(12);decimals(6);"`
	BonusDate     time.Time `json:"pledgeDate"  orm:"auto_now_add;type(datetime)"`
	Status        int       `json:"status" orm:"default(0);"`
	CreateUpdate
}

type UserMultiLogin struct {
	Id         int64  `json:"id" orm:"auto"`
	UserId     int64  `json:"userId" orm:"column(userId);"`
	ClientType string `json:"clientType" orm:" size(8)"`
	Status     int    `json:"status" `
	Token      string `json:"token" orm:" size(512)"`
	CreateUpdate
}

type SystemUser struct {
	Id       int64     `json:"id" orm:"auto"`
	Username string    `json:"username" orm:"column(username);size(128)"`
	Password string    `json:"password" orm:"column(password);size(128)"`
	Salt     string    `json:"salt"  orm:"column(salt);size(128)"`
	Mobile   string    `json:"mobile" orm:"size(32)"`
	RegTime  time.Time `json:"regTime"  orm:"auto_now_add;type(datetime)"`
	NickName string    `json:"nickName" orm:"size(32)"`
	Token    string    `json:"token" orm:"size(512)"`
	CreateUpdate
}
type SystemUserMultiLogin struct {
	Id         int64  `json:"id" orm:"auto"`
	UserId     int64  `json:"userId" orm:"column(userId);"`
	ClientType string `json:"clientType" orm:" size(8)"`
	Status     int    `json:"status" orm:" size(128)"`
	Token      string `json:"token" orm:" size(512)"`
	CreateUpdate
}

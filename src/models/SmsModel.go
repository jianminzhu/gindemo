package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"github.com/prometheus/common/log"
	"math/rand"
	"strings"
	"time"
)

type SmsSendResponse struct {
	Message   string
	RequestId string
	BizId     string
	Code      string

	BusinessType string
	SmsCode      string
}

type SmsRequest struct {
	Mobile       string
	UserId       int64
	BusinessType string
}

var BusinessTypeMaxGenTimesLimitMapping map[string]int = map[string]int{
	"reg":                 10,
	"forgetLoginPassword": 5,
	"forgetTradePassword": 5,
}

type SmsUtil struct{}

func CheckSmsCode(businessType string, mobile string, smsCode string) bool {
	count, _ := orm.NewOrm().QueryTable("sms").Filter("business_type", businessType).Filter("mobile", mobile).Filter("sms_code", smsCode).Count()
	return count > 0
}
func (su SmsUtil) SendSms(businessType string, mobile string) (Sms, error) {
	//1 存储到 短信验证码，发送时间  存储数据库，
	//2 根据业务类型及上次发送时间决定，要不要再次发送验证码，防止恶意频次 调用 验证码功能
	//3 调用 云短信验证码 接口发送验证码
	//根据商业类型在数据库中查找，如果存在，则查看是否超过发布次数限制
	o := orm.NewOrm()
	sms := Sms{
		BusinessType: businessType,
		Mobile:       mobile,
	}
	valid := validation.Validation{}
	valid.Required(mobile, "").Message("手机号不能为空")
	valid.Mobile(mobile, "").Message("手机号不正确")
	if !valid.HasErrors() {
		if businessType == "reg" {
			nums, _ := orm.NewOrm().QueryTable("user").Filter("mobile", mobile).Count()
			if nums > 0 {
				valid.Required("", "").Message("手机号已经注册")
			}
		}
	}
	if !valid.HasErrors() {
		code := su.GenValidateCode(6)
		if err := o.Read(&sms, "mobile", "business_type"); err == nil {
			//如果存在则检查是否超过限定次数
			limitTimes := BusinessTypeMaxGenTimesLimitMapping[businessType]
			if sms.GenTimes > limitTimes {
				err := errors.New("验证码请求次数太多，请稍后重试")
				return sms, err
			}
			//如果超过一天，则重置
			day := time.Now().Sub(sms.GenDate).Hours() / 24
			if day > 1 {
				sms.GenTimes = 0
			}
			sms.GenTimes += 1
			sms.SmsCode = code
			o.Update(&sms, "gen_times", "sms_code")
		} else {
			sms.SmsCode = code
			sms.GenTimes = 1
			if id, err := o.Insert(&sms); err == nil {
				sms.Id = id
			} else {
				return sms, errors.New("系统异常，请稍后重试")
			}
		}
		log.Debug("send sms for reg", mobile, code)
		res, err := su.ToMobile(mobile, code)
		log.Debug("send sms for reg Object ", res)
		if err != nil {
			return sms, errors.New("系统发送验证码失败")
		}
		return sms, err
	} else {
		return sms, errors.New(valid.Errors[0].Message)
	}
}

func (su SmsUtil) GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func (su SmsUtil) ToMobile(phoneNumbers string, code string) (SmsSendResponse, error) {
	var msg SmsSendResponse //阿里云返回的json信息对应的类
	regionId := "cn-hangzhou"
	accessKeyId := "LTAIKv6VnRQ5kgqz"
	accessKeySecret := "2nWB5IDutuwZk04ZDmdhZkI2Z3NP7q"
	templateCode := "SMS_111555002"

	client, err := sdk.NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret)
	if err != nil {
		return msg, err
	}
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = regionId
	request.QueryParams["PhoneNumbers"] = phoneNumbers //手机号
	signName := "古数信息注册验证"
	request.QueryParams["SignName"] = signName                       //阿里云验证过的项目名 自己设置
	request.QueryParams["TemplateCode"] = templateCode               //阿里云的短信模板号 自己设置
	request.QueryParams["TemplateParam"] = `{"code":"` + code + `"}` //短信模板中的验证码内容 自己生成
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return msg, err
	}
	err = json.Unmarshal(response.GetHttpContentBytes(), &msg)
	if err != nil {
		return msg, err
	} else {
		return msg, nil
	}
}

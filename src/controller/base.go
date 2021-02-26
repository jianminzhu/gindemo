package controller

import (
	"context"
	"gin/src/models"
	"github.com/gin-gonic/gin"
	"time"
)

type Base struct {
	context.Context
	Gin *gin.Context
}


func getCode(isSucc bool) int {
	code := 200
	if isSucc == false {
		code = 400
	}
	return code
}

var SYSTEM_CODE_NEED_LOGIN = 50001
var DefaultEmptyObj map[string]interface{} = map[string]interface{}{}
var DefaultEmptyArr []map[string]interface{} = []map[string]interface{}{}

func (c *Base) NeedLogin(msg string) {
	res := map[string]interface{}{
		"code":   SYSTEM_CODE_NEED_LOGIN,
		"isSucc": false,
		"msg":    msg,
		"data":   DefaultEmptyObj,
	}
	c.Gin.String(200, models.ToJson(res))
}
func (c *Base) FailDefObj(msg string) {
	res := map[string]interface{}{
		"code":   getCode(false),
		"isSucc": false,
		"msg":    msg,
		"data":   DefaultEmptyObj,
	}
	c.Gin.String(200, models.ToJson(res))
}
func (c *Base) FailDefArr(msg string) {
	res := map[string]interface{}{
		"code":   getCode(false),
		"isSucc": false,
		"msg":    msg,
		"data":   DefaultEmptyArr,
	}
	c.Gin.String(200, models.ToJson(res))
}

func (c *Base) Fail(data interface{}, msg string) {
	res := map[string]interface{}{
		"code":   getCode(false),
		"isSucc": false,
		"msg":    msg,
		"data":   data,
	}
	c.Gin.String(200, models.ToJson(res))
}

func (c *Base) Ok(data interface{}, msg string) {
	res := map[string]interface{}{
		"code":   getCode(true),
		"isSucc": true,
		"data":   data,
		"msg":    msg,
	}
	c.Gin.String(200, models.ToJson(res))
}

type BaseHandleFunc func(c *Base)

func WithFunc(baseHandle BaseHandleFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 可以在gin.Context中设置key-value
		c.Set("trace", "假设这是一个调用链追踪sdk")
		// 全局超时控制
		timeoutCtx, _ := context.WithTimeout(c, 5*time.Second)
		// ZDM上下文
		yuerCtx := Base{Context: timeoutCtx, Gin: c}
		// 回调接口
		baseHandle(&yuerCtx)
	}
}

func WithAuthFunc(baseHandle BaseHandleFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.DefaultQuery("token", "")
		if token == "" {
			base := Base{Gin: c}
			base.NeedLogin("needLogin")
		} else {
			// 可以在gin.Context中设置key-value
			c.Set("trace", "假设这是一个调用链追踪sdk")
			// 全局超时控制
			timeoutCtx, _ := context.WithTimeout(c, 5*time.Second)
			// ZDM上下文
			yuerCtx := Base{Context: timeoutCtx, Gin: c}
			// 回调接口
			baseHandle(&yuerCtx)
		}
	}
}

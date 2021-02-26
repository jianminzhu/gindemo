package user

import (
	"gin/src/util/tpl"
	"github.com/gin-gonic/gin"
	"net/http"
)
type User struct {

}
func GetTest(c *gin.Context) {
	id := c.Query("id")
	page := c.DefaultQuery("page", "0")
	name := c.DefaultPostForm("name", "")
	message := c.DefaultPostForm("message", "")
	res := tpl.ProcessString( "---- 22999 -----{{.id}} of {{.page}}-----------", map[string]interface{}{
		"id":      id,
		"page":     page,
		"name":    name,
		"message": message,
	})
	c.String(http.StatusOK, res)
}

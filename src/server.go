package main

import (
	"fmt"
	"gin/asset"
	c "gin/src/controller"
	"gin/src/controller/system"
	"gin/src/controller/user"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
	//"github.com/thinkerou/favicon"
	"gopkg.in/ini.v1"
	"net/http"
)

func main() {
	r := gin.Default()
	cfg, _ := ini.Load("app.conf")
	site := cfg.Section("site")
	dirForUpload := site.Key("dir_for_upload").String()
	tpl := assetfs.AssetFS{Asset: asset.Asset, AssetDir: asset.AssetDir, AssetInfo: asset.AssetInfo, Prefix: "static", Fallback: "index.html"}
	r.StaticFS("/static", &tpl)
	//r.Use(favicon.New("./favicon.ico"))
	fmt.Println("dir", dirForUpload)
	r.StaticFS("/upload", http.Dir(dirForUpload))
	apiRouter := site.Key("apiRouter").String()
	v1 := r.Group(apiRouter)
	{
		v1.GET("/health", c.WithFunc(system.Health))
		v1.GET("/system/test", c.WithAuthFunc(system.User))
		v1.GET("/user/test", user.GetTest)
	}
	r.Run(":9000")
}

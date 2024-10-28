package main

import (
	"scnu_acm_rank/biz/config"
	"scnu_acm_rank/biz/model"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/cors"
)

func main() {
	h := server.Default()
	h.Use(cors.New(cors.Config{
		// 允许使用所有的 Origin 头
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		// 允许使用以下的方法
		AllowMethods: []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		// 发起非简单请求时, 允许使用以下头信息
		AllowHeaders: []string{"Origin", "Authorization", "Content-Type"},
		// 允许暴露给客户端的响应头
		ExposeHeaders: []string{"Content-Type"},
		// 允许客户端请求携带用户凭证
		AllowCredentials: true,
		// 超时时间设定
		MaxAge: 24 * time.Hour,
	}))
	register(h)                 // 注册路由
	model.GetDB()               // 配置数据库信息
	config.Update <- struct{}{} // 初始化服务器配置信息

	h.Spin()
}

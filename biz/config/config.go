package config

import (
	"runtime"
	"scnu_acm_rank/biz/model"
)

var Conf model.Config        // 本项目运行所需要的配置 (Cookie, Email 等)
var Update chan struct{}     // 更新配置的信号量, 传入即更新配置
var updateList []*updateConf // 配置更新队列(更新那些依赖于 Conf 的变量)

// 初始化配置信息并将其放置于 runtime 中防止被回收
func init() {
	Conf = model.Config{}
	Update = make(chan struct{})
	updateList = make([]*updateConf, 0, 20)
	runtime.KeepAlive(Conf)
	runtime.KeepAlive(Update)
	runtime.KeepAlive(updateList)
	go update()
}

// 将那些依赖于 Conf 的其他配置变量添加到更新队列中
// 必须在 init() 中使用此函数
func Add(a updateConf) {
	updateList = append(updateList, &a)
}

// 更新项目中的配置
func update() {
	for range Update {
		// 从数据库中读取最新的配置
		model.DB.Model(&model.Config{}).First(&Conf)
		for _, v := range updateList { // 依次更新配置队列
			(*v).Update()
		}
	}
}

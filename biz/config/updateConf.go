package config

// 依赖于 Conf 的其他配置结构
// 必须实现 Update() 并在其中更新对于 Conf 的依赖关系
type updateConf interface {
	Update()
}

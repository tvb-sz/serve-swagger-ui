package initializer

import "github.com/tvb-sz/serve-swagger-ui/client"

// region 全局句柄初始化相关

// Init 初始化
//go:noinline
func Init() {
	client.Logger = iniLogger()            // 初始化logger，需要优先执行
	client.MemoryCache = initMemoryCache() // 初始化内存缓存
}

// endregion

// region 全局句柄热重载后初始化相关

// Reload 热重载再次初始化调用
//  - 当配置变更时又无需重新启动进程触发，监听调用
func Reload() {

}

// endregion

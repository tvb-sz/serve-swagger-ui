package client

// 定义全局句柄、需要初始化的变量，然后在initializer子包内单个文件进行初始化

import (
	"github.com/jjonline/go-lib-backend/logger"
	"github.com/jjonline/go-lib-backend/memory"
)

// 全局使用的句柄变量--client
//  ① 项目启动就要初始化的全局变量句柄
//  ② 按client包规则书写
var (
	Logger      *logger.Logger // 基于zap的日志记录器
	MemoryCache *memory.Cache  // 本地内存缓存客户端
)

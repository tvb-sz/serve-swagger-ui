package conf

import (
	"github.com/tvb-sz/serve-swagger-ui/utils/cfg"
)

// 项目config配置定义

// 暴露给全局使用的配置变量
var (
	Config config
	Cmd    cmdConfig
)

// config 项目配置上层结构
type config struct {
	Server       server  `json:"server"`  // server config
	Log          log     `json:"log"`     // log config
	Google       google  `json:"google"`  // google oauth config
	Swagger      swagger `json:"swagger"` // swagger json file config
	ConfigFile   string  `json:"-"`       // record config file path
	EnableGoogle bool    `json:"-"`       // record is set google client_id & client_secret
}

// cmdConfig command line args
type cmdConfig struct {
	ConfigFile         string // set config file path
	Host               string // set web server host ip
	Port               string // set web server ports
	SwaggerPath        string // set swagger file path
	GoogleClientID     string // set google oauth app-key
	GoogleClientSecret string // set google oauth app-secret
	LogLevel           string // set logger level
	LogPath            string // set logger path
}

// parseAfterLoad 配置项加载完成后的统一处理流程逻辑
func (c config) parseAfterLoad() {

}

// region 初始化

// Init 初始化
func Init() {
	// must load or panic quit
	var cfgLoader cfg.IFace
	cfgLoader = cfg.Viper{}
	if err := cfgLoader.Parse(Cmd.ConfigFile, Cmd.ConfigType, &Config); err != nil {
		panic(err)
	}

	// 配置加载并解析映射成功后统一处理逻辑：譬如Url统一处理后缀斜杠
	Config.parseAfterLoad()
}

// endregion

// region 热重载

// Reload 热重载
// 当配置变更时又无需重新启动进程触发，监听调用
func Reload() {

}

// endregion

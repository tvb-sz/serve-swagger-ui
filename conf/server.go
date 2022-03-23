package conf

// server 环境、端口、版本、http读写超时等server级别的配置
type server struct {
	SiteName     string `json:"SiteName"`      // site name
	Version      string `json:"-"`             // version
	Host         string `json:"host"`          // ip host
	Port         int    `json:"port"`          // port
	ReadTimeout  int    `json:"read_timeout"`  // request timeout waiting read
	WriteTimeout int    `json:"write_timeout"` // request timeout waiting write
	Cors         bool   `json:"cors"`          // Open or Close Cross-Origin Resource Sharing
}

package conf

// log 日志相关配置
type log struct {
	Level string `json:"level"` // log level,can use values: panic|fatal|error|warning|info|debug|trace
	Path  string `json:"path"`  // log storage path,can use values: stderr|stdout|file_path_only_dir
}

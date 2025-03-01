package config

import "github.com/labstack/echo/v4/middleware"

func (cfg *Config) GetEchoLogConfig() middleware.LoggerConfig {
	echoLogConf := middleware.DefaultLoggerConfig
	if cfg.isDev() {
		echoLogConf.Format = `[${time_rfc3339}] ${status} ${method} ${uri} ${latency_human}` +
			` ${remote_ip} ${host} ${error}` + "\n"
	}
	return echoLogConf
}

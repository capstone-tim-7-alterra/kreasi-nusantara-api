package logger

import "github.com/labstack/echo/v4/middleware"

func GetLoggerConfig() *middleware.LoggerConfig {
	return &middleware.LoggerConfig{
		Format: "method=${method}, url=${host}${uri}, status=${status}, latency=${latency_human}\n",
	}
}
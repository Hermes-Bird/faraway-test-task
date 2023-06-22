package config

import "time"

var ChallengeCondition = []byte{0, 0, 0}
var ServerShutdownTimeout = 5 * time.Second

type ServerConfig struct {
	MaxConn     int32
	ConnTimeout int
	Port        string
}

func DefaultConfig() ServerConfig {
	return ServerConfig{
		MaxConn:     5000,
		ConnTimeout: 7000,
		Port:        "8080",
	}
}

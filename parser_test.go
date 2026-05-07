package config

import (
	"strings"
	"testing"
	"time"
)

/*
server_addr: "ws://localhost:8443/sidecar/connect"  # 服务端 WebSocket 地址（生产使用 wss://）
token: "your-auth-token"                            # 鉴权 Token
reconnect_interval: 5s                              # 重连间隔基础值
max_reconnect_interval: 60s                         # 最大重连间隔
conn_alias: "abc"
*/
type SidecarConfig struct {
	ServerAddr           string        `mapstructure:"server_addr"`
	AppSn                string        `mapstructure:"app_sn"`
	AppId                string        `mapstructure:"app_id"`
	AppSec               string        `mapstructure:"app_sec"`
	ReconnectInterval    time.Duration `mapstructure:"reconnect_interval"`
	MaxReconnectInterval time.Duration `mapstructure:"max_reconnect_interval"`
	ConnAlias            string        `mapstructure:"conn_alias"`
}

type ServerConfig struct {
	HTTPAddr     string        `mapstructure:"http_addr"`
	GRPCAddr     string        `mapstructure:"grpc_addr"`
	ProxyTimeout time.Duration `mapstructure:"proxy_timeout"`
}

type Profile struct {
	Name    string         `mapstructure:"name"`
	Sidecar *SidecarConfig `mapstructure:"sidecar"`
	Server  *ServerConfig  `mapstructure:"server"`
}

func TestParser(t *testing.T) {
	var parser Unmarshal
	configFile := "configs"
	cs := strings.Split(configFile, ",")
	for _, f := range cs {
		parser = FromFiles(f, YAML)
	}
	profile := Profile{}
	parser.Unmarshal(&profile)
	t.Logf("sidecarConfig: %v", profile.Sidecar.ServerAddr)
	t.Logf("serverConfig: %v", profile.Server.HTTPAddr)
	t.Logf("name: %v", profile.Name)

}

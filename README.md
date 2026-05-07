# config

基于 Viper 的轻量配置加载器：支持从目录或单文件读取配置，并按环境（RUN_MODE）自动选择配置来源；同一目录下的多个配置文件会合并到同一个配置对象中。

## 安装

```bash
go get github.com/w6xian/config
```

## 快速开始

目录模式（推荐）：按环境分目录组织配置，例如 `configs/dev/*.yaml`、`configs/prod/*.yaml`。

```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/w6xian/config"
)

type SidecarConfig struct {
	ServerAddr           string        `mapstructure:"server_addr"`
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

func main() {
	_ = os.Setenv("RUN_MODE", "dev")

	var p Profile
	config.FromFiles("configs", config.YAML).Unmarshal(&p)

	fmt.Println("name:", p.Name)
	fmt.Println("server.http_addr:", p.Server.HTTPAddr)
	fmt.Println("sidecar.server_addr:", p.Sidecar.ServerAddr)
}
```

## 配置结构

### 目录模式

传入目录名 `configs` 时，会自动拼出 `configs/<RUN_MODE>/`，并读取该目录下所有后缀为指定类型的文件，然后依次合并。

```text
configs/
  dev/
    server.yaml
    sidecar.yaml
  prod/
    server.yaml
    sidecar.yaml
```

### 单文件模式

如果传入的路径不存在，会按 `<path>.<RUN_MODE>.<type>` 去尝试读取。例如：

```text
app.dev.yaml
app.prod.yaml
```

对应调用：

```go
config.FromFiles("app", config.YAML)
```

## API

- `FromFiles(path, type) Unmarshal`
  - `path`：目录路径或配置基名（单文件模式会补全环境与后缀）
  - `type`：配置类型，使用常量 `config.YAML` / `config.JSON` / `config.TOML`
- `GetMode() string`
  - 从环境变量 `RUN_MODE` 读取运行模式；未设置时默认 `dev`
- `Unmarshal(rawVal any, opts ...viper.DecoderConfigOption) error`
  - 将合并后的配置反序列化到结构体（基于 `mapstructure` 标签）

## 注意事项

- 本库通过 Viper 的全局实例合并配置；同一进程内多次加载会继续叠加合并结果。
- 同一进程内如果需要加载多套互相隔离的配置，建议在业务侧规划为单例加载，避免并发/隔离需求。

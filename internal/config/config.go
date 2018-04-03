package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/naoina/toml"
	"github.com/shiguanghuxian/face-login/internal/common"
)

var CFG *Config

// Config 系统根配置
type Config struct {
	Debug        bool         `toml:"debug"`
	APIKey       string       `toml:"api_key"`
	APISecret    string       `toml:"api_secret"`
	FacesetToken string       `toml:"faceset_token"`
	Http         *HTTPConfig  `toml:"http"`
	MySQL        *MySQLConfig `toml:"mysql"`
}

// HTTPConfig http监听配置
type HTTPConfig struct {
	Address string `toml:"address"`
	Port    int    `toml:"port"`
}

// Init http配置
func (hc *HTTPConfig) Init() {
	if hc.Address == "" {
		hc.Address = "0.0.0.0"
	}
	if hc.Port == 0 {
		hc.Port = 8080
	}
}

// NewConfig 初始化一个配置文件对象
func NewConfig(cfgPath string) (*Config, error) {
	if cfgPath == "" {
		rootPath := common.GetRootDir()
		cfgPath = fmt.Sprintf("%s%sconfig%sconfig.toml", rootPath, string(os.PathSeparator), string(os.PathSeparator))
	}
	log.Println(cfgPath)

	f, err := os.Open(cfgPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	config := new(Config)
	if err := toml.NewDecoder(f).Decode(config); err != nil {
		return nil, err
	}

	if config.Http == nil {
		return nil, errors.New("http配置不能为空")
	}
	config.Http.Init()
	CFG = config
	return config, nil
}

// Duration 用于日志文件解析出时间段
type Duration struct {
	time.Duration
}

// UnmarshalText implements encoding.TextUnmarshaler
func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

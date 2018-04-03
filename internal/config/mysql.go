package config

import "time"

// MySQLConfig 数据库配置
type MySQLConfig struct {
	Debug        bool     `toml:"debug"`          // 是否调试模式
	Address      string   `toml:"address"`        // 数据库连接地址
	Port         int      `toml:"port"`           // 数据库端口
	MaxIdleConns int      `toml:"max_idle_conns"` // 连接池最大连接数
	MaxOpenConns int      `toml:"max_open_conns"` // 默认打开连接数
	User         string   `toml:"user"`           // 数据库用户名
	Passwd       string   `toml:"passwd"`         // 数据库密码
	DbName       string   `toml:"db_name"`        // 数据库名
	Prefix       string   `toml:"prefix"`         // 数据库表前缀
	PingInterval Duration `toml:"ping_interval"`  // 定时保活定
}

// Init 初始化数据库配置
func (db *MySQLConfig) Init() {
	if db.Address == "" {
		db.Address = "127.0.0.1"
	}
	if db.Port == 0 {
		db.Port = 3306
	}
	if db.MaxIdleConns == 0 {
		db.MaxIdleConns = 64
	}
	if db.MaxOpenConns == 0 {
		db.MaxOpenConns = 16
	}
	if db.User == "" {
		db.User = "root"
	}
	if db.PingInterval.Seconds() == 0 {
		duration, _ := time.ParseDuration("10s")
		db.PingInterval.Duration = duration
	}
}

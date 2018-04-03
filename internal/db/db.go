package db

import (
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/logger"
	"github.com/jinzhu/gorm"
	"github.com/shiguanghuxian/face-login/internal/config"
)

// 数据库操作对象
var (
	DB *gorm.DB
)

// InitDB 初始化数据库
func InitDB(config *config.MySQLConfig) error {
	mySqlDB, err := New(config)
	if err != nil {
		return err
	}
	DB, err = mySqlDB.GetMySQLDB()
	return err
}

// MySQLDB mysql 结构体
type MySQLDB struct {
	db     *gorm.DB
	config *config.MySQLConfig
}

// GetMySQLDB 获取mysql操作对象
func (mysql *MySQLDB) GetMySQLDB() (*gorm.DB, error) {
	if mysql.db == nil {
		return nil, errors.New("连接对象为nil")
	}
	return mysql.db, nil
}

// Close 关闭连接
func (mysql *MySQLDB) Close() error {
	if mysql.db != nil {
		return mysql.db.Close()
	}
	return nil
}

// New 实例化连接MySQL数据库
func New(config *config.MySQLConfig) (*MySQLDB, error) {
	if config == nil {
		return nil, errors.New("配置文件不能为nil")
	}
	if config.Address == "" {
		return nil, errors.New("需要提供一个MySql连接地址")
	}
	logger.Infoln("正在与MySql建立连接")

	// 拼接连接数据库字符串
	connStr := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=UTC",
		config.User,
		config.Passwd,
		config.Address,
		config.Port,
		config.DbName)

	// 连接数据库
	db, err := gorm.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}
	// 设置表名前缀
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return config.Prefix + defaultTableName
	}
	// 禁用表名多元化
	db.SingularTable(true)
	// 是否开启debug模式
	if config.Debug {
		// debug 模式
		db = db.Debug()
	}

	// 连接池最大连接数
	db.DB().SetMaxIdleConns(config.MaxIdleConns)
	// 默认打开连接数
	db.DB().SetMaxOpenConns(config.MaxOpenConns)

	// 开启协程ping MySQL数据库查看连接状态
	go func() {
		for {
			// ping
			err = db.DB().Ping()
			if err != nil {
				logger.Infoln(err)
			}
			// 间隔5s ping一次
			time.Sleep(config.PingInterval.Duration)
		}
	}()
	mysqlDB := &MySQLDB{db: db, config: config}
	logger.Infoln("与MySQL建立连接成功")
	return mysqlDB, err
}

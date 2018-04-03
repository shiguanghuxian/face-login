package faceserver

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/logger"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/shiguanghuxian/face-login/internal/common"
	"github.com/shiguanghuxian/face-login/internal/config"
	"github.com/shiguanghuxian/face-login/internal/context"
	"github.com/shiguanghuxian/face-login/internal/db"
)

// FaceServer 程序服务对象
type FaceServer struct {
	e   *echo.Echo
	cfg *config.Config
}

// New 创建服务端对象
func New() *FaceServer {
	// 系统日志显示文件和行号
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	// 初始化配置文件
	cfg, err := config.NewConfig("")
	if err != nil {
		logger.Fatalln("配置文件读取失败:", err)
	}
	js, _ := json.Marshal(cfg)
	log.Println(string(js))
	// echo对象
	e := echo.New()
	e.Use(context.InitZContext())
	// 注册中间件
	e.Use(middleware.Logger()) // 根据配置将日志输出到哪里
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	e.Static("/", common.GetRootDir()+"/public")

	// 初始化mysql
	err = db.InitDB(cfg.MySQL)
	if err != nil {
		logger.Fatalln("mysql连接错误:", err)
	}

	return &FaceServer{
		e:   e,
		cfg: cfg,
	}
}

// Run 启动服务
func (fs *FaceServer) Run() {
	// 路由
	fs.Route()
	// 启动服务
	address := fmt.Sprintf("%s:%d", fs.cfg.Http.Address, fs.cfg.Http.Port)
	err := fs.e.Start(address)
	if err != nil {
		fs.e.Logger.Fatal(err)
	}
	fs.e.Logger.Info()

}

// Stop 停止服务
func (fs *FaceServer) Stop() {

}

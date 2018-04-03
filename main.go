/*
 * @Author: 时光弧线
 * @Date: 2018-03-28 10:05:04
 * @Last Modified by: 时光弧线
 * @Last Modified time: 2018-03-28 11:07:58
 */
package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/google/logger"
	"github.com/kardianos/service"
	"github.com/shiguanghuxian/face-login/faceserver"
	"github.com/shiguanghuxian/face-login/internal/common"
)

// 程序版本
var (
	VERSION    string
	BUILD_TIME string
	GO_VERSION string
	GIT_HASH   string
)

var serviceLogger service.Logger
var appServer *faceserver.FaceServer

type program struct{}

func (p *program) Start(s service.Service) error {
	// 开启异步任务，开启服务
	go p.run()
	return nil
}

func (p *program) run() {
	// 存储pid
	err := common.WritePidToFile("faceserver")
	if err != nil {
		log.Println("写pid文件错误")
	}
	// 启动服务
	appServer = faceserver.New()
	appServer.Run()
}

func (p *program) Stop(s service.Service) error {
	// 删除pid文件
	common.RemovePidFile("faceserver")
	// 停止服务
	appServer.Stop()
	// 停止任务，3秒后返回
	<-time.After(time.Second * 1)
	return nil
}

func main() {
	// 全部核心运行程序
	runtime.GOMAXPROCS(runtime.NumCPU())

	svcConfig := &service.Config{
		Name:        "faceserver",
		DisplayName: "faceserver",
		Description: "面部识别登录系统服务端",
	}
	// 实例化
	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		logger.Fatalln(err)
	}
	// 接收一个参数，install | uninstall | start | stop | restart
	if len(os.Args) > 1 {
		if os.Args[1] == "-v" || os.Args[1] == "-version" {
			ver := fmt.Sprintf("Version: %s\nBuilt: %s\nGo version: %s\nGit commit: %s", VERSION, BUILD_TIME, GO_VERSION, GIT_HASH)
			fmt.Println(ver)
			return
		}
		err = service.Control(s, os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	serviceLogger, err = s.Logger(nil)
	if err != nil {
		logger.Fatalln(err)
	}
	err = s.Run()
	if err != nil {
		serviceLogger.Error(err)
	}
}

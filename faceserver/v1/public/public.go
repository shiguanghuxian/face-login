package public

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	sdk "53it.net/face-golang-sdk"
	"github.com/google/logger"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	"github.com/shiguanghuxian/face-login/internal/common"
	"github.com/shiguanghuxian/face-login/internal/config"
	"github.com/shiguanghuxian/face-login/internal/db"
	"github.com/shiguanghuxian/face-login/model"
)

// PublicController 公共可访问接口
type PublicController struct {
}

// Login 登录
func (pc *PublicController) Login(c echo.Context) error {
	// 根据字段名获取表单文件
	formFile, header, err := c.Request().FormFile("file")
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	defer formFile.Close()

	// 获取文件后缀
	postfix := "png"
	index := strings.LastIndex(header.Filename, ".")
	if index > 0 {
		postfix = header.Filename[index+1:]
	}

	// 拼接保存路径
	pathSeparator := string(os.PathSeparator)
	uuid, _ := uuid.NewV4()
	fileName := uuid.String()
	savePath := fmt.Sprintf("%scache%s%s.%s", pathSeparator, pathSeparator, fileName, postfix)
	picPath := fmt.Sprintf("%s%spublic%s", common.GetRootDir(), pathSeparator, savePath)
	// 创建保存文件
	destFile, err := os.Create(picPath)
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	defer os.Remove(picPath)
	defer destFile.Close()

	// log.Println(picPath)

	// 读取表单文件，写入保存文件
	_, err = io.Copy(destFile, formFile)
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// 验证用户脸部信息
	faceSDK, err := sdk.NewFaceSDK(config.CFG.APIKey, config.CFG.APISecret, config.CFG.Debug)
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	search, err := faceSDK.Search()
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	r, _, err := search.SetFace(picPath, "image_file").
		SetFaceSet(config.CFG.FacesetToken, "faceset_token").
		SetOption("return_result_count", 1).End()
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if len(r.Results) == 0 {
		return c.JSON(http.StatusBadRequest, "用户不存在")
	}
	// 和万分之一做比较
	e4 := r.Thresholds["1e-4"]
	if r.Results[0].Confidence < e4 {
		return c.JSON(http.StatusBadRequest, "未识别出用户")
	}
	// 查询用户信息
	user := new(model.UserModel)
	err = db.DB.Where("face_token = ?", r.Results[0].FaceToken).Limit(1).Find(user).Error
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user.Password = ""

	return c.JSON(http.StatusOK, user)
}

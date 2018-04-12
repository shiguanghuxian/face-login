package public

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/logger"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	"github.com/shiguanghuxian/face-login/internal/common"
	"github.com/shiguanghuxian/face-login/internal/config"
	"github.com/shiguanghuxian/face-login/internal/db"
	"github.com/shiguanghuxian/face-login/internal/face"
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
	var faceToken string
	if config.CFG.FaceType == "face++" {
		faceToken, err = face.SearchFaceFaceToken(picPath)
	} else if config.CFG.FaceType == "seeta" {
		faceToken, err = face.SearchSeetaFaceToken(picPath)
	} else {
		return c.JSON(http.StatusBadRequest, "服务端未配置人脸检测方式")
	}
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// 查询用户信息
	user := new(model.UserModel)
	err = db.DB.Where("face_token = ?", faceToken).Limit(1).Find(user).Error
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user.Password = ""

	return c.JSON(http.StatusOK, user)
}

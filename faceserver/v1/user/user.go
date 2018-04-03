package user

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	sdk "git.53it.net/face-golang-sdk"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	"github.com/shiguanghuxian/face-login/internal/common"
	"github.com/shiguanghuxian/face-login/internal/config"
	"github.com/shiguanghuxian/face-login/internal/db"
	"github.com/shiguanghuxian/face-login/model"
	"github.com/shiguanghuxian/logger"
)

// var (
// 	APIKey    = "HiVoOzKxm9kRLcKG2ZaJS1I0414P2TJ5"
// 	APISecret = "RgGYzg3g4iuyw4B4t2z2pTvobj7KYHhl"
// )

// UserController 用户管理，即FaceSet管理
type UserController struct {
}

// AddUser 添加用户
func (uc *UserController) AddUser(c echo.Context) error {
	user := new(model.UserModel)
	user.Username = c.FormValue("username")
	if user.Username == "" {
		return c.JSON(http.StatusBadRequest, "用户名不能为空")
	}
	user.Password = common.UserPwdEncrypt(c.FormValue("password"))

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
	savePath := fmt.Sprintf("%sfaces%s%s.%s", pathSeparator, pathSeparator, fileName, postfix)
	picPath := fmt.Sprintf("%s%spublic%s", common.GetRootDir(), pathSeparator, savePath)
	// 创建保存文件
	destFile, err := os.Create(picPath)
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	defer destFile.Close()

	// log.Println(picPath)

	// 读取表单文件，写入保存文件
	_, err = io.Copy(destFile, formFile)
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user.FaceUrl = savePath
	user.FaceToken = ""
	// 如果未添加到数据库，则删除图片
	defer func() {
		log.Println(user.Id)
		if user.FaceToken == "" || user.Id == 0 {
			os.Remove(picPath)
		}
	}()

	// 上传脸部信息到face++
	// faceset_token a6f362d1977068f590ebec924114cd39   d497255b9c0dddc971356e2552ce1c81
	faceSDK, err := sdk.NewFaceSDK(config.CFG.APIKey, config.CFG.APISecret, config.CFG.Debug)
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// 人脸检测
	detect, err := faceSDK.Detect()
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	dr, _, err := detect.SetImage(picPath, "image_file").End()
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if len(dr.Faces) == 0 {
		return c.JSON(http.StatusBadRequest, "未检测到人脸信息")
	}
	if len(dr.Faces) > 1 {
		return c.JSON(http.StatusBadRequest, "请保证人脸照片中只包含一个人脸")
	}

	// 添加到faceset
	faceSet, err := faceSDK.FaceSet(map[string]interface{}{
		"faceset_token": config.CFG.FacesetToken,
		"face_tokens":   dr.Faces[0].FaceToken,
	})
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	cr, _, err := faceSet.AddFace().End()
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	crData := cr.(*sdk.FaceSetAddFaceFaceResponse)
	user.FaceToken = dr.Faces[0].FaceToken
	user.FacesetToken = crData.FacesetToken
	user.CreateTime = time.Now().Unix()

	err = db.DB.Create(user).Error
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

// UserList 用户列表
func (uc *UserController) UserList(c echo.Context) error {
	list := make([]*model.UserModel, 0)
	err := db.DB.Model(&model.UserModel{}).Order("id desc").Scan(&list).Error
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, list)
}

// DelUser 删除用户
func (uc *UserController) DelUser(c echo.Context) error {
	id := c.FormValue("id")
	// 查询用户信息
	user := new(model.UserModel)
	err := db.DB.Where("id = ?", id).Find(user).Limit(1).Error
	if err != nil {
		logger.Errorln(err)
		if err.Error() == "record not found" {
			return c.JSON(http.StatusBadRequest, "用户信息不存在")
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if user.Id == 0 {
		return c.JSON(http.StatusBadRequest, "用户信息不存在")
	}

	// 删除face++
	faceSDK, err := sdk.NewFaceSDK(config.CFG.APIKey, config.CFG.APISecret, config.CFG.Debug)
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	faceSet, err := faceSDK.FaceSet()
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	_, b, err := faceSet.RemoveFace().SetOptionMap(map[string]interface{}{
		"faceset_token": user.FacesetToken,
		"face_tokens":   user.FaceToken,
	}).End()
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	log.Println(b)
	// 删除用户
	err = db.DB.Where("id = ?", id).Delete(model.UserModel{}).Error
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	pathSeparator := string(os.PathSeparator)
	picPath := fmt.Sprintf("%s%spublic%s", common.GetRootDir(), pathSeparator, user.FaceUrl)
	os.Remove(picPath)

	return c.JSON(http.StatusOK, "ok")
}

// DelAll 删除全部用户
func (uc *UserController) DelAll(c echo.Context) error {
	// 删除用户
	err := db.DB.Delete(model.UserModel{}).Error
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// 清空face_token
	// 删除face++
	faceSDK, err := sdk.NewFaceSDK(config.CFG.APIKey, config.CFG.APISecret, config.CFG.Debug)
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	faceSet, err := faceSDK.FaceSet()
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	_, _, err = faceSet.RemoveFace().SetOptionMap(map[string]interface{}{
		"faceset_token": config.CFG.FacesetToken,
		"face_tokens":   "RemoveAllFaceTokens",
	}).End()
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// 删除所有图片
	pathSeparator := string(os.PathSeparator)
	picPath := fmt.Sprintf("%s%spublic%sfaces", common.GetRootDir(), pathSeparator, pathSeparator)
	err = os.RemoveAll(picPath)
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	err = os.MkdirAll(picPath, os.ModePerm)
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "ok")
}

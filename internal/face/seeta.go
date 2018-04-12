package face

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	sh "github.com/codeskyblue/go-sh"
	"github.com/shiguanghuxian/face-login/internal/common"
	"github.com/shiguanghuxian/face-login/internal/config"
	"github.com/shiguanghuxian/face-login/internal/db"
	"github.com/shiguanghuxian/face-login/model"
)

/* 本地人脸识别操作函数 */

// GetSeetaFaceCount 获取人脸数
func GetSeetaFaceCount(picPath string) (int, string, error) {
	log.Println(picPath)
	// 执行命令
	c, err := RunCmd("detect", picPath)
	if err != nil {
		return 0, "", err
	}
	cc, err := strconv.Atoi(c)
	if err != nil {
		return 0, "", err
	}
	// 获取文件md5
	md5, err := common.GetFileMd5(picPath)
	if err != nil {
		return 0, "", err
	}
	return cc, md5, nil
}

var wg sync.WaitGroup // 定义一个同步等待的组
type result struct {
	Per float64
	Key int
}

// SearchSeetaFaceToken 查找人员
func SearchSeetaFaceToken(picPath string) (string, error) {
	log.Println("本地验证用户")
	// 查询人员列表
	list := make([]*model.UserModel, 0)
	err := db.DB.Model(&model.UserModel{}).Order("id desc").Scan(&list).Error
	if err != nil {
		return "", err
	}
	// 比对结果
	results := make([]*result, 0)
	// 遍历所有用户对比
	for k, v := range list {
		v := v
		wg.Add(1)
		go func(key int, faceUrl string) {
			defer wg.Done()
			faceUrl = fmt.Sprintf("%s/public%s", common.GetRootDir(), faceUrl)
			val, err := RunCmd("compare", picPath, faceUrl)
			log.Println(val)
			if err == nil {
				fv, err := strconv.ParseFloat(val, 64)
				if err == nil {
					results = append(results, &result{
						Per: fv,
						Key: key,
					})
				}
			}
		}(k, v.FaceUrl)
	}
	wg.Wait()
	vvv, _ := json.Marshal(results)
	log.Println(string(vvv))

	if len(results) == 0 {
		return "", errors.New("未查询到相同人脸信息")
	}
	var per float64 = 0
	var key int = -1
	// 查找可能性最大的人脸
	for _, v := range results {
		if v.Per > 0.5 && v.Per > per {
			per = v.Per
			key = v.Key
		}
	}
	if per == 0 || key == -1 {
		return "", errors.New("未查询到相同人脸信息")
	}
	log.Println(key)
	// 用户人脸文件路径
	faceUrl := fmt.Sprintf("%s/public%s", common.GetRootDir(), list[key].FaceUrl)
	// 获取文件md5
	md5, err := common.GetFileMd5(faceUrl)
	if err != nil {
		return "", err
	}
	return md5, nil
}

// RunCmd 执行cmd
func RunCmd(args ...interface{}) (string, error) {
	// 执行的指令
	command := "./face"
	// 创建一个ssh session
	session := sh.NewSession()
	session.ShowCMD = false
	session.SetDir(config.CFG.Cmd)
	// session.SetEnv("PATH", appPath)
	c, err := session.Command(command, args...).Output()
	session.Command("exit").Run()
	if err != nil {
		return "", err
	}
	log.Println(string(c))
	log.Println(strings.TrimSpace(string(c)))
	return strings.TrimSpace(string(c)), nil
}

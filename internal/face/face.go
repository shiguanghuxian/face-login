package face

import (
	"errors"

	sdk "github.com/shiguanghuxian/face-golang-sdk"
	"github.com/shiguanghuxian/face-login/internal/config"
)

/* 人脸识别相关处理函数 */

// GetFaceFaceCount 获取人脸数 Face++
func GetFaceFaceCount(picPath string) (int, string, error) {
	// 上传脸部信息到face++
	// faceset_token a6f362d1977068f590ebec924114cd39   d497255b9c0dddc971356e2552ce1c81
	faceSDK, err := sdk.NewFaceSDK(config.CFG.APIKey, config.CFG.APISecret, config.CFG.Debug)
	if err != nil {
		return 0, "", err
	}
	// 人脸检测
	detect, err := faceSDK.Detect()
	if err != nil {
		return 0, "", err
	}
	dr, _, err := detect.SetImage(picPath, "image_file").End()
	if err != nil {
		return 0, "", err
	}

	faceToken := ""
	if len(dr.Faces) > 0 {
		faceToken = dr.Faces[0].FaceToken
	}

	return len(dr.Faces), faceToken, nil
}

// AddFaceTokenToFaceSet 添加到FaceSet Face++
func AddFaceTokenToFaceSet(faceToken string) (string, error) {
	faceSDK, err := sdk.NewFaceSDK(config.CFG.APIKey, config.CFG.APISecret, config.CFG.Debug)
	if err != nil {
		return "", err
	}
	// 添加到faceset
	faceSet, err := faceSDK.FaceSet(map[string]interface{}{
		"faceset_token": config.CFG.FacesetToken,
		"face_tokens":   faceToken,
	})
	if err != nil {
		return "", err
	}
	cr, _, err := faceSet.AddFace().End()
	if err != nil {
		return "", err
	}
	crData := cr.(*sdk.FaceSetAddFaceFaceResponse)

	return crData.FacesetToken, nil
}

// RemoveFaceFace 删除FaceSet中人脸
func RemoveFaceFace(faceToken, facesetToken string) error {
	// 删除face++
	faceSDK, err := sdk.NewFaceSDK(config.CFG.APIKey, config.CFG.APISecret, config.CFG.Debug)
	if err != nil {
		return err
	}
	faceSet, err := faceSDK.FaceSet()
	if err != nil {
		return err
	}
	_, _, err = faceSet.RemoveFace().SetOptionMap(map[string]interface{}{
		"faceset_token": facesetToken,
		"face_tokens":   faceToken,
	}).End()
	// log.Println(b)
	return err
}

// SearchFaceFaceToken 搜索人脸
func SearchFaceFaceToken(picPath string) (string, error) {
	// 验证用户脸部信息
	faceSDK, err := sdk.NewFaceSDK(config.CFG.APIKey, config.CFG.APISecret, config.CFG.Debug)
	if err != nil {
		return "", err
	}
	search, err := faceSDK.Search()
	if err != nil {
		return "", err
	}
	r, _, err := search.SetFace(picPath, "image_file").
		SetFaceSet(config.CFG.FacesetToken, "faceset_token").
		SetOption("return_result_count", 1).End()
	if err != nil {
		return "", err
	}
	if len(r.Results) == 0 {
		return "", errors.New("用户不存在")
	}
	// 和万分之一做比较
	e4 := r.Thresholds["1e-4"]
	if r.Results[0].Confidence < e4 {
		return "", errors.New("未识别出用户")
	}
	return r.Results[0].FaceToken, nil
}

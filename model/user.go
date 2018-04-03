package model

import "github.com/jinzhu/gorm"

// UserModel 用户表
type UserModel struct {
	Id           int    `gorm:"column:id;primary_key" json:"id"`
	Username     string `gorm:"column:username" json:"username"`
	Password     string `gorm:"column:password" json:"password"`
	FaceToken    string `gorm:"column:face_token" json:"face_token"`
	FaceUrl      string `gorm:"column:face_url" json:"face_url"`
	FacesetToken string `gorm:"column:faceset_token" json:"faceset_token"`
	CreateTime   int64  `gorm:"column:create_time" json:"create_time"`
}

// TableName 返回asc_door 表名称
func (UserModel) TableName() string {
	return gorm.DefaultTableNameHandler(nil, "user")
}

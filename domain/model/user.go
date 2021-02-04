package model

type User struct {
	//主键
	ID int64 `gorm:"primaryKey;not null;autoIncrement"`
	//用户名称
	UserName string `gorm:"uniqueIndex;not null;size:255;"`
	//
	FirstName string
	//密码
	HashPassword string
}

package auth

import "github.com/isd-sgcu/rnkm65-auth/src/app/model"

type Auth struct {
	model.Base
	StudentID    string `json:"student_id" gorm:"index:,unique"`
	Year         string `json:"year" gorm:"type:tinytext"`
	Faculty      string `json:"faculty" gorm:"type:tinytext"`
	UserID       string `json:"user_id" gorm:"index:,unique"`
	RefreshToken string `json:"refresh_token" gorm:"index"`
}

package dtos

import "time"

// RegisterInput 用户注册的输入结构
type RegisterInput struct {
	Username         string    `json:"Username" binding:"required,min=3,max=20"` // 用户名
	Email            string    `json:"Email" binding:"required,email"`           // 邮箱
	Password         string    `json:"Password" binding:"required,min=6,max=20"` // 密码
	FirstName        string    `json:"FirstName" binding:"required"`             // 名字
	LastName         string    `json:"LastName" binding:"required"`              // 姓氏
	ProfilePicture   string    `json:"ProfilePicture,omitempty"`                 // 头像（选填）
	DateOfBirth      time.Time `json:"DateOfBirth,omitempty"`                    // 出生日期（选填）
	Gender           string    `json:"Gender,omitempty"`                         // 性别（选填）
	FieldOfExpertise string    `json:"FieldOfExpertise,omitempty"`               // 专业领域（选填）
	EducationLevel   string    `json:"EducationLevel,omitempty"`                 // 受教育水平（选填）
	Skills           []string  `json:"Skills,omitempty"`                         // 技能（选填）
}

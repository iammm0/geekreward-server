package dtos

// LoginInput 用于用户登录的输入结构
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"` // 邮箱
	Password string `json:"password" binding:"required"`    // 密码
}

package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword 使用 bcrypt 生成密码的哈希值
func HashPassword(password string) (string, error) {
	// 使用 bcrypt 默认成本生成密码哈希
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash 验证密码是否与哈希匹配
func CheckPasswordHash(password, hash string) bool {
	// 使用 bcrypt 进行密码哈希验证
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

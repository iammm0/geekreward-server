package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
)

// SaveProfilePicture 保存用户头像到指定路径，并返回存储的文件路径
func SaveProfilePicture(file *multipart.FileHeader, uploadDir string) (string, error) {
	// 获取文件名并拼接成存储路径
	filename := filepath.Base(file.Filename)
	path := filepath.Join(uploadDir, filename)

	// 如果目录不存在则创建
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("无法创建上传目录: %v", err)
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("无法打开上传的文件: %v", err)
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("无法创建目标文件: %v", err)
	}
	defer dst.Close()

	// 保存文件内容到目标路径
	if _, err := dst.ReadFrom(src); err != nil {
		return "", fmt.Errorf("无法保存文件内容: %v", err)
	}

	return path, nil
}

package services

import (
	"GeekReward/inernal/app/models/dtos"
	"GeekReward/inernal/app/models/tables"
	"GeekReward/inernal/app/repositories"
	"GeekReward/pkg/utils"
	"errors"
	"gorm.io/gorm"
	"mime/multipart"
	"os"
)

type AuthService interface {
	Register(input dtos.RegisterInput, file *multipart.FileHeader) (*tables.User, error)
	Login(input dtos.LoginInput) (string, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(input dtos.RegisterInput, file *multipart.FileHeader) (*tables.User, error) {
	// 检查邮箱是否已被使用
	if _, err := s.userRepo.FindByEmail(input.Email); err == nil {
		return nil, errors.New("email already in use")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 检查用户名是否已被使用
	if _, err := s.userRepo.FindByUsername(input.Username); err == nil {
		return nil, errors.New("username already in use")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 保存头像文件
	var filePath string
	if file != nil {
		uploadDir := "uploads/avatars"
		var err error
		filePath, err = utils.SaveProfilePicture(file, uploadDir)
		if err != nil {
			return nil, err
		}
	}

	// 哈希用户密码
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &tables.User{
		Username:         input.Username,
		Email:            input.Email,
		Password:         hashedPassword,
		FirstName:        input.FirstName,
		LastName:         input.LastName,
		DateOfBirth:      input.DateOfBirth,
		Gender:           input.Gender,
		ProfilePicture:   filePath,
		FieldOfExpertise: input.FieldOfExpertise,
		EducationLevel:   input.EducationLevel,
		Skills:           input.Skills,
	}

	if err := s.userRepo.Create(user); err != nil {
		if filePath != "" {
			os.Remove(filePath)
		}
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(input dtos.LoginInput) (string, error) {
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil || !utils.CheckPasswordHash(input.Password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

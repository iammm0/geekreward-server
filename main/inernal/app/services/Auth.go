package services

import (
	dtos2 "GeekReward/main/inernal/app/models/dtos"
	"GeekReward/main/inernal/app/models/tables"
	"GeekReward/main/inernal/app/repositories"
	utils2 "GeekReward/main/pkg/utils"
	"errors"
	"gorm.io/gorm"
	"os"
)

type AuthService interface {
	Register(input dtos2.RegisterInput) (*tables.User, error)
	Login(input dtos2.LoginInput) (string, *tables.User, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(input dtos2.RegisterInput) (*tables.User, error) {
	// 1. 检查邮箱是否已被使用
	if _, err := s.userRepo.FindByEmail(input.Email); err == nil {
		return nil, errors.New("邮箱已被注册")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 2. 检查用户名是否已被使用
	if _, err := s.userRepo.FindByUsername(input.Username); err == nil {
		return nil, errors.New("用户名已被使用")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 3. 哈希密码
	hashedPassword, err := utils2.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	// 4. 创建用户模型
	user := &tables.User{
		Username:         input.Username,
		Email:            input.Email,
		Password:         hashedPassword,
		FirstName:        input.FirstName,
		LastName:         input.LastName,
		DateOfBirth:      input.DateOfBirth,
		Gender:           input.Gender,
		FieldOfExpertise: input.FieldOfExpertise,
		EducationLevel:   input.EducationLevel,
		Skills:           input.Skills,
		ProfilePicture:   input.ProfilePicture, // 直接使用Controller传来的URL
	}

	// 5. 创建数据库记录
	if err := s.userRepo.CreateUser(user); err != nil {
		// 如果数据库写入失败 & 用户上传了头像, 可能考虑删除文件
		if input.ProfilePicture != "" {
			// 这里的 input.ProfilePicture 是 "/uploads/avatars/xxx.jpg"
			// 如果想删除物理文件，需要组合 ./uploads/avatars/xxx.jpg
			filePath := "." + input.ProfilePicture // 变成 ./uploads/avatars/xxx.jpg
			_ = os.Remove(filePath)
		}
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(input dtos2.LoginInput) (string, *tables.User, error) {
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return "", nil, errors.New("用户不存在或数据库错误")
	}

	// 校验密码
	if !utils2.CheckPasswordHash(input.Password, user.Password) {
		return "", nil, errors.New("密码不正确")
	}

	// 生成JWT
	token, err := utils2.GenerateJWT(user.ID) // 传入 user.ID
	if err != nil {
		return "", nil, err
	}

	// 4. 返回 token, user
	return token, user, nil
}

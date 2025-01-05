package services

import (
	"GeekReward/inernal/app/models/dtos"
	"GeekReward/inernal/app/models/tables"
	"GeekReward/inernal/app/repositories"
	"github.com/google/uuid"
	"time"
)

type UserService interface {
	GetUserByID(id uuid.UUID) (*tables.User, error)
	UpdateUser(id uuid.UUID, input dtos.UpdateUserProfile) (*tables.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetUserByID(id uuid.UUID) (*tables.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *userService) UpdateUser(id uuid.UUID, input dtos.UpdateUserProfile) (*tables.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}
	if input.LastName != "" {
		user.LastName = input.LastName
	}
	if input.DateOfBirth != "" {
		dateOfBirth, err := time.Parse("2006-01-02", input.DateOfBirth)
		if err == nil {
			user.DateOfBirth = dateOfBirth
		}
	}
	user.Gender = input.Gender
	user.PhoneNumber = input.PhoneNumber
	user.Address = input.Address
	user.City = input.City
	user.State = input.State
	user.Country = input.Country
	user.PostalCode = input.PostalCode
	user.Institution = input.Institution
	user.Department = input.Department
	user.JobTitle = input.JobTitle
	user.EducationLevel = input.EducationLevel
	user.FieldOfExpertise = input.FieldOfExpertise
	user.YearsOfExperience = input.YearsOfExperience
	user.ProfilePicture = input.ProfilePicture
	user.Biography = input.Biography
	user.GitHubProfile = input.GitHubProfile
	user.Goals = input.Goals
	user.Skills = input.Skills
	user.Interests = input.Interests
	user.Languages = input.Languages
	user.Certifications = input.Certifications
	user.Awards = input.Awards
	user.Publications = input.Publications
	user.Projects = input.Projects
	user.Hobbies = input.Hobbies
	user.SocialMediaHandles = input.SocialMediaHandles
	user.Preferences = input.Preferences
	user.EmailPreferences = input.EmailPreferences
	user.ContactPreferences = input.ContactPreferences
	user.SecurityQuestions = input.SecurityQuestions
	user.TwoFactorEnabled = input.TwoFactorEnabled
	user.Timezone = input.Timezone
	user.PreferredLanguage = input.PreferredLanguage
	user.SolvedCount = input.SolvedCount
	user.MaxDifficulty = input.MaxDifficulty
	user.Reputation = input.Reputation

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

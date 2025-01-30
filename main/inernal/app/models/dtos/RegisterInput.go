package dtos

import "time"

type RegisterInput struct {
	Username         string    `form:"username"`
	Email            string    `form:"email"`
	Password         string    `form:"password"`
	FirstName        string    `form:"first_name"`
	LastName         string    `form:"last_name"`
	DateOfBirth      time.Time `form:"date_of_birth"`
	Gender           string    `form:"gender"`
	FieldOfExpertise string    `form:"field_of_expertise"`
	EducationLevel   string    `form:"education_level"`
	Skills           []string  `form:"skills"`

	ProfilePicture string `form:"profilePicture"`
}

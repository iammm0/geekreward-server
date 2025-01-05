package dtos

import "time"

type RegisterInput struct {
	Username         string    `json:"username"`
	Email            string    `json:"email"`
	Password         string    `json:"password"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	DateOfBirth      time.Time `json:"date_of_birth"`
	Gender           string    `json:"gender"`
	FieldOfExpertise string    `json:"field_of_expertise"`
	EducationLevel   string    `json:"education_level"`
	Skills           []string  `json:"skills"`
}

package dtos

import (
	"GeekReward/inernal/app/models/tables"
	"github.com/lib/pq"
)

type UpdateUserProfile struct {
	FirstName                string                  `json:"firstName"`
	LastName                 string                  `json:"lastName"`
	DateOfBirth              string                  `json:"dateOfBirth"` // 使用字符串格式，可以在服务端解析为时间格式
	Gender                   string                  `json:"gender"`
	PhoneNumber              string                  `json:"phoneNumber"`
	Address                  string                  `json:"address"`
	City                     string                  `json:"city"`
	State                    string                  `json:"state"`
	Country                  string                  `json:"country"`
	PostalCode               string                  `json:"postalCode"`
	Institution              string                  `json:"institution"`
	Department               string                  `json:"department"`
	JobTitle                 string                  `json:"jobTitle"`
	EducationLevel           string                  `json:"educationLevel"`
	FieldOfExpertise         string                  `json:"fieldOfExpertise"`
	YearsOfExperience        int                     `json:"yearsOfExperience"`
	ProfilePicture           string                  `json:"profilePicture"`
	Biography                string                  `json:"biography"`
	GitHubProfile            string                  `json:"gitHubProfile"`
	Website                  string                  `json:"website"`
	Goals                    string                  `json:"goals"`
	Skills                   pq.StringArray          `json:"skills"`
	Interests                pq.StringArray          `json:"interests"`
	Languages                pq.StringArray          `json:"languages"`
	Certifications           pq.StringArray          `json:"certifications"`
	Awards                   pq.StringArray          `json:"awards"`
	Publications             pq.StringArray          `json:"publications"`
	Projects                 pq.StringArray          `json:"projects"`
	Hobbies                  pq.StringArray          `json:"hobbies"`
	SocialMediaHandles       map[string]string       `json:"socialMediaHandles"`
	ProfessionalAffiliations string                  `json:"professionalAffiliations"`
	Memberships              pq.StringArray          `json:"memberships"`
	Preferences              map[string]string       `json:"preferences"`
	EmailPreferences         map[string]bool         `json:"emailPreferences"`
	ContactPreferences       map[string]bool         `json:"contactPreferences"`
	SecurityQuestions        map[string]string       `json:"securityQuestions"`
	TwoFactorEnabled         bool                    `json:"twoFactorEnabled"`
	Timezone                 string                  `json:"timezone"`
	PreferredLanguage        string                  `json:"preferredLanguage"`
	WorkExperience           []tables.WorkExperience `json:"workExperience"`
	SolvedCount              int                     `json:"solvedCount"`
	MaxDifficulty            string                  `json:"maxDifficulty"`
	Reputation               float64                 `json:"reputation"`
}

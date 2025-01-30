package tables

import (
	"github.com/lib/pq"
	"time"
)

type User struct {
	BaseModel
	Username             string    `gorm:"uniqueIndex;not null"`
	Email                string    `gorm:"uniqueIndex;not null"`
	Password             string    `gorm:"not null"`
	LastLogin            time.Time `gorm:"index"`
	AccountStatus        string    `gorm:"default:'active'"` // "active", "suspended", "deleted"
	Verified             bool      `gorm:"default:false"`
	NotificationsEnabled bool      `gorm:"default:true"`
	ProfilePicture       string    `gorm:"type:text"`

	Preferences        map[string]string `gorm:"type:jsonb"`
	EmailPreferences   map[string]bool   `gorm:"type:jsonb"`
	ContactPreferences map[string]bool   `gorm:"type:jsonb"`
	SecurityQuestions  map[string]string `gorm:"type:jsonb"`
	TwoFactorEnabled   bool              `gorm:"default:false"`
	LoginAttempts      int               `gorm:"default:0"`
	LastPasswordChange time.Time
	Timezone           string `gorm:"default:'UTC'"`
	PreferredLanguage  string `gorm:"default:'en'"`

	FirstName         string
	LastName          string
	DateOfBirth       time.Time
	Gender            string
	PhoneNumber       string
	Address           string
	City              string
	State             string
	Country           string
	PostalCode        string
	Institution       string
	Department        string
	JobTitle          string
	EducationLevel    string
	FieldOfExpertise  string
	YearsOfExperience int
	Biography         string
	GitHubProfile     string
	Goals             string
	Bio               string

	Skills             pq.StringArray    `gorm:"type:text[]"`
	Interests          pq.StringArray    `gorm:"type:text[]"`
	Languages          pq.StringArray    `gorm:"type:text[]"`
	Certifications     pq.StringArray    `gorm:"type:text[]"`
	Awards             pq.StringArray    `gorm:"type:text[]"`
	Publications       pq.StringArray    `gorm:"type:text[]"`
	Projects           pq.StringArray    `gorm:"type:text[]"`
	Hobbies            pq.StringArray    `gorm:"type:text[]"`
	SocialMediaHandles map[string]string `gorm:"type:jsonb"`
	WorkExperience     []WorkExperience  `gorm:"foreignKey:UserID;references:ID"`

	SolvedCount   int     `gorm:"default:0"`
	MaxDifficulty string  `gorm:"default:'medium'"`
	Reputation    float64 `gorm:"default:0"`

	Bounties     []Bounty      `gorm:"foreignKey:UserID;references:ID"`
	Applications []Application `gorm:"foreignKey:UserID;references:ID"`
	Comments     []Comment     `gorm:"foreignKey:UserID;references:ID"`
	Likes        []Like        `gorm:"foreignKey:UserID;references:ID"`
	Ratings      []Rating      `gorm:"foreignKey:UserID;references:ID"`
}

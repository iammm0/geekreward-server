package tables

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username             string    `gorm:"uniqueIndex;not null"`
	Email                string    `gorm:"uniqueIndex;not null"`
	Password             string    `gorm:"not null"`
	LastLogin            time.Time `gorm:"index"`
	AccountStatus        string    `gorm:"default:'active'"` // "active", "suspended", "deleted"
	Verified             bool      `gorm:"default:false"`
	NotificationsEnabled bool      `gorm:"default:true"`

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
	ProfilePicture    string
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
	WorkExperience     []WorkExperience  `gorm:"foreignKey:UserID"`

	SolvedCount   int     `gorm:"default:0"`
	MaxDifficulty string  `gorm:"default:'medium'"`
	Reputation    float64 `gorm:"default:0"`

	Bounties     []Bounty      `gorm:"foreignKey:UserID"`
	Applications []Application `gorm:"foreignKey:UserID"`
	Comments     []Comment     `gorm:"foreignKey:UserID"`
	Likes        []Like        `gorm:"foreignKey:UserID"`
	Ratings      []Rating      `gorm:"foreignKey:UserID"`
}

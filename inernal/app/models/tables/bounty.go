package tables

import (
	"GeekReward/inernal/app/models/common"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

type Bounty struct {
	gorm.Model
	UserID     common.UserID  `gorm:"not null;index"`
	ReceiverID *common.UserID `gorm:"index"`

	Title           string         `gorm:"not null"`
	Description     string         `gorm:"not null"`
	Reward          float64        `gorm:"not null"`
	Status          string         `gorm:"not null;default:'open';index"` // "open", "closed", "completed"
	Deadline        time.Time      `gorm:"index"`
	DifficultyLevel string         // "easy", "medium", "hard"
	Category        string         // 可通过预定义类别或数据库表进行规范
	Tags            pq.StringArray `gorm:"type:text[]"`
	Location        string
	AttachmentURLs  pq.StringArray `gorm:"type:text[]"`

	Anonymous               bool           `gorm:"default:false"`
	Priority                string         `gorm:"default:'normal'"` // "low", "normal", "high"
	PaymentStatus           string         `gorm:"default:'unpaid'"` // "unpaid", "paid"
	PreferredSolutionType   string         // "document", "code", etc.
	RequiredSkills          pq.StringArray `gorm:"type:text[]"`
	RequiredExperience      int            // 可选，年限
	RequiredCertifications  pq.StringArray `gorm:"type:text[]"`
	Visibility              string         `gorm:"default:'public'"` // "public", "private"
	Confidentiality         string         `gorm:"default:'non-confidential'"`
	ContractType            string         // "fixed", "hourly"
	EstimatedHours          float64
	ToolsRequired           pq.StringArray `gorm:"type:text[]"`
	CommunicationPreference string
	FeedbackRequired        bool `gorm:"default:false"`
	CompletionCriteria      string
	SubmissionGuidelines    string
	EvaluationCriteria      string
	ReferenceMaterials      string
	ExternalLinks           pq.StringArray `gorm:"type:text[]"`
	AdditionalNotes         string
	NDARequired             bool `gorm:"default:false"`
	AcceptanceCriteria      string
	PaymentMethod           string // "PayPal", "Bank Transfer", etc.

	ActualHours   float64
	LikesCount    int `gorm:"default:0"`
	CommentsCount int `gorm:"default:0"`
	ViewCount     int `gorm:"default:0"`
	AverageRating float64

	Milestones   []Milestone   `gorm:"foreignKey:BountyID"`
	Comments     []Comment     `gorm:"foreignKey:BountyID"`
	Applications []Application `gorm:"foreignKey:BountyID"`
	Likes        []Like        `gorm:"foreignKey:BountyID"`
	Ratings      []Rating      `gorm:"foreignKey:BountyID"`
}

package dtos

import (
	"github.com/google/uuid"
)

type BountyDTO struct {
	Title                   string    `json:"title"`
	Description             string    `json:"description"`
	Reward                  float64   `json:"reward"`
	Deadline                string    `json:"deadline"` // 使用字符串形式传递日期
	DifficultyLevel         string    `json:"difficulty_level"`
	Category                string    `json:"category"`
	Tags                    []string  `json:"tags"`
	Location                string    `json:"location"`
	AttachmentUrls          []string  `json:"attachment_urls"`
	Anonymous               bool      `json:"anonymous"`
	Priority                string    `json:"priority"`
	PaymentStatus           string    `json:"payment_status"`
	PreferredSolutionType   string    `json:"preferred_solution_type"`
	RequiredSkills          []string  `json:"required_skills"`
	RequiredExperience      int       `json:"required_experience"`
	RequiredCertifications  []string  `json:"required_certifications"`
	Visibility              string    `json:"visibility"`
	Confidentiality         string    `json:"confidentiality"`
	ContractType            string    `json:"contract_type"`
	EstimatedHours          float64   `json:"estimated_hours"`
	ToolsRequired           []string  `json:"tools_required"`
	CommunicationPreference string    `json:"communication_preference"`
	FeedbackRequired        bool      `json:"feedback_required"`
	CompletionCriteria      string    `json:"completion_criteria"`
	SubmissionGuidelines    string    `json:"submission_guidelines"`
	EvaluationCriteria      string    `json:"evaluation_criteria"`
	ReferenceMaterials      string    `json:"reference_materials"`
	ExternalLinks           []string  `json:"external_links"`
	AdditionalNotes         string    `json:"additional_notes"`
	NdaRequired             bool      `json:"nda_required"`
	AcceptanceCriteria      string    `json:"acceptance_criteria"`
	PaymentMethod           string    `json:"payment_method"`
	UserID                  uuid.UUID `json:"user_id"`
}

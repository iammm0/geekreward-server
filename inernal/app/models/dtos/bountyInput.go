package dtos

import "github.com/lib/pq"

type BountyDTO struct {
	Title                   string         `json:"Title" binding:"required"`
	Description             string         `json:"Description" binding:"required"`
	Reward                  float64        `json:"Reward" binding:"required"`
	Deadline                string         `json:"Deadline" binding:"required"` // 使用字符串格式，可以在服务端解析为时间格式
	DifficultyLevel         string         `json:"DifficultyLevel"`
	Category                string         `json:"Category"`
	Tags                    pq.StringArray `json:"Tags"`
	Location                string         `json:"Location"`
	AttachmentUrls          pq.StringArray `json:"AttachmentUrls"`
	Anonymous               bool           `json:"Anonymous"`
	Priority                string         `json:"Priority"`
	Budget                  float64        `json:"Budget"`
	PaymentStatus           string         `json:"PaymentStatus"`
	PreferredSolutionType   string         `json:"PreferredSolutionType"`
	RequiredSkills          pq.StringArray `json:"RequiredSkills"`
	RequiredExperience      int            `json:"RequiredExperience"`
	RequiredCertifications  pq.StringArray `json:"RequiredCertifications"`
	Visibility              string         `json:"Visibility"`
	Confidentiality         string         `json:"Confidentiality"`
	ContractType            string         `json:"ContractType"`
	EstimatedHours          float64        `json:"EstimatedHours"`
	ActualHours             float64        `json:"ActualHours"`
	ToolsRequired           pq.StringArray `json:"ToolsRequired"`
	CommunicationPreference string         `json:"CommunicationPreference"`
	FeedbackRequired        bool           `json:"FeedbackRequired"`
	CompletionCriteria      string         `json:"CompletionCriteria"`
	SubmissionGuidelines    string         `json:"SubmissionGuidelines"`
	EvaluationCriteria      string         `json:"EvaluationCriteria"`
	ReferenceMaterials      string         `json:"ReferenceMaterials"`
	ExternalLinks           pq.StringArray `json:"ExternalLinks"`
	AdditionalNotes         string         `json:"AdditionalNotes"`
	NdaRequired             bool           `json:"NdaRequired"`
	AcceptanceCriteria      string         `json:"AcceptanceCriteria"`
	PaymentMethod           string         `json:"PaymentMethod"`
}

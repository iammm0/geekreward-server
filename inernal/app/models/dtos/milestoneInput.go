package dtos

type MilestoneDTO struct {
	Title       string `json:"title" binding:"required"`       // 里程碑标题
	Description string `json:"description" binding:"required"` // 里程碑描述
	DueDate     string `json:"due_date" binding:"required"`    // 截止日期 (格式: YYYY-MM-DD)
	BountyID    uint   `json:"bounty_id" binding:"required"`   // 关联的悬赏令ID
}

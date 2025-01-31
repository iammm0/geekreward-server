package dtos

import (
	"GeekReward/inernal/app/models/tables"
	"github.com/google/uuid"
)

// BountyFilter Status, PublisherID, ReceiverID 都是指针，以便分辨 “有没有传” 与 “传了空” 的区别。
// 若只想查看某一方的悬赏令，可以设置 PublisherID=xxxx 或 ReceiverID=xxxx。
// 若您只想要 status + 分页，可以不定义 PublisherID, ReceiverID。
type BountyFilter struct {
	Status      *tables.BountyStatus // 可选状态
	PublisherID *uuid.UUID           // 可选发布者ID
	ReceiverID  *uuid.UUID           // 可选接收者ID
	Limit       int
	Offset      int
}

package entity

import (
	"github.com/google/uuid"
)

// PolicyCategory ポリシーカテゴリを表す構造体。
type PolicyCategory struct {
	PolicyCategoryID uuid.UUID `json:"policy_category_id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Key              string    `json:"key"`
}

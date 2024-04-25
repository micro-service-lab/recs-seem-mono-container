package entity

import (
	"github.com/google/uuid"
)

// Policy ポリシーを表す構造体。
type Policy struct {
	PolicyID         uuid.UUID `json:"policy_id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Key              string    `json:"key"`
	PolicyCategoryID uuid.UUID `json:"policy_category_id"`
}

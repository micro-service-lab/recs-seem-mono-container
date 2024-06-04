package entity

import "github.com/google/uuid"

// Group グループを表す構造体。
type Group struct {
	GroupID        uuid.UUID `json:"group_id"`
	Key            string    `json:"key"`
	OrganizationID uuid.UUID `json:"organization_id"`
}

// GroupWithOrganization グループと組織を表す構造体。
type GroupWithOrganization struct {
	GroupID      uuid.UUID    `json:"group_id"`
	Key          string       `json:"key"`
	Organization Organization `json:"organization"`
}

// GroupWithOrganizationForQuery グループと組織を表す構造体(クエリー用)。
type GroupWithOrganizationForQuery struct {
	Pkey Int `json:"-"`
	GroupWithOrganization
}

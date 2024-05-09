// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: group.sql

package query

import (
	"context"

	"github.com/google/uuid"
)

const countGroups = `-- name: CountGroups :one
SELECT COUNT(*) FROM m_groups
`

func (q *Queries) CountGroups(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countGroups)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createGroup = `-- name: CreateGroup :one
INSERT INTO m_groups (key, organization_id) VALUES ($1, $2) RETURNING m_groups_pkey, group_id, key, organization_id
`

type CreateGroupParams struct {
	Key            string    `json:"key"`
	OrganizationID uuid.UUID `json:"organization_id"`
}

func (q *Queries) CreateGroup(ctx context.Context, arg CreateGroupParams) (Group, error) {
	row := q.db.QueryRow(ctx, createGroup, arg.Key, arg.OrganizationID)
	var i Group
	err := row.Scan(
		&i.MGroupsPkey,
		&i.GroupID,
		&i.Key,
		&i.OrganizationID,
	)
	return i, err
}

type CreateGroupsParams struct {
	Key            string    `json:"key"`
	OrganizationID uuid.UUID `json:"organization_id"`
}

const deleteGroup = `-- name: DeleteGroup :exec
DELETE FROM m_groups WHERE group_id = $1
`

func (q *Queries) DeleteGroup(ctx context.Context, groupID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteGroup, groupID)
	return err
}

const deleteGroupByKey = `-- name: DeleteGroupByKey :exec
DELETE FROM m_groups WHERE key = $1
`

func (q *Queries) DeleteGroupByKey(ctx context.Context, key string) error {
	_, err := q.db.Exec(ctx, deleteGroupByKey, key)
	return err
}

const findGroupByID = `-- name: FindGroupByID :one
SELECT m_groups_pkey, group_id, key, organization_id FROM m_groups WHERE group_id = $1
`

func (q *Queries) FindGroupByID(ctx context.Context, groupID uuid.UUID) (Group, error) {
	row := q.db.QueryRow(ctx, findGroupByID, groupID)
	var i Group
	err := row.Scan(
		&i.MGroupsPkey,
		&i.GroupID,
		&i.Key,
		&i.OrganizationID,
	)
	return i, err
}

const findGroupByIDWithOrganization = `-- name: FindGroupByIDWithOrganization :one
SELECT m_groups.m_groups_pkey, m_groups.group_id, m_groups.key, m_groups.organization_id, m_organizations.m_organizations_pkey, m_organizations.organization_id, m_organizations.name, m_organizations.description, m_organizations.color, m_organizations.is_personal, m_organizations.is_whole, m_organizations.created_at, m_organizations.updated_at, m_organizations.chat_room_id FROM m_groups
LEFT JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
WHERE group_id = $1
`

type FindGroupByIDWithOrganizationRow struct {
	Group        Group        `json:"group"`
	Organization Organization `json:"organization"`
}

func (q *Queries) FindGroupByIDWithOrganization(ctx context.Context, groupID uuid.UUID) (FindGroupByIDWithOrganizationRow, error) {
	row := q.db.QueryRow(ctx, findGroupByIDWithOrganization, groupID)
	var i FindGroupByIDWithOrganizationRow
	err := row.Scan(
		&i.Group.MGroupsPkey,
		&i.Group.GroupID,
		&i.Group.Key,
		&i.Group.OrganizationID,
		&i.Organization.MOrganizationsPkey,
		&i.Organization.OrganizationID,
		&i.Organization.Name,
		&i.Organization.Description,
		&i.Organization.Color,
		&i.Organization.IsPersonal,
		&i.Organization.IsWhole,
		&i.Organization.CreatedAt,
		&i.Organization.UpdatedAt,
		&i.Organization.ChatRoomID,
	)
	return i, err
}

const findGroupByKey = `-- name: FindGroupByKey :one
SELECT m_groups_pkey, group_id, key, organization_id FROM m_groups WHERE key = $1
`

func (q *Queries) FindGroupByKey(ctx context.Context, key string) (Group, error) {
	row := q.db.QueryRow(ctx, findGroupByKey, key)
	var i Group
	err := row.Scan(
		&i.MGroupsPkey,
		&i.GroupID,
		&i.Key,
		&i.OrganizationID,
	)
	return i, err
}

const findGroupByKeyWithOrganization = `-- name: FindGroupByKeyWithOrganization :one
SELECT m_groups.m_groups_pkey, m_groups.group_id, m_groups.key, m_groups.organization_id, m_organizations.m_organizations_pkey, m_organizations.organization_id, m_organizations.name, m_organizations.description, m_organizations.color, m_organizations.is_personal, m_organizations.is_whole, m_organizations.created_at, m_organizations.updated_at, m_organizations.chat_room_id FROM m_groups
LEFT JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
WHERE key = $1
`

type FindGroupByKeyWithOrganizationRow struct {
	Group        Group        `json:"group"`
	Organization Organization `json:"organization"`
}

func (q *Queries) FindGroupByKeyWithOrganization(ctx context.Context, key string) (FindGroupByKeyWithOrganizationRow, error) {
	row := q.db.QueryRow(ctx, findGroupByKeyWithOrganization, key)
	var i FindGroupByKeyWithOrganizationRow
	err := row.Scan(
		&i.Group.MGroupsPkey,
		&i.Group.GroupID,
		&i.Group.Key,
		&i.Group.OrganizationID,
		&i.Organization.MOrganizationsPkey,
		&i.Organization.OrganizationID,
		&i.Organization.Name,
		&i.Organization.Description,
		&i.Organization.Color,
		&i.Organization.IsPersonal,
		&i.Organization.IsWhole,
		&i.Organization.CreatedAt,
		&i.Organization.UpdatedAt,
		&i.Organization.ChatRoomID,
	)
	return i, err
}

const getGroups = `-- name: GetGroups :many
SELECT m_groups_pkey, group_id, key, organization_id FROM m_groups
ORDER BY
	m_groups_pkey ASC
`

func (q *Queries) GetGroups(ctx context.Context) ([]Group, error) {
	rows, err := q.db.Query(ctx, getGroups)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Group{}
	for rows.Next() {
		var i Group
		if err := rows.Scan(
			&i.MGroupsPkey,
			&i.GroupID,
			&i.Key,
			&i.OrganizationID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getGroupsUseKeysetPaginate = `-- name: GetGroupsUseKeysetPaginate :many
SELECT m_groups_pkey, group_id, key, organization_id FROM m_groups
WHERE
	CASE $2::text
		WHEN 'next' THEN
			m_groups_pkey > $3::int
		WHEN 'prev' THEN
			m_groups_pkey < $3::int
	END
ORDER BY
	CASE WHEN $2::text = 'next' THEN m_groups_pkey END ASC,
	CASE WHEN $2::text = 'prev' THEN m_groups_pkey END DESC
LIMIT $1
`

type GetGroupsUseKeysetPaginateParams struct {
	Limit           int32  `json:"limit"`
	CursorDirection string `json:"cursor_direction"`
	Cursor          int32  `json:"cursor"`
}

func (q *Queries) GetGroupsUseKeysetPaginate(ctx context.Context, arg GetGroupsUseKeysetPaginateParams) ([]Group, error) {
	rows, err := q.db.Query(ctx, getGroupsUseKeysetPaginate, arg.Limit, arg.CursorDirection, arg.Cursor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Group{}
	for rows.Next() {
		var i Group
		if err := rows.Scan(
			&i.MGroupsPkey,
			&i.GroupID,
			&i.Key,
			&i.OrganizationID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getGroupsUseNumberedPaginate = `-- name: GetGroupsUseNumberedPaginate :many
SELECT m_groups_pkey, group_id, key, organization_id FROM m_groups
ORDER BY
	m_groups_pkey ASC
LIMIT $1 OFFSET $2
`

type GetGroupsUseNumberedPaginateParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetGroupsUseNumberedPaginate(ctx context.Context, arg GetGroupsUseNumberedPaginateParams) ([]Group, error) {
	rows, err := q.db.Query(ctx, getGroupsUseNumberedPaginate, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Group{}
	for rows.Next() {
		var i Group
		if err := rows.Scan(
			&i.MGroupsPkey,
			&i.GroupID,
			&i.Key,
			&i.OrganizationID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getGroupsWithOrganization = `-- name: GetGroupsWithOrganization :many
SELECT m_groups.m_groups_pkey, m_groups.group_id, m_groups.key, m_groups.organization_id, m_organizations.m_organizations_pkey, m_organizations.organization_id, m_organizations.name, m_organizations.description, m_organizations.color, m_organizations.is_personal, m_organizations.is_whole, m_organizations.created_at, m_organizations.updated_at, m_organizations.chat_room_id FROM m_groups
LEFT JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
ORDER BY
	CASE WHEN $1::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN $1::text = 'r_name' THEN m_organizations.name END DESC,
	m_groups_pkey ASC
`

type GetGroupsWithOrganizationRow struct {
	Group        Group        `json:"group"`
	Organization Organization `json:"organization"`
}

func (q *Queries) GetGroupsWithOrganization(ctx context.Context, orderMethod string) ([]GetGroupsWithOrganizationRow, error) {
	rows, err := q.db.Query(ctx, getGroupsWithOrganization, orderMethod)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetGroupsWithOrganizationRow{}
	for rows.Next() {
		var i GetGroupsWithOrganizationRow
		if err := rows.Scan(
			&i.Group.MGroupsPkey,
			&i.Group.GroupID,
			&i.Group.Key,
			&i.Group.OrganizationID,
			&i.Organization.MOrganizationsPkey,
			&i.Organization.OrganizationID,
			&i.Organization.Name,
			&i.Organization.Description,
			&i.Organization.Color,
			&i.Organization.IsPersonal,
			&i.Organization.IsWhole,
			&i.Organization.CreatedAt,
			&i.Organization.UpdatedAt,
			&i.Organization.ChatRoomID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getGroupsWithOrganizationUseKeysetPaginate = `-- name: GetGroupsWithOrganizationUseKeysetPaginate :many
SELECT m_groups.m_groups_pkey, m_groups.group_id, m_groups.key, m_groups.organization_id, m_organizations.m_organizations_pkey, m_organizations.organization_id, m_organizations.name, m_organizations.description, m_organizations.color, m_organizations.is_personal, m_organizations.is_whole, m_organizations.created_at, m_organizations.updated_at, m_organizations.chat_room_id FROM m_groups
LEFT JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
WHERE
	CASE $2::text
		WHEN 'next' THEN
			CASE $3::text
				WHEN 'name' THEN name > $4 OR (name = $4 AND m_groups_pkey > $5::int)
				WHEN 'r_name' THEN name < $4 OR (name = $4 AND m_groups_pkey > $5::int)
				ELSE m_groups_pkey > $5::int
			END
		WHEN 'prev' THEN
			CASE $3::text
				WHEN 'name' THEN name > $4 OR (name = $4 AND m_groups_pkey > $5::int)
				WHEN 'r_name' THEN name < $4 OR (name = $4 AND m_groups_pkey > $5::int)
				ELSE m_groups_pkey > $5::int
			END
	END
ORDER BY
	CASE WHEN $3::text = 'name' AND $2::text = 'next' THEN m_organizations.name END ASC,
	CASE WHEN $3::text = 'name' AND $2::text = 'prev' THEN m_organizations.name END DESC,
	CASE WHEN $3::text = 'r_name' AND $2::text = 'next' THEN m_organizations.name END ASC,
	CASE WHEN $3::text = 'r_name' AND $2::text = 'prev' THEN m_organizations.name END DESC,
	CASE WHEN $2::text = 'next' THEN m_groups_pkey END ASC,
	CASE WHEN $2::text = 'prev' THEN m_groups_pkey END DESC
LIMIT $1
`

type GetGroupsWithOrganizationUseKeysetPaginateParams struct {
	Limit           int32  `json:"limit"`
	CursorDirection string `json:"cursor_direction"`
	OrderMethod     string `json:"order_method"`
	NameCursor      string `json:"name_cursor"`
	Cursor          int32  `json:"cursor"`
}

type GetGroupsWithOrganizationUseKeysetPaginateRow struct {
	Group        Group        `json:"group"`
	Organization Organization `json:"organization"`
}

func (q *Queries) GetGroupsWithOrganizationUseKeysetPaginate(ctx context.Context, arg GetGroupsWithOrganizationUseKeysetPaginateParams) ([]GetGroupsWithOrganizationUseKeysetPaginateRow, error) {
	rows, err := q.db.Query(ctx, getGroupsWithOrganizationUseKeysetPaginate,
		arg.Limit,
		arg.CursorDirection,
		arg.OrderMethod,
		arg.NameCursor,
		arg.Cursor,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetGroupsWithOrganizationUseKeysetPaginateRow{}
	for rows.Next() {
		var i GetGroupsWithOrganizationUseKeysetPaginateRow
		if err := rows.Scan(
			&i.Group.MGroupsPkey,
			&i.Group.GroupID,
			&i.Group.Key,
			&i.Group.OrganizationID,
			&i.Organization.MOrganizationsPkey,
			&i.Organization.OrganizationID,
			&i.Organization.Name,
			&i.Organization.Description,
			&i.Organization.Color,
			&i.Organization.IsPersonal,
			&i.Organization.IsWhole,
			&i.Organization.CreatedAt,
			&i.Organization.UpdatedAt,
			&i.Organization.ChatRoomID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getGroupsWithOrganizationUseNumberedPaginate = `-- name: GetGroupsWithOrganizationUseNumberedPaginate :many
SELECT m_groups.m_groups_pkey, m_groups.group_id, m_groups.key, m_groups.organization_id, m_organizations.m_organizations_pkey, m_organizations.organization_id, m_organizations.name, m_organizations.description, m_organizations.color, m_organizations.is_personal, m_organizations.is_whole, m_organizations.created_at, m_organizations.updated_at, m_organizations.chat_room_id FROM m_groups
LEFT JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
ORDER BY
	CASE WHEN $3::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN $3::text = 'r_name' THEN m_organizations.name END DESC,
	m_groups_pkey ASC
LIMIT $1 OFFSET $2
`

type GetGroupsWithOrganizationUseNumberedPaginateParams struct {
	Limit       int32  `json:"limit"`
	Offset      int32  `json:"offset"`
	OrderMethod string `json:"order_method"`
}

type GetGroupsWithOrganizationUseNumberedPaginateRow struct {
	Group        Group        `json:"group"`
	Organization Organization `json:"organization"`
}

func (q *Queries) GetGroupsWithOrganizationUseNumberedPaginate(ctx context.Context, arg GetGroupsWithOrganizationUseNumberedPaginateParams) ([]GetGroupsWithOrganizationUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getGroupsWithOrganizationUseNumberedPaginate, arg.Limit, arg.Offset, arg.OrderMethod)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetGroupsWithOrganizationUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetGroupsWithOrganizationUseNumberedPaginateRow
		if err := rows.Scan(
			&i.Group.MGroupsPkey,
			&i.Group.GroupID,
			&i.Group.Key,
			&i.Group.OrganizationID,
			&i.Organization.MOrganizationsPkey,
			&i.Organization.OrganizationID,
			&i.Organization.Name,
			&i.Organization.Description,
			&i.Organization.Color,
			&i.Organization.IsPersonal,
			&i.Organization.IsWhole,
			&i.Organization.CreatedAt,
			&i.Organization.UpdatedAt,
			&i.Organization.ChatRoomID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPluralGroups = `-- name: GetPluralGroups :many
SELECT m_groups_pkey, group_id, key, organization_id FROM m_groups
WHERE organization_id = ANY($3::uuid[])
ORDER BY
	m_groups_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralGroupsParams struct {
	Limit           int32       `json:"limit"`
	Offset          int32       `json:"offset"`
	OrganizationIds []uuid.UUID `json:"organization_ids"`
}

func (q *Queries) GetPluralGroups(ctx context.Context, arg GetPluralGroupsParams) ([]Group, error) {
	rows, err := q.db.Query(ctx, getPluralGroups, arg.Limit, arg.Offset, arg.OrganizationIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Group{}
	for rows.Next() {
		var i Group
		if err := rows.Scan(
			&i.MGroupsPkey,
			&i.GroupID,
			&i.Key,
			&i.OrganizationID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPluralGroupsWithOrganization = `-- name: GetPluralGroupsWithOrganization :many
SELECT m_groups.m_groups_pkey, m_groups.group_id, m_groups.key, m_groups.organization_id, m_organizations.m_organizations_pkey, m_organizations.organization_id, m_organizations.name, m_organizations.description, m_organizations.color, m_organizations.is_personal, m_organizations.is_whole, m_organizations.created_at, m_organizations.updated_at, m_organizations.chat_room_id FROM m_groups
LEFT JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
WHERE group_id = ANY($3::uuid[])
ORDER BY
	m_groups_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralGroupsWithOrganizationParams struct {
	Limit    int32       `json:"limit"`
	Offset   int32       `json:"offset"`
	GroupIds []uuid.UUID `json:"group_ids"`
}

type GetPluralGroupsWithOrganizationRow struct {
	Group        Group        `json:"group"`
	Organization Organization `json:"organization"`
}

func (q *Queries) GetPluralGroupsWithOrganization(ctx context.Context, arg GetPluralGroupsWithOrganizationParams) ([]GetPluralGroupsWithOrganizationRow, error) {
	rows, err := q.db.Query(ctx, getPluralGroupsWithOrganization, arg.Limit, arg.Offset, arg.GroupIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralGroupsWithOrganizationRow{}
	for rows.Next() {
		var i GetPluralGroupsWithOrganizationRow
		if err := rows.Scan(
			&i.Group.MGroupsPkey,
			&i.Group.GroupID,
			&i.Group.Key,
			&i.Group.OrganizationID,
			&i.Organization.MOrganizationsPkey,
			&i.Organization.OrganizationID,
			&i.Organization.Name,
			&i.Organization.Description,
			&i.Organization.Color,
			&i.Organization.IsPersonal,
			&i.Organization.IsWhole,
			&i.Organization.CreatedAt,
			&i.Organization.UpdatedAt,
			&i.Organization.ChatRoomID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const pluralDeleteGroups = `-- name: PluralDeleteGroups :exec
DELETE FROM m_groups WHERE group_id = ANY($1::uuid[])
`

func (q *Queries) PluralDeleteGroups(ctx context.Context, dollar_1 []uuid.UUID) error {
	_, err := q.db.Exec(ctx, pluralDeleteGroups, dollar_1)
	return err
}

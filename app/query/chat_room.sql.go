// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: chat_room.sql

package query

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const countChatRooms = `-- name: CountChatRooms :one
SELECT COUNT(*) FROM m_chat_rooms
WHERE
	CASE WHEN $1::boolean = true THEN owner_id = ANY($2) ELSE TRUE END
AND
	CASE WHEN $3::boolean = true THEN is_private = $4 ELSE TRUE END
AND
	CASE WHEN $5::boolean = true THEN from_organization = $6 ELSE TRUE END
AND
	CASE WHEN $7::boolean = true THEN (SELECT chat_room_id FROM m_organizations WHERE organization_id = ANY($8)) = chat_room_id ELSE TRUE END
`

type CountChatRoomsParams struct {
	WhereInOwner            bool        `json:"where_in_owner"`
	InOwner                 pgtype.UUID `json:"in_owner"`
	WhereIsPrivate          bool        `json:"where_is_private"`
	IsPrivate               bool        `json:"is_private"`
	WhereIsFromOrganization bool        `json:"where_is_from_organization"`
	IsFromOrganization      bool        `json:"is_from_organization"`
	WhereFromOrganizations  bool        `json:"where_from_organizations"`
	InOrganizations         uuid.UUID   `json:"in_organizations"`
}

func (q *Queries) CountChatRooms(ctx context.Context, arg CountChatRoomsParams) (int64, error) {
	row := q.db.QueryRow(ctx, countChatRooms,
		arg.WhereInOwner,
		arg.InOwner,
		arg.WhereIsPrivate,
		arg.IsPrivate,
		arg.WhereIsFromOrganization,
		arg.IsFromOrganization,
		arg.WhereFromOrganizations,
		arg.InOrganizations,
	)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createChatRoom = `-- name: CreateChatRoom :one
INSERT INTO m_chat_rooms (name, is_private, cover_image_url, owner_id, from_organization, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING m_chat_rooms_pkey, chat_room_id, name, is_private, cover_image_url, owner_id, from_organization, created_at, updated_at
`

type CreateChatRoomParams struct {
	Name             pgtype.Text `json:"name"`
	IsPrivate        bool        `json:"is_private"`
	CoverImageUrl    pgtype.Text `json:"cover_image_url"`
	OwnerID          pgtype.UUID `json:"owner_id"`
	FromOrganization bool        `json:"from_organization"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
}

func (q *Queries) CreateChatRoom(ctx context.Context, arg CreateChatRoomParams) (ChatRoom, error) {
	row := q.db.QueryRow(ctx, createChatRoom,
		arg.Name,
		arg.IsPrivate,
		arg.CoverImageUrl,
		arg.OwnerID,
		arg.FromOrganization,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i ChatRoom
	err := row.Scan(
		&i.MChatRoomsPkey,
		&i.ChatRoomID,
		&i.Name,
		&i.IsPrivate,
		&i.CoverImageUrl,
		&i.OwnerID,
		&i.FromOrganization,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

type CreateChatRoomsParams struct {
	Name             pgtype.Text `json:"name"`
	IsPrivate        bool        `json:"is_private"`
	CoverImageUrl    pgtype.Text `json:"cover_image_url"`
	OwnerID          pgtype.UUID `json:"owner_id"`
	FromOrganization bool        `json:"from_organization"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
}

const deleteChatRoom = `-- name: DeleteChatRoom :exec
DELETE FROM m_chat_rooms WHERE chat_room_id = $1
`

func (q *Queries) DeleteChatRoom(ctx context.Context, chatRoomID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteChatRoom, chatRoomID)
	return err
}

const findChatRoomByID = `-- name: FindChatRoomByID :one
SELECT m_chat_rooms_pkey, chat_room_id, name, is_private, cover_image_url, owner_id, from_organization, created_at, updated_at FROM m_chat_rooms WHERE chat_room_id = $1
`

func (q *Queries) FindChatRoomByID(ctx context.Context, chatRoomID uuid.UUID) (ChatRoom, error) {
	row := q.db.QueryRow(ctx, findChatRoomByID, chatRoomID)
	var i ChatRoom
	err := row.Scan(
		&i.MChatRoomsPkey,
		&i.ChatRoomID,
		&i.Name,
		&i.IsPrivate,
		&i.CoverImageUrl,
		&i.OwnerID,
		&i.FromOrganization,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findChatRoomByIDWithOwner = `-- name: FindChatRoomByIDWithOwner :one
SELECT m_chat_rooms.m_chat_rooms_pkey, m_chat_rooms.chat_room_id, m_chat_rooms.name, m_chat_rooms.is_private, m_chat_rooms.cover_image_url, m_chat_rooms.owner_id, m_chat_rooms.from_organization, m_chat_rooms.created_at, m_chat_rooms.updated_at, m_members.m_members_pkey, m_members.member_id, m_members.login_id, m_members.password, m_members.email, m_members.name, m_members.attend_status_id, m_members.profile_image_url, m_members.grade_id, m_members.group_id, m_members.personal_organization_id, m_members.role_id, m_members.created_at, m_members.updated_at FROM m_chat_rooms
LEFT JOIN m_members ON m_chat_rooms.owner_id = m_members.member_id
WHERE chat_room_id = $1
`

type FindChatRoomByIDWithOwnerRow struct {
	ChatRoom ChatRoom `json:"chat_room"`
	Member   Member   `json:"member"`
}

func (q *Queries) FindChatRoomByIDWithOwner(ctx context.Context, chatRoomID uuid.UUID) (FindChatRoomByIDWithOwnerRow, error) {
	row := q.db.QueryRow(ctx, findChatRoomByIDWithOwner, chatRoomID)
	var i FindChatRoomByIDWithOwnerRow
	err := row.Scan(
		&i.ChatRoom.MChatRoomsPkey,
		&i.ChatRoom.ChatRoomID,
		&i.ChatRoom.Name,
		&i.ChatRoom.IsPrivate,
		&i.ChatRoom.CoverImageUrl,
		&i.ChatRoom.OwnerID,
		&i.ChatRoom.FromOrganization,
		&i.ChatRoom.CreatedAt,
		&i.ChatRoom.UpdatedAt,
		&i.Member.MMembersPkey,
		&i.Member.MemberID,
		&i.Member.LoginID,
		&i.Member.Password,
		&i.Member.Email,
		&i.Member.Name,
		&i.Member.AttendStatusID,
		&i.Member.ProfileImageUrl,
		&i.Member.GradeID,
		&i.Member.GroupID,
		&i.Member.PersonalOrganizationID,
		&i.Member.RoleID,
		&i.Member.CreatedAt,
		&i.Member.UpdatedAt,
	)
	return i, err
}

const getChatRooms = `-- name: GetChatRooms :many
SELECT m_chat_rooms_pkey, chat_room_id, name, is_private, cover_image_url, owner_id, from_organization, created_at, updated_at FROM m_chat_rooms
WHERE
	CASE WHEN $1::boolean = true THEN owner_id = ANY($2) ELSE TRUE END
AND
	CASE WHEN $3::boolean = true THEN is_private = $4 ELSE TRUE END
AND
	CASE WHEN $5::boolean = true THEN from_organization = $6 ELSE TRUE END
AND
	CASE WHEN $7::boolean = true THEN (SELECT chat_room_id FROM m_organizations WHERE organization_id = ANY($8)) = chat_room_id ELSE TRUE END
ORDER BY
	m_chat_rooms_pkey ASC
`

type GetChatRoomsParams struct {
	WhereInOwner            bool        `json:"where_in_owner"`
	InOwner                 pgtype.UUID `json:"in_owner"`
	WhereIsPrivate          bool        `json:"where_is_private"`
	IsPrivate               bool        `json:"is_private"`
	WhereIsFromOrganization bool        `json:"where_is_from_organization"`
	IsFromOrganization      bool        `json:"is_from_organization"`
	WhereFromOrganizations  bool        `json:"where_from_organizations"`
	InOrganizations         uuid.UUID   `json:"in_organizations"`
}

func (q *Queries) GetChatRooms(ctx context.Context, arg GetChatRoomsParams) ([]ChatRoom, error) {
	rows, err := q.db.Query(ctx, getChatRooms,
		arg.WhereInOwner,
		arg.InOwner,
		arg.WhereIsPrivate,
		arg.IsPrivate,
		arg.WhereIsFromOrganization,
		arg.IsFromOrganization,
		arg.WhereFromOrganizations,
		arg.InOrganizations,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ChatRoom{}
	for rows.Next() {
		var i ChatRoom
		if err := rows.Scan(
			&i.MChatRoomsPkey,
			&i.ChatRoomID,
			&i.Name,
			&i.IsPrivate,
			&i.CoverImageUrl,
			&i.OwnerID,
			&i.FromOrganization,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getChatRoomsUseKeysetPaginate = `-- name: GetChatRoomsUseKeysetPaginate :many
SELECT m_chat_rooms_pkey, chat_room_id, name, is_private, cover_image_url, owner_id, from_organization, created_at, updated_at FROM m_chat_rooms
WHERE
	CASE WHEN $2::boolean = true THEN owner_id = ANY($3) ELSE TRUE END
AND
	CASE WHEN $4::boolean = true THEN is_private = $5 ELSE TRUE END
AND
	CASE WHEN $6::boolean = true THEN from_organization = $7 ELSE TRUE END
AND
	CASE WHEN $8::boolean = true THEN (SELECT chat_room_id FROM m_organizations WHERE organization_id = ANY($9)) = chat_room_id ELSE TRUE END
AND
	CASE $10::text
		WHEN 'next' THEN
			m_chat_rooms_pkey > $11::int
		WHEN 'prev' THEN
			m_chat_rooms_pkey < $11::int
	END
ORDER BY
	CASE WHEN $10::text = 'next' THEN m_chat_rooms_pkey END ASC,
	CASE WHEN $10::text = 'prev' THEN m_chat_rooms_pkey END DESC
LIMIT $1
`

type GetChatRoomsUseKeysetPaginateParams struct {
	Limit                   int32       `json:"limit"`
	WhereInOwner            bool        `json:"where_in_owner"`
	InOwner                 pgtype.UUID `json:"in_owner"`
	WhereIsPrivate          bool        `json:"where_is_private"`
	IsPrivate               bool        `json:"is_private"`
	WhereIsFromOrganization bool        `json:"where_is_from_organization"`
	IsFromOrganization      bool        `json:"is_from_organization"`
	WhereFromOrganizations  bool        `json:"where_from_organizations"`
	InOrganizations         uuid.UUID   `json:"in_organizations"`
	CursorDirection         string      `json:"cursor_direction"`
	Cursor                  int32       `json:"cursor"`
}

func (q *Queries) GetChatRoomsUseKeysetPaginate(ctx context.Context, arg GetChatRoomsUseKeysetPaginateParams) ([]ChatRoom, error) {
	rows, err := q.db.Query(ctx, getChatRoomsUseKeysetPaginate,
		arg.Limit,
		arg.WhereInOwner,
		arg.InOwner,
		arg.WhereIsPrivate,
		arg.IsPrivate,
		arg.WhereIsFromOrganization,
		arg.IsFromOrganization,
		arg.WhereFromOrganizations,
		arg.InOrganizations,
		arg.CursorDirection,
		arg.Cursor,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ChatRoom{}
	for rows.Next() {
		var i ChatRoom
		if err := rows.Scan(
			&i.MChatRoomsPkey,
			&i.ChatRoomID,
			&i.Name,
			&i.IsPrivate,
			&i.CoverImageUrl,
			&i.OwnerID,
			&i.FromOrganization,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getChatRoomsUseNumberedPaginate = `-- name: GetChatRoomsUseNumberedPaginate :many
SELECT m_chat_rooms_pkey, chat_room_id, name, is_private, cover_image_url, owner_id, from_organization, created_at, updated_at FROM m_chat_rooms
WHERE
	CASE WHEN $3::boolean = true THEN owner_id = ANY($4) ELSE TRUE END
AND
	CASE WHEN $5::boolean = true THEN is_private = $6 ELSE TRUE END
AND
	CASE WHEN $7::boolean = true THEN from_organization = $8 ELSE TRUE END
AND
	CASE WHEN $9::boolean = true THEN (SELECT chat_room_id FROM m_organizations WHERE organization_id = ANY($10)) = chat_room_id ELSE TRUE END
ORDER BY
	m_chat_rooms_pkey ASC
LIMIT $1 OFFSET $2
`

type GetChatRoomsUseNumberedPaginateParams struct {
	Limit                   int32       `json:"limit"`
	Offset                  int32       `json:"offset"`
	WhereInOwner            bool        `json:"where_in_owner"`
	InOwner                 pgtype.UUID `json:"in_owner"`
	WhereIsPrivate          bool        `json:"where_is_private"`
	IsPrivate               bool        `json:"is_private"`
	WhereIsFromOrganization bool        `json:"where_is_from_organization"`
	IsFromOrganization      bool        `json:"is_from_organization"`
	WhereFromOrganizations  bool        `json:"where_from_organizations"`
	InOrganizations         uuid.UUID   `json:"in_organizations"`
}

func (q *Queries) GetChatRoomsUseNumberedPaginate(ctx context.Context, arg GetChatRoomsUseNumberedPaginateParams) ([]ChatRoom, error) {
	rows, err := q.db.Query(ctx, getChatRoomsUseNumberedPaginate,
		arg.Limit,
		arg.Offset,
		arg.WhereInOwner,
		arg.InOwner,
		arg.WhereIsPrivate,
		arg.IsPrivate,
		arg.WhereIsFromOrganization,
		arg.IsFromOrganization,
		arg.WhereFromOrganizations,
		arg.InOrganizations,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ChatRoom{}
	for rows.Next() {
		var i ChatRoom
		if err := rows.Scan(
			&i.MChatRoomsPkey,
			&i.ChatRoomID,
			&i.Name,
			&i.IsPrivate,
			&i.CoverImageUrl,
			&i.OwnerID,
			&i.FromOrganization,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getChatRoomsWithOwner = `-- name: GetChatRoomsWithOwner :many
SELECT m_chat_rooms.m_chat_rooms_pkey, m_chat_rooms.chat_room_id, m_chat_rooms.name, m_chat_rooms.is_private, m_chat_rooms.cover_image_url, m_chat_rooms.owner_id, m_chat_rooms.from_organization, m_chat_rooms.created_at, m_chat_rooms.updated_at, m_members.m_members_pkey, m_members.member_id, m_members.login_id, m_members.password, m_members.email, m_members.name, m_members.attend_status_id, m_members.profile_image_url, m_members.grade_id, m_members.group_id, m_members.personal_organization_id, m_members.role_id, m_members.created_at, m_members.updated_at FROM m_chat_rooms
LEFT JOIN m_members ON m_chat_rooms.owner_id = m_members.member_id
WHERE
	CASE WHEN $1::boolean THEN owner_id = ANY($2) ELSE TRUE END
AND
	CASE WHEN $3::boolean THEN is_private = $4 ELSE TRUE END
AND
	CASE WHEN $5::boolean = true THEN from_organization = $6 ELSE TRUE END
AND
	CASE WHEN $7::boolean = true THEN (SELECT chat_room_id FROM m_organizations WHERE organization_id = ANY($8)) = chat_room_id ELSE TRUE END
ORDER BY
	m_chat_rooms_pkey ASC
`

type GetChatRoomsWithOwnerParams struct {
	WhereInOwner            bool        `json:"where_in_owner"`
	InOwner                 pgtype.UUID `json:"in_owner"`
	WhereIsPrivate          bool        `json:"where_is_private"`
	IsPrivate               bool        `json:"is_private"`
	WhereIsFromOrganization bool        `json:"where_is_from_organization"`
	IsFromOrganization      bool        `json:"is_from_organization"`
	WhereFromOrganizations  bool        `json:"where_from_organizations"`
	InOrganizations         uuid.UUID   `json:"in_organizations"`
}

type GetChatRoomsWithOwnerRow struct {
	ChatRoom ChatRoom `json:"chat_room"`
	Member   Member   `json:"member"`
}

func (q *Queries) GetChatRoomsWithOwner(ctx context.Context, arg GetChatRoomsWithOwnerParams) ([]GetChatRoomsWithOwnerRow, error) {
	rows, err := q.db.Query(ctx, getChatRoomsWithOwner,
		arg.WhereInOwner,
		arg.InOwner,
		arg.WhereIsPrivate,
		arg.IsPrivate,
		arg.WhereIsFromOrganization,
		arg.IsFromOrganization,
		arg.WhereFromOrganizations,
		arg.InOrganizations,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetChatRoomsWithOwnerRow{}
	for rows.Next() {
		var i GetChatRoomsWithOwnerRow
		if err := rows.Scan(
			&i.ChatRoom.MChatRoomsPkey,
			&i.ChatRoom.ChatRoomID,
			&i.ChatRoom.Name,
			&i.ChatRoom.IsPrivate,
			&i.ChatRoom.CoverImageUrl,
			&i.ChatRoom.OwnerID,
			&i.ChatRoom.FromOrganization,
			&i.ChatRoom.CreatedAt,
			&i.ChatRoom.UpdatedAt,
			&i.Member.MMembersPkey,
			&i.Member.MemberID,
			&i.Member.LoginID,
			&i.Member.Password,
			&i.Member.Email,
			&i.Member.Name,
			&i.Member.AttendStatusID,
			&i.Member.ProfileImageUrl,
			&i.Member.GradeID,
			&i.Member.GroupID,
			&i.Member.PersonalOrganizationID,
			&i.Member.RoleID,
			&i.Member.CreatedAt,
			&i.Member.UpdatedAt,
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

const getChatRoomsWithOwnerUseKeysetPaginate = `-- name: GetChatRoomsWithOwnerUseKeysetPaginate :many
SELECT m_chat_rooms.m_chat_rooms_pkey, m_chat_rooms.chat_room_id, m_chat_rooms.name, m_chat_rooms.is_private, m_chat_rooms.cover_image_url, m_chat_rooms.owner_id, m_chat_rooms.from_organization, m_chat_rooms.created_at, m_chat_rooms.updated_at, m_members.m_members_pkey, m_members.member_id, m_members.login_id, m_members.password, m_members.email, m_members.name, m_members.attend_status_id, m_members.profile_image_url, m_members.grade_id, m_members.group_id, m_members.personal_organization_id, m_members.role_id, m_members.created_at, m_members.updated_at FROM m_chat_rooms
LEFT JOIN m_members ON m_chat_rooms.owner_id = m_members.member_id
WHERE
	CASE WHEN $2::boolean = true THEN owner_id = ANY($3) ELSE TRUE END
AND
	CASE WHEN $4::boolean = true THEN is_private = $5 ELSE TRUE END
AND
	CASE WHEN $6::boolean = true THEN from_organization = $7 ELSE TRUE END
AND
	CASE WHEN $8::boolean = true THEN (SELECT chat_room_id FROM m_organizations WHERE organization_id = ANY($9)) = chat_room_id ELSE TRUE END
AND
	CASE $10::text
		WHEN 'next' THEN
			m_chat_rooms_pkey > $11::int
		WHEN 'prev' THEN
			m_chat_rooms_pkey < $11::int
	END
ORDER BY
	CASE WHEN $10::text = 'next' THEN m_chat_rooms_pkey END ASC,
	CASE WHEN $10::text = 'prev' THEN m_chat_rooms_pkey END DESC
LIMIT $1
`

type GetChatRoomsWithOwnerUseKeysetPaginateParams struct {
	Limit                   int32       `json:"limit"`
	WhereInOwner            bool        `json:"where_in_owner"`
	InOwner                 pgtype.UUID `json:"in_owner"`
	WhereIsPrivate          bool        `json:"where_is_private"`
	IsPrivate               bool        `json:"is_private"`
	WhereIsFromOrganization bool        `json:"where_is_from_organization"`
	IsFromOrganization      bool        `json:"is_from_organization"`
	WhereFromOrganizations  bool        `json:"where_from_organizations"`
	InOrganizations         uuid.UUID   `json:"in_organizations"`
	CursorDirection         string      `json:"cursor_direction"`
	Cursor                  int32       `json:"cursor"`
}

type GetChatRoomsWithOwnerUseKeysetPaginateRow struct {
	ChatRoom ChatRoom `json:"chat_room"`
	Member   Member   `json:"member"`
}

func (q *Queries) GetChatRoomsWithOwnerUseKeysetPaginate(ctx context.Context, arg GetChatRoomsWithOwnerUseKeysetPaginateParams) ([]GetChatRoomsWithOwnerUseKeysetPaginateRow, error) {
	rows, err := q.db.Query(ctx, getChatRoomsWithOwnerUseKeysetPaginate,
		arg.Limit,
		arg.WhereInOwner,
		arg.InOwner,
		arg.WhereIsPrivate,
		arg.IsPrivate,
		arg.WhereIsFromOrganization,
		arg.IsFromOrganization,
		arg.WhereFromOrganizations,
		arg.InOrganizations,
		arg.CursorDirection,
		arg.Cursor,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetChatRoomsWithOwnerUseKeysetPaginateRow{}
	for rows.Next() {
		var i GetChatRoomsWithOwnerUseKeysetPaginateRow
		if err := rows.Scan(
			&i.ChatRoom.MChatRoomsPkey,
			&i.ChatRoom.ChatRoomID,
			&i.ChatRoom.Name,
			&i.ChatRoom.IsPrivate,
			&i.ChatRoom.CoverImageUrl,
			&i.ChatRoom.OwnerID,
			&i.ChatRoom.FromOrganization,
			&i.ChatRoom.CreatedAt,
			&i.ChatRoom.UpdatedAt,
			&i.Member.MMembersPkey,
			&i.Member.MemberID,
			&i.Member.LoginID,
			&i.Member.Password,
			&i.Member.Email,
			&i.Member.Name,
			&i.Member.AttendStatusID,
			&i.Member.ProfileImageUrl,
			&i.Member.GradeID,
			&i.Member.GroupID,
			&i.Member.PersonalOrganizationID,
			&i.Member.RoleID,
			&i.Member.CreatedAt,
			&i.Member.UpdatedAt,
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

const getChatRoomsWithOwnerUseNumberedPaginate = `-- name: GetChatRoomsWithOwnerUseNumberedPaginate :many
SELECT m_chat_rooms.m_chat_rooms_pkey, m_chat_rooms.chat_room_id, m_chat_rooms.name, m_chat_rooms.is_private, m_chat_rooms.cover_image_url, m_chat_rooms.owner_id, m_chat_rooms.from_organization, m_chat_rooms.created_at, m_chat_rooms.updated_at, m_members.m_members_pkey, m_members.member_id, m_members.login_id, m_members.password, m_members.email, m_members.name, m_members.attend_status_id, m_members.profile_image_url, m_members.grade_id, m_members.group_id, m_members.personal_organization_id, m_members.role_id, m_members.created_at, m_members.updated_at FROM m_chat_rooms
LEFT JOIN m_members ON m_chat_rooms.owner_id = m_members.member_id
WHERE
	CASE WHEN $3::boolean THEN owner_id = ANY($4) ELSE TRUE END
AND
	CASE WHEN $5::boolean THEN is_private = $6 ELSE TRUE END
AND
	CASE WHEN $7::boolean = true THEN from_organization = $8 ELSE TRUE END
AND
	CASE WHEN $9::boolean = true THEN (SELECT chat_room_id FROM m_organizations WHERE organization_id = ANY($10)) = chat_room_id ELSE TRUE END
ORDER BY
	m_chat_rooms_pkey ASC
LIMIT $1 OFFSET $2
`

type GetChatRoomsWithOwnerUseNumberedPaginateParams struct {
	Limit                   int32       `json:"limit"`
	Offset                  int32       `json:"offset"`
	WhereInOwner            bool        `json:"where_in_owner"`
	InOwner                 pgtype.UUID `json:"in_owner"`
	WhereIsPrivate          bool        `json:"where_is_private"`
	IsPrivate               bool        `json:"is_private"`
	WhereIsFromOrganization bool        `json:"where_is_from_organization"`
	IsFromOrganization      bool        `json:"is_from_organization"`
	WhereFromOrganizations  bool        `json:"where_from_organizations"`
	InOrganizations         uuid.UUID   `json:"in_organizations"`
}

type GetChatRoomsWithOwnerUseNumberedPaginateRow struct {
	ChatRoom ChatRoom `json:"chat_room"`
	Member   Member   `json:"member"`
}

func (q *Queries) GetChatRoomsWithOwnerUseNumberedPaginate(ctx context.Context, arg GetChatRoomsWithOwnerUseNumberedPaginateParams) ([]GetChatRoomsWithOwnerUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getChatRoomsWithOwnerUseNumberedPaginate,
		arg.Limit,
		arg.Offset,
		arg.WhereInOwner,
		arg.InOwner,
		arg.WhereIsPrivate,
		arg.IsPrivate,
		arg.WhereIsFromOrganization,
		arg.IsFromOrganization,
		arg.WhereFromOrganizations,
		arg.InOrganizations,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetChatRoomsWithOwnerUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetChatRoomsWithOwnerUseNumberedPaginateRow
		if err := rows.Scan(
			&i.ChatRoom.MChatRoomsPkey,
			&i.ChatRoom.ChatRoomID,
			&i.ChatRoom.Name,
			&i.ChatRoom.IsPrivate,
			&i.ChatRoom.CoverImageUrl,
			&i.ChatRoom.OwnerID,
			&i.ChatRoom.FromOrganization,
			&i.ChatRoom.CreatedAt,
			&i.ChatRoom.UpdatedAt,
			&i.Member.MMembersPkey,
			&i.Member.MemberID,
			&i.Member.LoginID,
			&i.Member.Password,
			&i.Member.Email,
			&i.Member.Name,
			&i.Member.AttendStatusID,
			&i.Member.ProfileImageUrl,
			&i.Member.GradeID,
			&i.Member.GroupID,
			&i.Member.PersonalOrganizationID,
			&i.Member.RoleID,
			&i.Member.CreatedAt,
			&i.Member.UpdatedAt,
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

const getPluralChatRooms = `-- name: GetPluralChatRooms :many
SELECT m_chat_rooms_pkey, chat_room_id, name, is_private, cover_image_url, owner_id, from_organization, created_at, updated_at FROM m_chat_rooms
WHERE chat_room_id = ANY($3::uuid[])
ORDER BY
	m_chat_rooms_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralChatRoomsParams struct {
	Limit       int32       `json:"limit"`
	Offset      int32       `json:"offset"`
	ChatRoomIds []uuid.UUID `json:"chat_room_ids"`
}

func (q *Queries) GetPluralChatRooms(ctx context.Context, arg GetPluralChatRoomsParams) ([]ChatRoom, error) {
	rows, err := q.db.Query(ctx, getPluralChatRooms, arg.Limit, arg.Offset, arg.ChatRoomIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ChatRoom{}
	for rows.Next() {
		var i ChatRoom
		if err := rows.Scan(
			&i.MChatRoomsPkey,
			&i.ChatRoomID,
			&i.Name,
			&i.IsPrivate,
			&i.CoverImageUrl,
			&i.OwnerID,
			&i.FromOrganization,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getPluralChatRoomsWithOwner = `-- name: GetPluralChatRoomsWithOwner :many
SELECT m_chat_rooms.m_chat_rooms_pkey, m_chat_rooms.chat_room_id, m_chat_rooms.name, m_chat_rooms.is_private, m_chat_rooms.cover_image_url, m_chat_rooms.owner_id, m_chat_rooms.from_organization, m_chat_rooms.created_at, m_chat_rooms.updated_at, m_members.m_members_pkey, m_members.member_id, m_members.login_id, m_members.password, m_members.email, m_members.name, m_members.attend_status_id, m_members.profile_image_url, m_members.grade_id, m_members.group_id, m_members.personal_organization_id, m_members.role_id, m_members.created_at, m_members.updated_at FROM m_chat_rooms
LEFT JOIN m_members ON m_chat_rooms.owner_id = m_members.member_id
WHERE chat_room_id = ANY($3::uuid[])
ORDER BY
	m_chat_rooms_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralChatRoomsWithOwnerParams struct {
	Limit       int32       `json:"limit"`
	Offset      int32       `json:"offset"`
	ChatRoomIds []uuid.UUID `json:"chat_room_ids"`
}

type GetPluralChatRoomsWithOwnerRow struct {
	ChatRoom ChatRoom `json:"chat_room"`
	Member   Member   `json:"member"`
}

func (q *Queries) GetPluralChatRoomsWithOwner(ctx context.Context, arg GetPluralChatRoomsWithOwnerParams) ([]GetPluralChatRoomsWithOwnerRow, error) {
	rows, err := q.db.Query(ctx, getPluralChatRoomsWithOwner, arg.Limit, arg.Offset, arg.ChatRoomIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralChatRoomsWithOwnerRow{}
	for rows.Next() {
		var i GetPluralChatRoomsWithOwnerRow
		if err := rows.Scan(
			&i.ChatRoom.MChatRoomsPkey,
			&i.ChatRoom.ChatRoomID,
			&i.ChatRoom.Name,
			&i.ChatRoom.IsPrivate,
			&i.ChatRoom.CoverImageUrl,
			&i.ChatRoom.OwnerID,
			&i.ChatRoom.FromOrganization,
			&i.ChatRoom.CreatedAt,
			&i.ChatRoom.UpdatedAt,
			&i.Member.MMembersPkey,
			&i.Member.MemberID,
			&i.Member.LoginID,
			&i.Member.Password,
			&i.Member.Email,
			&i.Member.Name,
			&i.Member.AttendStatusID,
			&i.Member.ProfileImageUrl,
			&i.Member.GradeID,
			&i.Member.GroupID,
			&i.Member.PersonalOrganizationID,
			&i.Member.RoleID,
			&i.Member.CreatedAt,
			&i.Member.UpdatedAt,
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

const pluralDeleteChatRooms = `-- name: PluralDeleteChatRooms :exec
DELETE FROM m_chat_rooms WHERE chat_room_id = ANY($1::uuid[])
`

func (q *Queries) PluralDeleteChatRooms(ctx context.Context, dollar_1 []uuid.UUID) error {
	_, err := q.db.Exec(ctx, pluralDeleteChatRooms, dollar_1)
	return err
}

const updateChatRoom = `-- name: UpdateChatRoom :one
UPDATE m_chat_rooms SET name = $2, is_private = $3, cover_image_url = $4, owner_id = $5, updated_at = $6 WHERE chat_room_id = $1 RETURNING m_chat_rooms_pkey, chat_room_id, name, is_private, cover_image_url, owner_id, from_organization, created_at, updated_at
`

type UpdateChatRoomParams struct {
	ChatRoomID    uuid.UUID   `json:"chat_room_id"`
	Name          pgtype.Text `json:"name"`
	IsPrivate     bool        `json:"is_private"`
	CoverImageUrl pgtype.Text `json:"cover_image_url"`
	OwnerID       pgtype.UUID `json:"owner_id"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

func (q *Queries) UpdateChatRoom(ctx context.Context, arg UpdateChatRoomParams) (ChatRoom, error) {
	row := q.db.QueryRow(ctx, updateChatRoom,
		arg.ChatRoomID,
		arg.Name,
		arg.IsPrivate,
		arg.CoverImageUrl,
		arg.OwnerID,
		arg.UpdatedAt,
	)
	var i ChatRoom
	err := row.Scan(
		&i.MChatRoomsPkey,
		&i.ChatRoomID,
		&i.Name,
		&i.IsPrivate,
		&i.CoverImageUrl,
		&i.OwnerID,
		&i.FromOrganization,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

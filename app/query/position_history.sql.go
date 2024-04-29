// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: position_history.sql

package query

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const countPositionHistories = `-- name: CountPositionHistories :one
SELECT COUNT(*) FROM t_position_histories
WHERE
	CASE WHEN $1::boolean = true THEN member_id = ANY($2::uuid[]) ELSE TRUE END
AND
	CASE WHEN $3::boolean = true THEN sent_at >= $4 ELSE TRUE END
AND
	CASE WHEN $5::boolean = true THEN sent_at <= $6 ELSE TRUE END
`

type CountPositionHistoriesParams struct {
	WhereInMember      bool        `json:"where_in_member"`
	InMemberIds        []uuid.UUID `json:"in_member_ids"`
	WhereEarlierSentAt bool        `json:"where_earlier_sent_at"`
	EarlierSentAt      time.Time   `json:"earlier_sent_at"`
	WhereLaterSentAt   bool        `json:"where_later_sent_at"`
	LaterSentAt        time.Time   `json:"later_sent_at"`
}

func (q *Queries) CountPositionHistories(ctx context.Context, arg CountPositionHistoriesParams) (int64, error) {
	row := q.db.QueryRow(ctx, countPositionHistories,
		arg.WhereInMember,
		arg.InMemberIds,
		arg.WhereEarlierSentAt,
		arg.EarlierSentAt,
		arg.WhereLaterSentAt,
		arg.LaterSentAt,
	)
	var count int64
	err := row.Scan(&count)
	return count, err
}

type CreatePositionHistoriesParams struct {
	MemberID uuid.UUID `json:"member_id"`
	XPos     float64   `json:"x_pos"`
	YPos     float64   `json:"y_pos"`
	SentAt   time.Time `json:"sent_at"`
}

const createPositionHistory = `-- name: CreatePositionHistory :one
INSERT INTO t_position_histories (member_id, x_pos, y_pos, sent_at) VALUES ($1, $2, $3, $4) RETURNING t_position_histories_pkey, position_history_id, member_id, x_pos, y_pos, sent_at
`

type CreatePositionHistoryParams struct {
	MemberID uuid.UUID `json:"member_id"`
	XPos     float64   `json:"x_pos"`
	YPos     float64   `json:"y_pos"`
	SentAt   time.Time `json:"sent_at"`
}

func (q *Queries) CreatePositionHistory(ctx context.Context, arg CreatePositionHistoryParams) (PositionHistory, error) {
	row := q.db.QueryRow(ctx, createPositionHistory,
		arg.MemberID,
		arg.XPos,
		arg.YPos,
		arg.SentAt,
	)
	var i PositionHistory
	err := row.Scan(
		&i.TPositionHistoriesPkey,
		&i.PositionHistoryID,
		&i.MemberID,
		&i.XPos,
		&i.YPos,
		&i.SentAt,
	)
	return i, err
}

const deletePositionHistory = `-- name: DeletePositionHistory :exec
DELETE FROM t_position_histories WHERE position_history_id = $1
`

func (q *Queries) DeletePositionHistory(ctx context.Context, positionHistoryID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deletePositionHistory, positionHistoryID)
	return err
}

const findPositionHistoryByID = `-- name: FindPositionHistoryByID :one
SELECT t_position_histories_pkey, position_history_id, member_id, x_pos, y_pos, sent_at FROM t_position_histories WHERE position_history_id = $1
`

func (q *Queries) FindPositionHistoryByID(ctx context.Context, positionHistoryID uuid.UUID) (PositionHistory, error) {
	row := q.db.QueryRow(ctx, findPositionHistoryByID, positionHistoryID)
	var i PositionHistory
	err := row.Scan(
		&i.TPositionHistoriesPkey,
		&i.PositionHistoryID,
		&i.MemberID,
		&i.XPos,
		&i.YPos,
		&i.SentAt,
	)
	return i, err
}

const findPositionHistoryByIDWithMember = `-- name: FindPositionHistoryByIDWithMember :one
SELECT t_position_histories.t_position_histories_pkey, t_position_histories.position_history_id, t_position_histories.member_id, t_position_histories.x_pos, t_position_histories.y_pos, t_position_histories.sent_at, m_members.m_members_pkey, m_members.member_id, m_members.login_id, m_members.password, m_members.email, m_members.name, m_members.attend_status_id, m_members.profile_image_id, m_members.grade_id, m_members.group_id, m_members.personal_organization_id, m_members.role_id, m_members.created_at, m_members.updated_at FROM t_position_histories
LEFT JOIN m_members ON t_position_histories.member_id = m_members.member_id
WHERE position_history_id = $1
`

type FindPositionHistoryByIDWithMemberRow struct {
	PositionHistory PositionHistory `json:"position_history"`
	Member          Member          `json:"member"`
}

func (q *Queries) FindPositionHistoryByIDWithMember(ctx context.Context, positionHistoryID uuid.UUID) (FindPositionHistoryByIDWithMemberRow, error) {
	row := q.db.QueryRow(ctx, findPositionHistoryByIDWithMember, positionHistoryID)
	var i FindPositionHistoryByIDWithMemberRow
	err := row.Scan(
		&i.PositionHistory.TPositionHistoriesPkey,
		&i.PositionHistory.PositionHistoryID,
		&i.PositionHistory.MemberID,
		&i.PositionHistory.XPos,
		&i.PositionHistory.YPos,
		&i.PositionHistory.SentAt,
		&i.Member.MMembersPkey,
		&i.Member.MemberID,
		&i.Member.LoginID,
		&i.Member.Password,
		&i.Member.Email,
		&i.Member.Name,
		&i.Member.AttendStatusID,
		&i.Member.ProfileImageID,
		&i.Member.GradeID,
		&i.Member.GroupID,
		&i.Member.PersonalOrganizationID,
		&i.Member.RoleID,
		&i.Member.CreatedAt,
		&i.Member.UpdatedAt,
	)
	return i, err
}

const getPluralPositionHistories = `-- name: GetPluralPositionHistories :many
SELECT t_position_histories_pkey, position_history_id, member_id, x_pos, y_pos, sent_at FROM t_position_histories WHERE position_history_id = ANY($3::uuid[])
ORDER BY
	t_position_histories_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralPositionHistoriesParams struct {
	Limit              int32       `json:"limit"`
	Offset             int32       `json:"offset"`
	PositionHistoryIds []uuid.UUID `json:"position_history_ids"`
}

func (q *Queries) GetPluralPositionHistories(ctx context.Context, arg GetPluralPositionHistoriesParams) ([]PositionHistory, error) {
	rows, err := q.db.Query(ctx, getPluralPositionHistories, arg.Limit, arg.Offset, arg.PositionHistoryIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PositionHistory{}
	for rows.Next() {
		var i PositionHistory
		if err := rows.Scan(
			&i.TPositionHistoriesPkey,
			&i.PositionHistoryID,
			&i.MemberID,
			&i.XPos,
			&i.YPos,
			&i.SentAt,
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

const getPluralPositionHistoriesWithMember = `-- name: GetPluralPositionHistoriesWithMember :many
SELECT t_position_histories.t_position_histories_pkey, t_position_histories.position_history_id, t_position_histories.member_id, t_position_histories.x_pos, t_position_histories.y_pos, t_position_histories.sent_at, m_members.m_members_pkey, m_members.member_id, m_members.login_id, m_members.password, m_members.email, m_members.name, m_members.attend_status_id, m_members.profile_image_id, m_members.grade_id, m_members.group_id, m_members.personal_organization_id, m_members.role_id, m_members.created_at, m_members.updated_at FROM t_position_histories
LEFT JOIN m_members ON t_position_histories.member_id = m_members.member_id
WHERE position_history_id = ANY($3::uuid[])
ORDER BY
	t_position_histories_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralPositionHistoriesWithMemberParams struct {
	Limit              int32       `json:"limit"`
	Offset             int32       `json:"offset"`
	PositionHistoryIds []uuid.UUID `json:"position_history_ids"`
}

type GetPluralPositionHistoriesWithMemberRow struct {
	PositionHistory PositionHistory `json:"position_history"`
	Member          Member          `json:"member"`
}

func (q *Queries) GetPluralPositionHistoriesWithMember(ctx context.Context, arg GetPluralPositionHistoriesWithMemberParams) ([]GetPluralPositionHistoriesWithMemberRow, error) {
	rows, err := q.db.Query(ctx, getPluralPositionHistoriesWithMember, arg.Limit, arg.Offset, arg.PositionHistoryIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralPositionHistoriesWithMemberRow{}
	for rows.Next() {
		var i GetPluralPositionHistoriesWithMemberRow
		if err := rows.Scan(
			&i.PositionHistory.TPositionHistoriesPkey,
			&i.PositionHistory.PositionHistoryID,
			&i.PositionHistory.MemberID,
			&i.PositionHistory.XPos,
			&i.PositionHistory.YPos,
			&i.PositionHistory.SentAt,
			&i.Member.MMembersPkey,
			&i.Member.MemberID,
			&i.Member.LoginID,
			&i.Member.Password,
			&i.Member.Email,
			&i.Member.Name,
			&i.Member.AttendStatusID,
			&i.Member.ProfileImageID,
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

const getPositionHistories = `-- name: GetPositionHistories :many
SELECT t_position_histories_pkey, position_history_id, member_id, x_pos, y_pos, sent_at FROM t_position_histories
WHERE
	CASE WHEN $1::boolean = true THEN member_id = ANY($2::uuid[]) ELSE TRUE END
AND
	CASE WHEN $3::boolean = true THEN sent_at >= $4 ELSE TRUE END
AND
	CASE WHEN $5::boolean = true THEN sent_at <= $6 ELSE TRUE END
ORDER BY
	CASE WHEN $7::text = 'old_send' THEN sent_at END ASC,
	CASE WHEN $7::text = 'late_send' THEN sent_at END DESC,
	t_position_histories_pkey ASC
`

type GetPositionHistoriesParams struct {
	WhereInMember      bool        `json:"where_in_member"`
	InMemberIds        []uuid.UUID `json:"in_member_ids"`
	WhereEarlierSentAt bool        `json:"where_earlier_sent_at"`
	EarlierSentAt      time.Time   `json:"earlier_sent_at"`
	WhereLaterSentAt   bool        `json:"where_later_sent_at"`
	LaterSentAt        time.Time   `json:"later_sent_at"`
	OrderMethod        string      `json:"order_method"`
}

func (q *Queries) GetPositionHistories(ctx context.Context, arg GetPositionHistoriesParams) ([]PositionHistory, error) {
	rows, err := q.db.Query(ctx, getPositionHistories,
		arg.WhereInMember,
		arg.InMemberIds,
		arg.WhereEarlierSentAt,
		arg.EarlierSentAt,
		arg.WhereLaterSentAt,
		arg.LaterSentAt,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PositionHistory{}
	for rows.Next() {
		var i PositionHistory
		if err := rows.Scan(
			&i.TPositionHistoriesPkey,
			&i.PositionHistoryID,
			&i.MemberID,
			&i.XPos,
			&i.YPos,
			&i.SentAt,
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

const getPositionHistoriesUseKeysetPaginate = `-- name: GetPositionHistoriesUseKeysetPaginate :many
SELECT t_position_histories_pkey, position_history_id, member_id, x_pos, y_pos, sent_at FROM t_position_histories
WHERE
	CASE WHEN $2::boolean = true THEN member_id = ANY($3::uuid[]) ELSE TRUE END
AND
	CASE WHEN $4::boolean = true THEN sent_at >= $5 ELSE TRUE END
AND
	CASE WHEN $6::boolean = true THEN sent_at <= $7 ELSE TRUE END
AND
	CASE $8::text
		WHEN 'next' THEN
			CASE $9::text
				WHEN 'old_send' THEN sent_at > $10 OR (sent_at = $10 AND t_position_histories_pkey > $11::int)
				WHEN 'late_send' THEN sent_at < $10 OR (sent_at = $10 AND t_position_histories_pkey > $11::int)
				ELSE t_position_histories_pkey > $11::int
			END
		WHEN 'prev' THEN
			CASE $9::text
				WHEN 'old_send' THEN sent_at < $10 OR (sent_at = $10 AND t_position_histories_pkey < $11::int)
				WHEN 'late_send' THEN sent_at > $10 OR (sent_at = $10 AND t_position_histories_pkey < $11::int)
				ELSE t_position_histories_pkey < $11::int
		END
	END
ORDER BY
	CASE WHEN $9::text = 'old_send' THEN sent_at END ASC,
	CASE WHEN $9::text = 'late_send' THEN sent_at END DESC,
	t_position_histories_pkey ASC
LIMIT $1
`

type GetPositionHistoriesUseKeysetPaginateParams struct {
	Limit              int32       `json:"limit"`
	WhereInMember      bool        `json:"where_in_member"`
	InMemberIds        []uuid.UUID `json:"in_member_ids"`
	WhereEarlierSentAt bool        `json:"where_earlier_sent_at"`
	EarlierSentAt      time.Time   `json:"earlier_sent_at"`
	WhereLaterSentAt   bool        `json:"where_later_sent_at"`
	LaterSentAt        time.Time   `json:"later_sent_at"`
	CursorDirection    string      `json:"cursor_direction"`
	OrderMethod        string      `json:"order_method"`
	SendCursor         time.Time   `json:"send_cursor"`
	Cursor             int32       `json:"cursor"`
}

func (q *Queries) GetPositionHistoriesUseKeysetPaginate(ctx context.Context, arg GetPositionHistoriesUseKeysetPaginateParams) ([]PositionHistory, error) {
	rows, err := q.db.Query(ctx, getPositionHistoriesUseKeysetPaginate,
		arg.Limit,
		arg.WhereInMember,
		arg.InMemberIds,
		arg.WhereEarlierSentAt,
		arg.EarlierSentAt,
		arg.WhereLaterSentAt,
		arg.LaterSentAt,
		arg.CursorDirection,
		arg.OrderMethod,
		arg.SendCursor,
		arg.Cursor,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PositionHistory{}
	for rows.Next() {
		var i PositionHistory
		if err := rows.Scan(
			&i.TPositionHistoriesPkey,
			&i.PositionHistoryID,
			&i.MemberID,
			&i.XPos,
			&i.YPos,
			&i.SentAt,
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

const getPositionHistoriesUseNumberedPaginate = `-- name: GetPositionHistoriesUseNumberedPaginate :many
SELECT t_position_histories_pkey, position_history_id, member_id, x_pos, y_pos, sent_at FROM t_position_histories
WHERE
	CASE WHEN $3::boolean = true THEN member_id = ANY($4::uuid[]) ELSE TRUE END
AND
	CASE WHEN $5::boolean = true THEN sent_at >= $6 ELSE TRUE END
AND
	CASE WHEN $7::boolean = true THEN sent_at <= $8 ELSE TRUE END
ORDER BY
	CASE WHEN $9::text = 'old_send' THEN sent_at END ASC,
	CASE WHEN $9::text = 'late_send' THEN sent_at END DESC,
	t_position_histories_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPositionHistoriesUseNumberedPaginateParams struct {
	Limit              int32       `json:"limit"`
	Offset             int32       `json:"offset"`
	WhereInMember      bool        `json:"where_in_member"`
	InMemberIds        []uuid.UUID `json:"in_member_ids"`
	WhereEarlierSentAt bool        `json:"where_earlier_sent_at"`
	EarlierSentAt      time.Time   `json:"earlier_sent_at"`
	WhereLaterSentAt   bool        `json:"where_later_sent_at"`
	LaterSentAt        time.Time   `json:"later_sent_at"`
	OrderMethod        string      `json:"order_method"`
}

func (q *Queries) GetPositionHistoriesUseNumberedPaginate(ctx context.Context, arg GetPositionHistoriesUseNumberedPaginateParams) ([]PositionHistory, error) {
	rows, err := q.db.Query(ctx, getPositionHistoriesUseNumberedPaginate,
		arg.Limit,
		arg.Offset,
		arg.WhereInMember,
		arg.InMemberIds,
		arg.WhereEarlierSentAt,
		arg.EarlierSentAt,
		arg.WhereLaterSentAt,
		arg.LaterSentAt,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PositionHistory{}
	for rows.Next() {
		var i PositionHistory
		if err := rows.Scan(
			&i.TPositionHistoriesPkey,
			&i.PositionHistoryID,
			&i.MemberID,
			&i.XPos,
			&i.YPos,
			&i.SentAt,
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

const getPositionHistoriesWithMember = `-- name: GetPositionHistoriesWithMember :many
SELECT t_position_histories.t_position_histories_pkey, t_position_histories.position_history_id, t_position_histories.member_id, t_position_histories.x_pos, t_position_histories.y_pos, t_position_histories.sent_at, m_members.m_members_pkey, m_members.member_id, m_members.login_id, m_members.password, m_members.email, m_members.name, m_members.attend_status_id, m_members.profile_image_id, m_members.grade_id, m_members.group_id, m_members.personal_organization_id, m_members.role_id, m_members.created_at, m_members.updated_at FROM t_position_histories
LEFT JOIN m_members ON t_position_histories.member_id = m_members.member_id
WHERE
	CASE WHEN $1::boolean = true THEN member_id = ANY($2::uuid[]) ELSE TRUE END
AND
	CASE WHEN $3::boolean = true THEN sent_at >= $4 ELSE TRUE END
AND
	CASE WHEN $5::boolean = true THEN sent_at <= $6 ELSE TRUE END
ORDER BY
	CASE WHEN $7::text = 'old_send' THEN sent_at END ASC,
	CASE WHEN $7::text = 'late_send' THEN sent_at END DESC,
	t_position_histories_pkey ASC
`

type GetPositionHistoriesWithMemberParams struct {
	WhereInMember      bool        `json:"where_in_member"`
	InMemberIds        []uuid.UUID `json:"in_member_ids"`
	WhereEarlierSentAt bool        `json:"where_earlier_sent_at"`
	EarlierSentAt      time.Time   `json:"earlier_sent_at"`
	WhereLaterSentAt   bool        `json:"where_later_sent_at"`
	LaterSentAt        time.Time   `json:"later_sent_at"`
	OrderMethod        string      `json:"order_method"`
}

type GetPositionHistoriesWithMemberRow struct {
	PositionHistory PositionHistory `json:"position_history"`
	Member          Member          `json:"member"`
}

func (q *Queries) GetPositionHistoriesWithMember(ctx context.Context, arg GetPositionHistoriesWithMemberParams) ([]GetPositionHistoriesWithMemberRow, error) {
	rows, err := q.db.Query(ctx, getPositionHistoriesWithMember,
		arg.WhereInMember,
		arg.InMemberIds,
		arg.WhereEarlierSentAt,
		arg.EarlierSentAt,
		arg.WhereLaterSentAt,
		arg.LaterSentAt,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPositionHistoriesWithMemberRow{}
	for rows.Next() {
		var i GetPositionHistoriesWithMemberRow
		if err := rows.Scan(
			&i.PositionHistory.TPositionHistoriesPkey,
			&i.PositionHistory.PositionHistoryID,
			&i.PositionHistory.MemberID,
			&i.PositionHistory.XPos,
			&i.PositionHistory.YPos,
			&i.PositionHistory.SentAt,
			&i.Member.MMembersPkey,
			&i.Member.MemberID,
			&i.Member.LoginID,
			&i.Member.Password,
			&i.Member.Email,
			&i.Member.Name,
			&i.Member.AttendStatusID,
			&i.Member.ProfileImageID,
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

const getPositionHistoriesWithMemberUseKeysetPaginate = `-- name: GetPositionHistoriesWithMemberUseKeysetPaginate :many
SELECT t_position_histories.t_position_histories_pkey, t_position_histories.position_history_id, t_position_histories.member_id, t_position_histories.x_pos, t_position_histories.y_pos, t_position_histories.sent_at, m_members.m_members_pkey, m_members.member_id, m_members.login_id, m_members.password, m_members.email, m_members.name, m_members.attend_status_id, m_members.profile_image_id, m_members.grade_id, m_members.group_id, m_members.personal_organization_id, m_members.role_id, m_members.created_at, m_members.updated_at FROM t_position_histories
LEFT JOIN m_members ON t_position_histories.member_id = m_members.member_id
WHERE
	CASE WHEN $2::boolean = true THEN member_id = ANY($3::uuid[]) ELSE TRUE END
AND
	CASE WHEN $4::boolean = true THEN sent_at >= $5 ELSE TRUE END
AND
	CASE WHEN $6::boolean = true THEN sent_at <= $7 ELSE TRUE END
AND
	CASE $8::text
		WHEN 'next' THEN
			CASE $9::text
				WHEN 'old_send' THEN sent_at > $10 OR (sent_at = $10 AND t_position_histories_pkey > $11::int)
				WHEN 'late_send' THEN sent_at < $10 OR (sent_at = $10 AND t_position_histories_pkey > $11::int)
				ELSE t_position_histories_pkey > $11::int
			END
		WHEN 'prev' THEN
			CASE $9::text
				WHEN 'old_send' THEN sent_at < $10 OR (sent_at = $10 AND t_position_histories_pkey < $11::int)
				WHEN 'late_send' THEN sent_at > $10 OR (sent_at = $10 AND t_position_histories_pkey < $11::int)
				ELSE t_position_histories_pkey < $11::int
			END
	END
ORDER BY
	CASE WHEN $9::text = 'old_send' THEN sent_at END ASC,
	CASE WHEN $9::text = 'late_send' THEN sent_at END DESC,
	t_position_histories_pkey ASC
LIMIT $1
`

type GetPositionHistoriesWithMemberUseKeysetPaginateParams struct {
	Limit              int32       `json:"limit"`
	WhereInMember      bool        `json:"where_in_member"`
	InMemberIds        []uuid.UUID `json:"in_member_ids"`
	WhereEarlierSentAt bool        `json:"where_earlier_sent_at"`
	EarlierSentAt      time.Time   `json:"earlier_sent_at"`
	WhereLaterSentAt   bool        `json:"where_later_sent_at"`
	LaterSentAt        time.Time   `json:"later_sent_at"`
	CursorDirection    string      `json:"cursor_direction"`
	OrderMethod        string      `json:"order_method"`
	SendCursor         time.Time   `json:"send_cursor"`
	Cursor             int32       `json:"cursor"`
}

type GetPositionHistoriesWithMemberUseKeysetPaginateRow struct {
	PositionHistory PositionHistory `json:"position_history"`
	Member          Member          `json:"member"`
}

func (q *Queries) GetPositionHistoriesWithMemberUseKeysetPaginate(ctx context.Context, arg GetPositionHistoriesWithMemberUseKeysetPaginateParams) ([]GetPositionHistoriesWithMemberUseKeysetPaginateRow, error) {
	rows, err := q.db.Query(ctx, getPositionHistoriesWithMemberUseKeysetPaginate,
		arg.Limit,
		arg.WhereInMember,
		arg.InMemberIds,
		arg.WhereEarlierSentAt,
		arg.EarlierSentAt,
		arg.WhereLaterSentAt,
		arg.LaterSentAt,
		arg.CursorDirection,
		arg.OrderMethod,
		arg.SendCursor,
		arg.Cursor,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPositionHistoriesWithMemberUseKeysetPaginateRow{}
	for rows.Next() {
		var i GetPositionHistoriesWithMemberUseKeysetPaginateRow
		if err := rows.Scan(
			&i.PositionHistory.TPositionHistoriesPkey,
			&i.PositionHistory.PositionHistoryID,
			&i.PositionHistory.MemberID,
			&i.PositionHistory.XPos,
			&i.PositionHistory.YPos,
			&i.PositionHistory.SentAt,
			&i.Member.MMembersPkey,
			&i.Member.MemberID,
			&i.Member.LoginID,
			&i.Member.Password,
			&i.Member.Email,
			&i.Member.Name,
			&i.Member.AttendStatusID,
			&i.Member.ProfileImageID,
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

const getPositionHistoriesWithMemberUseNumberedPaginate = `-- name: GetPositionHistoriesWithMemberUseNumberedPaginate :many
SELECT t_position_histories.t_position_histories_pkey, t_position_histories.position_history_id, t_position_histories.member_id, t_position_histories.x_pos, t_position_histories.y_pos, t_position_histories.sent_at, m_members.m_members_pkey, m_members.member_id, m_members.login_id, m_members.password, m_members.email, m_members.name, m_members.attend_status_id, m_members.profile_image_id, m_members.grade_id, m_members.group_id, m_members.personal_organization_id, m_members.role_id, m_members.created_at, m_members.updated_at FROM t_position_histories
LEFT JOIN m_members ON t_position_histories.member_id = m_members.member_id
WHERE
	CASE WHEN $3::boolean = true THEN member_id = ANY($4::uuid[]) ELSE TRUE END
AND
	CASE WHEN $5::boolean = true THEN sent_at >= $6 ELSE TRUE END
AND
	CASE WHEN $7::boolean = true THEN sent_at <= $8 ELSE TRUE END
ORDER BY
	CASE WHEN $9::text = 'old_send' THEN sent_at END ASC,
	CASE WHEN $9::text = 'late_send' THEN sent_at END DESC,
	t_position_histories_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPositionHistoriesWithMemberUseNumberedPaginateParams struct {
	Limit              int32       `json:"limit"`
	Offset             int32       `json:"offset"`
	WhereInMember      bool        `json:"where_in_member"`
	InMemberIds        []uuid.UUID `json:"in_member_ids"`
	WhereEarlierSentAt bool        `json:"where_earlier_sent_at"`
	EarlierSentAt      time.Time   `json:"earlier_sent_at"`
	WhereLaterSentAt   bool        `json:"where_later_sent_at"`
	LaterSentAt        time.Time   `json:"later_sent_at"`
	OrderMethod        string      `json:"order_method"`
}

type GetPositionHistoriesWithMemberUseNumberedPaginateRow struct {
	PositionHistory PositionHistory `json:"position_history"`
	Member          Member          `json:"member"`
}

func (q *Queries) GetPositionHistoriesWithMemberUseNumberedPaginate(ctx context.Context, arg GetPositionHistoriesWithMemberUseNumberedPaginateParams) ([]GetPositionHistoriesWithMemberUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getPositionHistoriesWithMemberUseNumberedPaginate,
		arg.Limit,
		arg.Offset,
		arg.WhereInMember,
		arg.InMemberIds,
		arg.WhereEarlierSentAt,
		arg.EarlierSentAt,
		arg.WhereLaterSentAt,
		arg.LaterSentAt,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPositionHistoriesWithMemberUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetPositionHistoriesWithMemberUseNumberedPaginateRow
		if err := rows.Scan(
			&i.PositionHistory.TPositionHistoriesPkey,
			&i.PositionHistory.PositionHistoryID,
			&i.PositionHistory.MemberID,
			&i.PositionHistory.XPos,
			&i.PositionHistory.YPos,
			&i.PositionHistory.SentAt,
			&i.Member.MMembersPkey,
			&i.Member.MemberID,
			&i.Member.LoginID,
			&i.Member.Password,
			&i.Member.Email,
			&i.Member.Name,
			&i.Member.AttendStatusID,
			&i.Member.ProfileImageID,
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

const updatePositionHistory = `-- name: UpdatePositionHistory :one
UPDATE t_position_histories SET member_id = $2, x_pos = $3, y_pos = $4, sent_at = $5 WHERE position_history_id = $1 RETURNING t_position_histories_pkey, position_history_id, member_id, x_pos, y_pos, sent_at
`

type UpdatePositionHistoryParams struct {
	PositionHistoryID uuid.UUID `json:"position_history_id"`
	MemberID          uuid.UUID `json:"member_id"`
	XPos              float64   `json:"x_pos"`
	YPos              float64   `json:"y_pos"`
	SentAt            time.Time `json:"sent_at"`
}

func (q *Queries) UpdatePositionHistory(ctx context.Context, arg UpdatePositionHistoryParams) (PositionHistory, error) {
	row := q.db.QueryRow(ctx, updatePositionHistory,
		arg.PositionHistoryID,
		arg.MemberID,
		arg.XPos,
		arg.YPos,
		arg.SentAt,
	)
	var i PositionHistory
	err := row.Scan(
		&i.TPositionHistoriesPkey,
		&i.PositionHistoryID,
		&i.MemberID,
		&i.XPos,
		&i.YPos,
		&i.SentAt,
	)
	return i, err
}

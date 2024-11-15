// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: lab_io_history.sql

package query

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const countLabIOHistories = `-- name: CountLabIOHistories :one
SELECT COUNT(*) FROM t_lab_io_histories
where
	CASE WHEN $1::boolean = true THEN member_id = ANY($2) ELSE TRUE END
AND
	CASE WHEN $3::boolean = true THEN entered_at >= $4 ELSE TRUE END
AND
	CASE WHEN $5::boolean = true THEN entered_at <= $6 ELSE TRUE END
AND
	CASE WHEN $7::boolean = true THEN exited_at >= $8 ELSE TRUE END
AND
	CASE WHEN $9::boolean = true THEN exited_at <= $10 ELSE TRUE END
`

type CountLabIOHistoriesParams struct {
	WhereInMember         bool               `json:"where_in_member"`
	InMember              uuid.UUID          `json:"in_member"`
	WhereEarlierEnteredAt bool               `json:"where_earlier_entered_at"`
	EarlierEnteredAt      time.Time          `json:"earlier_entered_at"`
	WhereLaterEnteredAt   bool               `json:"where_later_entered_at"`
	LaterEnteredAt        time.Time          `json:"later_entered_at"`
	WhereEarlierExitedAt  bool               `json:"where_earlier_exited_at"`
	EarlierExitedAt       pgtype.Timestamptz `json:"earlier_exited_at"`
	WhereLaterExitedAt    bool               `json:"where_later_exited_at"`
	LaterExitedAt         pgtype.Timestamptz `json:"later_exited_at"`
}

func (q *Queries) CountLabIOHistories(ctx context.Context, arg CountLabIOHistoriesParams) (int64, error) {
	row := q.db.QueryRow(ctx, countLabIOHistories,
		arg.WhereInMember,
		arg.InMember,
		arg.WhereEarlierEnteredAt,
		arg.EarlierEnteredAt,
		arg.WhereLaterEnteredAt,
		arg.LaterEnteredAt,
		arg.WhereEarlierExitedAt,
		arg.EarlierExitedAt,
		arg.WhereLaterExitedAt,
		arg.LaterExitedAt,
	)
	var count int64
	err := row.Scan(&count)
	return count, err
}

type CreateLabIOHistoriesParams struct {
	MemberID  uuid.UUID          `json:"member_id"`
	EnteredAt time.Time          `json:"entered_at"`
	ExitedAt  pgtype.Timestamptz `json:"exited_at"`
}

const createLabIOHistory = `-- name: CreateLabIOHistory :one
INSERT INTO t_lab_io_histories (member_id, entered_at, exited_at) VALUES ($1, $2, $3) RETURNING t_lab_io_histories_pkey, lab_io_history_id, member_id, entered_at, exited_at
`

type CreateLabIOHistoryParams struct {
	MemberID  uuid.UUID          `json:"member_id"`
	EnteredAt time.Time          `json:"entered_at"`
	ExitedAt  pgtype.Timestamptz `json:"exited_at"`
}

func (q *Queries) CreateLabIOHistory(ctx context.Context, arg CreateLabIOHistoryParams) (LabIOHistory, error) {
	row := q.db.QueryRow(ctx, createLabIOHistory, arg.MemberID, arg.EnteredAt, arg.ExitedAt)
	var i LabIOHistory
	err := row.Scan(
		&i.TLabIoHistoriesPkey,
		&i.LabIoHistoryID,
		&i.MemberID,
		&i.EnteredAt,
		&i.ExitedAt,
	)
	return i, err
}

const deleteLabIOHistory = `-- name: DeleteLabIOHistory :execrows
DELETE FROM t_lab_io_histories WHERE lab_io_history_id = $1
`

func (q *Queries) DeleteLabIOHistory(ctx context.Context, labIoHistoryID uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, deleteLabIOHistory, labIoHistoryID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const deleteLabIOHistoryOnMember = `-- name: DeleteLabIOHistoryOnMember :execrows
DELETE FROM t_lab_io_histories WHERE member_id = $1
`

func (q *Queries) DeleteLabIOHistoryOnMember(ctx context.Context, memberID uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, deleteLabIOHistoryOnMember, memberID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const exitLabIOHistory = `-- name: ExitLabIOHistory :one
UPDATE t_lab_io_histories SET exited_at = $2 WHERE lab_io_history_id = $1 RETURNING t_lab_io_histories_pkey, lab_io_history_id, member_id, entered_at, exited_at
`

type ExitLabIOHistoryParams struct {
	LabIoHistoryID uuid.UUID          `json:"lab_io_history_id"`
	ExitedAt       pgtype.Timestamptz `json:"exited_at"`
}

func (q *Queries) ExitLabIOHistory(ctx context.Context, arg ExitLabIOHistoryParams) (LabIOHistory, error) {
	row := q.db.QueryRow(ctx, exitLabIOHistory, arg.LabIoHistoryID, arg.ExitedAt)
	var i LabIOHistory
	err := row.Scan(
		&i.TLabIoHistoriesPkey,
		&i.LabIoHistoryID,
		&i.MemberID,
		&i.EnteredAt,
		&i.ExitedAt,
	)
	return i, err
}

const findLabIOHistoryByID = `-- name: FindLabIOHistoryByID :one
SELECT t_lab_io_histories_pkey, lab_io_history_id, member_id, entered_at, exited_at FROM t_lab_io_histories WHERE lab_io_history_id = $1
`

func (q *Queries) FindLabIOHistoryByID(ctx context.Context, labIoHistoryID uuid.UUID) (LabIOHistory, error) {
	row := q.db.QueryRow(ctx, findLabIOHistoryByID, labIoHistoryID)
	var i LabIOHistory
	err := row.Scan(
		&i.TLabIoHistoriesPkey,
		&i.LabIoHistoryID,
		&i.MemberID,
		&i.EnteredAt,
		&i.ExitedAt,
	)
	return i, err
}

const findLabIOHistoryWithMember = `-- name: FindLabIOHistoryWithMember :one
SELECT t_lab_io_histories.t_lab_io_histories_pkey, t_lab_io_histories.lab_io_history_id, t_lab_io_histories.member_id, t_lab_io_histories.entered_at, t_lab_io_histories.exited_at, m_members.m_members_pkey, m_members.member_id, m_members.login_id, m_members.password, m_members.email, m_members.name, m_members.first_name, m_members.last_name, m_members.attend_status_id, m_members.profile_image_id, m_members.grade_id, m_members.group_id, m_members.personal_organization_id, m_members.role_id, m_members.created_at, m_members.updated_at FROM t_lab_io_histories
LEFT JOIN m_members ON t_lab_io_histories.member_id = m_members.member_id
WHERE lab_io_history_id = $1
`

type FindLabIOHistoryWithMemberRow struct {
	LabIOHistory LabIOHistory `json:"lab_iohistory"`
	Member       Member       `json:"member"`
}

func (q *Queries) FindLabIOHistoryWithMember(ctx context.Context, labIoHistoryID uuid.UUID) (FindLabIOHistoryWithMemberRow, error) {
	row := q.db.QueryRow(ctx, findLabIOHistoryWithMember, labIoHistoryID)
	var i FindLabIOHistoryWithMemberRow
	err := row.Scan(
		&i.LabIOHistory.TLabIoHistoriesPkey,
		&i.LabIOHistory.LabIoHistoryID,
		&i.LabIOHistory.MemberID,
		&i.LabIOHistory.EnteredAt,
		&i.LabIOHistory.ExitedAt,
		&i.Member.MMembersPkey,
		&i.Member.MemberID,
		&i.Member.LoginID,
		&i.Member.Password,
		&i.Member.Email,
		&i.Member.Name,
		&i.Member.FirstName,
		&i.Member.LastName,
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

const getLabIOHistories = `-- name: GetLabIOHistories :many
SELECT t_lab_io_histories_pkey, lab_io_history_id, member_id, entered_at, exited_at FROM t_lab_io_histories
WHERE
	CASE WHEN $1::boolean = true THEN member_id = ANY($2) ELSE TRUE END
AND
	CASE WHEN $3::boolean = true THEN entered_at >= $4 ELSE TRUE END
AND
	CASE WHEN $5::boolean = true THEN entered_at <= $6 ELSE TRUE END
AND
	CASE WHEN $7::boolean = true THEN exited_at >= $8 ELSE TRUE END
AND
	CASE WHEN $9::boolean = true THEN exited_at <= $10 ELSE TRUE END
ORDER BY
	CASE WHEN $11::text = 'old_enter' THEN entered_at END ASC NULLS LAST,
	CASE WHEN $11::text = 'late_enter' THEN entered_at END DESC NULLS LAST,
	CASE WHEN $11::text = 'old_exit' THEN exited_at END ASC NULLS LAST,
	CASE WHEN $11::text = 'late_exit' THEN exited_at END DESC NULLS LAST,
	t_lab_io_histories_pkey ASC
`

type GetLabIOHistoriesParams struct {
	WhereInMember         bool               `json:"where_in_member"`
	InMember              uuid.UUID          `json:"in_member"`
	WhereEarlierEnteredAt bool               `json:"where_earlier_entered_at"`
	EarlierEnteredAt      time.Time          `json:"earlier_entered_at"`
	WhereLaterEnteredAt   bool               `json:"where_later_entered_at"`
	LaterEnteredAt        time.Time          `json:"later_entered_at"`
	WhereEarlierExitedAt  bool               `json:"where_earlier_exited_at"`
	EarlierExitedAt       pgtype.Timestamptz `json:"earlier_exited_at"`
	WhereLaterExitedAt    bool               `json:"where_later_exited_at"`
	LaterExitedAt         pgtype.Timestamptz `json:"later_exited_at"`
	OrderMethod           string             `json:"order_method"`
}

func (q *Queries) GetLabIOHistories(ctx context.Context, arg GetLabIOHistoriesParams) ([]LabIOHistory, error) {
	rows, err := q.db.Query(ctx, getLabIOHistories,
		arg.WhereInMember,
		arg.InMember,
		arg.WhereEarlierEnteredAt,
		arg.EarlierEnteredAt,
		arg.WhereLaterEnteredAt,
		arg.LaterEnteredAt,
		arg.WhereEarlierExitedAt,
		arg.EarlierExitedAt,
		arg.WhereLaterExitedAt,
		arg.LaterExitedAt,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []LabIOHistory{}
	for rows.Next() {
		var i LabIOHistory
		if err := rows.Scan(
			&i.TLabIoHistoriesPkey,
			&i.LabIoHistoryID,
			&i.MemberID,
			&i.EnteredAt,
			&i.ExitedAt,
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

const getLabIOHistoriesUseKeysetPaginate = `-- name: GetLabIOHistoriesUseKeysetPaginate :many
SELECT t_lab_io_histories_pkey, lab_io_history_id, member_id, entered_at, exited_at FROM t_lab_io_histories
WHERE
	CASE WHEN $2::boolean = true THEN member_id = ANY($3) ELSE TRUE END
AND
	CASE WHEN $4::boolean = true THEN entered_at >= $5 ELSE TRUE END
AND
	CASE WHEN $6::boolean = true THEN entered_at <= $7 ELSE TRUE END
AND
	CASE WHEN $8::boolean = true THEN exited_at >= $9 ELSE TRUE END
AND
	CASE WHEN $10::boolean = true THEN exited_at <= $11 ELSE TRUE END
AND
	CASE $12::text
		WHEN 'next' THEN
			CASE $13::text
				WHEN 'old_enter' THEN entered_at > $14 OR (entered_at = $14 AND t_lab_io_histories_pkey > $15::int)
				WHEN 'late_enter' THEN entered_at < $14 OR (entered_at = $14 AND t_lab_io_histories_pkey > $15::int)
				WHEN 'old_exit' THEN exited_at > $16 OR (exited_at = $16 AND t_lab_io_histories_pkey > $15::int)
				WHEN 'late_exit' THEN exited_at < $16 OR (exited_at = $16 AND t_lab_io_histories_pkey > $15::int)
				ELSE t_lab_io_histories_pkey > $15::int
			END
		WHEN 'prev' THEN
			CASE $13::text
				WHEN 'old_enter' THEN entered_at < $14 OR (entered_at = $14 AND t_lab_io_histories_pkey < $15::int)
				WHEN 'late_enter' THEN entered_at > $14 OR (entered_at = $14 AND t_lab_io_histories_pkey < $15::int)
				WHEN 'old_exit' THEN exited_at < $16 OR (exited_at = $16 AND t_lab_io_histories_pkey < $15::int)
				WHEN 'late_exit' THEN exited_at > $16 OR (exited_at = $16 AND t_lab_io_histories_pkey < $15::int)
				ELSE t_lab_io_histories_pkey < $15::int
		END
	END
ORDER BY
	CASE WHEN $13::text = 'old_enter' AND $12::text = 'next' THEN entered_at END ASC NULLS LAST,
	CASE WHEN $13::text = 'old_enter' AND $12::text = 'prev' THEN entered_at END DESC NULLS LAST,
	CASE WHEN $13::text = 'late_enter' AND $12::text = 'next' THEN entered_at END DESC NULLS LAST,
	CASE WHEN $13::text = 'late_enter' AND $12::text = 'prev' THEN entered_at END ASC NULLS LAST,
	CASE WHEN $13::text = 'old_exit' AND $12::text = 'next' THEN exited_at END ASC NULLS LAST,
	CASE WHEN $13::text = 'old_exit' AND $12::text = 'prev' THEN exited_at END DESC NULLS LAST,
	CASE WHEN $13::text = 'late_exit' AND $12::text = 'next' THEN exited_at END DESC NULLS LAST,
	CASE WHEN $13::text = 'late_exit' AND $12::text = 'prev' THEN exited_at END ASC NULLS LAST,
	CASE WHEN $12::text = 'next' THEN t_lab_io_histories_pkey END ASC,
	CASE WHEN $12::text = 'prev' THEN t_lab_io_histories_pkey END DESC
LIMIT $1
`

type GetLabIOHistoriesUseKeysetPaginateParams struct {
	Limit                 int32              `json:"limit"`
	WhereInMember         bool               `json:"where_in_member"`
	InMember              uuid.UUID          `json:"in_member"`
	WhereEarlierEnteredAt bool               `json:"where_earlier_entered_at"`
	EarlierEnteredAt      time.Time          `json:"earlier_entered_at"`
	WhereLaterEnteredAt   bool               `json:"where_later_entered_at"`
	LaterEnteredAt        time.Time          `json:"later_entered_at"`
	WhereEarlierExitedAt  bool               `json:"where_earlier_exited_at"`
	EarlierExitedAt       pgtype.Timestamptz `json:"earlier_exited_at"`
	WhereLaterExitedAt    bool               `json:"where_later_exited_at"`
	LaterExitedAt         pgtype.Timestamptz `json:"later_exited_at"`
	CursorDirection       string             `json:"cursor_direction"`
	OrderMethod           string             `json:"order_method"`
	EnterCursor           time.Time          `json:"enter_cursor"`
	Cursor                int32              `json:"cursor"`
	ExitCursor            pgtype.Timestamptz `json:"exit_cursor"`
}

func (q *Queries) GetLabIOHistoriesUseKeysetPaginate(ctx context.Context, arg GetLabIOHistoriesUseKeysetPaginateParams) ([]LabIOHistory, error) {
	rows, err := q.db.Query(ctx, getLabIOHistoriesUseKeysetPaginate,
		arg.Limit,
		arg.WhereInMember,
		arg.InMember,
		arg.WhereEarlierEnteredAt,
		arg.EarlierEnteredAt,
		arg.WhereLaterEnteredAt,
		arg.LaterEnteredAt,
		arg.WhereEarlierExitedAt,
		arg.EarlierExitedAt,
		arg.WhereLaterExitedAt,
		arg.LaterExitedAt,
		arg.CursorDirection,
		arg.OrderMethod,
		arg.EnterCursor,
		arg.Cursor,
		arg.ExitCursor,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []LabIOHistory{}
	for rows.Next() {
		var i LabIOHistory
		if err := rows.Scan(
			&i.TLabIoHistoriesPkey,
			&i.LabIoHistoryID,
			&i.MemberID,
			&i.EnteredAt,
			&i.ExitedAt,
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

const getLabIOHistoriesUseNumberedPaginate = `-- name: GetLabIOHistoriesUseNumberedPaginate :many
SELECT t_lab_io_histories_pkey, lab_io_history_id, member_id, entered_at, exited_at FROM t_lab_io_histories
WHERE
	CASE WHEN $3::boolean = true THEN member_id = ANY($4) ELSE TRUE END
AND
	CASE WHEN $5::boolean = true THEN entered_at >= $6 ELSE TRUE END
AND
	CASE WHEN $7::boolean = true THEN entered_at <= $8 ELSE TRUE END
AND
	CASE WHEN $9::boolean = true THEN exited_at >= $10 ELSE TRUE END
AND
	CASE WHEN $11::boolean = true THEN exited_at <= $12 ELSE TRUE END
ORDER BY
	CASE WHEN $13::text = 'old_enter' THEN entered_at END ASC NULLS LAST,
	CASE WHEN $13::text = 'late_enter' THEN entered_at END DESC NULLS LAST,
	CASE WHEN $13::text = 'old_exit' THEN exited_at END ASC NULLS LAST,
	CASE WHEN $13::text = 'late_exit' THEN exited_at END DESC NULLS LAST,
	t_lab_io_histories_pkey ASC
LIMIT $1 OFFSET $2
`

type GetLabIOHistoriesUseNumberedPaginateParams struct {
	Limit                 int32              `json:"limit"`
	Offset                int32              `json:"offset"`
	WhereInMember         bool               `json:"where_in_member"`
	InMember              uuid.UUID          `json:"in_member"`
	WhereEarlierEnteredAt bool               `json:"where_earlier_entered_at"`
	EarlierEnteredAt      time.Time          `json:"earlier_entered_at"`
	WhereLaterEnteredAt   bool               `json:"where_later_entered_at"`
	LaterEnteredAt        time.Time          `json:"later_entered_at"`
	WhereEarlierExitedAt  bool               `json:"where_earlier_exited_at"`
	EarlierExitedAt       pgtype.Timestamptz `json:"earlier_exited_at"`
	WhereLaterExitedAt    bool               `json:"where_later_exited_at"`
	LaterExitedAt         pgtype.Timestamptz `json:"later_exited_at"`
	OrderMethod           string             `json:"order_method"`
}

func (q *Queries) GetLabIOHistoriesUseNumberedPaginate(ctx context.Context, arg GetLabIOHistoriesUseNumberedPaginateParams) ([]LabIOHistory, error) {
	rows, err := q.db.Query(ctx, getLabIOHistoriesUseNumberedPaginate,
		arg.Limit,
		arg.Offset,
		arg.WhereInMember,
		arg.InMember,
		arg.WhereEarlierEnteredAt,
		arg.EarlierEnteredAt,
		arg.WhereLaterEnteredAt,
		arg.LaterEnteredAt,
		arg.WhereEarlierExitedAt,
		arg.EarlierExitedAt,
		arg.WhereLaterExitedAt,
		arg.LaterExitedAt,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []LabIOHistory{}
	for rows.Next() {
		var i LabIOHistory
		if err := rows.Scan(
			&i.TLabIoHistoriesPkey,
			&i.LabIoHistoryID,
			&i.MemberID,
			&i.EnteredAt,
			&i.ExitedAt,
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

const getLabIOHistoriesWithMember = `-- name: GetLabIOHistoriesWithMember :many
SELECT t_lab_io_histories.t_lab_io_histories_pkey, t_lab_io_histories.lab_io_history_id, t_lab_io_histories.member_id, t_lab_io_histories.entered_at, t_lab_io_histories.exited_at, m_members.m_members_pkey, m_members.member_id, m_members.login_id, m_members.password, m_members.email, m_members.name, m_members.first_name, m_members.last_name, m_members.attend_status_id, m_members.profile_image_id, m_members.grade_id, m_members.group_id, m_members.personal_organization_id, m_members.role_id, m_members.created_at, m_members.updated_at FROM t_lab_io_histories
LEFT JOIN m_members ON t_lab_io_histories.member_id = m_members.member_id
WHERE
	CASE WHEN $1::boolean = true THEN t_lab_io_histories.member_id = ANY($2) ELSE TRUE END
AND
	CASE WHEN $3::boolean = true THEN entered_at >= $4 ELSE TRUE END
AND
	CASE WHEN $5::boolean = true THEN entered_at <= $6 ELSE TRUE END
AND
	CASE WHEN $7::boolean = true THEN exited_at >= $8 ELSE TRUE END
AND
	CASE WHEN $9::boolean = true THEN exited_at <= $10 ELSE TRUE END
ORDER BY
	CASE WHEN $11::text = 'old_enter' THEN entered_at END ASC NULLS LAST,
	CASE WHEN $11::text = 'late_enter' THEN entered_at END DESC NULLS LAST,
	CASE WHEN $11::text = 'old_exit' THEN exited_at END ASC NULLS LAST,
	CASE WHEN $11::text = 'late_exit' THEN exited_at END DESC NULLS LAST,
	t_lab_io_histories_pkey ASC
`

type GetLabIOHistoriesWithMemberParams struct {
	WhereInMember         bool               `json:"where_in_member"`
	InMember              uuid.UUID          `json:"in_member"`
	WhereEarlierEnteredAt bool               `json:"where_earlier_entered_at"`
	EarlierEnteredAt      time.Time          `json:"earlier_entered_at"`
	WhereLaterEnteredAt   bool               `json:"where_later_entered_at"`
	LaterEnteredAt        time.Time          `json:"later_entered_at"`
	WhereEarlierExitedAt  bool               `json:"where_earlier_exited_at"`
	EarlierExitedAt       pgtype.Timestamptz `json:"earlier_exited_at"`
	WhereLaterExitedAt    bool               `json:"where_later_exited_at"`
	LaterExitedAt         pgtype.Timestamptz `json:"later_exited_at"`
	OrderMethod           string             `json:"order_method"`
}

type GetLabIOHistoriesWithMemberRow struct {
	LabIOHistory LabIOHistory `json:"lab_iohistory"`
	Member       Member       `json:"member"`
}

func (q *Queries) GetLabIOHistoriesWithMember(ctx context.Context, arg GetLabIOHistoriesWithMemberParams) ([]GetLabIOHistoriesWithMemberRow, error) {
	rows, err := q.db.Query(ctx, getLabIOHistoriesWithMember,
		arg.WhereInMember,
		arg.InMember,
		arg.WhereEarlierEnteredAt,
		arg.EarlierEnteredAt,
		arg.WhereLaterEnteredAt,
		arg.LaterEnteredAt,
		arg.WhereEarlierExitedAt,
		arg.EarlierExitedAt,
		arg.WhereLaterExitedAt,
		arg.LaterExitedAt,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetLabIOHistoriesWithMemberRow{}
	for rows.Next() {
		var i GetLabIOHistoriesWithMemberRow
		if err := rows.Scan(
			&i.LabIOHistory.TLabIoHistoriesPkey,
			&i.LabIOHistory.LabIoHistoryID,
			&i.LabIOHistory.MemberID,
			&i.LabIOHistory.EnteredAt,
			&i.LabIOHistory.ExitedAt,
			&i.Member.MMembersPkey,
			&i.Member.MemberID,
			&i.Member.LoginID,
			&i.Member.Password,
			&i.Member.Email,
			&i.Member.Name,
			&i.Member.FirstName,
			&i.Member.LastName,
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

const getLabIOHistoriesWithMemberUseKeysetPaginate = `-- name: GetLabIOHistoriesWithMemberUseKeysetPaginate :many
SELECT t_lab_io_histories.t_lab_io_histories_pkey, t_lab_io_histories.lab_io_history_id, t_lab_io_histories.member_id, t_lab_io_histories.entered_at, t_lab_io_histories.exited_at, m_members.m_members_pkey, m_members.member_id, m_members.login_id, m_members.password, m_members.email, m_members.name, m_members.first_name, m_members.last_name, m_members.attend_status_id, m_members.profile_image_id, m_members.grade_id, m_members.group_id, m_members.personal_organization_id, m_members.role_id, m_members.created_at, m_members.updated_at FROM t_lab_io_histories
LEFT JOIN m_members ON t_lab_io_histories.member_id = m_members.member_id
WHERE
	CASE WHEN $2::boolean = true THEN t_lab_io_histories.member_id = ANY($3) ELSE TRUE END
AND
	CASE WHEN $4::boolean = true THEN entered_at >= $5 ELSE TRUE END
AND
	CASE WHEN $6::boolean = true THEN entered_at <= $7 ELSE TRUE END
AND
	CASE WHEN $8::boolean = true THEN exited_at >= $9 ELSE TRUE END
AND
	CASE WHEN $10::boolean = true THEN exited_at <= $11 ELSE TRUE END
AND
	CASE $12::text
		WHEN 'next' THEN
			CASE $13::text
				WHEN 'old_enter' THEN entered_at > $14 OR (entered_at = $14 AND t_lab_io_histories_pkey > $15::int)
				WHEN 'late_enter' THEN entered_at < $14 OR (entered_at = $14 AND t_lab_io_histories_pkey > $15::int)
				WHEN 'old_exit' THEN exited_at > $16 OR (exited_at = $16 AND t_lab_io_histories_pkey > $15::int)
				WHEN 'late_exit' THEN exited_at < $16 OR (exited_at = $16 AND t_lab_io_histories_pkey > $15::int)
				ELSE t_lab_io_histories_pkey > $15::int
			END
		WHEN 'prev' THEN
			CASE $13::text
				WHEN 'old_enter' THEN entered_at < $14 OR (entered_at = $14 AND t_lab_io_histories_pkey < $15::int)
				WHEN 'late_enter' THEN entered_at > $14 OR (entered_at = $14 AND t_lab_io_histories_pkey < $15::int)
				WHEN 'old_exit' THEN exited_at < $16 OR (exited_at = $16 AND t_lab_io_histories_pkey < $15::int)
				WHEN 'late_exit' THEN exited_at > $16 OR (exited_at = $16 AND t_lab_io_histories_pkey < $15::int)
				ELSE t_lab_io_histories_pkey < $15::int
		END
	END
ORDER BY
	CASE WHEN $13::text = 'old_enter' AND $12::text = 'next' THEN entered_at END ASC NULLS LAST,
	CASE WHEN $13::text = 'old_enter' AND $12::text = 'prev' THEN entered_at END DESC NULLS LAST,
	CASE WHEN $13::text = 'late_enter' AND $12::text = 'next' THEN entered_at END DESC NULLS LAST,
	CASE WHEN $13::text = 'late_enter' AND $12::text = 'prev' THEN entered_at END ASC NULLS LAST,
	CASE WHEN $13::text = 'old_exit' AND $12::text = 'next' THEN exited_at END ASC NULLS LAST,
	CASE WHEN $13::text = 'old_exit' AND $12::text = 'prev' THEN exited_at END DESC NULLS LAST,
	CASE WHEN $13::text = 'late_exit' AND $12::text = 'next' THEN exited_at END DESC NULLS LAST,
	CASE WHEN $13::text = 'late_exit' AND $12::text = 'prev' THEN exited_at END ASC NULLS LAST,
	CASE WHEN $12::text = 'next' THEN t_lab_io_histories_pkey END ASC,
	CASE WHEN $12::text = 'prev' THEN t_lab_io_histories_pkey END DESC
LIMIT $1
`

type GetLabIOHistoriesWithMemberUseKeysetPaginateParams struct {
	Limit                 int32              `json:"limit"`
	WhereInMember         bool               `json:"where_in_member"`
	InMember              uuid.UUID          `json:"in_member"`
	WhereEarlierEnteredAt bool               `json:"where_earlier_entered_at"`
	EarlierEnteredAt      time.Time          `json:"earlier_entered_at"`
	WhereLaterEnteredAt   bool               `json:"where_later_entered_at"`
	LaterEnteredAt        time.Time          `json:"later_entered_at"`
	WhereEarlierExitedAt  bool               `json:"where_earlier_exited_at"`
	EarlierExitedAt       pgtype.Timestamptz `json:"earlier_exited_at"`
	WhereLaterExitedAt    bool               `json:"where_later_exited_at"`
	LaterExitedAt         pgtype.Timestamptz `json:"later_exited_at"`
	CursorDirection       string             `json:"cursor_direction"`
	OrderMethod           string             `json:"order_method"`
	EnterCursor           time.Time          `json:"enter_cursor"`
	Cursor                int32              `json:"cursor"`
	ExitCursor            pgtype.Timestamptz `json:"exit_cursor"`
}

type GetLabIOHistoriesWithMemberUseKeysetPaginateRow struct {
	LabIOHistory LabIOHistory `json:"lab_iohistory"`
	Member       Member       `json:"member"`
}

func (q *Queries) GetLabIOHistoriesWithMemberUseKeysetPaginate(ctx context.Context, arg GetLabIOHistoriesWithMemberUseKeysetPaginateParams) ([]GetLabIOHistoriesWithMemberUseKeysetPaginateRow, error) {
	rows, err := q.db.Query(ctx, getLabIOHistoriesWithMemberUseKeysetPaginate,
		arg.Limit,
		arg.WhereInMember,
		arg.InMember,
		arg.WhereEarlierEnteredAt,
		arg.EarlierEnteredAt,
		arg.WhereLaterEnteredAt,
		arg.LaterEnteredAt,
		arg.WhereEarlierExitedAt,
		arg.EarlierExitedAt,
		arg.WhereLaterExitedAt,
		arg.LaterExitedAt,
		arg.CursorDirection,
		arg.OrderMethod,
		arg.EnterCursor,
		arg.Cursor,
		arg.ExitCursor,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetLabIOHistoriesWithMemberUseKeysetPaginateRow{}
	for rows.Next() {
		var i GetLabIOHistoriesWithMemberUseKeysetPaginateRow
		if err := rows.Scan(
			&i.LabIOHistory.TLabIoHistoriesPkey,
			&i.LabIOHistory.LabIoHistoryID,
			&i.LabIOHistory.MemberID,
			&i.LabIOHistory.EnteredAt,
			&i.LabIOHistory.ExitedAt,
			&i.Member.MMembersPkey,
			&i.Member.MemberID,
			&i.Member.LoginID,
			&i.Member.Password,
			&i.Member.Email,
			&i.Member.Name,
			&i.Member.FirstName,
			&i.Member.LastName,
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

const getLabIOHistoriesWithMemberUseNumberedPaginate = `-- name: GetLabIOHistoriesWithMemberUseNumberedPaginate :many
SELECT t_lab_io_histories.t_lab_io_histories_pkey, t_lab_io_histories.lab_io_history_id, t_lab_io_histories.member_id, t_lab_io_histories.entered_at, t_lab_io_histories.exited_at, m_members.m_members_pkey, m_members.member_id, m_members.login_id, m_members.password, m_members.email, m_members.name, m_members.first_name, m_members.last_name, m_members.attend_status_id, m_members.profile_image_id, m_members.grade_id, m_members.group_id, m_members.personal_organization_id, m_members.role_id, m_members.created_at, m_members.updated_at FROM t_lab_io_histories
LEFT JOIN m_members ON t_lab_io_histories.member_id = m_members.member_id
WHERE
	CASE WHEN $3::boolean = true THEN t_lab_io_histories.member_id = ANY($4) ELSE TRUE END
AND
	CASE WHEN $5::boolean = true THEN entered_at >= $6 ELSE TRUE END
AND
	CASE WHEN $7::boolean = true THEN entered_at <= $8 ELSE TRUE END
AND
	CASE WHEN $9::boolean = true THEN exited_at >= $10 ELSE TRUE END
AND
	CASE WHEN $11::boolean = true THEN exited_at <= $12 ELSE TRUE END
ORDER BY
	CASE WHEN $13::text = 'old_enter' THEN entered_at END ASC NULLS LAST,
	CASE WHEN $13::text = 'late_enter' THEN entered_at END DESC NULLS LAST,
	CASE WHEN $13::text = 'old_exit' THEN exited_at END ASC NULLS LAST,
	CASE WHEN $13::text = 'late_exit' THEN exited_at END DESC NULLS LAST,
	t_lab_io_histories_pkey ASC
LIMIT $1 OFFSET $2
`

type GetLabIOHistoriesWithMemberUseNumberedPaginateParams struct {
	Limit                 int32              `json:"limit"`
	Offset                int32              `json:"offset"`
	WhereInMember         bool               `json:"where_in_member"`
	InMember              uuid.UUID          `json:"in_member"`
	WhereEarlierEnteredAt bool               `json:"where_earlier_entered_at"`
	EarlierEnteredAt      time.Time          `json:"earlier_entered_at"`
	WhereLaterEnteredAt   bool               `json:"where_later_entered_at"`
	LaterEnteredAt        time.Time          `json:"later_entered_at"`
	WhereEarlierExitedAt  bool               `json:"where_earlier_exited_at"`
	EarlierExitedAt       pgtype.Timestamptz `json:"earlier_exited_at"`
	WhereLaterExitedAt    bool               `json:"where_later_exited_at"`
	LaterExitedAt         pgtype.Timestamptz `json:"later_exited_at"`
	OrderMethod           string             `json:"order_method"`
}

type GetLabIOHistoriesWithMemberUseNumberedPaginateRow struct {
	LabIOHistory LabIOHistory `json:"lab_iohistory"`
	Member       Member       `json:"member"`
}

func (q *Queries) GetLabIOHistoriesWithMemberUseNumberedPaginate(ctx context.Context, arg GetLabIOHistoriesWithMemberUseNumberedPaginateParams) ([]GetLabIOHistoriesWithMemberUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getLabIOHistoriesWithMemberUseNumberedPaginate,
		arg.Limit,
		arg.Offset,
		arg.WhereInMember,
		arg.InMember,
		arg.WhereEarlierEnteredAt,
		arg.EarlierEnteredAt,
		arg.WhereLaterEnteredAt,
		arg.LaterEnteredAt,
		arg.WhereEarlierExitedAt,
		arg.EarlierExitedAt,
		arg.WhereLaterExitedAt,
		arg.LaterExitedAt,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetLabIOHistoriesWithMemberUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetLabIOHistoriesWithMemberUseNumberedPaginateRow
		if err := rows.Scan(
			&i.LabIOHistory.TLabIoHistoriesPkey,
			&i.LabIOHistory.LabIoHistoryID,
			&i.LabIOHistory.MemberID,
			&i.LabIOHistory.EnteredAt,
			&i.LabIOHistory.ExitedAt,
			&i.Member.MMembersPkey,
			&i.Member.MemberID,
			&i.Member.LoginID,
			&i.Member.Password,
			&i.Member.Email,
			&i.Member.Name,
			&i.Member.FirstName,
			&i.Member.LastName,
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

const getPluralLabIOHistories = `-- name: GetPluralLabIOHistories :many
SELECT t_lab_io_histories_pkey, lab_io_history_id, member_id, entered_at, exited_at FROM t_lab_io_histories WHERE lab_io_history_id = ANY($1::uuid[])
ORDER BY
	CASE WHEN $2::text = 'old_enter' THEN entered_at END ASC NULLS LAST,
	CASE WHEN $2::text = 'late_enter' THEN entered_at END DESC NULLS LAST,
	CASE WHEN $2::text = 'old_exit' THEN exited_at END ASC NULLS LAST,
	CASE WHEN $2::text = 'late_exit' THEN exited_at END DESC NULLS LAST,
	t_lab_io_histories_pkey ASC
`

type GetPluralLabIOHistoriesParams struct {
	LabIoHistoryIds []uuid.UUID `json:"lab_io_history_ids"`
	OrderMethod     string      `json:"order_method"`
}

func (q *Queries) GetPluralLabIOHistories(ctx context.Context, arg GetPluralLabIOHistoriesParams) ([]LabIOHistory, error) {
	rows, err := q.db.Query(ctx, getPluralLabIOHistories, arg.LabIoHistoryIds, arg.OrderMethod)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []LabIOHistory{}
	for rows.Next() {
		var i LabIOHistory
		if err := rows.Scan(
			&i.TLabIoHistoriesPkey,
			&i.LabIoHistoryID,
			&i.MemberID,
			&i.EnteredAt,
			&i.ExitedAt,
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

const getPluralLabIOHistoriesUseNumberedPaginate = `-- name: GetPluralLabIOHistoriesUseNumberedPaginate :many
SELECT t_lab_io_histories_pkey, lab_io_history_id, member_id, entered_at, exited_at FROM t_lab_io_histories WHERE lab_io_history_id = ANY($3::uuid[])
ORDER BY
	CASE WHEN $4::text = 'old_enter' THEN entered_at END ASC NULLS LAST,
	CASE WHEN $4::text = 'late_enter' THEN entered_at END DESC NULLS LAST,
	CASE WHEN $4::text = 'old_exit' THEN exited_at END ASC NULLS LAST,
	CASE WHEN $4::text = 'late_exit' THEN exited_at END DESC NULLS LAST,
	t_lab_io_histories_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralLabIOHistoriesUseNumberedPaginateParams struct {
	Limit           int32       `json:"limit"`
	Offset          int32       `json:"offset"`
	LabIoHistoryIds []uuid.UUID `json:"lab_io_history_ids"`
	OrderMethod     string      `json:"order_method"`
}

func (q *Queries) GetPluralLabIOHistoriesUseNumberedPaginate(ctx context.Context, arg GetPluralLabIOHistoriesUseNumberedPaginateParams) ([]LabIOHistory, error) {
	rows, err := q.db.Query(ctx, getPluralLabIOHistoriesUseNumberedPaginate,
		arg.Limit,
		arg.Offset,
		arg.LabIoHistoryIds,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []LabIOHistory{}
	for rows.Next() {
		var i LabIOHistory
		if err := rows.Scan(
			&i.TLabIoHistoriesPkey,
			&i.LabIoHistoryID,
			&i.MemberID,
			&i.EnteredAt,
			&i.ExitedAt,
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

const getPluralLabIOHistoriesWithMember = `-- name: GetPluralLabIOHistoriesWithMember :many
SELECT t_lab_io_histories.t_lab_io_histories_pkey, t_lab_io_histories.lab_io_history_id, t_lab_io_histories.member_id, t_lab_io_histories.entered_at, t_lab_io_histories.exited_at, m_members.m_members_pkey, m_members.member_id, m_members.login_id, m_members.password, m_members.email, m_members.name, m_members.first_name, m_members.last_name, m_members.attend_status_id, m_members.profile_image_id, m_members.grade_id, m_members.group_id, m_members.personal_organization_id, m_members.role_id, m_members.created_at, m_members.updated_at FROM t_lab_io_histories
LEFT JOIN m_members ON t_lab_io_histories.member_id = m_members.member_id
WHERE lab_io_history_id = ANY($1::uuid[])
ORDER BY
	CASE WHEN $2::text = 'old_enter' THEN entered_at END ASC NULLS LAST,
	CASE WHEN $2::text = 'late_enter' THEN entered_at END DESC NULLS LAST,
	CASE WHEN $2::text = 'old_exit' THEN exited_at END ASC NULLS LAST,
	CASE WHEN $2::text = 'late_exit' THEN exited_at END DESC NULLS LAST,
	t_lab_io_histories_pkey ASC
`

type GetPluralLabIOHistoriesWithMemberParams struct {
	LabIoHistoryIds []uuid.UUID `json:"lab_io_history_ids"`
	OrderMethod     string      `json:"order_method"`
}

type GetPluralLabIOHistoriesWithMemberRow struct {
	LabIOHistory LabIOHistory `json:"lab_iohistory"`
	Member       Member       `json:"member"`
}

func (q *Queries) GetPluralLabIOHistoriesWithMember(ctx context.Context, arg GetPluralLabIOHistoriesWithMemberParams) ([]GetPluralLabIOHistoriesWithMemberRow, error) {
	rows, err := q.db.Query(ctx, getPluralLabIOHistoriesWithMember, arg.LabIoHistoryIds, arg.OrderMethod)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralLabIOHistoriesWithMemberRow{}
	for rows.Next() {
		var i GetPluralLabIOHistoriesWithMemberRow
		if err := rows.Scan(
			&i.LabIOHistory.TLabIoHistoriesPkey,
			&i.LabIOHistory.LabIoHistoryID,
			&i.LabIOHistory.MemberID,
			&i.LabIOHistory.EnteredAt,
			&i.LabIOHistory.ExitedAt,
			&i.Member.MMembersPkey,
			&i.Member.MemberID,
			&i.Member.LoginID,
			&i.Member.Password,
			&i.Member.Email,
			&i.Member.Name,
			&i.Member.FirstName,
			&i.Member.LastName,
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

const getPluralLabIOHistoriesWithMemberUseNumberedPaginate = `-- name: GetPluralLabIOHistoriesWithMemberUseNumberedPaginate :many
SELECT t_lab_io_histories.t_lab_io_histories_pkey, t_lab_io_histories.lab_io_history_id, t_lab_io_histories.member_id, t_lab_io_histories.entered_at, t_lab_io_histories.exited_at, m_members.m_members_pkey, m_members.member_id, m_members.login_id, m_members.password, m_members.email, m_members.name, m_members.first_name, m_members.last_name, m_members.attend_status_id, m_members.profile_image_id, m_members.grade_id, m_members.group_id, m_members.personal_organization_id, m_members.role_id, m_members.created_at, m_members.updated_at FROM t_lab_io_histories
LEFT JOIN m_members ON t_lab_io_histories.member_id = m_members.member_id
WHERE lab_io_history_id = ANY($3::uuid[])
ORDER BY
	CASE WHEN $4::text = 'old_enter' THEN entered_at END ASC NULLS LAST,
	CASE WHEN $4::text = 'late_enter' THEN entered_at END DESC NULLS LAST,
	CASE WHEN $4::text = 'old_exit' THEN exited_at END ASC NULLS LAST,
	CASE WHEN $4::text = 'late_exit' THEN exited_at END DESC NULLS LAST,
	t_lab_io_histories_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralLabIOHistoriesWithMemberUseNumberedPaginateParams struct {
	Limit           int32       `json:"limit"`
	Offset          int32       `json:"offset"`
	LabIoHistoryIds []uuid.UUID `json:"lab_io_history_ids"`
	OrderMethod     string      `json:"order_method"`
}

type GetPluralLabIOHistoriesWithMemberUseNumberedPaginateRow struct {
	LabIOHistory LabIOHistory `json:"lab_iohistory"`
	Member       Member       `json:"member"`
}

func (q *Queries) GetPluralLabIOHistoriesWithMemberUseNumberedPaginate(ctx context.Context, arg GetPluralLabIOHistoriesWithMemberUseNumberedPaginateParams) ([]GetPluralLabIOHistoriesWithMemberUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getPluralLabIOHistoriesWithMemberUseNumberedPaginate,
		arg.Limit,
		arg.Offset,
		arg.LabIoHistoryIds,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralLabIOHistoriesWithMemberUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetPluralLabIOHistoriesWithMemberUseNumberedPaginateRow
		if err := rows.Scan(
			&i.LabIOHistory.TLabIoHistoriesPkey,
			&i.LabIOHistory.LabIoHistoryID,
			&i.LabIOHistory.MemberID,
			&i.LabIOHistory.EnteredAt,
			&i.LabIOHistory.ExitedAt,
			&i.Member.MMembersPkey,
			&i.Member.MemberID,
			&i.Member.LoginID,
			&i.Member.Password,
			&i.Member.Email,
			&i.Member.Name,
			&i.Member.FirstName,
			&i.Member.LastName,
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

const pluralDeleteLabIOHistories = `-- name: PluralDeleteLabIOHistories :execrows
DELETE FROM t_lab_io_histories WHERE lab_io_history_id = ANY($1::uuid[])
`

func (q *Queries) PluralDeleteLabIOHistories(ctx context.Context, labIoHistoryIds []uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, pluralDeleteLabIOHistories, labIoHistoryIds)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const updateLabIOHistory = `-- name: UpdateLabIOHistory :one
UPDATE t_lab_io_histories SET member_id = $2, entered_at = $3, exited_at = $4 WHERE lab_io_history_id = $1 RETURNING t_lab_io_histories_pkey, lab_io_history_id, member_id, entered_at, exited_at
`

type UpdateLabIOHistoryParams struct {
	LabIoHistoryID uuid.UUID          `json:"lab_io_history_id"`
	MemberID       uuid.UUID          `json:"member_id"`
	EnteredAt      time.Time          `json:"entered_at"`
	ExitedAt       pgtype.Timestamptz `json:"exited_at"`
}

func (q *Queries) UpdateLabIOHistory(ctx context.Context, arg UpdateLabIOHistoryParams) (LabIOHistory, error) {
	row := q.db.QueryRow(ctx, updateLabIOHistory,
		arg.LabIoHistoryID,
		arg.MemberID,
		arg.EnteredAt,
		arg.ExitedAt,
	)
	var i LabIOHistory
	err := row.Scan(
		&i.TLabIoHistoriesPkey,
		&i.LabIoHistoryID,
		&i.MemberID,
		&i.EnteredAt,
		&i.ExitedAt,
	)
	return i, err
}

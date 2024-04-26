// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: organization.sql

package query

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const countOrganizations = `-- name: CountOrganizations :one
SELECT COUNT(*) FROM m_organizations
WHERE
	CASE WHEN $1::boolean = true THEN m_organizations.name LIKE '%' || $2::text || '%' END
AND
	CASE WHEN $3::boolean = true THEN m_organizations.is_whole = $4 END
AND
	CASE WHEN $5::boolean = true THEN m_organizations.is_personal = $6 AND EXISTS (SELECT m_members_pkey, member_id, login_id, password, email, name, attend_status_id, profile_image_id, grade_id, group_id, personal_organization_id, role_id, created_at, updated_at FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = $7::uuid) END
AND
	CASE WHEN $8::boolean = true THEN EXISTS (SELECT m_groups_pkey, group_id, key, organization_id FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) END
AND
	CASE WHEN $9::boolean = true THEN EXISTS (SELECT m_grades_pkey, grade_id, key, organization_id FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) END
`

type CountOrganizationsParams struct {
	WhereLikeName    bool      `json:"where_like_name"`
	SearchName       string    `json:"search_name"`
	WhereIsWhole     bool      `json:"where_is_whole"`
	IsWhole          bool      `json:"is_whole"`
	WhereIsPersonal  bool      `json:"where_is_personal"`
	IsPersonal       bool      `json:"is_personal"`
	PersonalMemberID uuid.UUID `json:"personal_member_id"`
	WhereIsGroup     bool      `json:"where_is_group"`
	WhereIsGrade     bool      `json:"where_is_grade"`
}

func (q *Queries) CountOrganizations(ctx context.Context, arg CountOrganizationsParams) (int64, error) {
	row := q.db.QueryRow(ctx, countOrganizations,
		arg.WhereLikeName,
		arg.SearchName,
		arg.WhereIsWhole,
		arg.IsWhole,
		arg.WhereIsPersonal,
		arg.IsPersonal,
		arg.PersonalMemberID,
		arg.WhereIsGroup,
		arg.WhereIsGrade,
	)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createOrganization = `-- name: CreateOrganization :one
INSERT INTO m_organizations (name, description, is_personal, is_whole, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING m_organizations_pkey, organization_id, name, description, is_personal, is_whole, created_at, updated_at
`

type CreateOrganizationParams struct {
	Name        string      `json:"name"`
	Description pgtype.Text `json:"description"`
	IsPersonal  bool        `json:"is_personal"`
	IsWhole     bool        `json:"is_whole"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

func (q *Queries) CreateOrganization(ctx context.Context, arg CreateOrganizationParams) (Organization, error) {
	row := q.db.QueryRow(ctx, createOrganization,
		arg.Name,
		arg.Description,
		arg.IsPersonal,
		arg.IsWhole,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Organization
	err := row.Scan(
		&i.MOrganizationsPkey,
		&i.OrganizationID,
		&i.Name,
		&i.Description,
		&i.IsPersonal,
		&i.IsWhole,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

type CreateOrganizationsParams struct {
	Name        string      `json:"name"`
	Description pgtype.Text `json:"description"`
	IsPersonal  bool        `json:"is_personal"`
	IsWhole     bool        `json:"is_whole"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

const deleteOrganization = `-- name: DeleteOrganization :exec
DELETE FROM m_organizations WHERE organization_id = $1
`

func (q *Queries) DeleteOrganization(ctx context.Context, organizationID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteOrganization, organizationID)
	return err
}

const findOrganizationByID = `-- name: FindOrganizationByID :one
SELECT m_organizations_pkey, organization_id, name, description, is_personal, is_whole, created_at, updated_at FROM m_organizations WHERE organization_id = $1
`

func (q *Queries) FindOrganizationByID(ctx context.Context, organizationID uuid.UUID) (Organization, error) {
	row := q.db.QueryRow(ctx, findOrganizationByID, organizationID)
	var i Organization
	err := row.Scan(
		&i.MOrganizationsPkey,
		&i.OrganizationID,
		&i.Name,
		&i.Description,
		&i.IsPersonal,
		&i.IsWhole,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findOrganizationByIDWithDetail = `-- name: FindOrganizationByIDWithDetail :one
SELECT m_organizations.m_organizations_pkey, m_organizations.organization_id, m_organizations.name, m_organizations.description, m_organizations.is_personal, m_organizations.is_whole, m_organizations.created_at, m_organizations.updated_at, m_groups.m_groups_pkey, m_groups.group_id, m_groups.key, m_groups.organization_id, m_grades.m_grades_pkey, m_grades.grade_id, m_grades.key, m_grades.organization_id FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
WHERE m_organizations.organization_id = $1
`

type FindOrganizationByIDWithDetailRow struct {
	Organization Organization `json:"organization"`
	Group        Group        `json:"group"`
	Grade        Grade        `json:"grade"`
}

func (q *Queries) FindOrganizationByIDWithDetail(ctx context.Context, organizationID uuid.UUID) (FindOrganizationByIDWithDetailRow, error) {
	row := q.db.QueryRow(ctx, findOrganizationByIDWithDetail, organizationID)
	var i FindOrganizationByIDWithDetailRow
	err := row.Scan(
		&i.Organization.MOrganizationsPkey,
		&i.Organization.OrganizationID,
		&i.Organization.Name,
		&i.Organization.Description,
		&i.Organization.IsPersonal,
		&i.Organization.IsWhole,
		&i.Organization.CreatedAt,
		&i.Organization.UpdatedAt,
		&i.Group.MGroupsPkey,
		&i.Group.GroupID,
		&i.Group.Key,
		&i.Group.OrganizationID,
		&i.Grade.MGradesPkey,
		&i.Grade.GradeID,
		&i.Grade.Key,
		&i.Grade.OrganizationID,
	)
	return i, err
}

const findPersonalOrganization = `-- name: FindPersonalOrganization :one
SELECT m_organizations_pkey, organization_id, m_organizations.name, description, is_personal, is_whole, m_organizations.created_at, m_organizations.updated_at, m_members_pkey, member_id, login_id, password, email, m_members.name, attend_status_id, profile_image_id, grade_id, group_id, personal_organization_id, role_id, m_members.created_at, m_members.updated_at FROM m_organizations
LEFT JOIN m_members ON m_organizations.organization_id = m_members.personal_organization_id
WHERE m_organizations.is_personal = true AND m_members.member_id = $1
`

type FindPersonalOrganizationRow struct {
	MOrganizationsPkey     pgtype.Int8        `json:"m_organizations_pkey"`
	OrganizationID         uuid.UUID          `json:"organization_id"`
	Name                   string             `json:"name"`
	Description            pgtype.Text        `json:"description"`
	IsPersonal             bool               `json:"is_personal"`
	IsWhole                bool               `json:"is_whole"`
	CreatedAt              time.Time          `json:"created_at"`
	UpdatedAt              time.Time          `json:"updated_at"`
	MMembersPkey           pgtype.Int8        `json:"m_members_pkey"`
	MemberID               pgtype.UUID        `json:"member_id"`
	LoginID                pgtype.Text        `json:"login_id"`
	Password               pgtype.Text        `json:"password"`
	Email                  pgtype.Text        `json:"email"`
	Name_2                 pgtype.Text        `json:"name_2"`
	AttendStatusID         pgtype.UUID        `json:"attend_status_id"`
	ProfileImageID         pgtype.UUID        `json:"profile_image_id"`
	GradeID                pgtype.UUID        `json:"grade_id"`
	GroupID                pgtype.UUID        `json:"group_id"`
	PersonalOrganizationID pgtype.UUID        `json:"personal_organization_id"`
	RoleID                 pgtype.UUID        `json:"role_id"`
	CreatedAt_2            pgtype.Timestamptz `json:"created_at_2"`
	UpdatedAt_2            pgtype.Timestamptz `json:"updated_at_2"`
}

func (q *Queries) FindPersonalOrganization(ctx context.Context, memberID uuid.UUID) (FindPersonalOrganizationRow, error) {
	row := q.db.QueryRow(ctx, findPersonalOrganization, memberID)
	var i FindPersonalOrganizationRow
	err := row.Scan(
		&i.MOrganizationsPkey,
		&i.OrganizationID,
		&i.Name,
		&i.Description,
		&i.IsPersonal,
		&i.IsWhole,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.MMembersPkey,
		&i.MemberID,
		&i.LoginID,
		&i.Password,
		&i.Email,
		&i.Name_2,
		&i.AttendStatusID,
		&i.ProfileImageID,
		&i.GradeID,
		&i.GroupID,
		&i.PersonalOrganizationID,
		&i.RoleID,
		&i.CreatedAt_2,
		&i.UpdatedAt_2,
	)
	return i, err
}

const findWholeOrganization = `-- name: FindWholeOrganization :one
SELECT m_organizations_pkey, organization_id, name, description, is_personal, is_whole, created_at, updated_at FROM m_organizations WHERE is_whole = true
`

func (q *Queries) FindWholeOrganization(ctx context.Context) (Organization, error) {
	row := q.db.QueryRow(ctx, findWholeOrganization)
	var i Organization
	err := row.Scan(
		&i.MOrganizationsPkey,
		&i.OrganizationID,
		&i.Name,
		&i.Description,
		&i.IsPersonal,
		&i.IsWhole,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getOrganizations = `-- name: GetOrganizations :many
SELECT m_organizations_pkey, organization_id, name, description, is_personal, is_whole, created_at, updated_at FROM m_organizations
WHERE
	CASE WHEN $1::boolean = true THEN m_organizations.name LIKE '%' || $2::text || '%' END
AND
	CASE WHEN $3::boolean = true THEN m_organizations.is_whole = $4 END
AND
	CASE WHEN $5::boolean = true THEN m_organizations.is_personal = $6 AND EXISTS (SELECT m_members_pkey, member_id, login_id, password, email, name, attend_status_id, profile_image_id, grade_id, group_id, personal_organization_id, role_id, created_at, updated_at FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = $7::uuid) END
AND
	CASE WHEN $8::boolean = true THEN EXISTS (SELECT m_groups_pkey, group_id, key, organization_id FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) END
AND
	CASE WHEN $9::boolean = true THEN EXISTS (SELECT m_grades_pkey, grade_id, key, organization_id FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) END
ORDER BY
	CASE WHEN $10::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN $10::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey DESC
`

type GetOrganizationsParams struct {
	WhereLikeName    bool      `json:"where_like_name"`
	SearchName       string    `json:"search_name"`
	WhereIsWhole     bool      `json:"where_is_whole"`
	IsWhole          bool      `json:"is_whole"`
	WhereIsPersonal  bool      `json:"where_is_personal"`
	IsPersonal       bool      `json:"is_personal"`
	PersonalMemberID uuid.UUID `json:"personal_member_id"`
	WhereIsGroup     bool      `json:"where_is_group"`
	WhereIsGrade     bool      `json:"where_is_grade"`
	OrderMethod      string    `json:"order_method"`
}

func (q *Queries) GetOrganizations(ctx context.Context, arg GetOrganizationsParams) ([]Organization, error) {
	rows, err := q.db.Query(ctx, getOrganizations,
		arg.WhereLikeName,
		arg.SearchName,
		arg.WhereIsWhole,
		arg.IsWhole,
		arg.WhereIsPersonal,
		arg.IsPersonal,
		arg.PersonalMemberID,
		arg.WhereIsGroup,
		arg.WhereIsGrade,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Organization{}
	for rows.Next() {
		var i Organization
		if err := rows.Scan(
			&i.MOrganizationsPkey,
			&i.OrganizationID,
			&i.Name,
			&i.Description,
			&i.IsPersonal,
			&i.IsWhole,
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

const getOrganizationsUseKeysetPaginate = `-- name: GetOrganizationsUseKeysetPaginate :many
SELECT m_organizations_pkey, organization_id, name, description, is_personal, is_whole, created_at, updated_at FROM m_organizations
WHERE
	CASE WHEN $2::boolean = true THEN m_organizations.name LIKE '%' || $3::text || '%' END
AND
	CASE WHEN $4::boolean = true THEN m_organizations.is_whole = $5 END
AND
	CASE WHEN $6::boolean = true THEN m_organizations.is_personal = $7 AND EXISTS (SELECT m_members_pkey, member_id, login_id, password, email, name, attend_status_id, profile_image_id, grade_id, group_id, personal_organization_id, role_id, created_at, updated_at FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = $8::uuid) END
AND
	CASE WHEN $9::boolean = true THEN EXISTS (SELECT m_groups_pkey, group_id, key, organization_id FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) END
AND
	CASE WHEN $10::boolean = true THEN EXISTS (SELECT m_grades_pkey, grade_id, key, organization_id FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) END
AND
	CASE $11::text
		WHEN 'next' THEN
			CASE $12::text
				WHEN 'name' THEN m_organizations.name > $13 OR (m_organizations.name = $13 AND m_organizations_pkey < $14::int)
				WHEN 'r_name' THEN m_organizations.name < $13 OR (m_organizations.name = $13 AND m_organizations_pkey < $14::int)
				ELSE m_organizations_pkey < $14::int
			END
		WHEN 'prev' THEN
			CASE $12::text
				WHEN 'name' THEN m_organizations.name < $13 OR (m_organizations.name = $13 AND m_organizations_pkey > $14::int)
				WHEN 'r_name' THEN m_organizations.name > $13 OR (m_organizations.name = $13 AND m_organizations_pkey > $14::int)
				ELSE m_organizations_pkey > $14::int
			END
	END
ORDER BY
	CASE WHEN $12::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN $12::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey DESC
LIMIT $1
`

type GetOrganizationsUseKeysetPaginateParams struct {
	Limit            int32     `json:"limit"`
	WhereLikeName    bool      `json:"where_like_name"`
	SearchName       string    `json:"search_name"`
	WhereIsWhole     bool      `json:"where_is_whole"`
	IsWhole          bool      `json:"is_whole"`
	WhereIsPersonal  bool      `json:"where_is_personal"`
	IsPersonal       bool      `json:"is_personal"`
	PersonalMemberID uuid.UUID `json:"personal_member_id"`
	WhereIsGroup     bool      `json:"where_is_group"`
	WhereIsGrade     bool      `json:"where_is_grade"`
	CursorDirection  string    `json:"cursor_direction"`
	OrderMethod      string    `json:"order_method"`
	NameCursor       string    `json:"name_cursor"`
	Cursor           int32     `json:"cursor"`
}

func (q *Queries) GetOrganizationsUseKeysetPaginate(ctx context.Context, arg GetOrganizationsUseKeysetPaginateParams) ([]Organization, error) {
	rows, err := q.db.Query(ctx, getOrganizationsUseKeysetPaginate,
		arg.Limit,
		arg.WhereLikeName,
		arg.SearchName,
		arg.WhereIsWhole,
		arg.IsWhole,
		arg.WhereIsPersonal,
		arg.IsPersonal,
		arg.PersonalMemberID,
		arg.WhereIsGroup,
		arg.WhereIsGrade,
		arg.CursorDirection,
		arg.OrderMethod,
		arg.NameCursor,
		arg.Cursor,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Organization{}
	for rows.Next() {
		var i Organization
		if err := rows.Scan(
			&i.MOrganizationsPkey,
			&i.OrganizationID,
			&i.Name,
			&i.Description,
			&i.IsPersonal,
			&i.IsWhole,
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

const getOrganizationsUseNumberedPaginate = `-- name: GetOrganizationsUseNumberedPaginate :many
SELECT m_organizations_pkey, organization_id, name, description, is_personal, is_whole, created_at, updated_at FROM m_organizations
WHERE
	CASE WHEN $3::boolean = true THEN m_organizations.name LIKE '%' || $4::text || '%' END
AND
	CASE WHEN $5::boolean = true THEN m_organizations.is_whole = $6 END
AND
	CASE WHEN $7::boolean = true THEN m_organizations.is_personal = $8 AND EXISTS (SELECT m_members_pkey, member_id, login_id, password, email, name, attend_status_id, profile_image_id, grade_id, group_id, personal_organization_id, role_id, created_at, updated_at FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = $9::uuid) END
AND
	CASE WHEN $10::boolean = true THEN EXISTS (SELECT m_groups_pkey, group_id, key, organization_id FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) END
AND
	CASE WHEN $11::boolean = true THEN EXISTS (SELECT m_grades_pkey, grade_id, key, organization_id FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) END
ORDER BY
	CASE WHEN $12::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN $12::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey DESC
LIMIT $1 OFFSET $2
`

type GetOrganizationsUseNumberedPaginateParams struct {
	Limit            int32     `json:"limit"`
	Offset           int32     `json:"offset"`
	WhereLikeName    bool      `json:"where_like_name"`
	SearchName       string    `json:"search_name"`
	WhereIsWhole     bool      `json:"where_is_whole"`
	IsWhole          bool      `json:"is_whole"`
	WhereIsPersonal  bool      `json:"where_is_personal"`
	IsPersonal       bool      `json:"is_personal"`
	PersonalMemberID uuid.UUID `json:"personal_member_id"`
	WhereIsGroup     bool      `json:"where_is_group"`
	WhereIsGrade     bool      `json:"where_is_grade"`
	OrderMethod      string    `json:"order_method"`
}

func (q *Queries) GetOrganizationsUseNumberedPaginate(ctx context.Context, arg GetOrganizationsUseNumberedPaginateParams) ([]Organization, error) {
	rows, err := q.db.Query(ctx, getOrganizationsUseNumberedPaginate,
		arg.Limit,
		arg.Offset,
		arg.WhereLikeName,
		arg.SearchName,
		arg.WhereIsWhole,
		arg.IsWhole,
		arg.WhereIsPersonal,
		arg.IsPersonal,
		arg.PersonalMemberID,
		arg.WhereIsGroup,
		arg.WhereIsGrade,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Organization{}
	for rows.Next() {
		var i Organization
		if err := rows.Scan(
			&i.MOrganizationsPkey,
			&i.OrganizationID,
			&i.Name,
			&i.Description,
			&i.IsPersonal,
			&i.IsWhole,
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

const getOrganizationsWithDetail = `-- name: GetOrganizationsWithDetail :many
SELECT m_organizations.m_organizations_pkey, m_organizations.organization_id, m_organizations.name, m_organizations.description, m_organizations.is_personal, m_organizations.is_whole, m_organizations.created_at, m_organizations.updated_at, m_groups.m_groups_pkey, m_groups.group_id, m_groups.key, m_groups.organization_id, m_grades.m_grades_pkey, m_grades.grade_id, m_grades.key, m_grades.organization_id FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
WHERE
	CASE WHEN $1::boolean = true THEN m_organizations.name LIKE '%' || $2::text || '%' END
AND
	CASE WHEN $3::boolean = true THEN m_organizations.is_whole = $4 END
AND
	CASE WHEN $5::boolean = true THEN m_organizations.is_personal = $6 AND EXISTS (SELECT m_members_pkey, member_id, login_id, password, email, name, attend_status_id, profile_image_id, grade_id, group_id, personal_organization_id, role_id, created_at, updated_at FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = $7::uuid) END
AND
	CASE WHEN $8::boolean = true THEN EXISTS (SELECT m_groups_pkey, group_id, key, organization_id FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) END
AND
	CASE WHEN $9::boolean = true THEN EXISTS (SELECT m_grades_pkey, grade_id, key, organization_id FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) END
ORDER BY
	CASE WHEN $10::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN $10::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey DESC
`

type GetOrganizationsWithDetailParams struct {
	WhereLikeName    bool      `json:"where_like_name"`
	SearchName       string    `json:"search_name"`
	WhereIsWhole     bool      `json:"where_is_whole"`
	IsWhole          bool      `json:"is_whole"`
	WhereIsPersonal  bool      `json:"where_is_personal"`
	IsPersonal       bool      `json:"is_personal"`
	PersonalMemberID uuid.UUID `json:"personal_member_id"`
	WhereIsGroup     bool      `json:"where_is_group"`
	WhereIsGrade     bool      `json:"where_is_grade"`
	OrderMethod      string    `json:"order_method"`
}

type GetOrganizationsWithDetailRow struct {
	Organization Organization `json:"organization"`
	Group        Group        `json:"group"`
	Grade        Grade        `json:"grade"`
}

func (q *Queries) GetOrganizationsWithDetail(ctx context.Context, arg GetOrganizationsWithDetailParams) ([]GetOrganizationsWithDetailRow, error) {
	rows, err := q.db.Query(ctx, getOrganizationsWithDetail,
		arg.WhereLikeName,
		arg.SearchName,
		arg.WhereIsWhole,
		arg.IsWhole,
		arg.WhereIsPersonal,
		arg.IsPersonal,
		arg.PersonalMemberID,
		arg.WhereIsGroup,
		arg.WhereIsGrade,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetOrganizationsWithDetailRow{}
	for rows.Next() {
		var i GetOrganizationsWithDetailRow
		if err := rows.Scan(
			&i.Organization.MOrganizationsPkey,
			&i.Organization.OrganizationID,
			&i.Organization.Name,
			&i.Organization.Description,
			&i.Organization.IsPersonal,
			&i.Organization.IsWhole,
			&i.Organization.CreatedAt,
			&i.Organization.UpdatedAt,
			&i.Group.MGroupsPkey,
			&i.Group.GroupID,
			&i.Group.Key,
			&i.Group.OrganizationID,
			&i.Grade.MGradesPkey,
			&i.Grade.GradeID,
			&i.Grade.Key,
			&i.Grade.OrganizationID,
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

const getOrganizationsWithDetailUseKeysetPaginate = `-- name: GetOrganizationsWithDetailUseKeysetPaginate :many
SELECT m_organizations.m_organizations_pkey, m_organizations.organization_id, m_organizations.name, m_organizations.description, m_organizations.is_personal, m_organizations.is_whole, m_organizations.created_at, m_organizations.updated_at, m_groups.m_groups_pkey, m_groups.group_id, m_groups.key, m_groups.organization_id, m_grades.m_grades_pkey, m_grades.grade_id, m_grades.key, m_grades.organization_id FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
WHERE
	CASE WHEN $2::boolean = true THEN m_organizations.name LIKE '%' || $3::text || '%' END
AND
	CASE WHEN $4::boolean = true THEN m_organizations.is_whole = $5 END
AND
	CASE WHEN $6::boolean = true THEN m_organizations.is_personal = $7 AND EXISTS (SELECT m_members_pkey, member_id, login_id, password, email, name, attend_status_id, profile_image_id, grade_id, group_id, personal_organization_id, role_id, created_at, updated_at FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = $8::uuid) END
AND
	CASE WHEN $9::boolean = true THEN EXISTS (SELECT m_groups_pkey, group_id, key, organization_id FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) END
AND
	CASE WHEN $10::boolean = true THEN EXISTS (SELECT m_grades_pkey, grade_id, key, organization_id FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) END
AND
	CASE $11::text
		WHEN 'next' THEN
			CASE $12::text
				WHEN 'name' THEN m_organizations.name > $13 OR (m_organizations.name = $13 AND m_organizations_pkey < $14::int)
				WHEN 'r_name' THEN m_organizations.name < $13 OR (m_organizations.name = $13 AND m_organizations_pkey < $14::int)
				ELSE m_organizations_pkey < $14::int
			END
		WHEN 'prev' THEN
			CASE $12::text
				WHEN 'name' THEN m_organizations.name < $13 OR (m_organizations.name = $13 AND m_organizations_pkey > $14::int)
				WHEN 'r_name' THEN m_organizations.name > $13 OR (m_organizations.name = $13 AND m_organizations_pkey > $14::int)
				ELSE m_organizations_pkey > $14::int
			END
	END
ORDER BY
	CASE WHEN $12::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN $12::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey DESC
LIMIT $1
`

type GetOrganizationsWithDetailUseKeysetPaginateParams struct {
	Limit            int32     `json:"limit"`
	WhereLikeName    bool      `json:"where_like_name"`
	SearchName       string    `json:"search_name"`
	WhereIsWhole     bool      `json:"where_is_whole"`
	IsWhole          bool      `json:"is_whole"`
	WhereIsPersonal  bool      `json:"where_is_personal"`
	IsPersonal       bool      `json:"is_personal"`
	PersonalMemberID uuid.UUID `json:"personal_member_id"`
	WhereIsGroup     bool      `json:"where_is_group"`
	WhereIsGrade     bool      `json:"where_is_grade"`
	CursorDirection  string    `json:"cursor_direction"`
	OrderMethod      string    `json:"order_method"`
	NameCursor       string    `json:"name_cursor"`
	Cursor           int32     `json:"cursor"`
}

type GetOrganizationsWithDetailUseKeysetPaginateRow struct {
	Organization Organization `json:"organization"`
	Group        Group        `json:"group"`
	Grade        Grade        `json:"grade"`
}

func (q *Queries) GetOrganizationsWithDetailUseKeysetPaginate(ctx context.Context, arg GetOrganizationsWithDetailUseKeysetPaginateParams) ([]GetOrganizationsWithDetailUseKeysetPaginateRow, error) {
	rows, err := q.db.Query(ctx, getOrganizationsWithDetailUseKeysetPaginate,
		arg.Limit,
		arg.WhereLikeName,
		arg.SearchName,
		arg.WhereIsWhole,
		arg.IsWhole,
		arg.WhereIsPersonal,
		arg.IsPersonal,
		arg.PersonalMemberID,
		arg.WhereIsGroup,
		arg.WhereIsGrade,
		arg.CursorDirection,
		arg.OrderMethod,
		arg.NameCursor,
		arg.Cursor,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetOrganizationsWithDetailUseKeysetPaginateRow{}
	for rows.Next() {
		var i GetOrganizationsWithDetailUseKeysetPaginateRow
		if err := rows.Scan(
			&i.Organization.MOrganizationsPkey,
			&i.Organization.OrganizationID,
			&i.Organization.Name,
			&i.Organization.Description,
			&i.Organization.IsPersonal,
			&i.Organization.IsWhole,
			&i.Organization.CreatedAt,
			&i.Organization.UpdatedAt,
			&i.Group.MGroupsPkey,
			&i.Group.GroupID,
			&i.Group.Key,
			&i.Group.OrganizationID,
			&i.Grade.MGradesPkey,
			&i.Grade.GradeID,
			&i.Grade.Key,
			&i.Grade.OrganizationID,
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

const getOrganizationsWithDetailUseNumberedPaginate = `-- name: GetOrganizationsWithDetailUseNumberedPaginate :many
SELECT m_organizations.m_organizations_pkey, m_organizations.organization_id, m_organizations.name, m_organizations.description, m_organizations.is_personal, m_organizations.is_whole, m_organizations.created_at, m_organizations.updated_at, m_groups.m_groups_pkey, m_groups.group_id, m_groups.key, m_groups.organization_id, m_grades.m_grades_pkey, m_grades.grade_id, m_grades.key, m_grades.organization_id FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
WHERE
	CASE WHEN $3::boolean = true THEN m_organizations.name LIKE '%' || $4::text || '%' END
AND
	CASE WHEN $5::boolean = true THEN m_organizations.is_whole = $6 END
AND
	CASE WHEN $7::boolean = true THEN m_organizations.is_personal = $8 AND EXISTS (SELECT m_members_pkey, member_id, login_id, password, email, name, attend_status_id, profile_image_id, grade_id, group_id, personal_organization_id, role_id, created_at, updated_at FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = $9::uuid) END
AND
	CASE WHEN $10::boolean = true THEN EXISTS (SELECT m_groups_pkey, group_id, key, organization_id FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) END
AND
	CASE WHEN $11::boolean = true THEN EXISTS (SELECT m_grades_pkey, grade_id, key, organization_id FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) END
ORDER BY
	CASE WHEN $12::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN $12::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey DESC
LIMIT $1 OFFSET $2
`

type GetOrganizationsWithDetailUseNumberedPaginateParams struct {
	Limit            int32     `json:"limit"`
	Offset           int32     `json:"offset"`
	WhereLikeName    bool      `json:"where_like_name"`
	SearchName       string    `json:"search_name"`
	WhereIsWhole     bool      `json:"where_is_whole"`
	IsWhole          bool      `json:"is_whole"`
	WhereIsPersonal  bool      `json:"where_is_personal"`
	IsPersonal       bool      `json:"is_personal"`
	PersonalMemberID uuid.UUID `json:"personal_member_id"`
	WhereIsGroup     bool      `json:"where_is_group"`
	WhereIsGrade     bool      `json:"where_is_grade"`
	OrderMethod      string    `json:"order_method"`
}

type GetOrganizationsWithDetailUseNumberedPaginateRow struct {
	Organization Organization `json:"organization"`
	Group        Group        `json:"group"`
	Grade        Grade        `json:"grade"`
}

func (q *Queries) GetOrganizationsWithDetailUseNumberedPaginate(ctx context.Context, arg GetOrganizationsWithDetailUseNumberedPaginateParams) ([]GetOrganizationsWithDetailUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getOrganizationsWithDetailUseNumberedPaginate,
		arg.Limit,
		arg.Offset,
		arg.WhereLikeName,
		arg.SearchName,
		arg.WhereIsWhole,
		arg.IsWhole,
		arg.WhereIsPersonal,
		arg.IsPersonal,
		arg.PersonalMemberID,
		arg.WhereIsGroup,
		arg.WhereIsGrade,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetOrganizationsWithDetailUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetOrganizationsWithDetailUseNumberedPaginateRow
		if err := rows.Scan(
			&i.Organization.MOrganizationsPkey,
			&i.Organization.OrganizationID,
			&i.Organization.Name,
			&i.Organization.Description,
			&i.Organization.IsPersonal,
			&i.Organization.IsWhole,
			&i.Organization.CreatedAt,
			&i.Organization.UpdatedAt,
			&i.Group.MGroupsPkey,
			&i.Group.GroupID,
			&i.Group.Key,
			&i.Group.OrganizationID,
			&i.Grade.MGradesPkey,
			&i.Grade.GradeID,
			&i.Grade.Key,
			&i.Grade.OrganizationID,
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

const getPluralOrganizations = `-- name: GetPluralOrganizations :many
SELECT m_organizations_pkey, organization_id, name, description, is_personal, is_whole, created_at, updated_at FROM m_organizations WHERE organization_id = ANY($3::uuid[])
ORDER BY
	m_organizations_pkey DESC
LIMIT $1 OFFSET $2
`

type GetPluralOrganizationsParams struct {
	Limit           int32       `json:"limit"`
	Offset          int32       `json:"offset"`
	OrganizationIds []uuid.UUID `json:"organization_ids"`
}

func (q *Queries) GetPluralOrganizations(ctx context.Context, arg GetPluralOrganizationsParams) ([]Organization, error) {
	rows, err := q.db.Query(ctx, getPluralOrganizations, arg.Limit, arg.Offset, arg.OrganizationIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Organization{}
	for rows.Next() {
		var i Organization
		if err := rows.Scan(
			&i.MOrganizationsPkey,
			&i.OrganizationID,
			&i.Name,
			&i.Description,
			&i.IsPersonal,
			&i.IsWhole,
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

const getPluralOrganizationsWithDetail = `-- name: GetPluralOrganizationsWithDetail :many
SELECT m_organizations.m_organizations_pkey, m_organizations.organization_id, m_organizations.name, m_organizations.description, m_organizations.is_personal, m_organizations.is_whole, m_organizations.created_at, m_organizations.updated_at, m_groups.m_groups_pkey, m_groups.group_id, m_groups.key, m_groups.organization_id, m_grades.m_grades_pkey, m_grades.grade_id, m_grades.key, m_grades.organization_id FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
WHERE organization_id = ANY($3::uuid[])
ORDER BY
	m_organizations_pkey DESC
LIMIT $1 OFFSET $2
`

type GetPluralOrganizationsWithDetailParams struct {
	Limit           int32       `json:"limit"`
	Offset          int32       `json:"offset"`
	OrganizationIds []uuid.UUID `json:"organization_ids"`
}

type GetPluralOrganizationsWithDetailRow struct {
	Organization Organization `json:"organization"`
	Group        Group        `json:"group"`
	Grade        Grade        `json:"grade"`
}

func (q *Queries) GetPluralOrganizationsWithDetail(ctx context.Context, arg GetPluralOrganizationsWithDetailParams) ([]GetPluralOrganizationsWithDetailRow, error) {
	rows, err := q.db.Query(ctx, getPluralOrganizationsWithDetail, arg.Limit, arg.Offset, arg.OrganizationIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralOrganizationsWithDetailRow{}
	for rows.Next() {
		var i GetPluralOrganizationsWithDetailRow
		if err := rows.Scan(
			&i.Organization.MOrganizationsPkey,
			&i.Organization.OrganizationID,
			&i.Organization.Name,
			&i.Organization.Description,
			&i.Organization.IsPersonal,
			&i.Organization.IsWhole,
			&i.Organization.CreatedAt,
			&i.Organization.UpdatedAt,
			&i.Group.MGroupsPkey,
			&i.Group.GroupID,
			&i.Group.Key,
			&i.Group.OrganizationID,
			&i.Grade.MGradesPkey,
			&i.Grade.GradeID,
			&i.Grade.Key,
			&i.Grade.OrganizationID,
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

const updateOrganization = `-- name: UpdateOrganization :one
UPDATE m_organizations SET name = $2, description = $3, updated_at = $4 WHERE organization_id = $1 RETURNING m_organizations_pkey, organization_id, name, description, is_personal, is_whole, created_at, updated_at
`

type UpdateOrganizationParams struct {
	OrganizationID uuid.UUID   `json:"organization_id"`
	Name           string      `json:"name"`
	Description    pgtype.Text `json:"description"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

func (q *Queries) UpdateOrganization(ctx context.Context, arg UpdateOrganizationParams) (Organization, error) {
	row := q.db.QueryRow(ctx, updateOrganization,
		arg.OrganizationID,
		arg.Name,
		arg.Description,
		arg.UpdatedAt,
	)
	var i Organization
	err := row.Scan(
		&i.MOrganizationsPkey,
		&i.OrganizationID,
		&i.Name,
		&i.Description,
		&i.IsPersonal,
		&i.IsWhole,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

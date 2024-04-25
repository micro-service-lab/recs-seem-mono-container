-- name: CreateOrganizations :copyfrom
INSERT INTO m_organizations (name, description, is_personal, is_whole, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6);

-- name: CreateOrganization :one
INSERT INTO m_organizations (name, description, is_personal, is_whole, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: UpdateOrganization :one
UPDATE m_organizations SET name = $2, description = $3, updated_at = $4 WHERE organization_id = $1 RETURNING *;

-- name: DeleteOrganization :exec
DELETE FROM m_organizations WHERE organization_id = $1;

-- name: FindOrganizationByID :one
SELECT * FROM m_organizations WHERE organization_id = $1;

-- name: FindOrganizationByIDWithDetail :one
SELECT sqlc.embed(m_organizations), sqlc.embed(m_groups), sqlc.embed(m_grades) FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
WHERE m_organizations.organization_id = $1;

-- name: FindWholeOrganization :one
SELECT * FROM m_organizations WHERE is_whole = true;

-- name: FindPersonalOrganization :one
SELECT * FROM m_organizations
LEFT JOIN m_members ON m_organizations.organization_id = m_members.personal_organization_id
WHERE m_organizations.is_personal = true AND m_members.member_id = $1;

-- name: GetOrganizations :many
SELECT * FROM m_organizations
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey DESC;

-- name: GetOrganizationsUseNumberedPaginate :many
SELECT * FROM m_organizations
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetOrganizationsUseKeysetPaginate :many
SELECT * FROM m_organizations
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) END
AND
	CASE @cursor_direction
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_organizations.name > @cursor_column OR (m_organizations.name = @cursor_column AND m_organizations_pkey < @cursor)
				WHEN 'r_name' THEN m_organizations.name < @cursor_column OR (m_organizations.name = @cursor_column AND m_organizations_pkey < @cursor)
				ELSE m_organizations_pkey < @cursor
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_organizations.name < @cursor_column OR (m_organizations.name = @cursor_column AND m_organizations_pkey > @cursor)
				WHEN 'r_name' THEN m_organizations.name > @cursor_column OR (m_organizations.name = @cursor_column AND m_organizations_pkey > @cursor)
				ELSE m_organizations_pkey > @cursor
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey DESC
LIMIT $1;

-- name: GetOrganizationsWithDetail :many
SELECT sqlc.embed(m_organizations), sqlc.embed(m_groups), sqlc.embed(m_grades) FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey DESC;

-- name: GetOrganizationsWithDetailUseNumberedPaginate :many
SELECT sqlc.embed(m_organizations), sqlc.embed(m_groups), sqlc.embed(m_grades) FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetOrganizationsWithDetailUseKeysetPaginate :many
SELECT sqlc.embed(m_organizations), sqlc.embed(m_groups), sqlc.embed(m_grades) FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) END
AND
	CASE @cursor_direction
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_organizations.name > @cursor_column OR (m_organizations.name = @cursor_column AND m_organizations_pkey < @cursor)
				WHEN 'r_name' THEN m_organizations.name < @cursor_column OR (m_organizations.name = @cursor_column AND m_organizations_pkey < @cursor)
				ELSE m_organizations_pkey < @cursor
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_organizations.name < @cursor_column OR (m_organizations.name = @cursor_column AND m_organizations_pkey > @cursor)
				WHEN 'r_name' THEN m_organizations.name > @cursor_column OR (m_organizations.name = @cursor_column AND m_organizations_pkey > @cursor)
				ELSE m_organizations_pkey > @cursor
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey DESC
LIMIT $1;

-- name: CountOrganizations :one
SELECT COUNT(*) FROM m_organizations
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) END;
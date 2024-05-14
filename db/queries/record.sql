-- name: CreateRecords :copyfrom
INSERT INTO t_records (record_type_id, title, body, organization_id, posted_by, last_edited_by, posted_at, last_edited_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: CreateRecord :one
INSERT INTO t_records (record_type_id, title, body, organization_id, posted_by, last_edited_by, posted_at, last_edited_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: UpdateRecord :one
UPDATE t_records SET record_type_id = $2, title = $3, body = $4, organization_id = $5, last_edited_by = $6, last_edited_at = $7 WHERE record_id = $1 RETURNING *;

-- name: DeleteRecord :execrows
DELETE FROM t_records WHERE record_id = $1;

-- name: DeleteRecordOnOrganization :execrows
DELETE FROM t_records WHERE organization_id = $1;

-- name: PluralDeleteRecords :execrows
DELETE FROM t_records WHERE record_id = ANY($1::uuid[]);

-- name: FindRecordByID :one
SELECT * FROM t_records WHERE record_id = $1;

-- name: FindRecordByIDWithRecordType :one
SELECT sqlc.embed(t_records), sqlc.embed(m_record_types) FROM t_records
LEFT JOIN m_record_types ON t_records.record_type_id = m_record_types.record_type_id
WHERE record_id = $1;

-- name: FindRecordByIDWithOrganization :one
SELECT sqlc.embed(t_records), sqlc.embed(m_organizations) FROM t_records
LEFT JOIN m_organizations ON t_records.organization_id = m_organizations.organization_id
WHERE record_id = $1;

-- name: FindRecordByIDWithPostedBy :one
SELECT sqlc.embed(t_records), sqlc.embed(m_members) FROM t_records
LEFT JOIN m_members ON t_records.posted_by = m_members.member_id
WHERE record_id = $1;

-- name: FindRecordByIDWithLastEditedBy :one
SELECT sqlc.embed(t_records), sqlc.embed(m_members) FROM t_records
LEFT JOIN m_members ON t_records.last_edited_by = m_members.member_id
WHERE record_id = $1;

-- name: FindRecordByIDWithAll :one
SELECT sqlc.embed(t_records), sqlc.embed(m_record_types), sqlc.embed(m_organizations), sqlc.embed(m_members) FROM t_records
LEFT JOIN m_record_types ON t_records.record_type_id = m_record_types.record_type_id
LEFT JOIN m_organizations ON t_records.organization_id = m_organizations.organization_id
LEFT JOIN m_members ON t_records.posted_by = m_members.member_id
LEFT JOIN m_members AS m_members_2 ON t_records.last_edited_by = m_members_2.member_id
WHERE record_id = $1;

-- name: GetRecords :many
SELECT * FROM t_records
WHERE
	CASE WHEN @where_in_record_type::boolean = true THEN record_type_id = ANY(@in_record_type) ELSE TRUE END
AND
	CASE WHEN @where_like_title::boolean = true THEN title LIKE '%' || @search_title::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_organization::boolean = true THEN organization_id = ANY(@in_organization) ELSE TRUE END
AND
	CASE WHEN @where_in_posted_by::boolean = true THEN posted_by = ANY(@in_posted_by) ELSE TRUE END
AND
	CASE WHEN @where_in_last_edited_by::boolean = true THEN last_edited_by = ANY(@in_last_edited_by) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'title' THEN title END ASC,
	CASE WHEN @order_method::text = 'r_title' THEN title END DESC,
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC,
	t_records_pkey ASC;

-- name: GetRecordsUseNumberedPaginate :many
SELECT * FROM t_records
WHERE
	CASE WHEN @where_in_record_type::boolean = true THEN record_type_id = ANY(@in_record_type) ELSE TRUE END
AND
	CASE WHEN @where_like_title::boolean = true THEN title LIKE '%' || @search_title::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_organization::boolean = true THEN organization_id = ANY(@in_organization) ELSE TRUE END
AND
	CASE WHEN @where_in_posted_by::boolean = true THEN posted_by = ANY(@in_posted_by) ELSE TRUE END
AND
	CASE WHEN @where_in_last_edited_by::boolean = true THEN last_edited_by = ANY(@in_last_edited_by) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'title' THEN title END ASC,
	CASE WHEN @order_method::text = 'r_title' THEN title END DESC,
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC,
	t_records_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetRecordsUseKeysetPaginate :many
SELECT * FROM t_records
WHERE
	CASE WHEN @where_in_record_type::boolean = true THEN record_type_id = ANY(@in_record_type) ELSE TRUE END
AND
	CASE WHEN @where_like_title::boolean = true THEN title LIKE '%' || @search_title::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_organization::boolean = true THEN organization_id = ANY(@in_organization) ELSE TRUE END
AND
	CASE WHEN @where_in_posted_by::boolean = true THEN posted_by = ANY(@in_posted_by) ELSE TRUE END
AND
	CASE WHEN @where_in_last_edited_by::boolean = true THEN last_edited_by = ANY(@in_last_edited_by) ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'title' THEN title > @title_cursor OR (title = @title_cursor AND t_records_pkey > @cursor::int)
				WHEN 'r_title' THEN title < @title_cursor OR (title = @title_cursor AND t_records_pkey > @cursor::int)
				WHEN 'posted_at' THEN posted_at > @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_records_pkey > @cursor::int)
				WHEN 'r_posted_at' THEN posted_at < @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_records_pkey > @cursor::int)
				WHEN 'last_edited_at' THEN last_edited_at > @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_records_pkey > @cursor::int)
				WHEN 'r_last_edited_at' THEN last_edited_at < @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_records_pkey > @cursor::int)
				ELSE t_records_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'title' THEN title < @title_cursor OR (title = @title_cursor AND t_records_pkey < @cursor::int)
				WHEN 'r_title' THEN title > @title_cursor OR (title = @title_cursor AND t_records_pkey < @cursor::int)
				WHEN 'posted_at' THEN posted_at < @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_records_pkey < @cursor::int)
				WHEN 'r_posted_at' THEN posted_at > @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_records_pkey < @cursor::int)
				WHEN 'last_edited_at' THEN last_edited_at < @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_records_pkey < @cursor::int)
				WHEN 'r_last_edited_at' THEN last_edited_at > @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_records_pkey < @cursor::int)
				ELSE t_records_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'title' AND @cursor_direction::text = 'next' THEN title END ASC,
	CASE WHEN @order_method::text = 'title' AND @cursor_direction::text = 'prev' THEN title END DESC,
	CASE WHEN @order_method::text = 'r_title' AND @cursor_direction::text = 'next' THEN title END DESC,
	CASE WHEN @order_method::text = 'r_title' AND @cursor_direction::text = 'prev' THEN title END ASC,
	CASE WHEN @order_method::text = 'posted_at' AND @cursor_direction::text = 'next' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'posted_at' AND @cursor_direction::text = 'prev' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'r_posted_at' AND @cursor_direction::text = 'next' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'r_posted_at' AND @cursor_direction::text = 'prev' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'last_edited_at' AND @cursor_direction::text = 'next' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'last_edited_at' AND @cursor_direction::text = 'prev' THEN last_edited_at END DESC,
	CASE WHEN @order_method::text = 'r_last_edited_at' AND @cursor_direction::text = 'next' THEN last_edited_at END DESC,
	CASE WHEN @order_method::text = 'r_last_edited_at' AND @cursor_direction::text = 'prev' THEN last_edited_at END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN t_records_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_records_pkey END DESC
LIMIT $1;

-- name: GetPluralRecords :many
SELECT * FROM t_records WHERE record_id = ANY(@record_ids::uuid[])
ORDER BY
	t_records_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetRecordsWithRecordType :many
SELECT sqlc.embed(t_records), sqlc.embed(m_record_types) FROM t_records
LEFT JOIN m_record_types ON t_records.record_type_id = m_record_types.record_type_id
WHERE
	CASE WHEN @where_in_record_type::boolean = true THEN t_records.record_type_id = ANY(@in_record_type) ELSE TRUE END
AND
	CASE WHEN @where_like_title::boolean = true THEN title LIKE '%' || @search_title::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_organization::boolean = true THEN organization_id = ANY(@in_organization) ELSE TRUE END
AND
	CASE WHEN @where_in_posted_by::boolean = true THEN posted_by = ANY(@in_posted_by) ELSE TRUE END
AND
	CASE WHEN @where_in_last_edited_by::boolean = true THEN last_edited_by = ANY(@in_last_edited_by) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'title' THEN title END ASC,
	CASE WHEN @order_method::text = 'r_title' THEN title END DESC,
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC,
	t_records_pkey ASC;

-- name: GetRecordsWithRecordTypeUseNumberedPaginate :many
SELECT sqlc.embed(t_records), sqlc.embed(m_record_types) FROM t_records
LEFT JOIN m_record_types ON t_records.record_type_id = m_record_types.record_type_id
WHERE
	CASE WHEN @where_in_record_type::boolean = true THEN t_records.record_type_id = ANY(@in_record_type) ELSE TRUE END
AND
	CASE WHEN @where_like_title::boolean = true THEN title LIKE '%' || @search_title::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_organization::boolean = true THEN organization_id = ANY(@in_organization) ELSE TRUE END
AND
	CASE WHEN @where_in_posted_by::boolean = true THEN posted_by = ANY(@in_posted_by) ELSE TRUE END
AND
	CASE WHEN @where_in_last_edited_by::boolean = true THEN last_edited_by = ANY(@in_last_edited_by) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'title' THEN title END ASC,
	CASE WHEN @order_method::text = 'r_title' THEN title END DESC,
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC,
	t_records_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetRecordsWithRecordTypeUseKeysetPaginate :many
SELECT sqlc.embed(t_records), sqlc.embed(m_record_types) FROM t_records
LEFT JOIN m_record_types ON t_records.record_type_id = m_record_types.record_type_id
WHERE
	CASE WHEN @where_in_record_type::boolean = true THEN t_records.record_type_id = ANY(@in_record_type) ELSE TRUE END
AND
	CASE WHEN @where_like_title::boolean = true THEN title LIKE '%' || @search_title::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_organization::boolean = true THEN organization_id = ANY(@in_organization) ELSE TRUE END
AND
	CASE WHEN @where_in_posted_by::boolean = true THEN posted_by = ANY(@in_posted_by) ELSE TRUE END
AND
	CASE WHEN @where_in_last_edited_by::boolean = true THEN last_edited_by = ANY(@in_last_edited_by) ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'title' THEN title > @title_cursor OR (title = @title_cursor AND t_records_pkey > @cursor::int)
				WHEN 'r_title' THEN title < @title_cursor OR (title = @title_cursor AND t_records_pkey > @cursor::int)
				WHEN 'posted_at' THEN posted_at > @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_records_pkey > @cursor::int)
				WHEN 'r_posted_at' THEN posted_at < @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_records_pkey > @cursor::int)
				WHEN 'last_edited_at' THEN last_edited_at > @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_records_pkey > @cursor::int)
				WHEN 'r_last_edited_at' THEN last_edited_at < @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_records_pkey > @cursor::int)
				ELSE t_records_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'title' THEN title < @title_cursor OR (title = @title_cursor AND t_records_pkey < @cursor::int)
				WHEN 'r_title' THEN title > @title_cursor OR (title = @title_cursor AND t_records_pkey < @cursor::int)
				WHEN 'posted_at' THEN posted_at < @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_records_pkey < @cursor::int)
				WHEN 'r_posted_at' THEN posted_at > @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_records_pkey < @cursor::int)
				WHEN 'last_edited_at' THEN last_edited_at < @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_records_pkey < @cursor::int)
				WHEN 'r_last_edited_at' THEN last_edited_at > @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_records_pkey < @cursor::int)
				ELSE t_records_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'title' AND @cursor_direction::text = 'next' THEN title END ASC,
	CASE WHEN @order_method::text = 'title' AND @cursor_direction::text = 'prev' THEN title END DESC,
	CASE WHEN @order_method::text = 'r_title' AND @cursor_direction::text = 'next' THEN title END DESC,
	CASE WHEN @order_method::text = 'r_title' AND @cursor_direction::text = 'prev' THEN title END ASC,
	CASE WHEN @order_method::text = 'posted_at' AND @cursor_direction::text = 'next' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'posted_at' AND @cursor_direction::text = 'prev' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'r_posted_at' AND @cursor_direction::text = 'next' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'r_posted_at' AND @cursor_direction::text = 'prev' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'last_edited_at' AND @cursor_direction::text = 'next' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'last_edited_at' AND @cursor_direction::text = 'prev' THEN last_edited_at END DESC,
	CASE WHEN @order_method::text = 'r_last_edited_at' AND @cursor_direction::text = 'next' THEN last_edited_at END DESC,
	CASE WHEN @order_method::text = 'r_last_edited_at' AND @cursor_direction::text = 'prev' THEN last_edited_at END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN t_records_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_records_pkey END DESC
LIMIT $1;

-- name: GetPluralRecordsWithRecordType :many
SELECT sqlc.embed(t_records), sqlc.embed(m_record_types) FROM t_records
LEFT JOIN m_record_types ON t_records.record_type_id = m_record_types.record_type_id
WHERE record_id = ANY(@record_ids::uuid[])
ORDER BY
	t_records_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetRecordsWithOrganization :many
SELECT sqlc.embed(t_records), sqlc.embed(m_organizations) FROM t_records
LEFT JOIN m_organizations ON t_records.organization_id = m_organizations.organization_id
WHERE
	CASE WHEN @where_in_record_type::boolean = true THEN record_type_id = ANY(@in_record_type) ELSE TRUE END
AND
	CASE WHEN @where_like_title::boolean = true THEN title LIKE '%' || @search_title::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_organization::boolean = true THEN t_records.organization_id = ANY(@in_organization) ELSE TRUE END
AND
	CASE WHEN @where_in_posted_by::boolean = true THEN posted_by = ANY(@in_posted_by) ELSE TRUE END
AND
	CASE WHEN @where_in_last_edited_by::boolean = true THEN last_edited_by = ANY(@in_last_edited_by) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'title' THEN title END ASC,
	CASE WHEN @order_method::text = 'r_title' THEN title END DESC,
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC,
	t_records_pkey ASC;

-- name: GetRecordsWithOrganizationUseNumberedPaginate :many
SELECT sqlc.embed(t_records), sqlc.embed(m_organizations) FROM t_records
LEFT JOIN m_organizations ON t_records.organization_id = m_organizations.organization_id
WHERE
	CASE WHEN @where_in_record_type::boolean = true THEN record_type_id = ANY(@in_record_type) ELSE TRUE END
AND
	CASE WHEN @where_like_title::boolean = true THEN title LIKE '%' || @search_title::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_organization::boolean = true THEN t_records.organization_id = ANY(@in_organization) ELSE TRUE END
AND
	CASE WHEN @where_in_posted_by::boolean = true THEN posted_by = ANY(@in_posted_by) ELSE TRUE END
AND
	CASE WHEN @where_in_last_edited_by::boolean = true THEN last_edited_by = ANY(@in_last_edited_by) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'title' THEN title END ASC,
	CASE WHEN @order_method::text = 'r_title' THEN title END DESC,
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC,
	t_records_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetRecordsWithOrganizationUseKeysetPaginate :many
SELECT sqlc.embed(t_records), sqlc.embed(m_organizations) FROM t_records
LEFT JOIN m_organizations ON t_records.organization_id = m_organizations.organization_id
WHERE
	CASE WHEN @where_in_record_type::boolean = true THEN record_type_id = ANY(@in_record_type) ELSE TRUE END
AND
	CASE WHEN @where_like_title::boolean = true THEN title LIKE '%' || @search_title::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_organization::boolean = true THEN t_records.organization_id = ANY(@in_organization) ELSE TRUE END
AND
	CASE WHEN @where_in_posted_by::boolean = true THEN posted_by = ANY(@in_posted_by) ELSE TRUE END
AND
	CASE WHEN @where_in_last_edited_by::boolean = true THEN last_edited_by = ANY(@in_last_edited_by) ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'title' THEN title > @title_cursor OR (title = @title_cursor AND t_records_pkey > @cursor::int)
				WHEN 'r_title' THEN title < @title_cursor OR (title = @title_cursor AND t_records_pkey > @cursor::int)
				WHEN 'posted_at' THEN posted_at > @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_records_pkey > @cursor::int)
				WHEN 'r_posted_at' THEN posted_at < @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_records_pkey > @cursor::int)
				WHEN 'last_edited_at' THEN last_edited_at > @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_records_pkey > @cursor::int)
				WHEN 'r_last_edited_at' THEN last_edited_at < @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_records_pkey > @cursor::int)
				ELSE t_records_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'title' THEN title < @title_cursor OR (title = @title_cursor AND t_records_pkey < @cursor::int)
				WHEN 'r_title' THEN title > @title_cursor OR (title = @title_cursor AND t_records_pkey < @cursor::int)
				WHEN 'posted_at' THEN posted_at < @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_records_pkey < @cursor::int)
				WHEN 'r_posted_at' THEN posted_at > @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_records_pkey < @cursor::int)
				WHEN 'last_edited_at' THEN last_edited_at < @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_records_pkey < @cursor::int)
				WHEN 'r_last_edited_at' THEN last_edited_at > @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_records_pkey < @cursor::int)
				ELSE t_records_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'title' AND @cursor_direction::text = 'next' THEN title END ASC,
	CASE WHEN @order_method::text = 'title' AND @cursor_direction::text = 'prev' THEN title END DESC,
	CASE WHEN @order_method::text = 'r_title' AND @cursor_direction::text = 'next' THEN title END DESC,
	CASE WHEN @order_method::text = 'r_title' AND @cursor_direction::text = 'prev' THEN title END ASC,
	CASE WHEN @order_method::text = 'posted_at' AND @cursor_direction::text = 'next' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'posted_at' AND @cursor_direction::text = 'prev' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'r_posted_at' AND @cursor_direction::text = 'next' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'r_posted_at' AND @cursor_direction::text = 'prev' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'last_edited_at' AND @cursor_direction::text = 'next' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'last_edited_at' AND @cursor_direction::text = 'prev' THEN last_edited_at END DESC,
	CASE WHEN @order_method::text = 'r_last_edited_at' AND @cursor_direction::text = 'next' THEN last_edited_at END DESC,
	CASE WHEN @order_method::text = 'r_last_edited_at' AND @cursor_direction::text = 'prev' THEN last_edited_at END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN t_records_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_records_pkey END DESC
LIMIT $1;

-- name: GetPluralRecordsWithOrganization :many
SELECT sqlc.embed(t_records), sqlc.embed(m_organizations) FROM t_records
LEFT JOIN m_organizations ON t_records.organization_id = m_organizations.organization_id
WHERE record_id = ANY(@record_ids::uuid[])
ORDER BY
	t_records_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetRecordsWithPostedBy :many
SELECT sqlc.embed(t_records), sqlc.embed(m_members) FROM t_records
LEFT JOIN m_members ON t_records.posted_by = m_members.member_id
WHERE
	CASE WHEN @where_in_record_type::boolean = true THEN record_type_id = ANY(@in_record_type) ELSE TRUE END
AND
	CASE WHEN @where_like_title::boolean = true THEN title LIKE '%' || @search_title::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_organization::boolean = true THEN organization_id = ANY(@in_organization) ELSE TRUE END
AND
	CASE WHEN @where_in_posted_by::boolean = true THEN t_records.posted_by = ANY(@in_posted_by) ELSE TRUE END
AND
	CASE WHEN @where_in_last_edited_by::boolean = true THEN last_edited_by = ANY(@in_last_edited_by) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'title' THEN title END ASC,
	CASE WHEN @order_method::text = 'r_title' THEN title END DESC,
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC,
	t_records_pkey ASC;

-- name: GetRecordsWithPostedByUseNumberedPaginate :many
SELECT sqlc.embed(t_records), sqlc.embed(m_members) FROM t_records
LEFT JOIN m_members ON t_records.posted_by = m_members.member_id
WHERE
	CASE WHEN @where_in_record_type::boolean = true THEN record_type_id = ANY(@in_record_type) ELSE TRUE END
AND
	CASE WHEN @where_like_title::boolean = true THEN title LIKE '%' || @search_title::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_organization::boolean = true THEN organization_id = ANY(@in_organization) ELSE TRUE END
AND
	CASE WHEN @where_in_posted_by::boolean = true THEN t_records.posted_by = ANY(@in_posted_by) ELSE TRUE END
AND
	CASE WHEN @where_in_last_edited_by::boolean = true THEN last_edited_by = ANY(@in_last_edited_by) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'title' THEN title END ASC,
	CASE WHEN @order_method::text = 'r_title' THEN title END DESC,
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC,
	t_records_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetRecordsWithPostedByUseKeysetPaginate :many
SELECT sqlc.embed(t_records), sqlc.embed(m_members) FROM t_records
LEFT JOIN m_members ON t_records.posted_by = m_members.member_id
WHERE
	CASE WHEN @where_in_record_type::boolean = true THEN record_type_id = ANY(@in_record_type) ELSE TRUE END
AND
	CASE WHEN @where_like_title::boolean = true THEN title LIKE '%' || @search_title::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_organization::boolean = true THEN organization_id = ANY(@in_organization) ELSE TRUE END
AND
	CASE WHEN @where_in_posted_by::boolean = true THEN t_records.posted_by = ANY(@in_posted_by) ELSE TRUE END
AND
	CASE WHEN @where_in_last_edited_by::boolean = true THEN last_edited_by = ANY(@in_last_edited_by) ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'title' THEN title > @title_cursor OR (title = @title_cursor AND t_records_pkey > @cursor::int)
				WHEN 'r_title' THEN title < @title_cursor OR (title = @title_cursor AND t_records_pkey > @cursor::int)
				WHEN 'posted_at' THEN posted_at > @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_records_pkey > @cursor::int)
				WHEN 'r_posted_at' THEN posted_at < @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_records_pkey > @cursor::int)
				WHEN 'last_edited_at' THEN last_edited_at > @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_records_pkey > @cursor::int)
				WHEN 'r_last_edited_at' THEN last_edited_at < @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_records_pkey > @cursor::int)
				ELSE t_records_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'title' THEN title < @title_cursor OR (title = @title_cursor AND t_records_pkey < @cursor::int)
				WHEN 'r_title' THEN title > @title_cursor OR (title = @title_cursor AND t_records_pkey < @cursor::int)
				WHEN 'posted_at' THEN posted_at < @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_records_pkey < @cursor::int)
				WHEN 'r_posted_at' THEN posted_at > @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_records_pkey < @cursor::int)
				WHEN 'last_edited_at' THEN last_edited_at < @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_records_pkey < @cursor::int)
				WHEN 'r_last_edited_at' THEN last_edited_at > @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_records_pkey < @cursor::int)
				ELSE t_records_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'title' AND @cursor_direction::text = 'next' THEN title END ASC,
	CASE WHEN @order_method::text = 'title' AND @cursor_direction::text = 'prev' THEN title END DESC,
	CASE WHEN @order_method::text = 'r_title' AND @cursor_direction::text = 'next' THEN title END DESC,
	CASE WHEN @order_method::text = 'r_title' AND @cursor_direction::text = 'prev' THEN title END ASC,
	CASE WHEN @order_method::text = 'posted_at' AND @cursor_direction::text = 'next' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'posted_at' AND @cursor_direction::text = 'prev' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'r_posted_at' AND @cursor_direction::text = 'next' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'r_posted_at' AND @cursor_direction::text = 'prev' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'last_edited_at' AND @cursor_direction::text = 'next' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'last_edited_at' AND @cursor_direction::text = 'prev' THEN last_edited_at END DESC,
	CASE WHEN @order_method::text = 'r_last_edited_at' AND @cursor_direction::text = 'next' THEN last_edited_at END DESC,
	CASE WHEN @order_method::text = 'r_last_edited_at' AND @cursor_direction::text = 'prev' THEN last_edited_at END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN t_records_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_records_pkey END DESC
LIMIT $1;

-- name: GetPluralRecordsWithPostedBy :many
SELECT sqlc.embed(t_records), sqlc.embed(m_members) FROM t_records
LEFT JOIN m_members ON t_records.posted_by = m_members.member_id
WHERE record_id = ANY(@record_ids::uuid[])
ORDER BY
	t_records_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetRecordsWithLastEditedBy :many
SELECT sqlc.embed(t_records), sqlc.embed(m_members) FROM t_records
LEFT JOIN m_members ON t_records.last_edited_by = m_members.member_id
WHERE
	CASE WHEN @where_in_record_type::boolean = true THEN record_type_id = ANY(@in_record_type) ELSE TRUE END
AND
	CASE WHEN @where_like_title::boolean = true THEN title LIKE '%' || @search_title::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_organization::boolean = true THEN organization_id = ANY(@in_organization) ELSE TRUE END
AND
	CASE WHEN @where_in_posted_by::boolean = true THEN posted_by = ANY(@in_posted_by) ELSE TRUE END
AND
	CASE WHEN @where_in_last_edited_by::boolean = true THEN t_records.last_edited_by = ANY(@in_last_edited_by) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'title' THEN title END ASC,
	CASE WHEN @order_method::text = 'r_title' THEN title END DESC,
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC,
	t_records_pkey ASC;

-- name: GetRecordsWithLastEditedByUseNumberedPaginate :many
SELECT sqlc.embed(t_records), sqlc.embed(m_members) FROM t_records
LEFT JOIN m_members ON t_records.last_edited_by = m_members.member_id
WHERE
	CASE WHEN @where_in_record_type::boolean = true THEN record_type_id = ANY(@in_record_type) ELSE TRUE END
AND
	CASE WHEN @where_like_title::boolean = true THEN title LIKE '%' || @search_title::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_organization::boolean = true THEN organization_id = ANY(@in_organization) ELSE TRUE END
AND
	CASE WHEN @where_in_posted_by::boolean = true THEN posted_by = ANY(@in_posted_by) ELSE TRUE END
AND
	CASE WHEN @where_in_last_edited_by::boolean = true THEN t_records.last_edited_by = ANY(@in_last_edited_by) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'title' THEN title END ASC,
	CASE WHEN @order_method::text = 'r_title' THEN title END DESC,
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC,
	t_records_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetRecordsWithLastEditedByUseKeysetPaginate :many
SELECT sqlc.embed(t_records), sqlc.embed(m_members) FROM t_records
LEFT JOIN m_members ON t_records.last_edited_by = m_members.member_id
WHERE
	CASE WHEN @where_in_record_type::boolean = true THEN record_type_id = ANY(@in_record_type) ELSE TRUE END
AND
	CASE WHEN @where_like_title::boolean = true THEN title LIKE '%' || @search_title::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_organization::boolean = true THEN organization_id = ANY(@in_organization) ELSE TRUE END
AND
	CASE WHEN @where_in_posted_by::boolean = true THEN posted_by = ANY(@in_posted_by) ELSE TRUE END
AND
	CASE WHEN @where_in_last_edited_by::boolean = true THEN t_records.last_edited_by = ANY(@in_last_edited_by) ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'title' THEN title > @title_cursor OR (title = @title_cursor AND t_records_pkey > @cursor::int)
				WHEN 'r_title' THEN title < @title_cursor OR (title = @title_cursor AND t_records_pkey > @cursor::int)
				WHEN 'posted_at' THEN posted_at > @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_records_pkey > @cursor::int)
				WHEN 'r_posted_at' THEN posted_at < @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_records_pkey > @cursor::int)
				WHEN 'last_edited_at' THEN last_edited_at > @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_records_pkey > @cursor::int)
				WHEN 'r_last_edited_at' THEN last_edited_at < @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_records_pkey > @cursor::int)
				ELSE t_records_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'title' THEN title < @title_cursor OR (title = @title_cursor AND t_records_pkey < @cursor::int)
				WHEN 'r_title' THEN title > @title_cursor OR (title = @title_cursor AND t_records_pkey < @cursor::int)
				WHEN 'posted_at' THEN posted_at < @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_records_pkey < @cursor::int)
				WHEN 'r_posted_at' THEN posted_at > @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_records_pkey < @cursor::int)
				WHEN 'last_edited_at' THEN last_edited_at < @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_records_pkey < @cursor::int)
				WHEN 'r_last_edited_at' THEN last_edited_at > @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_records_pkey < @cursor::int)
				ELSE t_records_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'title' AND @cursor_direction::text = 'next' THEN title END ASC,
	CASE WHEN @order_method::text = 'title' AND @cursor_direction::text = 'prev' THEN title END DESC,
	CASE WHEN @order_method::text = 'r_title' AND @cursor_direction::text = 'next' THEN title END DESC,
	CASE WHEN @order_method::text = 'r_title' AND @cursor_direction::text = 'prev' THEN title END ASC,
	CASE WHEN @order_method::text = 'posted_at' AND @cursor_direction::text = 'next' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'posted_at' AND @cursor_direction::text = 'prev' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'r_posted_at' AND @cursor_direction::text = 'next' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'r_posted_at' AND @cursor_direction::text = 'prev' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'last_edited_at' AND @cursor_direction::text = 'next' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'last_edited_at' AND @cursor_direction::text = 'prev' THEN last_edited_at END DESC,
	CASE WHEN @order_method::text = 'r_last_edited_at' AND @cursor_direction::text = 'next' THEN last_edited_at END DESC,
	CASE WHEN @order_method::text = 'r_last_edited_at' AND @cursor_direction::text = 'prev' THEN last_edited_at END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN t_records_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_records_pkey END DESC
LIMIT $1;

-- name: GetPluralRecordsWithLastEditedBy :many
SELECT sqlc.embed(t_records), sqlc.embed(m_members) FROM t_records
LEFT JOIN m_members ON t_records.last_edited_by = m_members.member_id
WHERE record_id = ANY(@record_ids::uuid[])
ORDER BY
	t_records_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetRecordsWithAll :many
SELECT sqlc.embed(t_records), sqlc.embed(m_record_types), sqlc.embed(m_organizations), sqlc.embed(m_members), sqlc.embed(m_members) FROM t_records
LEFT JOIN m_record_types ON t_records.record_type_id = m_record_types.record_type_id
LEFT JOIN m_organizations ON t_records.organization_id = m_organizations.organization_id
LEFT JOIN m_members ON t_records.posted_by = m_members.member_id
LEFT JOIN m_members AS m_members_2 ON t_records.last_edited_by = m_members_2.member_id
WHERE
	CASE WHEN @where_in_record_type::boolean = true THEN t_records.record_type_id = ANY(@in_record_type) ELSE TRUE END
AND
	CASE WHEN @where_like_title::boolean = true THEN title LIKE '%' || @search_title::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_organization::boolean = true THEN t_records.organization_id = ANY(@in_organization) ELSE TRUE END
AND
	CASE WHEN @where_in_posted_by::boolean = true THEN t_records.posted_by = ANY(@in_posted_by) ELSE TRUE END
AND
	CASE WHEN @where_in_last_edited_by::boolean = true THEN t_records.last_edited_by = ANY(@in_last_edited_by) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'title' THEN title END ASC,
	CASE WHEN @order_method::text = 'r_title' THEN title END DESC,
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC,
	t_records_pkey ASC;

-- name: GetRecordsWithAllUseNumberedPaginate :many
SELECT sqlc.embed(t_records), sqlc.embed(m_record_types), sqlc.embed(m_organizations), sqlc.embed(m_members), sqlc.embed(m_members) FROM t_records
LEFT JOIN m_record_types ON t_records.record_type_id = m_record_types.record_type_id
LEFT JOIN m_organizations ON t_records.organization_id = m_organizations.organization_id
LEFT JOIN m_members ON t_records.posted_by = m_members.member_id
LEFT JOIN m_members AS m_members_2 ON t_records.last_edited_by = m_members_2.member_id
WHERE
	CASE WHEN @where_in_record_type::boolean = true THEN t_records.record_type_id = ANY(@in_record_type) ELSE TRUE END
AND
	CASE WHEN @where_like_title::boolean = true THEN title LIKE '%' || @search_title::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_organization::boolean = true THEN t_records.organization_id = ANY(@in_organization) ELSE TRUE END
AND
	CASE WHEN @where_in_posted_by::boolean = true THEN t_records.posted_by = ANY(@in_posted_by) ELSE TRUE END
AND
	CASE WHEN @where_in_last_edited_by::boolean = true THEN t_records.last_edited_by = ANY(@in_last_edited_by) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'title' THEN title END ASC,
	CASE WHEN @order_method::text = 'r_title' THEN title END DESC,
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC,
	t_records_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetRecordsWithAllUseKeysetPaginate :many
SELECT sqlc.embed(t_records), sqlc.embed(m_record_types), sqlc.embed(m_organizations), sqlc.embed(m_members), sqlc.embed(m_members) FROM t_records
LEFT JOIN m_record_types ON t_records.record_type_id = m_record_types.record_type_id
LEFT JOIN m_organizations ON t_records.organization_id = m_organizations.organization_id
LEFT JOIN m_members ON t_records.posted_by = m_members.member_id
LEFT JOIN m_members AS m_members_2 ON t_records.last_edited_by = m_members_2.member_id
WHERE
	CASE WHEN @where_in_record_type::boolean = true THEN t_records.record_type_id = ANY(@in_record_type) ELSE TRUE END
AND
	CASE WHEN @where_like_title::boolean = true THEN title LIKE '%' || @search_title::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_organization::boolean = true THEN t_records.organization_id = ANY(@in_organization) ELSE TRUE END
AND
	CASE WHEN @where_in_posted_by::boolean = true THEN t_records.posted_by = ANY(@in_posted_by) ELSE TRUE END
AND
	CASE WHEN @where_in_last_edited_by::boolean = true THEN t_records.last_edited_by = ANY(@in_last_edited_by) ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'title' THEN title > @title_cursor OR (title = @title_cursor AND t_records_pkey > @cursor::int)
				WHEN 'r_title' THEN title < @title_cursor OR (title = @title_cursor AND t_records_pkey > @cursor::int)
				WHEN 'posted_at' THEN posted_at > @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_records_pkey > @cursor::int)
				WHEN 'r_posted_at' THEN posted_at < @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_records_pkey > @cursor::int)
				WHEN 'last_edited_at' THEN last_edited_at > @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_records_pkey > @cursor::int)
				WHEN 'r_last_edited_at' THEN last_edited_at < @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_records_pkey > @cursor::int)
				ELSE t_records_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'title' THEN title < @title_cursor OR (title = @title_cursor AND t_records_pkey < @cursor::int)
				WHEN 'r_title' THEN title > @title_cursor OR (title = @title_cursor AND t_records_pkey < @cursor::int)
				WHEN 'posted_at' THEN posted_at < @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_records_pkey < @cursor::int)
				WHEN 'r_posted_at' THEN posted_at > @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_records_pkey < @cursor::int)
				WHEN 'last_edited_at' THEN last_edited_at < @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_records_pkey < @cursor::int)
				WHEN 'r_last_edited_at' THEN last_edited_at > @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_records_pkey < @cursor::int)
				ELSE t_records_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'title' AND @cursor_direction::text = 'next' THEN title END ASC,
	CASE WHEN @order_method::text = 'title' AND @cursor_direction::text = 'prev' THEN title END DESC,
	CASE WHEN @order_method::text = 'r_title' AND @cursor_direction::text = 'next' THEN title END DESC,
	CASE WHEN @order_method::text = 'r_title' AND @cursor_direction::text = 'prev' THEN title END ASC,
	CASE WHEN @order_method::text = 'posted_at' AND @cursor_direction::text = 'next' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'posted_at' AND @cursor_direction::text = 'prev' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'r_posted_at' AND @cursor_direction::text = 'next' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'r_posted_at' AND @cursor_direction::text = 'prev' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'last_edited_at' AND @cursor_direction::text = 'next' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'last_edited_at' AND @cursor_direction::text = 'prev' THEN last_edited_at END DESC,
	CASE WHEN @order_method::text = 'r_last_edited_at' AND @cursor_direction::text = 'next' THEN last_edited_at END DESC,
	CASE WHEN @order_method::text = 'r_last_edited_at' AND @cursor_direction::text = 'prev' THEN last_edited_at END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN t_records_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_records_pkey END DESC
LIMIT $1;

-- name: GetPluralRecordsWithAll :many
SELECT sqlc.embed(t_records), sqlc.embed(m_record_types), sqlc.embed(m_organizations), sqlc.embed(m_members), sqlc.embed(m_members) FROM t_records
LEFT JOIN m_record_types ON t_records.record_type_id = m_record_types.record_type_id
LEFT JOIN m_organizations ON t_records.organization_id = m_organizations.organization_id
LEFT JOIN m_members ON t_records.posted_by = m_members.member_id
LEFT JOIN m_members AS m_members_2 ON t_records.last_edited_by = m_members_2.member_id
WHERE record_id = ANY(@record_ids::uuid[])
ORDER BY
	t_records_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountRecords :one
SELECT COUNT(*) FROM t_records
WHERE
	CASE WHEN @where_in_record_type::boolean = true THEN record_type_id = ANY(@in_record_type) ELSE TRUE END
AND
	CASE WHEN @where_like_title::boolean = true THEN title LIKE '%' || @search_title::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_organization::boolean = true THEN organization_id = ANY(@in_organization) ELSE TRUE END
AND
	CASE WHEN @where_in_posted_by::boolean = true THEN posted_by = ANY(@in_posted_by) ELSE TRUE END
AND
	CASE WHEN @where_in_last_edited_by::boolean = true THEN last_edited_by = ANY(@in_last_edited_by) ELSE TRUE END;

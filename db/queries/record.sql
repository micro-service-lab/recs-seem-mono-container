-- name: CreateRecords :copyfrom
INSERT INTO t_records (record_type_id, title, body, organization_id, posted_by, last_edited_by, posted_at, last_edited_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: CreateRecord :one
INSERT INTO t_records (record_type_id, title, body, organization_id, posted_by, last_edited_by, posted_at, last_edited_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: UpdateRecord :one
UPDATE t_records SET record_type_id = $2, title = $3, body = $4, organization_id = $5, last_edited_by = $6, last_edited_at = $7 WHERE record_id = $1 RETURNING *;

-- name: DeleteRecord :exec
DELETE FROM t_records WHERE record_id = $1;

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
	CASE WHEN @where_in_organization::boolean = true THEN organization_id = ANY(@in_organization) ELSE TRUE END
AND
	CASE WHEN @where_in_posted_by::boolean = true THEN posted_by = ANY(@in_posted_by) ELSE TRUE END
AND
	CASE WHEN @where_in_last_edited_by::boolean = true THEN last_edited_by = ANY(@in_last_edited_by) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC,
	t_records_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetRecordsWithRecordType :many
SELECT sqlc.embed(t_records), sqlc.embed(m_record_types) FROM t_records
LEFT JOIN m_record_types ON t_records.record_type_id = m_record_types.record_type_id
WHERE
	CASE WHEN @where_in_record_type::boolean = true THEN t_records.record_type_id = ANY(@in_record_type) ELSE TRUE END
AND
	CASE WHEN @where_in_organization::boolean = true THEN organization_id = ANY(@in_organization) ELSE TRUE END
AND
	CASE WHEN @where_in_posted_by::boolean = true THEN posted_by = ANY(@in_posted_by) ELSE TRUE END
AND
	CASE WHEN @where_in_last_edited_by::boolean = true THEN last_edited_by = ANY(@in_last_edited_by) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC,
	t_records_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetRecordsWithOrganization :many
SELECT sqlc.embed(t_records), sqlc.embed(m_organizations) FROM t_records
LEFT JOIN m_organizations ON t_records.organization_id = m_organizations.organization_id
WHERE
	CASE WHEN @where_in_record_type::boolean = true THEN record_type_id = ANY(@in_record_type) ELSE TRUE END
AND
	CASE WHEN @where_in_organization::boolean = true THEN t_records.organization_id = ANY(@in_organization) ELSE TRUE END
AND
	CASE WHEN @where_in_posted_by::boolean = true THEN posted_by = ANY(@in_posted_by) ELSE TRUE END
AND
	CASE WHEN @where_in_last_edited_by::boolean = true THEN last_edited_by = ANY(@in_last_edited_by) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC,
	t_records_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetRecordsWithPostedBy :many
SELECT sqlc.embed(t_records), sqlc.embed(m_members) FROM t_records
LEFT JOIN m_members ON t_records.posted_by = m_members.member_id
WHERE
	CASE WHEN @where_in_record_type::boolean = true THEN record_type_id = ANY(@in_record_type) ELSE TRUE END
AND
	CASE WHEN @where_in_organization::boolean = true THEN organization_id = ANY(@in_organization) ELSE TRUE END
AND
	CASE WHEN @where_in_posted_by::boolean = true THEN t_records.posted_by = ANY(@in_posted_by) ELSE TRUE END
AND
	CASE WHEN @where_in_last_edited_by::boolean = true THEN last_edited_by = ANY(@in_last_edited_by) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC,
	t_records_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetRecordsWithLastEditedBy :many
SELECT sqlc.embed(t_records), sqlc.embed(m_members) FROM t_records
LEFT JOIN m_members ON t_records.last_edited_by = m_members.member_id
WHERE
	CASE WHEN @where_in_record_type::boolean = true THEN record_type_id = ANY(@in_record_type) ELSE TRUE END
AND
	CASE WHEN @where_in_organization::boolean = true THEN organization_id = ANY(@in_organization) ELSE TRUE END
AND
	CASE WHEN @where_in_posted_by::boolean = true THEN posted_by = ANY(@in_posted_by) ELSE TRUE END
AND
	CASE WHEN @where_in_last_edited_by::boolean = true THEN t_records.last_edited_by = ANY(@in_last_edited_by) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC,
	t_records_pkey DESC
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
	CASE WHEN @where_in_organization::boolean = true THEN t_records.organization_id = ANY(@in_organization) ELSE TRUE END
AND
	CASE WHEN @where_in_posted_by::boolean = true THEN t_records.posted_by = ANY(@in_posted_by) ELSE TRUE END
AND
	CASE WHEN @where_in_last_edited_by::boolean = true THEN t_records.last_edited_by = ANY(@in_last_edited_by) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC,
	t_records_pkey DESC
LIMIT $1 OFFSET $2;

-- name: CountRecords :one
SELECT COUNT(*) FROM t_records
WHERE
	CASE WHEN @where_in_record_type::boolean = true THEN record_type_id = ANY(@in_record_type) ELSE TRUE END
AND
	CASE WHEN @where_in_organization::boolean = true THEN organization_id = ANY(@in_organization) ELSE TRUE END
AND
	CASE WHEN @where_in_posted_by::boolean = true THEN posted_by = ANY(@in_posted_by) ELSE TRUE END
AND
	CASE WHEN @where_in_last_edited_by::boolean = true THEN last_edited_by = ANY(@in_last_edited_by) ELSE TRUE END;

-- name: CreateEvents :copyfrom
INSERT INTO t_events (event_type_id, title, description, organization_id, start_time, end_time, mail_send_flag, send_organization_id, posted_by, last_edited_by, posted_at, last_edited_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);

-- name: CreateEvent :one
INSERT INTO t_events (event_type_id, title, description, organization_id, start_time, end_time, mail_send_flag, send_organization_id, posted_by, last_edited_by, posted_at, last_edited_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING *;

-- name: UpdateEvent :one
UPDATE t_events SET event_type_id = $2, title = $3, description = $4, organization_id = $5, start_time = $6, end_time = $7, send_organization_id = $8, posted_by = $9, last_edited_by = $10, posted_at = $11, last_edited_at = $12 WHERE event_id = $1 RETURNING *;

-- name: DeleteEvent :exec
DELETE FROM t_events WHERE event_id = $1;

-- name: FindEventByID :one
SELECT * FROM t_events WHERE event_id = $1;

-- name: FindEventByIDWithOrganization :one
SELECT sqlc.embed(t_events), sqlc.embed(m_organizations) FROM t_events
INNER JOIN m_organizations ON t_events.organization_id = m_organizations.organization_id
WHERE event_id = $1;

-- name: GetEvents :many
SELECT * FROM t_events
WHERE
	CASE WHEN @where_like_title::boolean = true THEN title LIKE '%' || @search_title::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_organization::boolean = true THEN organization_id = @organization_id ELSE TRUE END
AND
	CASE WHEN @where_mail_send_flag::boolean = true THEN mail_send_flag = @mail_send_flag ELSE TRUE END
AND
	CASE WHEN @where_send_organization::boolean = true THEN send_organization_id = @send_organization_id ELSE TRUE END
AND
	CASE WHEN @where_earlier_start_time::boolean = true THEN start_time >= @earlier_start_time ELSE TRUE END
AND
	CASE WHEN @where_later_start_time::boolean = true THEN start_time <= @later_start_time ELSE TRUE END
AND
	CASE WHEN @where_earlier_end_time::boolean = true THEN end_time >= @earlier_end_time ELSE TRUE END
AND
	CASE WHEN @where_later_end_time::boolean = true THEN end_time <= @later_end_time ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'start_time' THEN start_time END ASC,
	t_events_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetEventsWithOrganization :many
SELECT sqlc.embed(t_events), sqlc.embed(m_organizations) FROM t_events
INNER JOIN m_organizations ON t_events.organization_id = m_organizations.organization_id
WHERE
	CASE WHEN @where_like_title::boolean = true THEN title LIKE '%' || @search_title::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_organization::boolean = true THEN t_events.organization_id = @organization_id ELSE TRUE END
AND
	CASE WHEN @where_mail_send_flag::boolean = true THEN mail_send_flag = @mail_send_flag ELSE TRUE END
AND
	CASE WHEN @where_send_organization::boolean = true THEN send_organization_id = @send_organization_id ELSE TRUE END
AND
	CASE WHEN @where_earlier_start_time::boolean = true THEN start_time >= @earlier_start_time ELSE TRUE END
AND
	CASE WHEN @where_later_start_time::boolean = true THEN start_time <= @later_start_time ELSE TRUE END
AND
	CASE WHEN @where_earlier_end_time::boolean = true THEN end_time >= @earlier_end_time ELSE TRUE END
AND
	CASE WHEN @where_later_end_time::boolean = true THEN end_time <= @later_end_time ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'start_time' THEN start_time END ASC,
	t_events_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetEventsWithSendOrganization :many
SELECT sqlc.embed(t_events), sqlc.embed(m_organizations) FROM t_events
INNER JOIN m_organizations ON t_events.send_organization_id = m_organizations.organization_id
WHERE
	CASE WHEN @where_like_title::boolean = true THEN title LIKE '%' || @search_title::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_organization::boolean = true THEN t_events.organization_id = @organization_id ELSE TRUE END
AND
	CASE WHEN @where_mail_send_flag::boolean = true THEN mail_send_flag = @mail_send_flag ELSE TRUE END
AND
	CASE WHEN @where_send_organization::boolean = true THEN send_organization_id = @send_organization_id ELSE TRUE END
AND
	CASE WHEN @where_earlier_start_time::boolean = true THEN start_time >= @earlier_start_time ELSE TRUE END
AND
	CASE WHEN @where_later_start_time::boolean = true THEN start_time <= @later_start_time ELSE TRUE END
AND
	CASE WHEN @where_earlier_end_time::boolean = true THEN end_time >= @earlier_end_time ELSE TRUE END
AND
	CASE WHEN @where_later_end_time::boolean = true THEN end_time <= @later_end_time ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'start_time' THEN start_time END ASC,
	t_events_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetEventsWithAll :many
SELECT sqlc.embed(t_events), sqlc.embed(o), sqlc.embed(s) FROM t_events
INNER JOIN m_organizations o ON t_events.organization_id = m_organizations.organization_id
INNER JOIN m_organizations s ON t_events.send_organization_id = s.organization_id
WHERE
	CASE WHEN @where_like_title::boolean = true THEN title LIKE '%' || @search_title::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_organization::boolean = true THEN t_events.organization_id = @organization_id ELSE TRUE END
AND
	CASE WHEN @where_mail_send_flag::boolean = true THEN mail_send_flag = @mail_send_flag ELSE TRUE END
AND
	CASE WHEN @where_send_organization::boolean = true THEN send_organization_id = @send_organization_id ELSE TRUE END
AND
	CASE WHEN @where_earlier_start_time::boolean = true THEN start_time >= @earlier_start_time ELSE TRUE END
AND
	CASE WHEN @where_later_start_time::boolean = true THEN start_time <= @later_start_time ELSE TRUE END
AND
	CASE WHEN @where_earlier_end_time::boolean = true THEN end_time >= @earlier_end_time ELSE TRUE END
AND
	CASE WHEN @where_later_end_time::boolean = true THEN end_time <= @later_end_time ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'start_time' THEN start_time END ASC,
	t_events_pkey DESC
LIMIT $1 OFFSET $2;

-- name: CountEvents :one
SELECT COUNT(*) FROM t_events
WHERE
	CASE WHEN @where_like_title::boolean = true THEN title LIKE '%' || @search_title::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_organization::boolean = true THEN t_events.organization_id = @organization_id ELSE TRUE END
AND
	CASE WHEN @where_mail_send_flag::boolean = true THEN mail_send_flag = @mail_send_flag ELSE TRUE END
AND
	CASE WHEN @where_send_organization::boolean = true THEN send_organization_id = @send_organization_id ELSE TRUE END
AND
	CASE WHEN @where_earlier_start_time::boolean = true THEN start_time >= @earlier_start_time ELSE TRUE END
AND
	CASE WHEN @where_later_start_time::boolean = true THEN start_time <= @later_start_time ELSE TRUE END
AND
	CASE WHEN @where_earlier_end_time::boolean = true THEN end_time >= @earlier_end_time ELSE TRUE END
AND
	CASE WHEN @where_later_end_time::boolean = true THEN end_time <= @later_end_time ELSE TRUE END;

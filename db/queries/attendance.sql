-- name: CreateAttendances :copyfrom
INSERT INTO t_attendances (attendance_type_id, member_id, description, date, mail_send_flag, send_organization_id, posted_at, last_edited_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: CreateAttendance :one
INSERT INTO t_attendances (attendance_type_id, member_id, description, date, mail_send_flag, send_organization_id, posted_at, last_edited_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: UpdateAttendance :one
UPDATE t_attendances SET attendance_type_id = $2, member_id = $3, description = $4, date = $5, mail_send_flag = $6, send_organization_id = $7, last_edited_at = $8 WHERE attendance_id = $1 RETURNING *;

-- name: DeleteAttendance :exec
DELETE FROM t_attendances WHERE attendance_id = $1;

-- name: FindAttendanceByID :one
SELECT * FROM t_attendances WHERE attendance_id = $1;

-- name: FindAttendanceByIDWithMember :one
SELECT sqlc.embed(t_attendances), sqlc.embed(m_members) FROM t_attendances
LEFT JOIN m_members ON t_attendances.member_id = m_members.member_id
WHERE attendance_id = $1;

-- name: FindAttendanceByIDWithAttendanceType :one
SELECT sqlc.embed(t_attendances), sqlc.embed(m_attendance_types) FROM t_attendances
LEFT JOIN m_attendance_types ON t_attendances.attendance_type_id = m_attendance_types.attendance_type_id
WHERE attendance_id = $1;

-- name: FindAttendanceByIDWithSendOrganization :one
SELECT sqlc.embed(t_attendances), sqlc.embed(m_organizations) FROM t_attendances
LEFT JOIN m_organizations ON t_attendances.send_organization_id = m_organizations.organization_id
WHERE attendance_id = $1;

-- name: FindAttendanceByIDWithDetails :one
SELECT sqlc.embed(t_attendances), sqlc.embed(t_early_leavings), sqlc.embed(t_late_arrivals), sqlc.embed(t_absences) FROM t_attendances
LEFT JOIN t_early_leavings ON t_attendances.attendance_id = t_early_leavings.attendance_id
LEFT JOIN t_late_arrivals ON t_attendances.attendance_id = t_late_arrivals.attendance_id
LEFT JOIN t_absences ON t_attendances.attendance_id = t_absences.attendance_id
WHERE t_attendances.attendance_id = $1;

-- name: FindAttendanceByIDWithAll :one
SELECT sqlc.embed(t_attendances), sqlc.embed(m_members), sqlc.embed(m_attendance_types), sqlc.embed(m_organizations), sqlc.embed(t_early_leavings), sqlc.embed(t_late_arrivals), sqlc.embed(t_absences) FROM t_attendances
LEFT JOIN m_members ON t_attendances.member_id = m_members.member_id
LEFT JOIN m_attendance_types ON t_attendances.attendance_type_id = m_attendance_types.attendance_type_id
LEFT JOIN m_organizations ON t_attendances.send_organization_id = m_organizations.organization_id
LEFT JOIN t_early_leavings ON t_attendances.attendance_id = t_early_leavings.attendance_id
LEFT JOIN t_late_arrivals ON t_attendances.attendance_id = t_late_arrivals.attendance_id
LEFT JOIN t_absences ON t_attendances.attendance_id = t_absences.attendance_id
WHERE t_attendances.attendance_id = $1;

-- name: GetAttendances :many
SELECT * FROM t_attendances
WHERE
	CASE WHEN @where_in_attendance_type::boolean = true THEN t_attendances.attendance_type_id = ANY(@in_attendance_type) ELSE TRUE END
AND
	CASE WHEN @where_in_member::boolean = true THEN t_attendances.member_id = ANY(@in_member) ELSE TRUE END
AND
	CASE WHEN @where_earlier_date::boolean = true THEN t_attendances.date >= @earlier_date ELSE TRUE END
AND
	CASE WHEN @where_later_date::boolean = true THEN t_attendances.date <= @later_date ELSE TRUE END
AND
	CASE WHEN @where_mail_send_flag::boolean = true THEN t_attendances.mail_send_flag = @mail_send_flag ELSE TRUE END
AND
	CASE WHEN @where_in_send_organization::boolean = true THEN t_attendances.send_organization_id = ANY(@in_send_organization) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'date' THEN t_attendances.date END ASC,
	CASE WHEN @order_method::text = 'r_date' THEN t_attendances.date END DESC,
	t_attendances_pkey DESC;

-- name: GetAttendanceUseNumberedPaginate :many
SELECT * FROM t_attendances
WHERE
	CASE WHEN @where_in_attendance_type::boolean = true THEN t_attendances.attendance_type_id = ANY(@in_attendance_type) ELSE TRUE END
AND
	CASE WHEN @where_in_member::boolean = true THEN t_attendances.member_id = ANY(@in_member) ELSE TRUE END
AND
	CASE WHEN @where_earlier_date::boolean = true THEN t_attendances.date >= @earlier_date ELSE TRUE END
AND
	CASE WHEN @where_later_date::boolean = true THEN t_attendances.date <= @later_date ELSE TRUE END
AND
	CASE WHEN @where_mail_send_flag::boolean = true THEN t_attendances.mail_send_flag = @mail_send_flag ELSE TRUE END
AND
	CASE WHEN @where_in_send_organization::boolean = true THEN t_attendances.send_organization_id = ANY(@in_send_organization) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'date' THEN t_attendances.date END ASC,
	CASE WHEN @order_method::text = 'r_date' THEN t_attendances.date END DESC,
	t_attendances_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetAttendanceUseKeysetPaginate :many
SELECT * FROM t_attendances
WHERE
	CASE WHEN @where_in_attendance_type::boolean = true THEN t_attendances.attendance_type_id = ANY(@in_attendance_type) ELSE TRUE END
AND
	CASE WHEN @where_in_member::boolean = true THEN t_attendances.member_id = ANY(@in_member) ELSE TRUE END
AND
	CASE WHEN @where_earlier_date::boolean = true THEN t_attendances.date >= @earlier_date ELSE TRUE END
AND
	CASE WHEN @where_later_date::boolean = true THEN t_attendances.date <= @later_date ELSE TRUE END
AND
	CASE WHEN @where_mail_send_flag::boolean = true THEN t_attendances.mail_send_flag = @mail_send_flag ELSE TRUE END
AND
	CASE WHEN @where_in_send_organization::boolean = true THEN t_attendances.send_organization_id = ANY(@in_send_organization) ELSE TRUE END
AND
	CASE @cursor_direction
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'date' THEN date > @cursor_column OR (date = @cursor_column AND t_attendances_pkey < @cursor)
				WHEN 'r_date' THEN date < @cursor_column OR (date = @cursor_column AND t_attendances_pkey < @cursor)
				ELSE t_attendances_pkey < @cursor
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'date' THEN date < @cursor_column OR (date = @cursor_column AND t_attendances_pkey > @cursor)
				WHEN 'r_date' THEN date > @cursor_column OR (date = @cursor_column AND t_attendances_pkey > @cursor)
				ELSE t_attendances_pkey > @cursor
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'date' THEN date END ASC,
	CASE WHEN @order_method::text = 'r_date' THEN date END DESC,
	t_attendances_pkey DESC
LIMIT $1;

-- name: GetAttendanceWithMember :many
SELECT sqlc.embed(t_attendances), sqlc.embed(m_members) FROM t_attendances
LEFT JOIN m_members ON t_attendances.member_id = m_members.member_id
WHERE
	CASE WHEN @where_in_attendance_type::boolean = true THEN t_attendances.attendance_type_id = ANY(@in_attendance_type) ELSE TRUE END
AND
	CASE WHEN @where_in_member::boolean = true THEN t_attendances.member_id = ANY(@in_member) ELSE TRUE END
AND
	CASE WHEN @where_earlier_date::boolean = true THEN t_attendances.date >= @earlier_date ELSE TRUE END
AND
	CASE WHEN @where_later_date::boolean = true THEN t_attendances.date <= @later_date ELSE TRUE END
AND
	CASE WHEN @where_mail_send_flag::boolean = true THEN t_attendances.mail_send_flag = @mail_send_flag ELSE TRUE END
AND
	CASE WHEN @where_in_send_organization::boolean = true THEN t_attendances.send_organization_id = ANY(@in_send_organization) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'date' THEN t_attendances.date END ASC,
	CASE WHEN @order_method::text = 'r_date' THEN t_attendances.date END DESC,
	t_attendances_pkey DESC;

-- name: GetAttendanceWithMemberUseNumberedPaginate :many
SELECT sqlc.embed(t_attendances), sqlc.embed(m_members) FROM t_attendances
LEFT JOIN m_members ON t_attendances.member_id = m_members.member_id
WHERE
	CASE WHEN @where_in_attendance_type::boolean = true THEN t_attendances.attendance_type_id = ANY(@in_attendance_type) ELSE TRUE END
AND
	CASE WHEN @where_in_member::boolean = true THEN t_attendances.member_id = ANY(@in_member) ELSE TRUE END
AND
	CASE WHEN @where_earlier_date::boolean = true THEN t_attendances.date >= @earlier_date ELSE TRUE END
AND
	CASE WHEN @where_later_date::boolean = true THEN t_attendances.date <= @later_date ELSE TRUE END
AND
	CASE WHEN @where_mail_send_flag::boolean = true THEN t_attendances.mail_send_flag = @mail_send_flag ELSE TRUE END
AND
	CASE WHEN @where_in_send_organization::boolean = true THEN t_attendances.send_organization_id = ANY(@in_send_organization) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'date' THEN t_attendances.date END ASC,
	CASE WHEN @order_method::text = 'r_date' THEN t_attendances.date END DESC,
	t_attendances_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetAttendanceWithMemberUseKeysetPaginate :many
SELECT sqlc.embed(t_attendances), sqlc.embed(m_members) FROM t_attendances
LEFT JOIN m_members ON t_attendances.member_id = m_members.member_id
WHERE
	CASE WHEN @where_in_attendance_type::boolean = true THEN t_attendances.attendance_type_id = ANY(@in_attendance_type) ELSE TRUE END
AND
	CASE WHEN @where_in_member::boolean = true THEN t_attendances.member_id = ANY(@in_member) ELSE TRUE END
AND
	CASE WHEN @where_earlier_date::boolean = true THEN t_attendances.date >= @earlier_date ELSE TRUE END
AND
	CASE WHEN @where_later_date::boolean = true THEN t_attendances.date <= @later_date ELSE TRUE END
AND
	CASE WHEN @where_mail_send_flag::boolean = true THEN t_attendances.mail_send_flag = @mail_send_flag ELSE TRUE END
AND
	CASE WHEN @where_in_send_organization::boolean = true THEN t_attendances.send_organization_id = ANY(@in_send_organization) ELSE TRUE END
AND
	CASE @cursor_direction
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'date' THEN t_attendances.date > @cursor_column OR (t_attendances.date = @cursor_column AND t_attendances_pkey < @cursor)
				WHEN 'r_date' THEN t_attendances.date < @cursor_column OR (t_attendances.date = @cursor_column AND t_attendances_pkey < @cursor)
				ELSE t_attendances_pkey < @cursor
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'date' THEN t_attendances.date < @cursor_column OR (t_attendances.date = @cursor_column AND t_attendances_pkey > @cursor)
				WHEN 'r_date' THEN t_attendances.date > @cursor_column OR (t_attendances.date = @cursor_column AND t_attendances_pkey > @cursor)
				ELSE t_attendances_pkey > @cursor
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'date' THEN t_attendances.date END ASC,
	CASE WHEN @order_method::text = 'r_date' THEN t_attendances.date END DESC,
	t_attendances_pkey DESC
LIMIT $1;

-- name: GetAttendanceWithAttendanceType :many
SELECT sqlc.embed(t_attendances), sqlc.embed(m_attendance_types) FROM t_attendances
LEFT JOIN m_attendance_types ON t_attendances.attendance_type_id = m_attendance_types.attendance_type_id
WHERE
	CASE WHEN @where_in_attendance_type::boolean = true THEN t_attendances.attendance_type_id = ANY(@in_attendance_type) ELSE TRUE END
AND
	CASE WHEN @where_in_member::boolean = true THEN t_attendances.member_id = ANY(@in_member) ELSE TRUE END
AND
	CASE WHEN @where_earlier_date::boolean = true THEN t_attendances.date >= @earlier_date ELSE TRUE END
AND
	CASE WHEN @where_later_date::boolean = true THEN t_attendances.date <= @later_date ELSE TRUE END
AND
	CASE WHEN @where_mail_send_flag::boolean = true THEN t_attendances.mail_send_flag = @mail_send_flag ELSE TRUE END
AND
	CASE WHEN @where_in_send_organization::boolean = true THEN t_attendances.send_organization_id = ANY(@in_send_organization) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'date' THEN t_attendances.date END ASC,
	CASE WHEN @order_method::text = 'r_date' THEN t_attendances.date END DESC,
	t_attendances_pkey DESC;

-- name: GetAttendanceWithAttendanceTypeUseNumberedPaginate :many
SELECT sqlc.embed(t_attendances), sqlc.embed(m_attendance_types) FROM t_attendances
LEFT JOIN m_attendance_types ON t_attendances.attendance_type_id = m_attendance_types.attendance_type_id
WHERE
	CASE WHEN @where_in_attendance_type::boolean = true THEN t_attendances.attendance_type_id = ANY(@in_attendance_type) ELSE TRUE END
AND
	CASE WHEN @where_in_member::boolean = true THEN t_attendances.member_id = ANY(@in_member) ELSE TRUE END
AND
	CASE WHEN @where_earlier_date::boolean = true THEN t_attendances.date >= @earlier_date ELSE TRUE END
AND
	CASE WHEN @where_later_date::boolean = true THEN t_attendances.date <= @later_date ELSE TRUE END
AND
	CASE WHEN @where_mail_send_flag::boolean = true THEN t_attendances.mail_send_flag = @mail_send_flag ELSE TRUE END
AND
	CASE WHEN @where_in_send_organization::boolean = true THEN t_attendances.send_organization_id = ANY(@in_send_organization) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'date' THEN t_attendances.date END ASC,
	CASE WHEN @order_method::text = 'r_date' THEN t_attendances.date END DESC,
	t_attendances_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetAttendanceWithAttendanceTypeUseKeysetPaginate :many
SELECT sqlc.embed(t_attendances), sqlc.embed(m_attendance_types) FROM t_attendances
LEFT JOIN m_attendance_types ON t_attendances.attendance_type_id = m_attendance_types.attendance_type_id
WHERE
	CASE WHEN @where_in_attendance_type::boolean = true THEN t_attendances.attendance_type_id = ANY(@in_attendance_type) ELSE TRUE END
AND
	CASE WHEN @where_in_member::boolean = true THEN t_attendances.member_id = ANY(@in_member) ELSE TRUE END
AND
	CASE WHEN @where_earlier_date::boolean = true THEN t_attendances.date >= @earlier_date ELSE TRUE END
AND
	CASE WHEN @where_later_date::boolean = true THEN t_attendances.date <= @later_date ELSE TRUE END
AND
	CASE WHEN @where_mail_send_flag::boolean = true THEN t_attendances.mail_send_flag = @mail_send_flag ELSE TRUE END
AND
	CASE WHEN @where_in_send_organization::boolean = true THEN t_attendances.send_organization_id = ANY(@in_send_organization) ELSE TRUE END
AND
	CASE @cursor_direction
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'date' THEN t_attendances.date > @cursor_column OR (t_attendances.date = @cursor_column AND t_attendances_pkey < @cursor)
				WHEN 'r_date' THEN t_attendances.date < @cursor_column OR (t_attendances.date = @cursor_column AND t_attendances_pkey < @cursor)
				ELSE t_attendances_pkey < @cursor
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'date' THEN t_attendances.date < @cursor_column OR (t_attendances.date = @cursor_column AND t_attendances_pkey > @cursor)
				WHEN 'r_date' THEN t_attendances.date > @cursor_column OR (t_attendances.date = @cursor_column AND t_attendances_pkey > @cursor)
				ELSE t_attendances_pkey > @cursor
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'date' THEN t_attendances.date END ASC,
	CASE WHEN @order_method::text = 'r_date' THEN t_attendances.date END DESC,
	t_attendances_pkey DESC
LIMIT $1;

-- name: GetAttendanceWithSendOrganization :many
SELECT sqlc.embed(t_attendances), sqlc.embed(m_organizations) FROM t_attendances
LEFT JOIN m_organizations ON t_attendances.send_organization_id = m_organizations.organization_id
WHERE
	CASE WHEN @where_in_attendance_type::boolean = true THEN t_attendances.attendance_type_id = ANY(@in_attendance_type) ELSE TRUE END
AND
	CASE WHEN @where_in_member::boolean = true THEN t_attendances.member_id = ANY(@in_member) ELSE TRUE END
AND
	CASE WHEN @where_earlier_date::boolean = true THEN t_attendances.date >= @earlier_date ELSE TRUE END
AND
	CASE WHEN @where_later_date::boolean = true THEN t_attendances.date <= @later_date ELSE TRUE END
AND
	CASE WHEN @where_mail_send_flag::boolean = true THEN t_attendances.mail_send_flag = @mail_send_flag ELSE TRUE END
AND
	CASE WHEN @where_in_send_organization::boolean = true THEN t_attendances.send_organization_id = ANY(@in_send_organization) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'date' THEN t_attendances.date END ASC,
	CASE WHEN @order_method::text = 'r_date' THEN t_attendances.date END DESC,
	t_attendances_pkey DESC;

-- name: GetAttendanceWithSendOrganizationUseNumberedPaginate :many
SELECT sqlc.embed(t_attendances), sqlc.embed(m_organizations) FROM t_attendances
LEFT JOIN m_organizations ON t_attendances.send_organization_id = m_organizations.organization_id
WHERE
	CASE WHEN @where_in_attendance_type::boolean = true THEN t_attendances.attendance_type_id = ANY(@in_attendance_type) ELSE TRUE END
AND
	CASE WHEN @where_in_member::boolean = true THEN t_attendances.member_id = ANY(@in_member) ELSE TRUE END
AND
	CASE WHEN @where_earlier_date::boolean = true THEN t_attendances.date >= @earlier_date ELSE TRUE END
AND
	CASE WHEN @where_later_date::boolean = true THEN t_attendances.date <= @later_date ELSE TRUE END
AND
	CASE WHEN @where_mail_send_flag::boolean = true THEN t_attendances.mail_send_flag = @mail_send_flag ELSE TRUE END
AND
	CASE WHEN @where_in_send_organization::boolean = true THEN t_attendances.send_organization_id = ANY(@in_send_organization) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'date' THEN t_attendances.date END ASC,
	CASE WHEN @order_method::text = 'r_date' THEN t_attendances.date END DESC,
	t_attendances_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetAttendanceWithSendOrganizationUseKeysetPaginate :many
SELECT sqlc.embed(t_attendances), sqlc.embed(m_organizations) FROM t_attendances
LEFT JOIN m_organizations ON t_attendances.send_organization_id = m_organizations.organization_id
WHERE
	CASE WHEN @where_in_attendance_type::boolean = true THEN t_attendances.attendance_type_id = ANY(@in_attendance_type) ELSE TRUE END
AND
	CASE WHEN @where_in_member::boolean = true THEN t_attendances.member_id = ANY(@in_member) ELSE TRUE END
AND
	CASE WHEN @where_earlier_date::boolean = true THEN t_attendances.date >= @earlier_date ELSE TRUE END
AND
	CASE WHEN @where_later_date::boolean = true THEN t_attendances.date <= @later_date ELSE TRUE END
AND
	CASE WHEN @where_mail_send_flag::boolean = true THEN t_attendances.mail_send_flag = @mail_send_flag ELSE TRUE END
AND
	CASE WHEN @where_in_send_organization::boolean = true THEN t_attendances.send_organization_id = ANY(@in_send_organization) ELSE TRUE END
AND
	CASE @cursor_direction
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'date' THEN t_attendances.date > @cursor_column OR (t_attendances.date = @cursor_column AND t_attendances_pkey < @cursor)
				WHEN 'r_date' THEN t_attendances.date < @cursor_column OR (t_attendances.date = @cursor_column AND t_attendances_pkey < @cursor)
				ELSE t_attendances_pkey < @cursor
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'date' THEN t_attendances.date < @cursor_column OR (t_attendances.date = @cursor_column AND t_attendances_pkey > @cursor)
				WHEN 'r_date' THEN t_attendances.date > @cursor_column OR (t_attendances.date = @cursor_column AND t_attendances_pkey > @cursor)
				ELSE t_attendances_pkey > @cursor
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'date' THEN t_attendances.date END ASC,
	CASE WHEN @order_method::text = 'r_date' THEN t_attendances.date END DESC,
	t_attendances_pkey DESC
LIMIT $1;

-- name: GetAttendanceWithDetails :many
SELECT sqlc.embed(t_attendances), sqlc.embed(t_early_leavings), sqlc.embed(t_late_arrivals), sqlc.embed(t_absences) FROM t_attendances
LEFT JOIN t_early_leavings ON t_attendances.attendance_id = t_early_leavings.attendance_id
LEFT JOIN t_late_arrivals ON t_attendances.attendance_id = t_late_arrivals.attendance_id
LEFT JOIN t_absences ON t_attendances.attendance_id = t_absences.attendance_id
WHERE
	CASE WHEN @where_in_attendance_type::boolean = true THEN t_attendances.attendance_type_id = ANY(@in_attendance_type) ELSE TRUE END
AND
	CASE WHEN @where_in_member::boolean = true THEN t_attendances.member_id = ANY(@in_member) ELSE TRUE END
AND
	CASE WHEN @where_earlier_date::boolean = true THEN t_attendances.date >= @earlier_date ELSE TRUE END
AND
	CASE WHEN @where_later_date::boolean = true THEN t_attendances.date <= @later_date ELSE TRUE END
AND
	CASE WHEN @where_mail_send_flag::boolean = true THEN t_attendances.mail_send_flag = @mail_send_flag ELSE TRUE END
AND
	CASE WHEN @where_in_send_organization::boolean = true THEN t_attendances.send_organization_id = ANY(@in_send_organization) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'date' THEN t_attendances.date END ASC,
	CASE WHEN @order_method::text = 'r_date' THEN t_attendances.date END DESC,
	t_attendances_pkey DESC;

-- name: GetAttendanceWithDetailsUseNumberedPaginate :many
SELECT sqlc.embed(t_attendances), sqlc.embed(t_early_leavings), sqlc.embed(t_late_arrivals), sqlc.embed(t_absences) FROM t_attendances
LEFT JOIN t_early_leavings ON t_attendances.attendance_id = t_early_leavings.attendance_id
LEFT JOIN t_late_arrivals ON t_attendances.attendance_id = t_late_arrivals.attendance_id
LEFT JOIN t_absences ON t_attendances.attendance_id = t_absences.attendance_id
WHERE
	CASE WHEN @where_in_attendance_type::boolean = true THEN t_attendances.attendance_type_id = ANY(@in_attendance_type) ELSE TRUE END
AND
	CASE WHEN @where_in_member::boolean = true THEN t_attendances.member_id = ANY(@in_member) ELSE TRUE END
AND
	CASE WHEN @where_earlier_date::boolean = true THEN t_attendances.date >= @earlier_date ELSE TRUE END
AND
	CASE WHEN @where_later_date::boolean = true THEN t_attendances.date <= @later_date ELSE TRUE END
AND
	CASE WHEN @where_mail_send_flag::boolean = true THEN t_attendances.mail_send_flag = @mail_send_flag ELSE TRUE END
AND
	CASE WHEN @where_in_send_organization::boolean = true THEN t_attendances.send_organization_id = ANY(@in_send_organization) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'date' THEN t_attendances.date END ASC,
	CASE WHEN @order_method::text = 'r_date' THEN t_attendances.date END DESC,
	t_attendances_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetAttendanceWithDetailsUseKeysetPaginate :many
SELECT sqlc.embed(t_attendances), sqlc.embed(t_early_leavings), sqlc.embed(t_late_arrivals), sqlc.embed(t_absences) FROM t_attendances
LEFT JOIN t_early_leavings ON t_attendances.attendance_id = t_early_leavings.attendance_id
LEFT JOIN t_late_arrivals ON t_attendances.attendance_id = t_late_arrivals.attendance_id
LEFT JOIN t_absences ON t_attendances.attendance_id = t_absences.attendance_id
WHERE
	CASE WHEN @where_in_attendance_type::boolean = true THEN t_attendances.attendance_type_id = ANY(@in_attendance_type) ELSE TRUE END
AND
	CASE WHEN @where_in_member::boolean = true THEN t_attendances.member_id = ANY(@in_member) ELSE TRUE END
AND
	CASE WHEN @where_earlier_date::boolean = true THEN t_attendances.date >= @earlier_date ELSE TRUE END
AND
	CASE WHEN @where_later_date::boolean = true THEN t_attendances.date <= @later_date ELSE TRUE END
AND
	CASE WHEN @where_mail_send_flag::boolean = true THEN t_attendances.mail_send_flag = @mail_send_flag ELSE TRUE END
AND
	CASE WHEN @where_in_send_organization::boolean = true THEN t_attendances.send_organization_id = ANY(@in_send_organization) ELSE TRUE END
AND
	CASE @cursor_direction
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'date' THEN t_attendances.date > @cursor_column OR (t_attendances.date = @cursor_column AND t_attendances_pkey < @cursor)
				WHEN 'r_date' THEN t_attendances.date < @cursor_column OR (t_attendances.date = @cursor_column AND t_attendances_pkey < @cursor)
				ELSE t_attendances_pkey < @cursor
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'date' THEN t_attendances.date < @cursor_column OR (t_attendances.date = @cursor_column AND t_attendances_pkey > @cursor)
				WHEN 'r_date' THEN t_attendances.date > @cursor_column OR (t_attendances.date = @cursor_column AND t_attendances_pkey > @cursor)
				ELSE t_attendances_pkey > @cursor
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'date' THEN t_attendances.date END ASC,
	CASE WHEN @order_method::text = 'r_date' THEN t_attendances.date END DESC,
	t_attendances_pkey DESC
LIMIT $1;

-- name: GetAttendanceWithAll :many
SELECT sqlc.embed(t_attendances), sqlc.embed(m_members), sqlc.embed(m_attendance_types), sqlc.embed(m_organizations), sqlc.embed(t_early_leavings), sqlc.embed(t_late_arrivals), sqlc.embed(t_absences) FROM t_attendances
LEFT JOIN t_early_leavings ON t_attendances.attendance_id = t_early_leavings.attendance_id
LEFT JOIN t_late_arrivals ON t_attendances.attendance_id = t_late_arrivals.attendance_id
LEFT JOIN t_absences ON t_attendances.attendance_id = t_absences.attendance_id
LEFT JOIN m_members ON t_attendances.member_id = m_members.member_id
LEFT JOIN m_attendance_types ON t_attendances.attendance_type_id = m_attendance_types.attendance_type_id
LEFT JOIN m_organizations ON t_attendances.send_organization_id = m_organizations.organization_id
WHERE
	CASE WHEN @where_in_attendance_type::boolean = true THEN t_attendances.attendance_type_id = ANY(@in_attendance_type) ELSE TRUE END
AND
	CASE WHEN @where_in_member::boolean = true THEN t_attendances.member_id = ANY(@in_member) ELSE TRUE END
AND
	CASE WHEN @where_earlier_date::boolean = true THEN t_attendances.date >= @earlier_date ELSE TRUE END
AND
	CASE WHEN @where_later_date::boolean = true THEN t_attendances.date <= @later_date ELSE TRUE END
AND
	CASE WHEN @where_mail_send_flag::boolean = true THEN t_attendances.mail_send_flag = @mail_send_flag ELSE TRUE END
AND
	CASE WHEN @where_in_send_organization::boolean = true THEN t_attendances.send_organization_id = ANY(@in_send_organization) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'date' THEN t_attendances.date END ASC,
	CASE WHEN @order_method::text = 'r_date' THEN t_attendances.date END DESC,
	t_attendances_pkey DESC;

-- name: GetAttendanceWithAllUseNumberedPaginate :many
SELECT sqlc.embed(t_attendances), sqlc.embed(m_members), sqlc.embed(m_attendance_types), sqlc.embed(m_organizations), sqlc.embed(t_early_leavings), sqlc.embed(t_late_arrivals), sqlc.embed(t_absences) FROM t_attendances
LEFT JOIN t_early_leavings ON t_attendances.attendance_id = t_early_leavings.attendance_id
LEFT JOIN t_late_arrivals ON t_attendances.attendance_id = t_late_arrivals.attendance_id
LEFT JOIN t_absences ON t_attendances.attendance_id = t_absences.attendance_id
LEFT JOIN m_members ON t_attendances.member_id = m_members.member_id
LEFT JOIN m_attendance_types ON t_attendances.attendance_type_id = m_attendance_types.attendance_type_id
LEFT JOIN m_organizations ON t_attendances.send_organization_id = m_organizations.organization_id
WHERE
	CASE WHEN @where_in_attendance_type::boolean = true THEN t_attendances.attendance_type_id = ANY(@in_attendance_type) ELSE TRUE END
AND
	CASE WHEN @where_in_member::boolean = true THEN t_attendances.member_id = ANY(@in_member) ELSE TRUE END
AND
	CASE WHEN @where_earlier_date::boolean = true THEN t_attendances.date >= @earlier_date ELSE TRUE END
AND
	CASE WHEN @where_later_date::boolean = true THEN t_attendances.date <= @later_date ELSE TRUE END
AND
	CASE WHEN @where_mail_send_flag::boolean = true THEN t_attendances.mail_send_flag = @mail_send_flag ELSE TRUE END
AND
	CASE WHEN @where_in_send_organization::boolean = true THEN t_attendances.send_organization_id = ANY(@in_send_organization) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'date' THEN t_attendances.date END ASC,
	CASE WHEN @order_method::text = 'r_date' THEN t_attendances.date END DESC,
	t_attendances_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetAttendanceWithAllUseKeysetPaginate :many
SELECT sqlc.embed(t_attendances), sqlc.embed(m_members), sqlc.embed(m_attendance_types), sqlc.embed(m_organizations), sqlc.embed(t_early_leavings), sqlc.embed(t_late_arrivals), sqlc.embed(t_absences) FROM t_attendances
LEFT JOIN t_early_leavings ON t_attendances.attendance_id = t_early_leavings.attendance_id
LEFT JOIN t_late_arrivals ON t_attendances.attendance_id = t_late_arrivals.attendance_id
LEFT JOIN t_absences ON t_attendances.attendance_id = t_absences.attendance_id
LEFT JOIN m_members ON t_attendances.member_id = m_members.member_id
LEFT JOIN m_attendance_types ON t_attendances.attendance_type_id = m_attendance_types.attendance_type_id
LEFT JOIN m_organizations ON t_attendances.send_organization_id = m_organizations.organization_id
WHERE
	CASE WHEN @where_in_attendance_type::boolean = true THEN t_attendances.attendance_type_id = ANY(@in_attendance_type) ELSE TRUE END
AND
	CASE WHEN @where_in_member::boolean = true THEN t_attendances.member_id = ANY(@in_member) ELSE TRUE END
AND
	CASE WHEN @where_earlier_date::boolean = true THEN t_attendances.date >= @earlier_date ELSE TRUE END
AND
	CASE WHEN @where_later_date::boolean = true THEN t_attendances.date <= @later_date ELSE TRUE END
AND
	CASE WHEN @where_mail_send_flag::boolean = true THEN t_attendances.mail_send_flag = @mail_send_flag ELSE TRUE END
AND
	CASE WHEN @where_in_send_organization::boolean = true THEN t_attendances.send_organization_id = ANY(@in_send_organization) ELSE TRUE END
AND
	CASE @cursor_direction
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'date' THEN t_attendances.date > @cursor_column OR (t_attendances.date = @cursor_column AND t_attendances_pkey < @cursor)
				WHEN 'r_date' THEN t_attendances.date < @cursor_column OR (t_attendances.date = @cursor_column AND t_attendances_pkey < @cursor)
				ELSE t_attendances_pkey < @cursor
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'date' THEN t_attendances.date < @cursor_column OR (t_attendances.date = @cursor_column AND t_attendances_pkey > @cursor)
				WHEN 'r_date' THEN t_attendances.date > @cursor_column OR (t_attendances.date = @cursor_column AND t_attendances_pkey > @cursor)
				ELSE t_attendances_pkey > @cursor
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'date' THEN t_attendances.date END ASC,
	CASE WHEN @order_method::text = 'r_date' THEN t_attendances.date END DESC,
	t_attendances_pkey DESC
LIMIT $1;

-- name: CountAttendances :one
SELECT COUNT(*) FROM t_attendances
WHERE
	CASE WHEN @where_in_attendance_type::boolean = true THEN t_attendances.attendance_type_id = ANY(@in_attendance_type) ELSE TRUE END
AND
	CASE WHEN @where_in_member::boolean = true THEN t_attendances.member_id = ANY(@in_member) ELSE TRUE END
AND
	CASE WHEN @where_earlier_date::boolean = true THEN t_attendances.date >= @earlier_date ELSE TRUE END
AND
	CASE WHEN @where_later_date::boolean = true THEN t_attendances.date <= @later_date ELSE TRUE END
AND
	CASE WHEN @where_mail_send_flag::boolean = true THEN t_attendances.mail_send_flag = @mail_send_flag ELSE TRUE END
AND
	CASE WHEN @where_in_send_organization::boolean = true THEN t_attendances.send_organization_id = ANY(@in_send_organization) ELSE TRUE END;

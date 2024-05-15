-- name: CreateAttendances :copyfrom
INSERT INTO t_attendances (attendance_type_id, member_id, description, date, mail_send_flag, send_organization_id, posted_at, last_edited_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: CreateAttendance :one
INSERT INTO t_attendances (attendance_type_id, member_id, description, date, mail_send_flag, send_organization_id, posted_at, last_edited_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: UpdateAttendance :one
UPDATE t_attendances SET attendance_type_id = $2, member_id = $3, description = $4, date = $5, mail_send_flag = $6, send_organization_id = $7, last_edited_at = $8 WHERE attendance_id = $1 RETURNING *;

-- name: DeleteAttendance :execrows
DELETE FROM t_attendances WHERE attendance_id = $1;

-- name: DeleteAttendancesOnMember :execrows
DELETE FROM t_attendances WHERE member_id = $1;

-- name: DeleteAttendancesOnMembers :execrows
DELETE FROM t_attendances WHERE member_id = ANY($1::uuid[]);

-- name: PluralDeleteAttendances :execrows
DELETE FROM t_attendances WHERE attendance_id = ANY($1::uuid[]);

-- name: FindAttendanceByID :one
SELECT * FROM t_attendances WHERE attendance_id = $1;

-- name: FindAttendanceByIDWithMember :one
SELECT t_attendances.*, sqlc.embed(m_members) FROM t_attendances
LEFT JOIN m_members ON t_attendances.member_id = m_members.member_id
LEFT JOIN m_attend_statuses ON m_members.attend_status_id = m_attend_statuses.attend_status_id
LEFT JOIN m_grades ON m_members.grade_id = m_grades.grade_id
LEFT JOIN m_groups ON m_members.group_id = m_groups.group_id
WHERE attendance_id = $1;

-- name: FindAttendanceByIDWithAttendanceType :one
SELECT t_attendances.*, m_attendance_types.attendance_type_id, m_attendance_types.name as attendance_type_name, m_attendance_types.key as attendance_type_key, m_attendance_types.color as attendance_type_color FROM t_attendances
LEFT JOIN m_attendance_types ON t_attendances.attendance_type_id = m_attendance_types.attendance_type_id
WHERE attendance_id = $1;

-- name: FindAttendanceByIDWithSendOrganization :one
SELECT t_attendances.*, sqlc.embed(m_organizations) FROM t_attendances
LEFT JOIN m_organizations ON t_attendances.send_organization_id = m_organizations.organization_id
WHERE attendance_id = $1;

-- name: FindAttendanceByIDWithDetails :one
SELECT t_attendances.*, sqlc.embed(t_early_leavings), sqlc.embed(t_late_arrivals), sqlc.embed(t_absences) FROM t_attendances
LEFT JOIN t_early_leavings ON t_attendances.attendance_id = t_early_leavings.attendance_id
LEFT JOIN t_late_arrivals ON t_attendances.attendance_id = t_late_arrivals.attendance_id
LEFT JOIN t_absences ON t_attendances.attendance_id = t_absences.attendance_id
WHERE t_attendances.attendance_id = $1;

-- name: FindAttendanceByIDWithAll :one
SELECT t_attendances.*, sqlc.embed(m_members), m_attendance_types.attendance_type_id, m_attendance_types.name as attendance_type_name, m_attendance_types.key as attendance_type_key, m_attendance_types.color as attendance_type_color, sqlc.embed(m_organizations), sqlc.embed(t_early_leavings), sqlc.embed(t_late_arrivals), sqlc.embed(t_absences) FROM t_attendances
LEFT JOIN m_members ON t_attendances.member_id = m_members.member_id
LEFT JOIN m_attend_statuses ON m_members.attend_status_id = m_attend_statuses.attend_status_id
LEFT JOIN m_grades ON m_members.grade_id = m_grades.grade_id
LEFT JOIN m_groups ON m_members.group_id = m_groups.group_id
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
	t_attendances_pkey ASC;

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
	t_attendances_pkey ASC
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
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'date' THEN date > @date_cursor OR (date = @date_cursor AND t_attendances_pkey > @cursor::int)
				WHEN 'r_date' THEN date < @date_cursor OR (date = @date_cursor AND t_attendances_pkey > @cursor::int)
				ELSE t_attendances_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'date' THEN date < @date_cursor OR (date = @date_cursor AND t_attendances_pkey < @cursor::int)
				WHEN 'r_date' THEN date > @date_cursor OR (date = @date_cursor AND t_attendances_pkey < @cursor::int)
				ELSE t_attendances_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'date' AND @cursor_direction::text = 'next' THEN date END ASC,
	CASE WHEN @order_method::text = 'date' AND @cursor_direction::text = 'prev' THEN date END DESC,
	CASE WHEN @order_method::text = 'r_date' AND @cursor_direction::text = 'next' THEN date END DESC,
	CASE WHEN @order_method::text = 'r_date' AND @cursor_direction::text = 'prev' THEN date END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN t_attendances_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_attendances_pkey END DESC
LIMIT $1;

-- name: GetPluralAttendances :many
SELECT * FROM t_attendances
WHERE attendance_id = ANY(@attendance_ids::uuid[])
ORDER BY
	t_attendances_pkey ASC;

-- name: GetPluralAttendancesUseNumberedPaginate :many
SELECT * FROM t_attendances
WHERE attendance_id = ANY(@attendance_ids::uuid[])
ORDER BY
	t_attendances_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetAttendanceWithMember :many
SELECT t_attendances.*, sqlc.embed(m_members) FROM t_attendances
LEFT JOIN m_members ON t_attendances.member_id = m_members.member_id
LEFT JOIN m_attend_statuses ON m_members.attend_status_id = m_attend_statuses.attend_status_id
LEFT JOIN m_grades ON m_members.grade_id = m_grades.grade_id
LEFT JOIN m_groups ON m_members.group_id = m_groups.group_id
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
	t_attendances_pkey ASC;

-- name: GetAttendanceWithMemberUseNumberedPaginate :many
SELECT t_attendances.*, sqlc.embed(m_members) FROM t_attendances
LEFT JOIN m_members ON t_attendances.member_id = m_members.member_id
LEFT JOIN m_attend_statuses ON m_members.attend_status_id = m_attend_statuses.attend_status_id
LEFT JOIN m_grades ON m_members.grade_id = m_grades.grade_id
LEFT JOIN m_groups ON m_members.group_id = m_groups.group_id
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
	t_attendances_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetAttendanceWithMemberUseKeysetPaginate :many
SELECT t_attendances.*, sqlc.embed(m_members) FROM t_attendances
LEFT JOIN m_members ON t_attendances.member_id = m_members.member_id
LEFT JOIN m_attend_statuses ON m_members.attend_status_id = m_attend_statuses.attend_status_id
LEFT JOIN m_grades ON m_members.grade_id = m_grades.grade_id
LEFT JOIN m_groups ON m_members.group_id = m_groups.group_id
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
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'date' THEN t_attendances.date > @date_cursor OR (t_attendances.date = @date_cursor AND t_attendances_pkey > @cursor::int)
				WHEN 'r_date' THEN t_attendances.date < @date_cursor OR (t_attendances.date = @date_cursor AND t_attendances_pkey > @cursor::int)
				ELSE t_attendances_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'date' THEN t_attendances.date < @date_cursor OR (t_attendances.date = @date_cursor AND t_attendances_pkey < @cursor::int)
				WHEN 'r_date' THEN t_attendances.date > @date_cursor OR (t_attendances.date = @date_cursor AND t_attendances_pkey < @cursor::int)
				ELSE t_attendances_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'date' AND @cursor_direction::text = 'next' THEN date END ASC,
	CASE WHEN @order_method::text = 'date' AND @cursor_direction::text = 'prev' THEN date END DESC,
	CASE WHEN @order_method::text = 'r_date' AND @cursor_direction::text = 'next' THEN date END DESC,
	CASE WHEN @order_method::text = 'r_date' AND @cursor_direction::text = 'prev' THEN date END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN t_attendances_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_attendances_pkey END DESC
LIMIT $1;

-- name: GetPluralAttendanceWithMember :many
SELECT t_attendances.*, sqlc.embed(m_members) FROM t_attendances
LEFT JOIN m_members ON t_attendances.member_id = m_members.member_id
LEFT JOIN m_attend_statuses ON m_members.attend_status_id = m_attend_statuses.attend_status_id
LEFT JOIN m_grades ON m_members.grade_id = m_grades.grade_id
LEFT JOIN m_groups ON m_members.group_id = m_groups.group_id
WHERE attendance_id = ANY(@attendance_ids::uuid[])
ORDER BY
	t_attendances_pkey ASC;

-- name: GetPluralAttendanceWithMemberUseNumberedPaginate :many
SELECT t_attendances.*, sqlc.embed(m_members) FROM t_attendances
LEFT JOIN m_members ON t_attendances.member_id = m_members.member_id
LEFT JOIN m_attend_statuses ON m_members.attend_status_id = m_attend_statuses.attend_status_id
LEFT JOIN m_grades ON m_members.grade_id = m_grades.grade_id
LEFT JOIN m_groups ON m_members.group_id = m_groups.group_id
WHERE attendance_id = ANY(@attendance_ids::uuid[])
ORDER BY
	t_attendances_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetAttendanceWithAttendanceType :many
SELECT t_attendances.*, m_attendance_types.attendance_type_id, m_attendance_types.name as attendance_type_name, m_attendance_types.key as attendance_type_key, m_attendance_types.color as attendance_type_color FROM t_attendances
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
	t_attendances_pkey ASC;

-- name: GetAttendanceWithAttendanceTypeUseNumberedPaginate :many
SELECT t_attendances.*, m_attendance_types.attendance_type_id, m_attendance_types.name as attendance_type_name, m_attendance_types.key as attendance_type_key, m_attendance_types.color as attendance_type_color FROM t_attendances
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
	t_attendances_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetAttendanceWithAttendanceTypeUseKeysetPaginate :many
SELECT t_attendances.*, m_attendance_types.attendance_type_id, m_attendance_types.name as attendance_type_name, m_attendance_types.key as attendance_type_key, m_attendance_types.color as attendance_type_color FROM t_attendances
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
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'date' THEN t_attendances.date > @date_cursor OR (t_attendances.date = @date_cursor AND t_attendances_pkey > @cursor::int)
				WHEN 'r_date' THEN t_attendances.date < @date_cursor OR (t_attendances.date = @date_cursor AND t_attendances_pkey > @cursor::int)
				ELSE t_attendances_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'date' THEN t_attendances.date < @date_cursor OR (t_attendances.date = @date_cursor AND t_attendances_pkey < @cursor::int)
				WHEN 'r_date' THEN t_attendances.date > @date_cursor OR (t_attendances.date = @date_cursor AND t_attendances_pkey < @cursor::int)
				ELSE t_attendances_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'date' AND @cursor_direction::text = 'next' THEN date END ASC,
	CASE WHEN @order_method::text = 'date' AND @cursor_direction::text = 'prev' THEN date END DESC,
	CASE WHEN @order_method::text = 'r_date' AND @cursor_direction::text = 'next' THEN date END DESC,
	CASE WHEN @order_method::text = 'r_date' AND @cursor_direction::text = 'prev' THEN date END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN t_attendances_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_attendances_pkey END DESC
LIMIT $1;

-- name: GetPluralAttendanceWithAttendanceType :many
SELECT t_attendances.*, m_attendance_types.attendance_type_id, m_attendance_types.name as attendance_type_name, m_attendance_types.key as attendance_type_key, m_attendance_types.color as attendance_type_color FROM t_attendances
LEFT JOIN m_attendance_types ON t_attendances.attendance_type_id = m_attendance_types.attendance_type_id
WHERE attendance_id = ANY(@attendance_ids::uuid[])
ORDER BY
	t_attendances_pkey ASC;

-- name: GetPluralAttendanceWithAttendanceTypeUseNumberedPaginate :many
SELECT t_attendances.*, m_attendance_types.attendance_type_id, m_attendance_types.name as attendance_type_name, m_attendance_types.key as attendance_type_key, m_attendance_types.color as attendance_type_color FROM t_attendances
LEFT JOIN m_attendance_types ON t_attendances.attendance_type_id = m_attendance_types.attendance_type_id
WHERE attendance_id = ANY(@attendance_ids::uuid[])
ORDER BY
	t_attendances_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetAttendanceWithSendOrganization :many
SELECT t_attendances.*, sqlc.embed(m_organizations) FROM t_attendances
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
	t_attendances_pkey ASC;

-- name: GetAttendanceWithSendOrganizationUseNumberedPaginate :many
SELECT t_attendances.*, sqlc.embed(m_organizations) FROM t_attendances
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
	t_attendances_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetAttendanceWithSendOrganizationUseKeysetPaginate :many
SELECT t_attendances.*, sqlc.embed(m_organizations) FROM t_attendances
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
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'date' THEN t_attendances.date > @date_cursor OR (t_attendances.date = @date_cursor AND t_attendances_pkey > @cursor::int)
				WHEN 'r_date' THEN t_attendances.date < @date_cursor OR (t_attendances.date = @date_cursor AND t_attendances_pkey > @cursor::int)
				ELSE t_attendances_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'date' THEN t_attendances.date < @date_cursor OR (t_attendances.date = @date_cursor AND t_attendances_pkey < @cursor::int)
				WHEN 'r_date' THEN t_attendances.date > @date_cursor OR (t_attendances.date = @date_cursor AND t_attendances_pkey < @cursor::int)
				ELSE t_attendances_pkey > @cursor
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'date' AND @cursor_direction::text = 'next' THEN date END ASC,
	CASE WHEN @order_method::text = 'date' AND @cursor_direction::text = 'prev' THEN date END DESC,
	CASE WHEN @order_method::text = 'r_date' AND @cursor_direction::text = 'next' THEN date END DESC,
	CASE WHEN @order_method::text = 'r_date' AND @cursor_direction::text = 'prev' THEN date END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN t_attendances_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_attendances_pkey END DESC
LIMIT $1;

-- name: GetPluralAttendanceWithSendOrganization :many
SELECT t_attendances.*, sqlc.embed(m_organizations) FROM t_attendances
LEFT JOIN m_organizations ON t_attendances.send_organization_id = m_organizations.organization_id
WHERE attendance_id = ANY(@attendance_ids::uuid[])
ORDER BY
	t_attendances_pkey ASC;

-- name: GetPluralAttendanceWithSendOrganizationUseNumberedPaginate :many
SELECT t_attendances.*, sqlc.embed(m_organizations) FROM t_attendances
LEFT JOIN m_organizations ON t_attendances.send_organization_id = m_organizations.organization_id
WHERE attendance_id = ANY(@attendance_ids::uuid[])
ORDER BY
	t_attendances_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetAttendanceWithDetails :many
SELECT t_attendances.*, sqlc.embed(t_early_leavings), sqlc.embed(t_late_arrivals), sqlc.embed(t_absences) FROM t_attendances
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
	t_attendances_pkey ASC;

-- name: GetAttendanceWithDetailsUseNumberedPaginate :many
SELECT t_attendances.*, sqlc.embed(t_early_leavings), sqlc.embed(t_late_arrivals), sqlc.embed(t_absences) FROM t_attendances
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
	t_attendances_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetAttendanceWithDetailsUseKeysetPaginate :many
SELECT t_attendances.*, sqlc.embed(t_early_leavings), sqlc.embed(t_late_arrivals), sqlc.embed(t_absences) FROM t_attendances
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
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'date' THEN t_attendances.date > @date_cursor OR (t_attendances.date = @date_cursor AND t_attendances_pkey > @cursor::int)
				WHEN 'r_date' THEN t_attendances.date < @date_cursor OR (t_attendances.date = @date_cursor AND t_attendances_pkey > @cursor::int)
				ELSE t_attendances_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'date' THEN t_attendances.date < @date_cursor OR (t_attendances.date = @date_cursor AND t_attendances_pkey < @cursor::int)
				WHEN 'r_date' THEN t_attendances.date > @date_cursor OR (t_attendances.date = @date_cursor AND t_attendances_pkey < @cursor::int)
				ELSE t_attendances_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'date' AND @cursor_direction::text = 'next' THEN date END ASC,
	CASE WHEN @order_method::text = 'date' AND @cursor_direction::text = 'prev' THEN date END DESC,
	CASE WHEN @order_method::text = 'r_date' AND @cursor_direction::text = 'next' THEN date END DESC,
	CASE WHEN @order_method::text = 'r_date' AND @cursor_direction::text = 'prev' THEN date END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN t_attendances_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_attendances_pkey END DESC
LIMIT $1;

-- name: GetPluralAttendanceWithDetails :many
SELECT t_attendances.*, sqlc.embed(t_early_leavings), sqlc.embed(t_late_arrivals), sqlc.embed(t_absences) FROM t_attendances
LEFT JOIN t_early_leavings ON t_attendances.attendance_id = t_early_leavings.attendance_id
LEFT JOIN t_late_arrivals ON t_attendances.attendance_id = t_late_arrivals.attendance_id
LEFT JOIN t_absences ON t_attendances.attendance_id = t_absences.attendance_id
WHERE attendance_id = ANY(@attendance_ids::uuid[])
ORDER BY
	t_attendances_pkey ASC;

-- name: GetPluralAttendanceWithDetailsUseNumberedPaginate :many
SELECT t_attendances.*, sqlc.embed(t_early_leavings), sqlc.embed(t_late_arrivals), sqlc.embed(t_absences) FROM t_attendances
LEFT JOIN t_early_leavings ON t_attendances.attendance_id = t_early_leavings.attendance_id
LEFT JOIN t_late_arrivals ON t_attendances.attendance_id = t_late_arrivals.attendance_id
LEFT JOIN t_absences ON t_attendances.attendance_id = t_absences.attendance_id
WHERE attendance_id = ANY(@attendance_ids::uuid[])
ORDER BY
	t_attendances_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetAttendanceWithAll :many
SELECT t_attendances.*, sqlc.embed(m_members), m_attendance_types.attendance_type_id, m_attendance_types.name as attendance_type_name, m_attendance_types.key as attendance_type_key, m_attendance_types.color as attendance_type_color, sqlc.embed(m_organizations), sqlc.embed(t_early_leavings), sqlc.embed(t_late_arrivals), sqlc.embed(t_absences) FROM t_attendances
LEFT JOIN t_early_leavings ON t_attendances.attendance_id = t_early_leavings.attendance_id
LEFT JOIN t_late_arrivals ON t_attendances.attendance_id = t_late_arrivals.attendance_id
LEFT JOIN t_absences ON t_attendances.attendance_id = t_absences.attendance_id
LEFT JOIN m_members ON t_attendances.member_id = m_members.member_id
LEFT JOIN m_attend_statuses ON m_members.attend_status_id = m_attend_statuses.attend_status_id
LEFT JOIN m_grades ON m_members.grade_id = m_grades.grade_id
LEFT JOIN m_groups ON m_members.group_id = m_groups.group_id
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
	t_attendances_pkey ASC;

-- name: GetAttendanceWithAllUseNumberedPaginate :many
SELECT t_attendances.*, sqlc.embed(m_members), m_attendance_types.attendance_type_id, m_attendance_types.name as attendance_type_name, m_attendance_types.key as attendance_type_key, m_attendance_types.color as attendance_type_color, sqlc.embed(m_organizations), sqlc.embed(t_early_leavings), sqlc.embed(t_late_arrivals), sqlc.embed(t_absences) FROM t_attendances
LEFT JOIN t_early_leavings ON t_attendances.attendance_id = t_early_leavings.attendance_id
LEFT JOIN t_late_arrivals ON t_attendances.attendance_id = t_late_arrivals.attendance_id
LEFT JOIN t_absences ON t_attendances.attendance_id = t_absences.attendance_id
LEFT JOIN m_members ON t_attendances.member_id = m_members.member_id
LEFT JOIN m_attend_statuses ON m_members.attend_status_id = m_attend_statuses.attend_status_id
LEFT JOIN m_grades ON m_members.grade_id = m_grades.grade_id
LEFT JOIN m_groups ON m_members.group_id = m_groups.group_id
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
	t_attendances_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetAttendanceWithAllUseKeysetPaginate :many
SELECT t_attendances.*, sqlc.embed(m_members), m_attendance_types.attendance_type_id, m_attendance_types.name as attendance_type_name, m_attendance_types.key as attendance_type_key, m_attendance_types.color as attendance_type_color, sqlc.embed(m_organizations), sqlc.embed(t_early_leavings), sqlc.embed(t_late_arrivals), sqlc.embed(t_absences) FROM t_attendances
LEFT JOIN t_early_leavings ON t_attendances.attendance_id = t_early_leavings.attendance_id
LEFT JOIN t_late_arrivals ON t_attendances.attendance_id = t_late_arrivals.attendance_id
LEFT JOIN t_absences ON t_attendances.attendance_id = t_absences.attendance_id
LEFT JOIN m_members ON t_attendances.member_id = m_members.member_id
LEFT JOIN m_attend_statuses ON m_members.attend_status_id = m_attend_statuses.attend_status_id
LEFT JOIN m_grades ON m_members.grade_id = m_grades.grade_id
LEFT JOIN m_groups ON m_members.group_id = m_groups.group_id
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
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'date' THEN t_attendances.date > @date_cursor OR (t_attendances.date = @date_cursor AND t_attendances_pkey > @cursor::int)
				WHEN 'r_date' THEN t_attendances.date < @date_cursor OR (t_attendances.date = @date_cursor AND t_attendances_pkey > @cursor::int)
				ELSE t_attendances_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'date' THEN t_attendances.date < @date_cursor OR (t_attendances.date = @date_cursor AND t_attendances_pkey < @cursor::int)
				WHEN 'r_date' THEN t_attendances.date > @date_cursor OR (t_attendances.date = @date_cursor AND t_attendances_pkey < @cursor::int)
				ELSE t_attendances_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'date' AND @cursor_direction::text = 'next' THEN date END ASC,
	CASE WHEN @order_method::text = 'date' AND @cursor_direction::text = 'prev' THEN date END DESC,
	CASE WHEN @order_method::text = 'r_date' AND @cursor_direction::text = 'next' THEN date END DESC,
	CASE WHEN @order_method::text = 'r_date' AND @cursor_direction::text = 'prev' THEN date END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN t_attendances_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_attendances_pkey END DESC
LIMIT $1;

-- name: GetPluralAttendanceWithAll :many
SELECT t_attendances.*, sqlc.embed(m_members), m_attendance_types.attendance_type_id, m_attendance_types.name as attendance_type_name, m_attendance_types.key as attendance_type_key, m_attendance_types.color as attendance_type_color, sqlc.embed(m_organizations), sqlc.embed(t_early_leavings), sqlc.embed(t_late_arrivals), sqlc.embed(t_absences) FROM t_attendances
LEFT JOIN t_early_leavings ON t_attendances.attendance_id = t_early_leavings.attendance_id
LEFT JOIN t_late_arrivals ON t_attendances.attendance_id = t_late_arrivals.attendance_id
LEFT JOIN t_absences ON t_attendances.attendance_id = t_absences.attendance_id
LEFT JOIN m_members ON t_attendances.member_id = m_members.member_id
LEFT JOIN m_attend_statuses ON m_members.attend_status_id = m_attend_statuses.attend_status_id
LEFT JOIN m_grades ON m_members.grade_id = m_grades.grade_id
LEFT JOIN m_groups ON m_members.group_id = m_groups.group_id
LEFT JOIN m_attendance_types ON t_attendances.attendance_type_id = m_attendance_types.attendance_type_id
LEFT JOIN m_organizations ON t_attendances.send_organization_id = m_organizations.organization_id
WHERE attendance_id = ANY(@attendance_ids::uuid[])
ORDER BY
	t_attendances_pkey ASC;

-- name: GetPluralAttendanceWithAllUseNumberedPaginate :many
SELECT t_attendances.*, sqlc.embed(m_members), m_attendance_types.attendance_type_id, m_attendance_types.name as attendance_type_name, m_attendance_types.key as attendance_type_key, m_attendance_types.color as attendance_type_color, sqlc.embed(m_organizations), sqlc.embed(t_early_leavings), sqlc.embed(t_late_arrivals), sqlc.embed(t_absences) FROM t_attendances
LEFT JOIN t_early_leavings ON t_attendances.attendance_id = t_early_leavings.attendance_id
LEFT JOIN t_late_arrivals ON t_attendances.attendance_id = t_late_arrivals.attendance_id
LEFT JOIN t_absences ON t_attendances.attendance_id = t_absences.attendance_id
LEFT JOIN m_members ON t_attendances.member_id = m_members.member_id
LEFT JOIN m_attend_statuses ON m_members.attend_status_id = m_attend_statuses.attend_status_id
LEFT JOIN m_grades ON m_members.grade_id = m_grades.grade_id
LEFT JOIN m_groups ON m_members.group_id = m_groups.group_id
LEFT JOIN m_attendance_types ON t_attendances.attendance_type_id = m_attendance_types.attendance_type_id
LEFT JOIN m_organizations ON t_attendances.send_organization_id = m_organizations.organization_id
WHERE attendance_id = ANY(@attendance_ids::uuid[])
ORDER BY
	t_attendances_pkey ASC
LIMIT $1 OFFSET $2;

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

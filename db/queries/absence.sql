-- name: CreateAbsences :copyfrom
INSERT INTO t_absences (attendance_id) VALUES ($1);

-- name: CreateAbsence :one
INSERT INTO t_absences (attendance_id) VALUES ($1) RETURNING *;

-- name: DeleteAbsence :exec
DELETE FROM t_absences WHERE absence_id = $1;

-- name: FindAbsenceByID :one
SELECT * FROM t_absences WHERE absence_id = $1;

-- name: GetAbsences :many
SELECT * FROM t_absences
ORDER BY
	t_absences_pkey ASC;

-- name: GetAbsencesUseNumberedPaginate :many
SELECT * FROM t_absences
ORDER BY
	t_absences_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetAbsencesUseKeysetPaginate :many
SELECT * FROM t_absences
WHERE
	CASE @cursor_direction::text
		WHEN 'next' THEN
			t_absences_pkey > @cursor::int
		WHEN 'prev' THEN
			t_absences_pkey < @cursor::int
	END
ORDER BY
	t_absences_pkey ASC
LIMIT $1;

-- name: GetPluralAbsences :many
SELECT * FROM t_absences
WHERE absence_id = ANY(@absence_ids::uuid[])
ORDER BY
	t_absences_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountAbsences :one
SELECT COUNT(*) FROM t_absences;

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
	t_absences_pkey DESC
LIMIT $1 OFFSET $2;

-- name: CountAbsences :one
SELECT COUNT(*) FROM t_absences;

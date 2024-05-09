// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: event_type.sql

package query

import (
	"context"

	"github.com/google/uuid"
)

const countEventTypes = `-- name: CountEventTypes :one
SELECT COUNT(*) FROM m_event_types
WHERE
	CASE WHEN $1::boolean = true THEN name LIKE '%' || $2::text || '%' ELSE TRUE END
`

type CountEventTypesParams struct {
	WhereLikeName bool   `json:"where_like_name"`
	SearchName    string `json:"search_name"`
}

func (q *Queries) CountEventTypes(ctx context.Context, arg CountEventTypesParams) (int64, error) {
	row := q.db.QueryRow(ctx, countEventTypes, arg.WhereLikeName, arg.SearchName)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createEventType = `-- name: CreateEventType :one
INSERT INTO m_event_types (name, key, color) VALUES ($1, $2, $3) RETURNING m_event_types_pkey, event_type_id, name, key, color
`

type CreateEventTypeParams struct {
	Name  string `json:"name"`
	Key   string `json:"key"`
	Color string `json:"color"`
}

func (q *Queries) CreateEventType(ctx context.Context, arg CreateEventTypeParams) (EventType, error) {
	row := q.db.QueryRow(ctx, createEventType, arg.Name, arg.Key, arg.Color)
	var i EventType
	err := row.Scan(
		&i.MEventTypesPkey,
		&i.EventTypeID,
		&i.Name,
		&i.Key,
		&i.Color,
	)
	return i, err
}

type CreateEventTypesParams struct {
	Name  string `json:"name"`
	Key   string `json:"key"`
	Color string `json:"color"`
}

const deleteEventType = `-- name: DeleteEventType :exec
DELETE FROM m_event_types WHERE event_type_id = $1
`

func (q *Queries) DeleteEventType(ctx context.Context, eventTypeID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteEventType, eventTypeID)
	return err
}

const deleteEventTypeByKey = `-- name: DeleteEventTypeByKey :exec
DELETE FROM m_event_types WHERE key = $1
`

func (q *Queries) DeleteEventTypeByKey(ctx context.Context, key string) error {
	_, err := q.db.Exec(ctx, deleteEventTypeByKey, key)
	return err
}

const findEventTypeByID = `-- name: FindEventTypeByID :one
SELECT m_event_types_pkey, event_type_id, name, key, color FROM m_event_types WHERE event_type_id = $1
`

func (q *Queries) FindEventTypeByID(ctx context.Context, eventTypeID uuid.UUID) (EventType, error) {
	row := q.db.QueryRow(ctx, findEventTypeByID, eventTypeID)
	var i EventType
	err := row.Scan(
		&i.MEventTypesPkey,
		&i.EventTypeID,
		&i.Name,
		&i.Key,
		&i.Color,
	)
	return i, err
}

const findEventTypeByKey = `-- name: FindEventTypeByKey :one
SELECT m_event_types_pkey, event_type_id, name, key, color FROM m_event_types WHERE key = $1
`

func (q *Queries) FindEventTypeByKey(ctx context.Context, key string) (EventType, error) {
	row := q.db.QueryRow(ctx, findEventTypeByKey, key)
	var i EventType
	err := row.Scan(
		&i.MEventTypesPkey,
		&i.EventTypeID,
		&i.Name,
		&i.Key,
		&i.Color,
	)
	return i, err
}

const getEventTypes = `-- name: GetEventTypes :many
SELECT m_event_types_pkey, event_type_id, name, key, color FROM m_event_types
WHERE
	CASE WHEN $1::boolean = true THEN name LIKE '%' || $2::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN $3::text = 'name' THEN name END ASC,
	CASE WHEN $3::text = 'r_name' THEN name END DESC,
	m_event_types_pkey ASC
`

type GetEventTypesParams struct {
	WhereLikeName bool   `json:"where_like_name"`
	SearchName    string `json:"search_name"`
	OrderMethod   string `json:"order_method"`
}

func (q *Queries) GetEventTypes(ctx context.Context, arg GetEventTypesParams) ([]EventType, error) {
	rows, err := q.db.Query(ctx, getEventTypes, arg.WhereLikeName, arg.SearchName, arg.OrderMethod)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []EventType{}
	for rows.Next() {
		var i EventType
		if err := rows.Scan(
			&i.MEventTypesPkey,
			&i.EventTypeID,
			&i.Name,
			&i.Key,
			&i.Color,
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

const getEventTypesUseKeysetPaginate = `-- name: GetEventTypesUseKeysetPaginate :many
SELECT m_event_types_pkey, event_type_id, name, key, color FROM m_event_types
WHERE
	CASE WHEN $2::boolean = true THEN name LIKE '%' || $3::text || '%' ELSE TRUE END
AND
	CASE $4::text
		WHEN 'next' THEN
			CASE $5::text
				WHEN 'name' THEN name > $6 OR (name = $6 AND m_event_types_pkey > $7::int)
				WHEN 'r_name' THEN name < $6 OR (name = $6 AND m_event_types_pkey > $7::int)
				ELSE m_event_types_pkey > $7::int
			END
		WHEN 'prev' THEN
			CASE $5::text
				WHEN 'name' THEN name < $6 OR (name = $6 AND m_event_types_pkey < $7::int)
				WHEN 'r_name' THEN name > $6 OR (name = $6 AND m_event_types_pkey < $7::int)
				ELSE m_event_types_pkey < $7::int
			END
	END
ORDER BY
	CASE WHEN $5::text = 'name' AND $4::text = 'next' THEN name END ASC,
	CASE WHEN $5::text = 'name' AND $4::text = 'prev' THEN name END DESC,
	CASE WHEN $5::text = 'r_name' AND $4::text = 'next' THEN name END ASC,
	CASE WHEN $5::text = 'r_name' AND $4::text = 'prev' THEN name END DESC,
	CASE WHEN $4::text = 'next' THEN m_event_types_pkey END ASC,
	CASE WHEN $4::text = 'prev' THEN m_event_types_pkey END DESC
LIMIT $1
`

type GetEventTypesUseKeysetPaginateParams struct {
	Limit           int32  `json:"limit"`
	WhereLikeName   bool   `json:"where_like_name"`
	SearchName      string `json:"search_name"`
	CursorDirection string `json:"cursor_direction"`
	OrderMethod     string `json:"order_method"`
	NameCursor      string `json:"name_cursor"`
	Cursor          int32  `json:"cursor"`
}

func (q *Queries) GetEventTypesUseKeysetPaginate(ctx context.Context, arg GetEventTypesUseKeysetPaginateParams) ([]EventType, error) {
	rows, err := q.db.Query(ctx, getEventTypesUseKeysetPaginate,
		arg.Limit,
		arg.WhereLikeName,
		arg.SearchName,
		arg.CursorDirection,
		arg.OrderMethod,
		arg.NameCursor,
		arg.Cursor,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []EventType{}
	for rows.Next() {
		var i EventType
		if err := rows.Scan(
			&i.MEventTypesPkey,
			&i.EventTypeID,
			&i.Name,
			&i.Key,
			&i.Color,
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

const getEventTypesUseNumberedPaginate = `-- name: GetEventTypesUseNumberedPaginate :many
SELECT m_event_types_pkey, event_type_id, name, key, color FROM m_event_types
WHERE
	CASE WHEN $3::boolean = true THEN name LIKE '%' || $4::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN $5::text = 'name' THEN name END ASC,
	CASE WHEN $5::text = 'r_name' THEN name END DESC,
	m_event_types_pkey ASC
LIMIT $1 OFFSET $2
`

type GetEventTypesUseNumberedPaginateParams struct {
	Limit         int32  `json:"limit"`
	Offset        int32  `json:"offset"`
	WhereLikeName bool   `json:"where_like_name"`
	SearchName    string `json:"search_name"`
	OrderMethod   string `json:"order_method"`
}

func (q *Queries) GetEventTypesUseNumberedPaginate(ctx context.Context, arg GetEventTypesUseNumberedPaginateParams) ([]EventType, error) {
	rows, err := q.db.Query(ctx, getEventTypesUseNumberedPaginate,
		arg.Limit,
		arg.Offset,
		arg.WhereLikeName,
		arg.SearchName,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []EventType{}
	for rows.Next() {
		var i EventType
		if err := rows.Scan(
			&i.MEventTypesPkey,
			&i.EventTypeID,
			&i.Name,
			&i.Key,
			&i.Color,
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

const getPluralEventTypes = `-- name: GetPluralEventTypes :many
SELECT m_event_types_pkey, event_type_id, name, key, color FROM m_event_types WHERE event_type_id = ANY($3::uuid[])
ORDER BY
	m_event_types_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralEventTypesParams struct {
	Limit        int32       `json:"limit"`
	Offset       int32       `json:"offset"`
	EventTypeIds []uuid.UUID `json:"event_type_ids"`
}

func (q *Queries) GetPluralEventTypes(ctx context.Context, arg GetPluralEventTypesParams) ([]EventType, error) {
	rows, err := q.db.Query(ctx, getPluralEventTypes, arg.Limit, arg.Offset, arg.EventTypeIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []EventType{}
	for rows.Next() {
		var i EventType
		if err := rows.Scan(
			&i.MEventTypesPkey,
			&i.EventTypeID,
			&i.Name,
			&i.Key,
			&i.Color,
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

const pluralDeleteEventTypes = `-- name: PluralDeleteEventTypes :exec
DELETE FROM m_event_types WHERE event_type_id = ANY($1::uuid[])
`

func (q *Queries) PluralDeleteEventTypes(ctx context.Context, dollar_1 []uuid.UUID) error {
	_, err := q.db.Exec(ctx, pluralDeleteEventTypes, dollar_1)
	return err
}

const updateEventType = `-- name: UpdateEventType :one
UPDATE m_event_types SET name = $2, key = $3, color = $4 WHERE event_type_id = $1 RETURNING m_event_types_pkey, event_type_id, name, key, color
`

type UpdateEventTypeParams struct {
	EventTypeID uuid.UUID `json:"event_type_id"`
	Name        string    `json:"name"`
	Key         string    `json:"key"`
	Color       string    `json:"color"`
}

func (q *Queries) UpdateEventType(ctx context.Context, arg UpdateEventTypeParams) (EventType, error) {
	row := q.db.QueryRow(ctx, updateEventType,
		arg.EventTypeID,
		arg.Name,
		arg.Key,
		arg.Color,
	)
	var i EventType
	err := row.Scan(
		&i.MEventTypesPkey,
		&i.EventTypeID,
		&i.Name,
		&i.Key,
		&i.Color,
	)
	return i, err
}

const updateEventTypeByKey = `-- name: UpdateEventTypeByKey :one
UPDATE m_event_types SET name = $2, color = $3 WHERE key = $1 RETURNING m_event_types_pkey, event_type_id, name, key, color
`

type UpdateEventTypeByKeyParams struct {
	Key   string `json:"key"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

func (q *Queries) UpdateEventTypeByKey(ctx context.Context, arg UpdateEventTypeByKeyParams) (EventType, error) {
	row := q.db.QueryRow(ctx, updateEventTypeByKey, arg.Key, arg.Name, arg.Color)
	var i EventType
	err := row.Scan(
		&i.MEventTypesPkey,
		&i.EventTypeID,
		&i.Name,
		&i.Key,
		&i.Color,
	)
	return i, err
}

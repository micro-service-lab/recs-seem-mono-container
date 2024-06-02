-- name: CreateAttachableItems :copyfrom
INSERT INTO t_attachable_items (url, size, alias, owner_id, from_outer, mime_type_id) VALUES ($1, $2, $3, $4, $5, $6);

-- name: CreateAttachableItem :one
INSERT INTO t_attachable_items (url, size, alias, owner_id, from_outer, mime_type_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: UpdateAttachableItem :one
UPDATE t_attachable_items SET url = $2, size = $3, alias = $4, mime_type_id = $5 WHERE attachable_item_id = $1 RETURNING *;

-- name: DeleteAttachableItem :execrows
DELETE FROM t_attachable_items WHERE attachable_item_id = $1;

-- name: PluralDeleteAttachableItems :execrows
DELETE FROM t_attachable_items WHERE t_attachable_items.attachable_item_id = ANY(@attachable_item_ids::uuid[]);

-- name: FindAttachableItemByID :one
SELECT t_attachable_items.*, t_images.image_id, t_images.height image_height,
t_images.width image_width, t_files.file_id
FROM t_attachable_items
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE t_attachable_items.attachable_item_id = $1;

-- name: FindAttachableItemByIDWithMimeType :one
SELECT t_attachable_items.*, m_mime_types.name mime_type_name, m_mime_types.kind mime_type_kind,
m_mime_types.key mime_type_key, t_images.height image_height, t_images.width image_width, t_images.image_id, t_files.file_id
FROM t_attachable_items
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE t_attachable_items.attachable_item_id = $1;

-- name: FindAttachableItemByURL :one
SELECT t_attachable_items.*, t_images.image_id, t_images.height image_height,
t_images.width image_width, t_files.file_id
FROM t_attachable_items
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE t_attachable_items.url = $1;

-- name: FindAttachableItemByURLWithMimeType :one
SELECT t_attachable_items.*, m_mime_types.name mime_type_name, m_mime_types.kind mime_type_kind,
m_mime_types.key mime_type_key, t_images.height image_height, t_images.width image_width, t_images.image_id, t_files.file_id
FROM t_attachable_items
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE t_attachable_items.url = $1;

-- name: GetAttachableItems :many
SELECT t_attachable_items.*, t_images.image_id, t_images.height image_height,
t_images.width image_width, t_files.file_id
FROM t_attachable_items
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE
	CASE WHEN @where_in_mime_type_ids::boolean = true THEN mime_type_id = ANY(@in_mime_type_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_in_owner_ids::boolean = true THEN owner_id = ANY(@in_owner_ids::uuid[]) ELSE TRUE END
ORDER BY
	t_attachable_items_pkey ASC;

-- name: GetAttachableItemsUseNumberedPaginate :many
SELECT t_attachable_items.*, t_images.image_id, t_images.height image_height,
t_images.width image_width, t_files.file_id
FROM t_attachable_items
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE
	CASE WHEN @where_in_mime_type_ids::boolean = true THEN mime_type_id = ANY(@in_mime_type_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_in_owner_ids::boolean = true THEN owner_id = ANY(@in_owner_ids::uuid[]) ELSE TRUE END
ORDER BY
	t_attachable_items_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetAttachableItemsUseKeysetPaginate :many
SELECT t_attachable_items.*, t_images.image_id, t_images.height image_height,
t_images.width image_width, t_files.file_id
FROM t_attachable_items
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE
	CASE WHEN @where_in_mime_type_ids::boolean = true THEN mime_type_id = ANY(@in_mime_type_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_in_owner_ids::boolean = true THEN owner_id = ANY(@in_owner_ids::uuid[]) ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			t_attachable_items_pkey < @cursor::int
		WHEN 'prev' THEN
			t_attachable_items_pkey > @cursor::int
	END
ORDER BY
	CASE WHEN @cursor_direction::text = 'next' THEN t_attachable_items_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_attachable_items_pkey END DESC
LIMIT $1;

-- name: GetPluralAttachableItems :many
SELECT t_attachable_items.*, t_images.image_id, t_images.height image_height,
t_images.width image_width, t_files.file_id
FROM t_attachable_items
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE t_attachable_items.attachable_item_id = ANY(@attachable_item_ids::uuid[])
ORDER BY
	t_attachable_items_pkey ASC;

-- name: GetPluralAttachableItemsUseNumberedPaginate :many
SELECT t_attachable_items.*, t_images.image_id, t_images.height image_height,
t_images.width image_width, t_files.file_id
FROM t_attachable_items
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE t_attachable_items.attachable_item_id = ANY(@attachable_item_ids::uuid[])
ORDER BY
	t_attachable_items_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetAttachableItemsWithMimeType :many
SELECT t_attachable_items.*, m_mime_types.name mime_type_name, m_mime_types.kind mime_type_kind,
m_mime_types.key mime_type_key, t_images.height image_height, t_images.width image_width, t_images.image_id, t_files.file_id
FROM t_attachable_items
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.where_mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE
	CASE WHEN @where_in_mime_type_ids::boolean = true THEN mime_type_id = ANY(@in_mime_type_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_in_owner_ids::boolean = true THEN owner_id = ANY(@in_owner_ids::uuid[]) ELSE TRUE END
ORDER BY
	t_attachable_items_pkey ASC;

-- name: GetAttachableItemsWithMimeTypeUseNumberedPaginate :many
SELECT t_attachable_items.*, m_mime_types.name mime_type_name, m_mime_types.kind mime_type_kind,
m_mime_types.key mime_type_key, t_images.height image_height, t_images.width image_width, t_images.image_id, t_files.file_id
FROM t_attachable_items
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE
	CASE WHEN @where_in_mime_type_ids::boolean = true THEN mime_type_id = ANY(@in_mime_type_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_in_owner_ids::boolean = true THEN owner_id = ANY(@in_owner_ids::uuid[]) ELSE TRUE END
ORDER BY
	t_attachable_items_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetAttachableItemsWithMimeTypeUseKeysetPaginate :many
SELECT t_attachable_items.*, m_mime_types.name mime_type_name, m_mime_types.kind mime_type_kind,
m_mime_types.key mime_type_key, t_images.height image_height, t_images.width image_width, t_images.image_id, t_files.file_id
FROM t_attachable_items
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE
	CASE WHEN @where_in_mime_type_ids::boolean = true THEN mime_type_id = ANY(@in_mime_type_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_in_owner_ids::boolean = true THEN owner_id = ANY(@in_owner_ids::uuid[]) ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			t_attachable_items_pkey > @cursor::int
		WHEN 'prev' THEN
			t_attachable_items_pkey < @cursor::int
	END
ORDER BY
	CASE WHEN @cursor_direction::text = 'next' THEN t_attachable_items_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_attachable_items_pkey END DESC
LIMIT $1;

-- name: GetPluralAttachableItemsWithMimeType :many
SELECT t_attachable_items.*, m_mime_types.name mime_type_name, m_mime_types.kind mime_type_kind,
m_mime_types.key mime_type_key, t_images.height image_height, t_images.width image_width, t_images.image_id, t_files.file_id
FROM t_attachable_items
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE t_attachable_items.attachable_item_id = ANY(@attachable_item_ids::uuid[])
ORDER BY
	t_attachable_items_pkey ASC;

-- name: GetPluralAttachableItemsWithMimeTypeUseNumberedPaginate :many
SELECT t_attachable_items.*, m_mime_types.name mime_type_name, m_mime_types.kind mime_type_kind,
m_mime_types.key mime_type_key, t_images.height image_height, t_images.width image_width, t_images.image_id, t_files.file_id
FROM t_attachable_items
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE t_attachable_items.attachable_item_id = ANY(@attachable_item_ids::uuid[])
ORDER BY
	t_attachable_items_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountAttachableItems :one
SELECT COUNT(*) FROM t_attachable_items
WHERE
	CASE WHEN @where_in_mime_type_ids::boolean = true THEN mime_type_id = ANY(@in_mime_type_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_in_owner_ids::boolean = true THEN owner_id = ANY(@in_owner_ids::uuid[]) ELSE TRUE END;

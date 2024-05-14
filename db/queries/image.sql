-- name: CreateImages :copyfrom
INSERT INTO t_images (height, width, attachable_item_id) VALUES ($1, $2, $3);

-- name: CreateImage :one
INSERT INTO t_images (height, width, attachable_item_id) VALUES ($1, $2, $3) RETURNING *;

-- name: DeleteImage :execrows
DELETE FROM t_images WHERE image_id = $1;

-- name: PluralDeleteImages :execrows
DELETE FROM t_images WHERE image_id = ANY($1::uuid[]);

-- name: FindImageByID :one
SELECT * FROM t_images WHERE image_id = $1;

-- name: FindImageByIDWithAttachableItem :one
SELECT sqlc.embed(t_images), sqlc.embed(t_attachable_items), sqlc.embed(m_mime_types) FROM t_images
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
WHERE image_id = $1;

-- name: GetImages :many
SELECT * FROM t_images
ORDER BY
	t_images_pkey ASC;

-- name: GetImagesUseNumberedPaginate :many
SELECT * FROM t_images
ORDER BY
	t_images_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetImagesUseKeysetPaginate :many
SELECT * FROM t_images
WHERE
	CASE @cursor_direction::text
		WHEN 'next' THEN
			t_images_pkey > @cursor::int
		WHEN 'prev' THEN
			t_images_pkey < @cursor::int
	END
ORDER BY
	CASE WHEN @cursor_direction::text = 'next' THEN t_images_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_images_pkey END DESC
LIMIT $1;

-- name: GetPluralImages :many
SELECT * FROM t_images
WHERE attachable_item_id = ANY(@attachable_item_ids::uuid[])
ORDER BY
	t_images_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetImagesWithAttachableItem :many
SELECT sqlc.embed(t_images), sqlc.embed(t_attachable_items), sqlc.embed(m_mime_types) FROM t_images
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
ORDER BY
	t_images_pkey ASC;

-- name: GetImagesWithAttachableItemUseNumberedPaginate :many
SELECT sqlc.embed(t_images), sqlc.embed(t_attachable_items), sqlc.embed(m_mime_types) FROM t_images
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
ORDER BY
	t_images_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetImagesWithAttachableItemUseKeysetPaginate :many
SELECT sqlc.embed(t_images), sqlc.embed(t_attachable_items), sqlc.embed(m_mime_types) FROM t_images
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
WHERE
	CASE @cursor_direction::text
		WHEN 'next' THEN
			t_images_pkey > @cursor::int
		WHEN 'prev' THEN
			t_images_pkey < @cursor::int
	END
ORDER BY
	CASE WHEN @cursor_direction::text = 'next' THEN t_images_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_images_pkey END DESC
LIMIT $1;

-- name: GetPluralImagesWithAttachableItem :many
SELECT sqlc.embed(t_images), sqlc.embed(t_attachable_items), sqlc.embed(m_mime_types) FROM t_images
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
WHERE attachable_item_id = ANY(@attachable_item_ids::uuid[])
ORDER BY
	t_images_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountImages :one
SELECT COUNT(*) FROM t_images;

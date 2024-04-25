-- name: CreateImages :copyfrom
INSERT INTO t_images (height, width, attachable_item_id) VALUES ($1, $2, $3);

-- name: CreateImage :one
INSERT INTO t_images (height, width, attachable_item_id) VALUES ($1, $2, $3) RETURNING *;

-- name: DeleteImage :exec
DELETE FROM t_images WHERE image_id = $1;

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
	t_images_pkey DESC;

-- name: GetImagesUseNumberedPaginate :many
SELECT * FROM t_images
ORDER BY
	t_images_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetImagesUseKeysetPaginate :many
SELECT * FROM t_images
WHERE
	CASE @cursor_direction
		WHEN 'next' THEN
			t_images_pkey < @cursor
		WHEN 'prev' THEN
			t_images_pkey > @cursor
	END
ORDER BY
	t_images_pkey DESC
LIMIT $1;

-- name: GetImagesWithAttachableItem :many
SELECT sqlc.embed(t_images), sqlc.embed(t_attachable_items), sqlc.embed(m_mime_types) FROM t_images
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
ORDER BY
	t_images_pkey DESC;

-- name: GetImagesWithAttachableItem :many
SELECT sqlc.embed(t_images), sqlc.embed(t_attachable_items), sqlc.embed(m_mime_types) FROM t_images
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
ORDER BY
	t_images_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetImagesWithAttachableItem :many
SELECT sqlc.embed(t_images), sqlc.embed(t_attachable_items), sqlc.embed(m_mime_types) FROM t_images
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
WHERE
	CASE @cursor_direction
		WHEN 'next' THEN
			t_images_pkey < @cursor
		WHEN 'prev' THEN
			t_images_pkey > @cursor
	END
ORDER BY
	t_images_pkey DESC
LIMIT $1;

-- name: CountImages :one
SELECT COUNT(*) FROM t_images;

-- name: CreateImages :copyfrom
INSERT INTO t_images (height, width, attachable_item_id) VALUES ($1, $2, $3);

-- name: CreateImage :one
INSERT INTO t_images (height, width, attachable_item_id) VALUES ($1, $2, $3) RETURNING *;

-- name: DeleteImage :execrows
DELETE FROM t_images WHERE image_id = $1;

-- name: PluralDeleteImages :execrows
DELETE FROM t_images WHERE image_id = ANY(@image_ids::uuid[]);

-- name: FindImageByID :one
SELECT * FROM t_images WHERE image_id = $1;

-- name: FindImageByIDWithAttachableItem :one
SELECT t_images.*, t_attachable_items.owner_id owner_id, t_attachable_items.from_outer from_outer, t_attachable_items.alias alias,
t_attachable_items.url url, t_attachable_items.size size, t_attachable_items.mime_type_id mime_type_id FROM t_images
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
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
WHERE image_id = ANY(@image_ids::uuid[])
ORDER BY
	t_images_pkey ASC;

-- name: GetPluralImagesUseNumberedPaginate :many
SELECT * FROM t_images
WHERE image_id = ANY(@image_ids::uuid[])
ORDER BY
	t_images_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetImagesWithAttachableItem :many
SELECT t_images.*, t_attachable_items.owner_id owner_id, t_attachable_items.from_outer from_outer, t_attachable_items.alias alias,
t_attachable_items.url url, t_attachable_items.size size, t_attachable_items.mime_type_id mime_type_id FROM t_images
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
ORDER BY
	t_images_pkey ASC;

-- name: GetImagesWithAttachableItemUseNumberedPaginate :many
SELECT t_images.*, t_attachable_items.owner_id owner_id, t_attachable_items.from_outer from_outer, t_attachable_items.alias alias,
t_attachable_items.url url, t_attachable_items.size size, t_attachable_items.mime_type_id mime_type_id FROM t_images
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
ORDER BY
	t_images_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetImagesWithAttachableItemUseKeysetPaginate :many
SELECT t_images.*, t_attachable_items.owner_id owner_id, t_attachable_items.from_outer from_outer, t_attachable_items.alias alias,
t_attachable_items.url url, t_attachable_items.size size, t_attachable_items.mime_type_id mime_type_id FROM t_images
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
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
SELECT t_images.*, t_attachable_items.owner_id owner_id, t_attachable_items.from_outer from_outer, t_attachable_items.alias alias,
t_attachable_items.url url, t_attachable_items.size size, t_attachable_items.mime_type_id mime_type_id FROM t_images
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE image_id = ANY(@image_ids::uuid[])
ORDER BY
	t_images_pkey ASC;

-- name: GetPluralImagesWithAttachableItemUseNumberedPaginate :many
SELECT t_images.*, t_attachable_items.owner_id owner_id, t_attachable_items.from_outer from_outer, t_attachable_items.alias alias,
t_attachable_items.url url, t_attachable_items.size size, t_attachable_items.mime_type_id mime_type_id FROM t_images
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE image_id = ANY(@image_ids::uuid[])
ORDER BY
	t_images_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountImages :one
SELECT COUNT(*) FROM t_images;

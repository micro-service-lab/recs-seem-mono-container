package service

//nolint:revive
import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/ccitt"
	_ "golang.org/x/image/riff"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/vp8"
	_ "golang.org/x/image/vp8l"
	_ "golang.org/x/image/webp"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/storage"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// ManageImage 画像管理サービス。
type ManageImage struct {
	DB      store.Store
	Storage storage.Storage
}

// DefaultImages デフォルト画像。
//
//go:embed static/images/*
var DefaultImages embed.FS

// DefaultImageKeys デフォルト画像のキー。
var DefaultImageKeys []string

func init() {
	entries, err := DefaultImages.ReadDir("static/images")
	if err != nil {
		panic(fmt.Sprintf("failed to read directory: %v", err))
	}
	for _, e := range entries {
		DefaultImageKeys = append(DefaultImageKeys, e.Name())
	}
}

// CreateImage 画像を作成する。
func (m *ManageImage) CreateImage(
	ctx context.Context,
	origin io.Reader,
	alias string,
	ownerID entity.UUID,
) (e entity.Image, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.Image{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	var storageKeys []string
	defer func() {
		if err != nil {
			if rerr := m.DB.Rollback(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to rollback transaction: %w", rerr)
			}
			if len(storageKeys) > 0 {
				if rerr := m.Storage.DeleteObjects(ctx, storageKeys); rerr != nil {
					err = fmt.Errorf("failed to delete objects: %w", rerr)
				}
			}
		} else {
			if rerr := m.DB.Commit(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", rerr)
			}
		}
	}()
	if ownerID.Valid {
		// 画像の所有者が存在する場合は、画像の所有者が存在するか確認する。
		_, err := m.DB.FindMemberByIDWithSd(ctx, sd, ownerID.Bytes)
		if err != nil {
			return entity.Image{}, errhandle.NewModelNotFoundError("member")
		}
	}
	mtype, err := mimetype.DetectReader(origin)
	if err != nil {
		return entity.Image{}, fmt.Errorf("failed to detect mimetype: %w", err)
	}
	mime, err := m.DB.FindMimeTypeByKindWithSd(ctx, sd, mtype.String())
	if err != nil {
		return entity.Image{}, errhandle.NewModelNotFoundError("mime type")
	}
	uid, err := uuid.NewRandom()
	if err != nil {
		return entity.Image{}, fmt.Errorf("failed to generate uuid: %w", err)
	}
	fname := fmt.Sprintf("%s.%s", uid.String(), mtype.Extension())
	url, err := m.Storage.UploadObject(ctx, origin, fname)
	if err != nil {
		return entity.Image{}, errhandle.NewCommonError(response.FailedUpload, nil)
	}
	storageKeys = append(storageKeys, fname)
	buf := new(bytes.Buffer)
	size, err := io.Copy(buf, origin)
	if err != nil {
		return entity.Image{}, fmt.Errorf("failed to copy file: %w", err)
	}
	fsize := entity.Float{
		Valid:   true,
		Float64: float64(size),
	}
	img, _, err := image.Decode(buf)
	var imgHeight, imgWidth entity.Float
	if err == nil {
		imgHeight = entity.Float{Valid: true, Float64: float64(img.Bounds().Dy())}
		imgWidth = entity.Float{Valid: true, Float64: float64(img.Bounds().Dx())}
	}
	caip := parameter.CreateAttachableItemParam{
		URL:        url,
		Size:       fsize,
		OwnerID:    ownerID,
		FromOuter:  false,
		Alias:      alias,
		MimeTypeID: mime.MimeTypeID,
	}
	ai, err := m.DB.CreateAttachableItemWithSd(ctx, sd, caip)
	if err != nil {
		return entity.Image{}, fmt.Errorf("failed to create attachable item: %w", err)
	}
	p := parameter.CreateImageParam{
		Height:           imgHeight,
		Width:            imgWidth,
		AttachableItemID: ai.AttachableItemID,
	}
	e, err = m.DB.CreateImageWithSd(ctx, sd, p)
	if err != nil {
		return entity.Image{}, fmt.Errorf("failed to create image: %w", err)
	}
	return e, nil
}

// CreateImages 画像を複数作成する。
func (m *ManageImage) CreateImages(
	ctx context.Context,
	ownerID entity.UUID,
	params []parameter.CreateImageServiceParam,
) (es []entity.Image, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	var storageKeys []string
	defer func() {
		if err != nil {
			if rerr := m.DB.Rollback(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to rollback transaction: %w", rerr)
			}
			if len(storageKeys) > 0 {
				if rerr := m.Storage.DeleteObjects(ctx, storageKeys); rerr != nil {
					err = fmt.Errorf("failed to delete objects: %w", rerr)
				}
			}
		} else {
			if rerr := m.DB.Commit(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", rerr)
			}
		}
	}()
	if ownerID.Valid {
		// 画像の所有者が存在する場合は、画像の所有者が存在するか確認する。
		_, err := m.DB.FindMemberByIDWithSd(ctx, sd, ownerID.Bytes)
		if err != nil {
			return nil, errhandle.NewModelNotFoundError("member")
		}
	}
	var images []entity.Image
	for _, p := range params {
		mtype, err := mimetype.DetectReader(p.Origin)
		if err != nil {
			return nil, fmt.Errorf("failed to detect mimetype: %w", err)
		}
		mime, err := m.DB.FindMimeTypeByKindWithSd(ctx, sd, mtype.String())
		if err != nil {
			return nil, errhandle.NewModelNotFoundError("mime type")
		}
		uid, err := uuid.NewRandom()
		if err != nil {
			return nil, fmt.Errorf("failed to generate uuid: %w", err)
		}
		fname := fmt.Sprintf("%s.%s", uid.String(), mtype.Extension())
		url, err := m.Storage.UploadObject(ctx, p.Origin, fname)
		if err != nil {
			return nil, errhandle.NewCommonError(response.FailedUpload, nil)
		}
		storageKeys = append(storageKeys, fname)
		buf := new(bytes.Buffer)
		size, err := io.Copy(buf, p.Origin)
		if err != nil {
			return nil, fmt.Errorf("failed to copy file: %w", err)
		}
		fsize := entity.Float{
			Valid:   true,
			Float64: float64(size),
		}
		img, _, err := image.Decode(buf)
		var imgHeight, imgWidth entity.Float
		if err == nil {
			imgHeight = entity.Float{Valid: true, Float64: float64(img.Bounds().Dy())}
			imgWidth = entity.Float{Valid: true, Float64: float64(img.Bounds().Dx())}
		}
		caip := parameter.CreateAttachableItemParam{
			URL:        url,
			Alias:      p.Alias,
			Size:       fsize,
			OwnerID:    ownerID,
			FromOuter:  false,
			MimeTypeID: mime.MimeTypeID,
		}
		ai, err := m.DB.CreateAttachableItemWithSd(ctx, sd, caip)
		if err != nil {
			return nil, fmt.Errorf("failed to create attachable item: %w", err)
		}
		p := parameter.CreateImageParam{
			Height:           imgHeight,
			Width:            imgWidth,
			AttachableItemID: ai.AttachableItemID,
		}
		e, err := m.DB.CreateImageWithSd(ctx, sd, p)
		if err != nil {
			return nil, fmt.Errorf("failed to create image: %w", err)
		}
		images = append(images, e)
	}
	return images, nil
}

// CreateImageSpecifyFilename 画像を作成する。
func (m *ManageImage) CreateImageSpecifyFilename(
	ctx context.Context,
	origin io.Reader,
	alias string,
	ownerID entity.UUID,
	filename string,
) (e entity.Image, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.Image{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	var storageKeys []string
	defer func() {
		if err != nil {
			if rerr := m.DB.Rollback(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to rollback transaction: %w", rerr)
			}
			if len(storageKeys) > 0 {
				if rerr := m.Storage.DeleteObjects(ctx, storageKeys); rerr != nil {
					err = fmt.Errorf("failed to delete objects: %w", rerr)
				}
			}
		} else {
			if rerr := m.DB.Commit(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", rerr)
			}
		}
	}()
	if ownerID.Valid {
		// 画像の所有者が存在する場合は、画像の所有者が存在するか確認する。
		_, err := m.DB.FindMemberByIDWithSd(ctx, sd, ownerID.Bytes)
		if err != nil {
			return entity.Image{}, errhandle.NewModelNotFoundError("member")
		}
	}
	mtype, err := mimetype.DetectReader(origin)
	if err != nil {
		return entity.Image{}, fmt.Errorf("failed to detect mimetype: %w", err)
	}
	mime, err := m.DB.FindMimeTypeByKindWithSd(ctx, sd, mtype.String())
	if err != nil {
		return entity.Image{}, errhandle.NewModelNotFoundError("mime type")
	}
	exist, err := m.Storage.ExistsObject(ctx, filename)
	if err != nil {
		return entity.Image{}, fmt.Errorf("failed to check object existence: %w", err)
	}
	if exist {
		return entity.Image{}, errhandle.NewCommonError(response.ConflictStorageKey, nil)
	}
	url, err := m.Storage.UploadObject(ctx, origin, filename)
	if err != nil {
		return entity.Image{}, errhandle.NewCommonError(response.FailedUpload, nil)
	}
	storageKeys = append(storageKeys, filename)
	buf := new(bytes.Buffer)
	size, err := io.Copy(buf, origin)
	if err != nil {
		return entity.Image{}, fmt.Errorf("failed to copy file: %w", err)
	}
	fsize := entity.Float{
		Valid:   true,
		Float64: float64(size),
	}
	img, _, err := image.Decode(buf)
	var imgHeight, imgWidth entity.Float
	if err == nil {
		imgHeight = entity.Float{Valid: true, Float64: float64(img.Bounds().Dy())}
		imgWidth = entity.Float{Valid: true, Float64: float64(img.Bounds().Dx())}
	}
	caip := parameter.CreateAttachableItemParam{
		URL:        url,
		Alias:      alias,
		Size:       fsize,
		OwnerID:    ownerID,
		FromOuter:  false,
		MimeTypeID: mime.MimeTypeID,
	}
	ai, err := m.DB.CreateAttachableItemWithSd(ctx, sd, caip)
	p := parameter.CreateImageParam{
		Height:           imgHeight,
		Width:            imgWidth,
		AttachableItemID: ai.AttachableItemID,
	}
	e, err = m.DB.CreateImageWithSd(ctx, sd, p)
	if err != nil {
		return entity.Image{}, fmt.Errorf("failed to create image: %w", err)
	}
	return e, nil
}

// CreateImagesSpecifyFilename 画像を複数作成する。
func (m *ManageImage) CreateImagesSpecifyFilename(
	ctx context.Context,
	ownerID entity.UUID,
	params []parameter.CreateImageSpecifyFilenameServiceParam,
) (es []entity.Image, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	var storageKeys []string
	defer func() {
		if err != nil {
			if rerr := m.DB.Rollback(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to rollback transaction: %w", rerr)
			}
			if len(storageKeys) > 0 {
				if rerr := m.Storage.DeleteObjects(ctx, storageKeys); rerr != nil {
					err = fmt.Errorf("failed to delete objects: %w", rerr)
				}
			}
		} else {
			if rerr := m.DB.Commit(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", rerr)
			}
		}
	}()
	if ownerID.Valid {
		// 画像の所有者が存在する場合は、画像の所有者が存在するか確認する。
		_, err := m.DB.FindMemberByIDWithSd(ctx, sd, ownerID.Bytes)
		if err != nil {
			return nil, errhandle.NewModelNotFoundError("member")
		}
	}
	var images []entity.Image
	for _, p := range params {
		mtype, err := mimetype.DetectReader(p.Origin)
		if err != nil {
			return nil, fmt.Errorf("failed to detect mimetype: %w", err)
		}
		mime, err := m.DB.FindMimeTypeByKindWithSd(ctx, sd, mtype.String())
		if err != nil {
			return nil, errhandle.NewModelNotFoundError("mime type")
		}
		exist, err := m.Storage.ExistsObject(ctx, p.Filename)
		if err != nil {
			return nil, fmt.Errorf("failed to check object existence: %w", err)
		}
		if exist {
			return nil, errhandle.NewCommonError(response.ConflictStorageKey, nil)
		}
		url, err := m.Storage.UploadObject(ctx, p.Origin, p.Filename)
		if err != nil {
			return nil, errhandle.NewCommonError(response.FailedUpload, nil)
		}
		storageKeys = append(storageKeys, p.Filename)
		buf := new(bytes.Buffer)
		size, err := io.Copy(buf, p.Origin)
		if err != nil {
			return nil, fmt.Errorf("failed to copy file: %w", err)
		}
		fsize := entity.Float{
			Valid:   true,
			Float64: float64(size),
		}
		img, _, err := image.Decode(buf)
		var imgHeight, imgWidth entity.Float
		if err == nil {
			imgHeight = entity.Float{Valid: true, Float64: float64(img.Bounds().Dy())}
			imgWidth = entity.Float{Valid: true, Float64: float64(img.Bounds().Dx())}
		}
		caip := parameter.CreateAttachableItemParam{
			URL:        url,
			Alias:      p.Alias,
			Size:       fsize,
			OwnerID:    ownerID,
			FromOuter:  false,
			MimeTypeID: mime.MimeTypeID,
		}
		ai, err := m.DB.CreateAttachableItemWithSd(ctx, sd, caip)
		if err != nil {
			return nil, fmt.Errorf("failed to create attachable item: %w", err)
		}
		p := parameter.CreateImageParam{
			Height:           imgHeight,
			Width:            imgWidth,
			AttachableItemID: ai.AttachableItemID,
		}
		e, err := m.DB.CreateImageWithSd(ctx, sd, p)
		if err != nil {
			return nil, fmt.Errorf("failed to create image: %w", err)
		}
		images = append(images, e)
	}
	return images, nil
}

// CreateImageFromOuter 外部画像を作成する。
func (m *ManageImage) CreateImageFromOuter(
	ctx context.Context,
	url,
	alias string,
	size entity.Float,
	ownerID entity.UUID,
	mimeTypeID uuid.UUID,
	height, width entity.Float,
) (e entity.Image, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.Image{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rerr := m.DB.Rollback(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to rollback transaction: %w", rerr)
			}
		} else {
			if rerr := m.DB.Commit(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", rerr)
			}
		}
	}()
	if ownerID.Valid {
		// 画像の所有者が存在する場合は、画像の所有者が存在するか確認する。
		_, err := m.DB.FindMemberByIDWithSd(ctx, sd, ownerID.Bytes)
		if err != nil {
			return entity.Image{}, errhandle.NewModelNotFoundError("member")
		}
	}
	// 画像のMIMEタイプが存在するか確認する。
	_, err = m.DB.FindMimeTypeByIDWithSd(ctx, sd, mimeTypeID)
	if err != nil {
		return entity.Image{}, errhandle.NewModelNotFoundError("mime type")
	}
	caip := parameter.CreateAttachableItemParam{
		URL:        url,
		Alias:      alias,
		Size:       size,
		OwnerID:    ownerID,
		FromOuter:  true,
		MimeTypeID: mimeTypeID,
	}
	ai, err := m.DB.CreateAttachableItemWithSd(ctx, sd, caip)
	p := parameter.CreateImageParam{
		Height:           height,
		Width:            width,
		AttachableItemID: ai.AttachableItemID,
	}
	e, err = m.DB.CreateImageWithSd(ctx, sd, p)
	if err != nil {
		return entity.Image{}, fmt.Errorf("failed to create image: %w", err)
	}
	return e, nil
}

// CreateImagesFromOuter 外部画像を複数作成する。
func (m *ManageImage) CreateImagesFromOuter(
	ctx context.Context,
	ownerID entity.UUID,
	params []parameter.CreateImageFromOuterServiceParam,
) (es []entity.Image, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rerr := m.DB.Rollback(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to rollback transaction: %w", rerr)
			}
		} else {
			if rerr := m.DB.Commit(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", rerr)
			}
		}
	}()
	if ownerID.Valid {
		// 画像の所有者が存在する場合は、画像の所有者が存在するか確認する。
		_, err := m.DB.FindMemberByIDWithSd(ctx, sd, ownerID.Bytes)
		if err != nil {
			return nil, errhandle.NewModelNotFoundError("member")
		}
	}
	var images []entity.Image
	for _, p := range params {
		// 画像のMIMEタイプが存在するか確認する。
		_, err = m.DB.FindMimeTypeByIDWithSd(ctx, sd, p.MimeTypeID)
		if err != nil {
			return nil, errhandle.NewModelNotFoundError("mime type")
		}
		caip := parameter.CreateAttachableItemParam{
			URL:        p.URL,
			Alias:      p.Alias,
			Size:       p.Size,
			OwnerID:    ownerID,
			FromOuter:  true,
			MimeTypeID: p.MimeTypeID,
		}
		ai, err := m.DB.CreateAttachableItemWithSd(ctx, sd, caip)
		if err != nil {
			return nil, fmt.Errorf("failed to create attachable item: %w", err)
		}
		p := parameter.CreateImageParam{
			Height:           p.Height,
			Width:            p.Width,
			AttachableItemID: ai.AttachableItemID,
		}
		e, err := m.DB.CreateImageWithSd(ctx, sd, p)
		if err != nil {
			return nil, fmt.Errorf("failed to create image: %w", err)
		}
		images = append(images, e)
	}
	return images, nil
}

// DeleteImage 画像を削除する。
func (m *ManageImage) DeleteImage(ctx context.Context, id uuid.UUID) (c int64, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rerr := m.DB.Rollback(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to rollback transaction: %w", rerr)
			}
		} else {
			if rerr := m.DB.Commit(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", rerr)
			}
		}
	}()
	image, err := m.DB.FindImageWithAttachableItemWithSd(ctx, sd, id)
	if err != nil {
		return 0, fmt.Errorf("failed to find image: %w", err)
	}
	if !image.AttachableItem.FromOuter {
		key, err := m.Storage.GetKeyFromURL(ctx, image.AttachableItem.URL)
		if err != nil {
			return 0, fmt.Errorf("failed to get key from url: %w", err)
		}
		err = m.Storage.DeleteObjects(ctx, []string{key})
		if err != nil {
			return 0, fmt.Errorf("failed to delete object: %w", err)
		}
	}
	c, err = m.DB.DeleteImageWithSd(ctx, sd, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete image: %w", err)
	}
	_, err = m.DB.DeleteAttachableItemWithSd(ctx, sd, image.AttachableItem.AttachableItemID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete attachable item: %w", err)
	}
	return c, nil
}

func pluralDeleteImages(
	ctx context.Context,
	sd store.Sd,
	db store.Store,
	stg storage.Storage,
	ids []uuid.UUID,
	ownerID entity.UUID,
) (c int64, err error) {
	image, err := db.GetPluralImagesWithAttachableItemWithSd(
		ctx, sd, ids, parameter.ImageOrderMethodDefault, store.NumberedPaginationParam{})
	if err != nil {
		return 0, fmt.Errorf("failed to get images: %w", err)
	}
	var keys []string
	var attachableItemIDs []uuid.UUID
	for _, i := range image.Data {
		if !i.AttachableItem.OwnerID.Valid {
			continue
		}
		if !ownerID.Valid || i.AttachableItem.OwnerID.Bytes != ownerID.Bytes {
			return 0, errhandle.NewCommonError(response.NotFileOwner, nil)
		}
		if !i.AttachableItem.FromOuter {
			key, err := stg.GetKeyFromURL(ctx, i.AttachableItem.URL)
			if err != nil {
				return 0, fmt.Errorf("failed to get key from url: %w", err)
			}
			keys = append(keys, key)
		}
		attachableItemIDs = append(attachableItemIDs, i.AttachableItem.AttachableItemID)
	}
	if len(keys) == 0 {
		return 0, nil
	}
	err = stg.DeleteObjects(ctx, keys)
	if err != nil {
		return 0, fmt.Errorf("failed to delete objects: %w", err)
	}
	c, err = db.PluralDeleteImagesWithSd(ctx, sd, ids)
	if err != nil {
		return 0, fmt.Errorf("failed to delete policy: %w", err)
	}
	_, err = db.PluralDeleteAttachableItemsWithSd(ctx, sd, attachableItemIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to delete attachable items: %w", err)
	}
	return c, nil
}

// PluralDeleteImages 画像を複数削除する。
func (m *ManageImage) PluralDeleteImages(
	ctx context.Context, ids []uuid.UUID,
) (c int64, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rerr := m.DB.Rollback(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to rollback transaction: %w", rerr)
			}
		} else {
			if rerr := m.DB.Commit(ctx, sd); rerr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", rerr)
			}
		}
	}()
	return pluralDeleteImages(ctx, sd, m.DB, m.Storage, ids, entity.UUID{})
}

// GetImages 画像を取得する。
func (m *ManageImage) GetImages(
	ctx context.Context,
	order parameter.ImageOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.Image], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereImageParam{}
	switch pg {
	case parameter.NumberedPagination:
		np = store.NumberedPaginationParam{
			Valid:  true,
			Offset: entity.Int{Int64: int64(offset)},
			Limit:  entity.Int{Int64: int64(limit)},
		}
	case parameter.CursorPagination:
		cp = store.CursorPaginationParam{
			Valid:  true,
			Cursor: string(cursor),
			Limit:  entity.Int{Int64: int64(limit)},
		}
	case parameter.NonePagination:
	}
	r, err := m.DB.GetImages(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.Image]{}, fmt.Errorf("failed to get images: %w", err)
	}
	return r, nil
}

// GetImagesCount 画像の数を取得する。
func (m *ManageImage) GetImagesCount(
	ctx context.Context,
) (int64, error) {
	p := parameter.WhereImageParam{}
	c, err := m.DB.CountImages(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to get images count: %w", err)
	}
	return c, nil
}

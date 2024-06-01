package service

//nolint:revive
import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"strings"

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
) (e entity.ImageWithAttachableItem, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.ImageWithAttachableItem{}, fmt.Errorf("failed to begin transaction: %w", err)
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
			return entity.ImageWithAttachableItem{}, errhandle.NewModelNotFoundError("member")
		}
	}
	data, size, err := readAll(origin)
	if err != nil {
		return entity.ImageWithAttachableItem{}, err
	}
	mtype := mimetype.Detect(data)
	if err != nil {
		return entity.ImageWithAttachableItem{}, fmt.Errorf("failed to detect mimetype: %w", err)
	}
	// "; "が含まれている場合は、"; "より前の文字列を取得する。
	mtypeStr := mtype.String()
	if i := bytes.Index([]byte(mtypeStr), []byte("; ")); i != -1 {
		mtypeStr = mtypeStr[:i]
	}
	mime, err := m.DB.FindMimeTypeByKindWithSd(ctx, sd, mtypeStr)
	if err != nil {
		var merr *errhandle.ModelNotFoundError
		if errors.As(err, &merr) {
			mime, err = m.DB.FindMimeTypeByKeyWithSd(ctx, sd, string(MimeTypeKeyUnknown))
			if err != nil {
				return entity.ImageWithAttachableItem{}, errhandle.NewModelNotFoundError("mime type")
			}
		} else {
			return entity.ImageWithAttachableItem{}, fmt.Errorf("failed to find mime type: %w", err)
		}
	}
	if !strings.HasPrefix(mime.Kind, "image/") {
		return entity.ImageWithAttachableItem{}, errhandle.NewCommonError(response.NotImageFile, nil)
	}
	uid, err := uuid.NewRandom()
	if err != nil {
		return entity.ImageWithAttachableItem{}, fmt.Errorf("failed to generate uuid: %w", err)
	}
	extension := mtype.Extension()
	if i := bytes.LastIndex([]byte(alias), []byte(".")); i != -1 {
		extension = "." + alias[i+1:]
	}
	fname := fmt.Sprintf("%s%s", uid.String(), extension)
	url, err := m.Storage.UploadObject(ctx, bytes.NewReader(data), fname)
	if err != nil {
		return entity.ImageWithAttachableItem{}, errhandle.NewCommonError(response.FailedUpload, nil)
	}
	storageKeys = append(storageKeys, fname)
	fsize := entity.Float{
		Valid:   true,
		Float64: float64(size),
	}
	img, _, err := image.DecodeConfig(bytes.NewReader(data))
	var imgHeight, imgWidth entity.Float
	if err == nil {
		imgHeight = entity.Float{Valid: true, Float64: float64(img.Height)}
		imgWidth = entity.Float{Valid: true, Float64: float64(img.Width)}
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
		return entity.ImageWithAttachableItem{}, fmt.Errorf("failed to create attachable item: %w", err)
	}
	p := parameter.CreateImageParam{
		Height:           imgHeight,
		Width:            imgWidth,
		AttachableItemID: ai.AttachableItemID,
	}
	image, err := m.DB.CreateImageWithSd(ctx, sd, p)
	if err != nil {
		return entity.ImageWithAttachableItem{}, fmt.Errorf("failed to create image: %w", err)
	}
	ai.ImageID = entity.UUID{Bytes: image.ImageID, Valid: true}
	return entity.ImageWithAttachableItem{
		ImageID:        image.ImageID,
		Height:         image.Height,
		Width:          image.Width,
		AttachableItem: ai,
	}, nil
}

// CreateImages 画像を複数作成する。
func (m *ManageImage) CreateImages(
	ctx context.Context,
	ownerID entity.UUID,
	params []parameter.CreateImageServiceParam,
) (es []entity.ImageWithAttachableItem, err error) {
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
	var images []entity.ImageWithAttachableItem
	for _, p := range params {
		data, size, err := readAll(p.Origin)
		if err != nil {
			return nil, err
		}
		mtype := mimetype.Detect(data)
		// "; "が含まれている場合は、"; "より前の文字列を取得する。
		mtypeStr := mtype.String()
		if i := bytes.Index([]byte(mtypeStr), []byte("; ")); i != -1 {
			mtypeStr = mtypeStr[:i]
		}
		mime, err := m.DB.FindMimeTypeByKindWithSd(ctx, sd, mtypeStr)
		if err != nil {
			var merr *errhandle.ModelNotFoundError
			if errors.As(err, &merr) {
				mime, err = m.DB.FindMimeTypeByKeyWithSd(ctx, sd, string(MimeTypeKeyUnknown))
				if err != nil {
					return nil, errhandle.NewModelNotFoundError("mime type")
				}
			} else {
				return nil, fmt.Errorf("failed to find mime type: %w", err)
			}
		}
		if !strings.HasPrefix(mime.Kind, "image/") {
			return nil, errhandle.NewCommonError(response.NotImageFile, nil)
		}
		uid, err := uuid.NewRandom()
		if err != nil {
			return nil, fmt.Errorf("failed to generate uuid: %w", err)
		}
		extension := mtype.Extension()
		if i := bytes.LastIndex([]byte(p.Alias), []byte(".")); i != -1 {
			extension = "." + p.Alias[i+1:]
		}
		fname := fmt.Sprintf("%s%s", uid.String(), extension)
		url, err := m.Storage.UploadObject(ctx, bytes.NewReader(data), fname)
		if err != nil {
			return nil, errhandle.NewCommonError(response.FailedUpload, nil)
		}
		storageKeys = append(storageKeys, fname)

		fsize := entity.Float{
			Valid:   true,
			Float64: float64(size),
		}
		img, _, err := image.DecodeConfig(bytes.NewReader(data))
		var imgHeight, imgWidth entity.Float
		if err == nil {
			imgHeight = entity.Float{Valid: true, Float64: float64(img.Height)}
			imgWidth = entity.Float{Valid: true, Float64: float64(img.Width)}
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
		ai.ImageID = entity.UUID{Bytes: e.ImageID, Valid: true}
		images = append(images, entity.ImageWithAttachableItem{
			ImageID:        e.ImageID,
			Height:         e.Height,
			Width:          e.Width,
			AttachableItem: ai,
		})
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
) (e entity.ImageWithAttachableItem, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.ImageWithAttachableItem{}, fmt.Errorf("failed to begin transaction: %w", err)
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
			return entity.ImageWithAttachableItem{}, errhandle.NewModelNotFoundError("member")
		}
	}
	data, size, err := readAll(origin)
	if err != nil {
		return entity.ImageWithAttachableItem{}, err
	}
	mtype := mimetype.Detect(data)
	if err != nil {
		return entity.ImageWithAttachableItem{}, fmt.Errorf("failed to detect mimetype: %w", err)
	}
	// "; "が含まれている場合は、"; "より前の文字列を取得する。
	mtypeStr := mtype.String()
	if i := bytes.Index([]byte(mtypeStr), []byte("; ")); i != -1 {
		mtypeStr = mtypeStr[:i]
	}
	mime, err := m.DB.FindMimeTypeByKindWithSd(ctx, sd, mtypeStr)
	if err != nil {
		var merr *errhandle.ModelNotFoundError
		if errors.As(err, &merr) {
			mime, err = m.DB.FindMimeTypeByKeyWithSd(ctx, sd, string(MimeTypeKeyUnknown))
			if err != nil {
				return entity.ImageWithAttachableItem{}, errhandle.NewModelNotFoundError("mime type")
			}
		} else {
			return entity.ImageWithAttachableItem{}, fmt.Errorf("failed to find mime type: %w", err)
		}
	}
	if !strings.HasPrefix(mime.Kind, "image/") {
		return entity.ImageWithAttachableItem{}, errhandle.NewCommonError(response.NotImageFile, nil)
	}
	exist, err := m.Storage.ExistsObject(ctx, filename)
	if err != nil {
		return entity.ImageWithAttachableItem{}, fmt.Errorf("failed to check object existence: %w", err)
	}
	if exist {
		return entity.ImageWithAttachableItem{}, errhandle.NewCommonError(response.ConflictStorageKey, nil)
	}
	url, err := m.Storage.UploadObject(ctx, bytes.NewReader(data), filename)
	if err != nil {
		return entity.ImageWithAttachableItem{}, errhandle.NewCommonError(response.FailedUpload, nil)
	}
	storageKeys = append(storageKeys, filename)
	if err != nil {
		return entity.ImageWithAttachableItem{}, fmt.Errorf("failed to copy file: %w", err)
	}
	fsize := entity.Float{
		Valid:   true,
		Float64: float64(size),
	}
	img, _, err := image.DecodeConfig(bytes.NewReader(data))
	var imgHeight, imgWidth entity.Float
	if err == nil {
		imgHeight = entity.Float{Valid: true, Float64: float64(img.Height)}
		imgWidth = entity.Float{Valid: true, Float64: float64(img.Width)}
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
	im, err := m.DB.CreateImageWithSd(ctx, sd, p)
	if err != nil {
		return entity.ImageWithAttachableItem{}, fmt.Errorf("failed to create image: %w", err)
	}
	ai.ImageID = entity.UUID{Bytes: im.ImageID, Valid: true}
	return entity.ImageWithAttachableItem{
		ImageID:        im.ImageID,
		Height:         im.Height,
		Width:          im.Width,
		AttachableItem: ai,
	}, nil
}

// CreateImagesSpecifyFilename 画像を複数作成する。
func (m *ManageImage) CreateImagesSpecifyFilename(
	ctx context.Context,
	ownerID entity.UUID,
	params []parameter.CreateImageSpecifyFilenameServiceParam,
) (es []entity.ImageWithAttachableItem, err error) {
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
	var images []entity.ImageWithAttachableItem
	for _, p := range params {
		data, size, err := readAll(p.Origin)
		if err != nil {
			return nil, err
		}
		mtype := mimetype.Detect(data)
		// "; "が含まれている場合は、"; "より前の文字列を取得する。
		mtypeStr := mtype.String()
		if i := bytes.Index([]byte(mtypeStr), []byte("; ")); i != -1 {
			mtypeStr = mtypeStr[:i]
		}
		mime, err := m.DB.FindMimeTypeByKindWithSd(ctx, sd, mtypeStr)
		if err != nil {
			var merr *errhandle.ModelNotFoundError
			if errors.As(err, &merr) {
				mime, err = m.DB.FindMimeTypeByKeyWithSd(ctx, sd, string(MimeTypeKeyUnknown))
				if err != nil {
					return nil, errhandle.NewModelNotFoundError("mime type")
				}
			} else {
				return nil, fmt.Errorf("failed to find mime type: %w", err)
			}
		}
		if !strings.HasPrefix(mime.Kind, "image/") {
			return nil, errhandle.NewCommonError(response.NotImageFile, nil)
		}
		exist, err := m.Storage.ExistsObject(ctx, p.Filename)
		if err != nil {
			return nil, fmt.Errorf("failed to check object existence: %w", err)
		}
		if exist {
			return nil, errhandle.NewCommonError(response.ConflictStorageKey, nil)
		}
		url, err := m.Storage.UploadObject(ctx, bytes.NewReader(data), p.Filename)
		if err != nil {
			return nil, errhandle.NewCommonError(response.FailedUpload, nil)
		}
		storageKeys = append(storageKeys, p.Filename)
		fsize := entity.Float{
			Valid:   true,
			Float64: float64(size),
		}
		img, _, err := image.DecodeConfig(bytes.NewReader(data))
		var imgHeight, imgWidth entity.Float
		if err == nil {
			imgHeight = entity.Float{Valid: true, Float64: float64(img.Height)}
			imgWidth = entity.Float{Valid: true, Float64: float64(img.Width)}
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
		ai.ImageID = entity.UUID{Bytes: e.ImageID, Valid: true}
		images = append(images, entity.ImageWithAttachableItem{
			ImageID:        e.ImageID,
			Height:         e.Height,
			Width:          e.Width,
			AttachableItem: ai,
		})
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
) (e entity.ImageWithAttachableItem, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.ImageWithAttachableItem{}, fmt.Errorf("failed to begin transaction: %w", err)
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
			return entity.ImageWithAttachableItem{}, errhandle.NewModelNotFoundError("member")
		}
	}
	// 画像のMIMEタイプが存在するか確認する。
	mime, err := m.DB.FindMimeTypeByIDWithSd(ctx, sd, mimeTypeID)
	if err != nil {
		return entity.ImageWithAttachableItem{}, errhandle.NewModelNotFoundError("mime type")
	}
	if !strings.HasPrefix(mime.Kind, "image/") {
		return entity.ImageWithAttachableItem{}, errhandle.NewCommonError(response.NotImageFile, nil)
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
	im, err := m.DB.CreateImageWithSd(ctx, sd, p)
	if err != nil {
		return entity.ImageWithAttachableItem{}, fmt.Errorf("failed to create image: %w", err)
	}
	ai.ImageID = entity.UUID{Bytes: im.ImageID, Valid: true}
	return entity.ImageWithAttachableItem{
		ImageID:        im.ImageID,
		Height:         im.Height,
		Width:          im.Width,
		AttachableItem: ai,
	}, nil
}

// CreateImagesFromOuter 外部画像を複数作成する。
func (m *ManageImage) CreateImagesFromOuter(
	ctx context.Context,
	ownerID entity.UUID,
	params []parameter.CreateImageFromOuterServiceParam,
) (es []entity.ImageWithAttachableItem, err error) {
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
	var images []entity.ImageWithAttachableItem
	for _, p := range params {
		// 画像のMIMEタイプが存在するか確認する。
		mime, err := m.DB.FindMimeTypeByIDWithSd(ctx, sd, p.MimeTypeID)
		if err != nil {
			return nil, errhandle.NewModelNotFoundError("mime type")
		}
		if !strings.HasPrefix(mime.Kind, "image/") {
			return nil, errhandle.NewCommonError(response.NotImageFile, nil)
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
		ai.ImageID = entity.UUID{Bytes: e.ImageID, Valid: true}
		images = append(images, entity.ImageWithAttachableItem{
			ImageID:        e.ImageID,
			Height:         e.Height,
			Width:          e.Width,
			AttachableItem: ai,
		})
	}
	return images, nil
}

// DeleteImage 画像を削除する。
func (m *ManageImage) DeleteImage(ctx context.Context, id uuid.UUID, ownerID entity.UUID) (c int64, err error) {
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
	if !image.AttachableItem.OwnerID.Valid {
		return 0, nil
	}
	if !ownerID.Valid || image.AttachableItem.OwnerID.Bytes != ownerID.Bytes {
		return 0, errhandle.NewCommonError(response.NotFileOwner, nil)
	}
	c, err = m.DB.DeleteImageWithSd(ctx, sd, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete image: %w", err)
	}
	_, err = m.DB.DeleteAttachableItemWithSd(ctx, sd, image.AttachableItem.AttachableItemID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete attachable item: %w", err)
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
	return c, nil
}

func pluralDeleteImages(
	ctx context.Context,
	sd store.Sd,
	db store.Store,
	stg storage.Storage,
	ids []uuid.UUID,
	ownerID entity.UUID,
	force bool,
) (c int64, err error) {
	image, err := db.GetPluralImagesWithAttachableItemWithSd(
		ctx, sd, ids, parameter.ImageOrderMethodDefault, store.NumberedPaginationParam{})
	if err != nil {
		return 0, fmt.Errorf("failed to get images: %w", err)
	}
	if len(image.Data) != len(ids) {
		return 0, errhandle.NewModelNotFoundError(ImageTargetImages)
	}
	var keys []string
	var attachableItemIDs []uuid.UUID
	for _, i := range image.Data {
		if !i.AttachableItem.OwnerID.Valid {
			if force {
				continue
			}
			return 0, errhandle.NewCommonError(response.CannotDeleteSystemFile, nil)
		}
		if !force && (!ownerID.Valid || i.AttachableItem.OwnerID.Bytes != ownerID.Bytes) {
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
	c, err = db.PluralDeleteImagesWithSd(ctx, sd, ids)
	if err != nil {
		return 0, fmt.Errorf("failed to delete policy: %w", err)
	}
	_, err = db.PluralDeleteAttachableItemsWithSd(ctx, sd, attachableItemIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to delete attachable items: %w", err)
	}
	err = stg.DeleteObjects(ctx, keys)
	if err != nil {
		return 0, fmt.Errorf("failed to delete objects: %w", err)
	}
	return c, nil
}

// PluralDeleteImages 画像を複数削除する。
func (m *ManageImage) PluralDeleteImages(
	ctx context.Context, ids []uuid.UUID, ownerID entity.UUID,
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
	return pluralDeleteImages(ctx, sd, m.DB, m.Storage, ids, ownerID, false)
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
			Offset: entity.Int{Int64: int64(offset), Valid: true},
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.CursorPagination:
		cp = store.CursorPaginationParam{
			Valid:  true,
			Cursor: string(cursor),
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
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

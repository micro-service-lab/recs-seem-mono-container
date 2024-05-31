package service

import (
	"context"
	"fmt"
	"io"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/storage"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

const bufSize = 1024

// ManageFile ファイル管理サービス。
type ManageFile struct {
	DB      store.Store
	Storage storage.Storage
}

// CreateFile ファイルを作成する。
func (m *ManageFile) CreateFile(
	ctx context.Context,
	origin io.Reader,
	alias string,
	ownerID entity.UUID,
) (e entity.File, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.File{}, fmt.Errorf("failed to begin transaction: %w", err)
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
		// ファイルの所有者が存在する場合は、ファイルの所有者が存在するか確認する。
		_, err := m.DB.FindMemberByIDWithSd(ctx, sd, ownerID.Bytes)
		if err != nil {
			return entity.File{}, errhandle.NewModelNotFoundError("member")
		}
	}
	mtype, err := mimetype.DetectReader(origin)
	if err != nil {
		return entity.File{}, fmt.Errorf("failed to detect mimetype: %w", err)
	}
	mime, err := m.DB.FindMimeTypeByKindWithSd(ctx, sd, mtype.String())
	if err != nil {
		return entity.File{}, errhandle.NewModelNotFoundError("mime type")
	}
	uid, err := uuid.NewRandom()
	if err != nil {
		return entity.File{}, fmt.Errorf("failed to generate uuid: %w", err)
	}
	fname := fmt.Sprintf("%s.%s", uid.String(), mtype.Extension())
	url, err := m.Storage.UploadObject(ctx, origin, fname)
	if err != nil {
		return entity.File{}, errhandle.NewCommonError(response.FailedUpload, nil)
	}
	storageKeys = append(storageKeys, fname)
	var size int64
	for {
		buf := make([]byte, bufSize)
		n, err := origin.Read(buf)
		if err != nil && err != io.EOF {
			return entity.File{}, fmt.Errorf("failed to read file: %w", err)
		}
		if n == 0 {
			break
		}
		size += int64(n)
	}
	fsize := entity.Float{
		Valid:   true,
		Float64: float64(size),
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
		return entity.File{}, fmt.Errorf("failed to create attachable item: %w", err)
	}
	p := parameter.CreateFileParam{
		AttachableItemID: ai.AttachableItemID,
	}
	e, err = m.DB.CreateFileWithSd(ctx, sd, p)
	if err != nil {
		return entity.File{}, fmt.Errorf("failed to create file: %w", err)
	}
	return e, nil
}

// CreateFiles ファイルを複数作成する。
func (m *ManageFile) CreateFiles(
	ctx context.Context,
	ownerID entity.UUID,
	params []parameter.CreateFileServiceParam,
) (es []entity.File, err error) {
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
		// ファイルの所有者が存在する場合は、ファイルの所有者が存在するか確認する。
		_, err := m.DB.FindMemberByIDWithSd(ctx, sd, ownerID.Bytes)
		if err != nil {
			return nil, errhandle.NewModelNotFoundError("member")
		}
	}
	var files []entity.File
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
		var size int64
		for {
			buf := make([]byte, bufSize)
			n, err := p.Origin.Read(buf)
			if err != nil && err != io.EOF {
				return nil, fmt.Errorf("failed to read file: %w", err)
			}
			if n == 0 {
				break
			}
			size += int64(n)
		}
		fsize := entity.Float{
			Valid:   true,
			Float64: float64(size),
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
		p := parameter.CreateFileParam{
			AttachableItemID: ai.AttachableItemID,
		}
		e, err := m.DB.CreateFileWithSd(ctx, sd, p)
		if err != nil {
			return nil, fmt.Errorf("failed to create file: %w", err)
		}
		files = append(files, e)
	}
	return files, nil
}

// CreateFileSpecifyFilename ファイルを作成する。
func (m *ManageFile) CreateFileSpecifyFilename(
	ctx context.Context,
	origin io.Reader,
	alias string,
	ownerID entity.UUID,
	filename string,
) (e entity.File, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.File{}, fmt.Errorf("failed to begin transaction: %w", err)
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
		// ファイルの所有者が存在する場合は、ファイルの所有者が存在するか確認する。
		_, err := m.DB.FindMemberByIDWithSd(ctx, sd, ownerID.Bytes)
		if err != nil {
			return entity.File{}, errhandle.NewModelNotFoundError("member")
		}
	}
	mtype, err := mimetype.DetectReader(origin)
	if err != nil {
		return entity.File{}, fmt.Errorf("failed to detect mimetype: %w", err)
	}
	mime, err := m.DB.FindMimeTypeByKindWithSd(ctx, sd, mtype.String())
	if err != nil {
		return entity.File{}, errhandle.NewModelNotFoundError("mime type")
	}
	exist, err := m.Storage.ExistsObject(ctx, filename)
	if err != nil {
		return entity.File{}, fmt.Errorf("failed to check object existence: %w", err)
	}
	if exist {
		return entity.File{}, errhandle.NewCommonError(response.ConflictStorageKey, nil)
	}
	url, err := m.Storage.UploadObject(ctx, origin, filename)
	if err != nil {
		return entity.File{}, errhandle.NewCommonError(response.FailedUpload, nil)
	}
	storageKeys = append(storageKeys, filename)
	var size int64
	for {
		buf := make([]byte, bufSize)
		n, err := origin.Read(buf)
		if err != nil && err != io.EOF {
			return entity.File{}, fmt.Errorf("failed to read file: %w", err)
		}
		if n == 0 {
			break
		}
		size += int64(n)
	}
	fsize := entity.Float{
		Valid:   true,
		Float64: float64(size),
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
	p := parameter.CreateFileParam{
		AttachableItemID: ai.AttachableItemID,
	}
	e, err = m.DB.CreateFileWithSd(ctx, sd, p)
	if err != nil {
		return entity.File{}, fmt.Errorf("failed to create file: %w", err)
	}
	return e, nil
}

// CreateFilesSpecifyFilename ファイルを複数作成する。
func (m *ManageFile) CreateFilesSpecifyFilename(
	ctx context.Context,
	ownerID entity.UUID,
	params []parameter.CreateFileSpecifyFilenameServiceParam,
) (es []entity.File, err error) {
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
		// ファイルの所有者が存在する場合は、ファイルの所有者が存在するか確認する。
		_, err := m.DB.FindMemberByIDWithSd(ctx, sd, ownerID.Bytes)
		if err != nil {
			return nil, errhandle.NewModelNotFoundError("member")
		}
	}
	var files []entity.File
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
		var size int64
		for {
			buf := make([]byte, bufSize)
			n, err := p.Origin.Read(buf)
			if err != nil && err != io.EOF {
				return nil, fmt.Errorf("failed to read file: %w", err)
			}
			if n == 0 {
				break
			}
			size += int64(n)
		}
		fsize := entity.Float{
			Valid:   true,
			Float64: float64(size),
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
		p := parameter.CreateFileParam{
			AttachableItemID: ai.AttachableItemID,
		}
		e, err := m.DB.CreateFileWithSd(ctx, sd, p)
		if err != nil {
			return nil, fmt.Errorf("failed to create file: %w", err)
		}
		files = append(files, e)
	}
	return files, nil
}

// CreateFileFromOuter 外部ファイルを作成する。
func (m *ManageFile) CreateFileFromOuter(
	ctx context.Context,
	url,
	alias string,
	size entity.Float,
	ownerID entity.UUID,
	mimeTypeID uuid.UUID,
) (e entity.File, err error) {
	sd, err := m.DB.Begin(ctx)
	if err != nil {
		return entity.File{}, fmt.Errorf("failed to begin transaction: %w", err)
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
		// ファイルの所有者が存在する場合は、ファイルの所有者が存在するか確認する。
		_, err := m.DB.FindMemberByIDWithSd(ctx, sd, ownerID.Bytes)
		if err != nil {
			return entity.File{}, errhandle.NewModelNotFoundError("member")
		}
	}
	// ファイルのMIMEタイプが存在するか確認する。
	_, err = m.DB.FindMimeTypeByIDWithSd(ctx, sd, mimeTypeID)
	if err != nil {
		return entity.File{}, errhandle.NewModelNotFoundError("mime type")
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
	p := parameter.CreateFileParam{
		AttachableItemID: ai.AttachableItemID,
	}
	e, err = m.DB.CreateFileWithSd(ctx, sd, p)
	if err != nil {
		return entity.File{}, fmt.Errorf("failed to create file: %w", err)
	}
	return e, nil
}

// CreateFilesFromOuter 外部ファイルを複数作成する。
func (m *ManageFile) CreateFilesFromOuter(
	ctx context.Context,
	ownerID entity.UUID,
	params []parameter.CreateFileFromOuterServiceParam,
) (es []entity.File, err error) {
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
		// ファイルの所有者が存在する場合は、ファイルの所有者が存在するか確認する。
		_, err := m.DB.FindMemberByIDWithSd(ctx, sd, ownerID.Bytes)
		if err != nil {
			return nil, errhandle.NewModelNotFoundError("member")
		}
	}
	var files []entity.File
	for _, p := range params {
		// ファイルのMIMEタイプが存在するか確認する。
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
		p := parameter.CreateFileParam{
			AttachableItemID: ai.AttachableItemID,
		}
		e, err := m.DB.CreateFileWithSd(ctx, sd, p)
		if err != nil {
			return nil, fmt.Errorf("failed to create file: %w", err)
		}
		files = append(files, e)
	}
	return files, nil
}

// DeleteFile ファイルを削除する。
func (m *ManageFile) DeleteFile(ctx context.Context, id uuid.UUID) (c int64, err error) {
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
	file, err := m.DB.FindFileWithAttachableItemWithSd(ctx, sd, id)
	if err != nil {
		return 0, fmt.Errorf("failed to find file: %w", err)
	}
	if !file.AttachableItem.OwnerID.Valid {
		return 0, nil
	}
	if !file.AttachableItem.FromOuter {
		key, err := m.Storage.GetKeyFromURL(ctx, file.AttachableItem.URL)
		if err != nil {
			return 0, fmt.Errorf("failed to get key from url: %w", err)
		}
		err = m.Storage.DeleteObjects(ctx, []string{key})
		if err != nil {
			return 0, fmt.Errorf("failed to delete object: %w", err)
		}
	}
	c, err = m.DB.DeleteFileWithSd(ctx, sd, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete file: %w", err)
	}
	_, err = m.DB.DeleteAttachableItemWithSd(ctx, sd, file.AttachableItem.AttachableItemID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete attachable item: %w", err)
	}
	return c, nil
}

func pluralDeleteFiles(
	ctx context.Context,
	sd store.Sd,
	db store.Store,
	stg storage.Storage,
	ids []uuid.UUID,
	ownerID entity.UUID,
	force bool,
) (c int64, err error) {
	file, err := db.GetPluralFilesWithAttachableItemWithSd(
		ctx, sd, ids, parameter.FileOrderMethodDefault, store.NumberedPaginationParam{})
	if err != nil {
		return 0, fmt.Errorf("failed to get files: %w", err)
	}
	var keys []string
	var attachableItemIDs []uuid.UUID
	for _, i := range file.Data {
		if !i.AttachableItem.OwnerID.Valid {
			continue
		}
		if !force || (!ownerID.Valid || i.AttachableItem.OwnerID.Bytes != ownerID.Bytes) {
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
	c, err = db.PluralDeleteFilesWithSd(ctx, sd, ids)
	if err != nil {
		return 0, fmt.Errorf("failed to delete policy: %w", err)
	}
	_, err = db.PluralDeleteAttachableItemsWithSd(ctx, sd, attachableItemIDs)
	if err != nil {
		return 0, fmt.Errorf("failed to delete attachable items: %w", err)
	}
	return c, nil
}

// PluralDeleteFiles ファイルを複数削除する。
func (m *ManageFile) PluralDeleteFiles(
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
	return pluralDeleteFiles(ctx, sd, m.DB, m.Storage, ids, entity.UUID{}, false)
}

// GetFiles ファイルを取得する。
func (m *ManageFile) GetFiles(
	ctx context.Context,
	order parameter.FileOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.File], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereFileParam{}
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
	r, err := m.DB.GetFiles(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.File]{}, fmt.Errorf("failed to get files: %w", err)
	}
	return r, nil
}

// GetFilesCount ファイルの数を取得する。
func (m *ManageFile) GetFilesCount(
	ctx context.Context,
) (int64, error) {
	p := parameter.WhereFileParam{}
	c, err := m.DB.CountFiles(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to get files count: %w", err)
	}
	return c, nil
}

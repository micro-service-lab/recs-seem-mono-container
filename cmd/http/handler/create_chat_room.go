package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/i18n"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/lang"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/validation"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
)

// CreateChatRoom is a handler for creating ChatRoom.
type CreateChatRoom struct {
	Service    service.ManagerInterface
	Validator  validation.Validator
	Translator i18n.Translation
}

// CreateChatRoomRequest is a request for CreateChatRoom.
type CreateChatRoomRequest struct {
	Name         string      `json:"name" validate:"required,max=255" ja:"名前" en:"Name"`
	CoverImageID uuid.UUID   `json:"cover_image_id" validate:"" ja:"カバー画像ID" en:"CoverImageID"`
	MemberIDS    []uuid.UUID `json:"member_ids" validate:"required,unique" ja:"メンバーID" en:"MemberIDs"`
}

func (h *CreateChatRoom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser := auth.FromContext(ctx)
	var err error
	var chatRoomReq CreateChatRoomRequest
	if err = json.NewDecoder(r.Body).Decode(&chatRoomReq); err == nil || errors.Is(err, io.EOF) {
		if errors.Is(err, io.EOF) {
			chatRoomReq = CreateChatRoomRequest{}
		}
		err = h.Validator.ValidateWithLocale(ctx, &chatRoomReq, lang.GetLocale(r.Context()))
	} else {
		err = errhandle.NewJSONFormatError()
	}
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
	var chatRoom entity.ChatRoom
	if chatRoom, err = h.Service.CreateChatRoom(
		ctx,
		chatRoomReq.Name,
		entity.UUID{
			Valid: chatRoomReq.CoverImageID != uuid.Nil,
			Bytes: chatRoomReq.CoverImageID,
		},
		authUser.MemberID,
		chatRoomReq.MemberIDS,
	); err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			switch e.Target() {
			case service.ChatRoomTargetCoverImages:
				coverImageStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "CoverImageID", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "CoverImageID",
						Other: "CoverImageID",
					},
				})
				msgStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "ModelNotExists", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "ModelNotExists",
						Other: "{{.ID}} not found",
					},
					TemplateData: map[string]any{
						"ID": coverImageStr,
					},
				})
				ve := errhandle.NewValidationError(nil)
				ve.Add("cover_image_id", msgStr)
				err = ve
			case service.ChatRoomTargetOwner:
				e.SetTarget("owner")
				err = e
			case service.ChatRoomTargetMembers:
				memberStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "MemberIDs", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "MemberIDs",
						Other: "MemberIDs",
					},
				})
				msgStr := h.Translator.TranslateWithOpts(
					lang.GetLocaleForTranslation(ctx), "PluralModelNotExists", i18n.Options{
						DefaultMessage: &i18n.Message{
							ID:    "PluralModelNotExists",
							Other: "{{.ID}} not found",
						},
						TemplateData: map[string]any{
							"ID":        memberStr,
							"ValueType": "ID",
						},
					})
				ve := errhandle.NewValidationError(nil)
				ve.Add("member_ids", msgStr)
				err = ve
			}
		}
		var ce errhandle.CommonError
		if errors.As(err, &ce) {
			if ce.Code.Code == response.CannotAddSelfToChatRoom.Code {
				msgStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "CannotAddSelfToMembers", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "CannotAddSelfToMembers",
						Other: "Cannot add self to members",
					},
				})
				ve := errhandle.NewValidationError(nil)
				ve.Add("member_ids", msgStr)
				err = ve
			}
		}
	} else {
		err = response.JSONResponseWriter(ctx, w, response.Success, chatRoom, nil)
	}
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
}

package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
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

// CreatePrivateMessage is a handler for creating message on private room.
type CreatePrivateMessage struct {
	Service    service.ManagerInterface
	Validator  validation.Validator
	Translator i18n.Translation
}

// CreatePrivateMessageRequest is a request for CreatePrivateMessage.
type CreatePrivateMessageRequest struct {
	Content           string      `json:"content" validate:"required,max=255" ja:"コンテンツ" en:"Content"`
	AttachableItemIDS []uuid.UUID `json:"attachable_item_ids" validate:"unique" ja:"添付可能アイテムID" en:"AttachableItemIDs"`
}

func (h *CreatePrivateMessage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser := auth.FromContext(ctx)
	receiverID := uuid.MustParse(chi.URLParam(r, "member_id"))
	var err error
	var messageReq CreatePrivateMessageRequest
	if err = json.NewDecoder(r.Body).Decode(&messageReq); err == nil || errors.Is(err, io.EOF) {
		if errors.Is(err, io.EOF) {
			messageReq = CreatePrivateMessageRequest{}
		}
		err = h.Validator.ValidateWithLocale(ctx, &messageReq, lang.GetLocale(r.Context()))
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
	var message entity.Message
	if message, err = h.Service.CreateMessageOnPrivateRoom(
		ctx,
		authUser.MemberID,
		receiverID,
		messageReq.Content,
		messageReq.AttachableItemIDS,
	); err != nil {
		var ce errhandle.CommonError
		if errors.As(err, &ce) {
			if ce.Code.Code == response.CannotAttachSystemFile.Code {
				attachableItemStr := h.Translator.TranslateWithOpts(
					lang.GetLocaleForTranslation(ctx), "AttachableItemIDs", i18n.Options{
						DefaultMessage: &i18n.Message{
							ID:    "AttachableItemIDs",
							Other: "AttachableItemIDs",
						},
					})
				msgStr := h.Translator.TranslateWithOpts(
					lang.GetLocaleForTranslation(ctx), "ContainSystemFile", i18n.Options{
						DefaultMessage: &i18n.Message{
							ID:    "ContainSystemFile",
							Other: "{{.Target}} contains system files",
						},
						TemplateData: map[string]any{
							"Target": attachableItemStr,
						},
					})
				ve := errhandle.NewValidationError(nil)
				ve.Add("attachable_item_ids", msgStr)
				err = ve
			} else if ce.Code.Code == response.NotFileOwner.Code {
				attachableItemStr := h.Translator.TranslateWithOpts(
					lang.GetLocaleForTranslation(ctx), "AttachableItemIDs", i18n.Options{
						DefaultMessage: &i18n.Message{
							ID:    "AttachableItemIDs",
							Other: "AttachableItemIDs",
						},
					})
				msgStr := h.Translator.TranslateWithOpts(
					lang.GetLocaleForTranslation(ctx), "ContainNotOwnerFile", i18n.Options{
						DefaultMessage: &i18n.Message{
							ID:    "ContainNotOwnerFile",
							Other: "{{.Target}} contains files that are not owned by the user",
						},
						TemplateData: map[string]any{
							"Target": attachableItemStr,
						},
					})
				ve := errhandle.NewValidationError(nil)
				ve.Add("attachable_item_ids", msgStr)
				err = ve
			}
		}
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			switch e.Target() {
			case service.MessageTargetAttachments:
				attachableItemStr := h.Translator.TranslateWithOpts(
					lang.GetLocaleForTranslation(ctx), "AttachableItemIDs", i18n.Options{
						DefaultMessage: &i18n.Message{
							ID:    "AttachableItemIDs",
							Other: "AttachableItemIDs",
						},
					})
				msgStr := h.Translator.TranslateWithOpts(
					lang.GetLocaleForTranslation(ctx), "PluralModelNotExists", i18n.Options{
						DefaultMessage: &i18n.Message{
							ID:    "PluralModelNotExists",
							Other: "{{.ID}} not found",
						},
						TemplateData: map[string]any{
							"ID":        attachableItemStr,
							"ValueType": "ID",
						},
					})
				ve := errhandle.NewValidationError(nil)
				ve.Add("attachable_item_ids", msgStr)
				err = ve
			}
		}
	} else {
		err = response.JSONResponseWriter(ctx, w, response.Success, message, nil)
	}
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
}

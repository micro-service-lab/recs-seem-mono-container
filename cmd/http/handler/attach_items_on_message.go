package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/errhandle"
	"github.com/micro-service-lab/recs-seem-mono-container/app/i18n"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/lang"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/validation"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
)

// AttachItemsOnMessage is a handler for attach items on message.
type AttachItemsOnMessage struct {
	Service    service.ManagerInterface
	Validator  validation.Validator
	Translator i18n.Translation
}

// AttachItemsOnMessageRequest is a request for AttachItemsOnMessage.
type AttachItemsOnMessageRequest struct {
	AttachableItemIDS []uuid.UUID `json:"attachable_item_ids" validate:"required,unique,min=1" ja:"添付可能アイテムID" en:"AttachableItemIDs"`
}

func (h *AttachItemsOnMessage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser := auth.FromContext(ctx)
	chatRoomID := uuid.MustParse(chi.URLParam(r, "chat_room_id"))
	messageID := uuid.MustParse(chi.URLParam(r, "message_id"))
	var err error
	var req AttachItemsOnMessageRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err == nil || errors.Is(err, io.EOF) {
		if errors.Is(err, io.EOF) {
			req = AttachItemsOnMessageRequest{}
		}
		err = h.Validator.ValidateWithLocale(ctx, &req, lang.GetLocale(r.Context()))
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
	if _, err = h.Service.AttachItemsOnMessage(
		ctx, chatRoomID, messageID, authUser.MemberID, req.AttachableItemIDS); err != nil {
		var de errhandle.ModelDuplicatedError
		if errors.As(err, &de) {
			switch de.Target() {
			case service.MessageTargetAttacheMessage:
				attachableItemStr := h.Translator.TranslateWithOpts(
					lang.GetLocaleForTranslation(ctx), "AttachableItemIDs", i18n.Options{
						DefaultMessage: &i18n.Message{
							ID:    "AttachableItemIDs",
							Other: "AttachableItemIDs",
						},
					})
				messageStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "Message", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "Message",
						Other: "Message",
					},
				})
				msgStr := h.Translator.TranslateWithOpts(
					lang.GetLocaleForTranslation(ctx), "PluralAlreadyAssociated", i18n.Options{
						DefaultMessage: &i18n.Message{
							ID:    "PluralAlreadyAssociated",
							Other: "{{.ID}} not found",
						},
						TemplateData: map[string]any{
							"ID":         attachableItemStr,
							"ValueType":  "ID",
							"Associated": messageStr,
						},
					})
				ve := errhandle.NewValidationError(nil)
				ve.Add("attachable_item_ids", msgStr)
				err = ve
			}
		}
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
			case service.MessageTargetMessages:
				e.SetTarget("message")
				err = e
			}
		}
	} else {
		err = response.JSONResponseWriter(ctx, w, response.Success, nil, nil)
	}

	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
}

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

// RemoveMembersOnChatRoom is a handler for remove members on chat room.
type RemoveMembersOnChatRoom struct {
	Service    service.ManagerInterface
	Validator  validation.Validator
	Translator i18n.Translation
}

// RemoveMembersOnChatRoomRequest is a request for RemoveMembersOnChatRoom.
type RemoveMembersOnChatRoomRequest struct {
	MemberIDS []uuid.UUID `json:"member_ids" validate:"required,unique,min=1" ja:"メンバーID" en:"MemberIDs"`
}

func (h *RemoveMembersOnChatRoom) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser := auth.FromContext(ctx)
	id := uuid.MustParse(chi.URLParam(r, "chat_room_id"))
	var err error
	var req RemoveMembersOnChatRoomRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err == nil || errors.Is(err, io.EOF) {
		if errors.Is(err, io.EOF) {
			req = RemoveMembersOnChatRoomRequest{}
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
	if _, err = h.Service.RemoveMembersFromChatRoom(ctx, id, authUser.MemberID, req.MemberIDS); err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			switch e.Target() {
			case service.ChatRoomBelongingTargetOwner:
				e.SetTarget("owner")
				err = e
			case service.ChatRoomBelongingTargetChatRoom:
				e.SetTarget("chat_room")
				err = e
			case service.ChatRoomBelongingTargetMembers:
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
			case service.ChatRoomBelongingTargetChatRoomBelongings:
				memberStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "MemberIDs", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "MemberIDs",
						Other: "MemberIDs",
					},
				})
				msgStr := h.Translator.TranslateWithOpts(
					lang.GetLocaleForTranslation(ctx), "PluralNotAssociated", i18n.Options{
						DefaultMessage: &i18n.Message{
							ID:    "PluralNotAssociated",
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

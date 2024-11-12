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

// RemoveMembersOnOrganization is a handler for remove members on organization.
type RemoveMembersOnOrganization struct {
	Service    service.ManagerInterface
	Validator  validation.Validator
	Translator i18n.Translation
}

// RemoveMembersOnOrganizationRequest is a request for RemoveMembersOnOrganization.
type RemoveMembersOnOrganizationRequest struct {
	MemberIDS []uuid.UUID `json:"member_ids" validate:"required,unique,min=1" ja:"メンバーID" en:"MemberIDs"`
}

func (h *RemoveMembersOnOrganization) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser := auth.FromContext(ctx)
	id := uuid.MustParse(chi.URLParam(r, "organization_id"))
	var err error
	var req RemoveMembersOnOrganizationRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err == nil || errors.Is(err, io.EOF) {
		if errors.Is(err, io.EOF) {
			req = RemoveMembersOnOrganizationRequest{}
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
	if _, err = h.Service.RemoveMembersFromOrganization(ctx, id, authUser.MemberID, req.MemberIDS); err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			switch e.Target() {
			case service.OrganizationBelongingTargetOwner:
				e.SetTarget("owner")
				err = e
			case service.OrganizationBelongingTargetOrganization:
				e.SetTarget("organization")
				err = e
			case service.ChatRoomBelongingTargetChatRoom:
				e.SetTarget("chat_room")
				err = e
			case service.OrganizationBelongingTargetMembers:
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
			case service.OrganizationBelongingTargetOrganizationBelongings:
				memberStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "MemberIDs", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "MemberIDs",
						Other: "MemberIDs",
					},
				})
				organizationStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "Organization", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "Organization",
						Other: "Organization",
					},
				})
				msgStr := h.Translator.TranslateWithOpts(
					lang.GetLocaleForTranslation(ctx), "PluralNotAssociated", i18n.Options{
						DefaultMessage: &i18n.Message{
							ID:    "PluralNotAssociated",
							Other: "{{.ID}} not found",
						},
						TemplateData: map[string]any{
							"ID":         memberStr,
							"ValueType":  "ID",
							"Associated": organizationStr,
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
				chatRoomStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "ChatRoom", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "ChatRoom",
						Other: "ChatRoom",
					},
				})
				msgStr := h.Translator.TranslateWithOpts(
					lang.GetLocaleForTranslation(ctx), "PluralNotAssociated", i18n.Options{
						DefaultMessage: &i18n.Message{
							ID:    "PluralNotAssociated",
							Other: "{{.ID}} not found",
						},
						TemplateData: map[string]any{
							"ID":         memberStr,
							"ValueType":  "ID",
							"Associated": chatRoomStr,
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

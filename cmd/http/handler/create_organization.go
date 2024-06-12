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

// CreateOrganization is a handler for creating organization.
type CreateOrganization struct {
	Service    service.ManagerInterface
	Validator  validation.Validator
	Translator i18n.Translation
}

// CreateOrganizationRequest is a request for CreateOrganization.
type CreateOrganizationRequest struct {
	Name                 string      `json:"name" validate:"required,max=255" ja:"名前" en:"Name"`
	Description          string      `json:"description,omitempty" validate:"omitempty" ja:"説明" en:"Description"`
	Color                string      `json:"color,omitempty" validate:"omitempty,hexcolor" ja:"色" en:"Color"`
	WithChatRoom         bool        `json:"with_chat_room" validate:"boolean" ja:"WithChatRoom" en:"WithChatRoom"`
	ChatRoomCoverImageID uuid.UUID   `json:"chat_room_cover_image_id" validate:"" ja:"カバー画像ID" en:"CoverImageID"`
	MemberIDS            []uuid.UUID `json:"member_ids" validate:"required,unique" ja:"メンバーID" en:"MemberIDs"`
}

func (h *CreateOrganization) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authUser := auth.FromContext(ctx)
	var err error
	var organizationReq CreateOrganizationRequest
	if err = json.NewDecoder(r.Body).Decode(&organizationReq); err == nil || errors.Is(err, io.EOF) {
		if errors.Is(err, io.EOF) {
			organizationReq = CreateOrganizationRequest{}
		}
		err = h.Validator.ValidateWithLocale(ctx, &organizationReq, lang.GetLocale(r.Context()))
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
	var organization entity.Organization
	if organization, err = h.Service.CreateOrganization(
		ctx,
		organizationReq.Name,
		entity.String{
			Valid:  organizationReq.Description != "",
			String: organizationReq.Description,
		},
		entity.String{
			Valid:  organizationReq.Color != "",
			String: organizationReq.Color,
		},
		authUser.MemberID,
		organizationReq.MemberIDS,
		organizationReq.WithChatRoom,
		entity.UUID{
			Valid: organizationReq.ChatRoomCoverImageID != uuid.Nil,
			Bytes: organizationReq.ChatRoomCoverImageID,
		},
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
				ve.Add("chat_room_cover_image_id", msgStr)
				err = ve
			case service.OrganizationTargetOwner:
				e.SetTarget("owner")
				err = e
			case service.OrganizationTargetMembers:
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
			if ce.Code.Code == response.CannotAddSelfToOrganization.Code {
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
		err = response.JSONResponseWriter(ctx, w, response.Success, organization, nil)
	}
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
}

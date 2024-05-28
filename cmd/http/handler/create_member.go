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
)

// CreateMember is a handler for creating member.
type CreateMember struct {
	Service    service.ManagerInterface
	Validator  validation.Validator
	Translator i18n.Translation
}

// CreateMemberRequest is a request for CreateMember.
//
//nolint:lll
type CreateMemberRequest struct {
	LoginID              string    `json:"login_id" validate:"required,max=255" ja:"ログインID" en:"LoginID"`
	Password             string    `json:"password" validate:"required,max=255" ja:"パスワード" en:"Password"`
	PasswordConfirmation string    `json:"password_confirmation" validate:"required,max=255,eqfield=Password" ja:"パスワード確認" en:"PasswordConfirmation"`
	Email                string    `json:"email" validate:"required,email,max=255" ja:"メールアドレス" en:"Email"`
	Name                 string    `json:"name" validate:"required,max=255" ja:"名前" en:"Name"`
	FirstName            string    `json:"first_name" validate:"max=255" ja:"名" en:"FirstName"`
	LastName             string    `json:"last_name" validate:"max=255" ja:"姓" en:"LastName"`
	GradeID              uuid.UUID `json:"grade_id" validate:"required" ja:"学年ID" en:"GradeID"`
	GroupID              uuid.UUID `json:"group_id" validate:"required" ja:"班ID" en:"GroupID"`
	ProfileImageID       uuid.UUID `json:"profile_image_id" validate:"" ja:"プロフィール画像ID" en:"ProfileImageID"`
	RoleID               uuid.UUID `json:"member_id" validate:"" ja:"ロールID" en:"RoleID"`
}

func (h *CreateMember) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	var memberReq CreateMemberRequest
	if err = json.NewDecoder(r.Body).Decode(&memberReq); err == nil || errors.Is(err, io.EOF) {
		if errors.Is(err, io.EOF) {
			memberReq = CreateMemberRequest{}
		}
		err = h.Validator.ValidateWithLocale(ctx, &memberReq, lang.GetLocale(r.Context()))
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
	var member entity.Member
	if member, err = h.Service.CreateMember(
		ctx,
		memberReq.LoginID,
		memberReq.Password,
		memberReq.Email,
		memberReq.Name,
		entity.String{
			Valid:  memberReq.FirstName != "",
			String: memberReq.FirstName,
		},
		entity.String{
			Valid:  memberReq.LastName != "",
			String: memberReq.LastName,
		},
		memberReq.GradeID,
		memberReq.GroupID,
		entity.UUID{
			Valid: memberReq.ProfileImageID != uuid.Nil,
			Bytes: memberReq.ProfileImageID,
		},
		entity.UUID{
			Valid: memberReq.RoleID != uuid.Nil,
			Bytes: memberReq.RoleID,
		},
	); err != nil {
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			switch e.Target() {
			case service.MemberTargetGrades:
				gradeStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "GradeID", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "GradeID",
						Other: "GradeID",
					},
				})
				msgStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "ModelNotExists", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "ModelNotExists",
						Other: "{{.ID}} not found",
					},
					TemplateData: map[string]any{
						"ID": gradeStr,
					},
				})
				ve := errhandle.NewValidationError(nil)
				ve.Add("grade_id", msgStr)
				err = ve
			case service.MemberTargetGroups:
				groupStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "GroupID", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "GroupID",
						Other: "GroupID",
					},
				})
				msgStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "ModelNotExists", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "ModelNotExists",
						Other: "{{.ID}} not found",
					},
					TemplateData: map[string]any{
						"ID": groupStr,
					},
				})
				ve := errhandle.NewValidationError(nil)
				ve.Add("group_id", msgStr)
				err = ve
			case service.MemberTargetProfileImages:
				profileImageStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "ProfileImageID", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "ProfileImageID",
						Other: "ProfileImageID",
					},
				})
				msgStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "ModelNotExists", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "ModelNotExists",
						Other: "{{.ID}} not found",
					},
					TemplateData: map[string]any{
						"ID": profileImageStr,
					},
				})
				ve := errhandle.NewValidationError(nil)
				ve.Add("profile_image_id", msgStr)
				err = ve
			case service.MemberTargetRoles:
				roleStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "RoleID", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "RoleID",
						Other: "RoleID",
					},
				})
				msgStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "ModelNotExists", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "ModelNotExists",
						Other: "{{.ID}} not found",
					},
					TemplateData: map[string]any{
						"ID": roleStr,
					},
				})
				ve := errhandle.NewValidationError(nil)
				ve.Add("role_id", msgStr)
				err = ve
			}
		}
	} else {
		err = response.JSONResponseWriter(ctx, w, response.Success, member, nil)
	}
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
}

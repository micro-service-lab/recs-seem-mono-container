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

// CreateProfessor is a handler for creating professor.
type CreateProfessor struct {
	Service    service.ManagerInterface
	Validator  validation.Validator
	Translator i18n.Translation
}

// CreateProfessorRequest is a request for CreateProfessor.
//
//nolint:lll
type CreateProfessorRequest struct {
	LoginID              string    `json:"login_id" validate:"required,max=255" ja:"ログインID" en:"LoginID"`
	Password             string    `json:"password" validate:"required,max=255" ja:"パスワード" en:"Password"`
	PasswordConfirmation string    `json:"password_confirmation" validate:"required,max=255,eqfield=Password" ja:"パスワード確認" en:"PasswordConfirmation"`
	Email                string    `json:"email" validate:"required,email,max=255" ja:"メールアドレス" en:"Email"`
	Name                 string    `json:"name" validate:"required,max=255" ja:"名前" en:"Name"`
	FirstName            string    `json:"first_name" validate:"max=255" ja:"名" en:"FirstName"`
	LastName             string    `json:"last_name" validate:"max=255" ja:"姓" en:"LastName"`
	RoleID               uuid.UUID `json:"professor_id" validate:"" ja:"ロールID" en:"RoleID"`
}

func (h *CreateProfessor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	var professorReq CreateProfessorRequest
	if err = json.NewDecoder(r.Body).Decode(&professorReq); err == nil || errors.Is(err, io.EOF) {
		if errors.Is(err, io.EOF) {
			professorReq = CreateProfessorRequest{}
		}
		err = h.Validator.ValidateWithLocale(ctx, &professorReq, lang.GetLocale(r.Context()))
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
	var professor entity.Professor
	if professor, err = h.Service.CreateProfessor(
		ctx,
		professorReq.LoginID,
		professorReq.Password,
		professorReq.Email,
		professorReq.Name,
		entity.String{
			Valid:  professorReq.FirstName != "",
			String: professorReq.FirstName,
		},
		entity.String{
			Valid:  professorReq.LastName != "",
			String: professorReq.LastName,
		},
		entity.UUID{
			Valid: professorReq.RoleID != uuid.Nil,
			Bytes: professorReq.RoleID,
		},
	); err != nil {
		var de errhandle.ModelDuplicatedError
		if errors.As(err, &de) {
			switch de.Target() {
			case service.MemberTargetLoginID:
				loginIDStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "LoginID", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "LoginID",
						Other: "LoginID",
					},
				})
				msgStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "ModelExists", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "ModelExists",
						Other: "{{.ID}} already exists",
					},
					TemplateData: map[string]any{
						"ID": loginIDStr,
					},
				})
				ve := errhandle.NewValidationError(nil)
				ve.Add("login_id", msgStr)
				err = ve
			}
		}
		var e errhandle.ModelNotFoundError
		if errors.As(err, &e) {
			switch e.Target() {
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
		err = response.JSONResponseWriter(ctx, w, response.Success, professor, nil)
	}
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
}

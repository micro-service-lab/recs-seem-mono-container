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

// CreateStudent is a handler for creating student.
type CreateStudent struct {
	Service    service.ManagerInterface
	Validator  validation.Validator
	Translator i18n.Translation
}

// CreateStudentRequest is a request for CreateStudent.
//
//nolint:lll
type CreateStudentRequest struct {
	LoginID              string    `json:"login_id" validate:"required,max=255" ja:"ログインID" en:"LoginID"`
	Password             string    `json:"password" validate:"required,max=255" ja:"パスワード" en:"Password"`
	PasswordConfirmation string    `json:"password_confirmation" validate:"required,max=255,eqfield=Password" ja:"パスワード確認" en:"PasswordConfirmation"`
	Email                string    `json:"email" validate:"required,email,max=255" ja:"メールアドレス" en:"Email"`
	Name                 string    `json:"name" validate:"required,max=255" ja:"名前" en:"Name"`
	FirstName            string    `json:"first_name" validate:"max=255" ja:"名" en:"FirstName"`
	LastName             string    `json:"last_name" validate:"max=255" ja:"姓" en:"LastName"`
	GradeID              uuid.UUID `json:"grade_id" validate:"required" ja:"学年ID" en:"GradeID"`
	GroupID              uuid.UUID `json:"group_id" validate:"required" ja:"班ID" en:"GroupID"`
	RoleID               uuid.UUID `json:"role_id" validate:"" ja:"ロールID" en:"RoleID"`
}

func (h *CreateStudent) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	var studentReq CreateStudentRequest
	if err = json.NewDecoder(r.Body).Decode(&studentReq); err == nil || errors.Is(err, io.EOF) {
		if errors.Is(err, io.EOF) {
			studentReq = CreateStudentRequest{}
		}
		err = h.Validator.ValidateWithLocale(ctx, &studentReq, lang.GetLocale(r.Context()))
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
	var student entity.Student
	if student, err = h.Service.CreateStudent(
		ctx,
		studentReq.LoginID,
		studentReq.Password,
		studentReq.Email,
		studentReq.Name,
		entity.String{
			Valid:  studentReq.FirstName != "",
			String: studentReq.FirstName,
		},
		entity.String{
			Valid:  studentReq.LastName != "",
			String: studentReq.LastName,
		},
		studentReq.GradeID,
		studentReq.GroupID,
		entity.UUID{
			Valid: studentReq.RoleID != uuid.Nil,
			Bytes: studentReq.RoleID,
		},
	); err != nil {
		var ce errhandle.CommonError
		if errors.As(err, &ce) && ce.Code.Code == response.OnlyProfessorAction.Code {
			switch ce.Target {
			case service.MemberTargetGrades:
				gradeStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "GradeID", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "GradeID",
						Other: "GradeID",
					},
				})
				msgStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "OnlyProfessorModel", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "OnlyProfessorModel",
						Other: "{{.ID}} professor only",
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
				msgStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "OnlyProfessorModel", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "OnlyProfessorModel",
						Other: "{{.ID}} professor only",
					},
					TemplateData: map[string]any{
						"ID": groupStr,
					},
				})
				ve := errhandle.NewValidationError(nil)
				ve.Add("group_id", msgStr)
				err = ve
			}
		}
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
		err = response.JSONResponseWriter(ctx, w, response.Success, student, nil)
	}
	if err != nil {
		handled, err := errhandle.ErrorHandle(ctx, w, err)
		if !handled || err != nil {
			log.Printf("failed to handle error: %v", err)
		}
		return
	}
}

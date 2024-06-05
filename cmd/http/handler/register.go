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

// Register is a handler for registering member.
type Register struct {
	Service    service.ManagerInterface
	Validator  validation.Validator
	Translator i18n.Translation
}

// RegisterRequest is a request for Register.
//
//nolint:lll
type RegisterRequest struct {
	LoginID              string    `json:"login_id" validate:"required,max=255" ja:"ログインID" en:"LoginID"`
	Password             string    `json:"password" validate:"required,max=255" ja:"パスワード" en:"Password"`
	PasswordConfirmation string    `json:"password_confirmation" validate:"required,max=255,eqfield=Password" ja:"パスワード確認" en:"PasswordConfirmation"`
	Email                string    `json:"email" validate:"required,email,max=255" ja:"メールアドレス" en:"Email"`
	Name                 string    `json:"name" validate:"required,max=255" ja:"名前" en:"Name"`
	FirstName            string    `json:"first_name" validate:"max=255" ja:"名" en:"FirstName"`
	LastName             string    `json:"last_name" validate:"max=255" ja:"姓" en:"LastName"`
	GradeID              uuid.UUID `json:"grade_id" validate:"required" ja:"学年ID" en:"GradeID"`
	GroupID              uuid.UUID `json:"group_id" validate:"required" ja:"班ID" en:"GroupID"`
}

func (h *Register) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	var memberReq RegisterRequest
	if err = json.NewDecoder(r.Body).Decode(&memberReq); err == nil || errors.Is(err, io.EOF) {
		if errors.Is(err, io.EOF) {
			memberReq = RegisterRequest{}
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
	var member entity.MemberWithDetail
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
		entity.UUID{},
	); err != nil {
		var ce errhandle.CommonError
		if errors.As(err, &ce) {
			if ce.Code.Code == response.OnlyProfessorAction.Code {
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
			} else if ce.Code.Code == response.MustBeGradeProfessorIfGroupProfessor.Code {
				gradeStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "Grade", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "Grade",
						Other: "Grade",
					},
				})
				groupStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "Group", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "Group",
						Other: "Group",
					},
				})
				professorStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "Professor", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "Professor",
						Other: "Professor",
					},
				})
				msgStr := h.Translator.TranslateWithOpts(
					lang.GetLocaleForTranslation(ctx), "MustBeEntityStatusIfEntityStatus", i18n.Options{
						DefaultMessage: &i18n.Message{
							ID:    "MustBeEntityStatusIfEntityStatus",
							Other: "If the {{.Assumption}} is a {{.AssumptionStatus}}, the {{.Target}} must be a {{.TargetStatus}}",
						},
						TemplateData: map[string]any{
							"Assumption":       groupStr,
							"AssumptionStatus": professorStr,
							"Target":           gradeStr,
							"TargetStatus":     professorStr,
						},
					})
				ve := errhandle.NewValidationError(nil)
				ve.Add("grade_id", msgStr)
				err = ve
			} else if ce.Code.Code == response.MustBeGroupProfessorIfGradeProfessor.Code {
				gradeStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "Grade", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "Grade",
						Other: "Grade",
					},
				})
				groupStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "Group", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "Group",
						Other: "Group",
					},
				})
				professorStr := h.Translator.TranslateWithOpts(lang.GetLocaleForTranslation(ctx), "Professor", i18n.Options{
					DefaultMessage: &i18n.Message{
						ID:    "Professor",
						Other: "Professor",
					},
				})
				msgStr := h.Translator.TranslateWithOpts(
					lang.GetLocaleForTranslation(ctx), "MustBeEntityStatusIfEntityStatus", i18n.Options{
						DefaultMessage: &i18n.Message{
							ID:    "MustBeEntityStatusIfEntityStatus",
							Other: "If the {{.Assumption}} is a {{.AssumptionStatus}}, the {{.Target}} must be a {{.TargetStatus}}",
						},
						TemplateData: map[string]any{
							"Assumption":       gradeStr,
							"AssumptionStatus": professorStr,
							"Target":           groupStr,
							"TargetStatus":     professorStr,
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

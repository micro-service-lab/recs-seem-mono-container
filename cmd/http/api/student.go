package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/micro-service-lab/recs-seem-mono-container/app/i18n"
	"github.com/micro-service-lab/recs-seem-mono-container/app/service"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler"
	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/validation"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/auth"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/clock"
	"github.com/micro-service-lab/recs-seem-mono-container/internal/session"
)

// StudentHandler is a handler for roles.
func StudentHandler(
	svc service.ManagerInterface,
	vd validation.Validator,
	t i18n.Translation,
	clk clock.Clock,
	auth auth.Auth,
	ssm session.Manager,
) http.Handler {
	createHandler := handler.CreateStudent{
		Service:    svc,
		Validator:  vd,
		Translator: t,
	}
	deleteHandler := handler.DeleteStudent{
		Service: svc,
	}
	updateGroupHandler := handler.UpdateStudentGroup{
		Service:    svc,
		Validator:  vd,
		Translator: t,
	}
	updateGradeHandler := handler.UpdateStudentGrade{
		Service:    svc,
		Validator:  vd,
		Translator: t,
	}

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware(clk.Now, auth, svc, ssm))

		r.Post("/", createHandler.ServeHTTP)
		r.Delete("/{student_id}", deleteHandler.ServeHTTP)
		r.Put("/{student_id}/group", updateGroupHandler.ServeHTTP)
		r.Put("/{student_id}/grade", updateGradeHandler.ServeHTTP)
	})

	return r
}

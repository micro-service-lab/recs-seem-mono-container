package recoverer

import (
	"log"
	"net/http"

	"github.com/micro-service-lab/recs-seem-mono-container/cmd/http/handler/response"
)

// Handler is a function that handles a panic.
func Handler(w http.ResponseWriter, r *http.Request, _ any) {
	e := response.JSONResponseWriter(r.Context(), w, response.System, nil, nil)
	if e != nil {
		log.Printf("failed to write response: %v", e)
	}
}

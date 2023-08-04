package rotavator

import (
	"fmt"
	"net/http"

	"rotorvator/rotavator/security"
	"rotorvator/rotavator/templates"
)

type HomeHandler struct {
}

// HomeHandler main user info
func (h HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	rCtx := security.CookieContext(r.Context())
	s := templates.RenderPage("home.html", map[string]any{
		"title":    "Home",
		"email":    rCtx["email"],
		"name":     "Bob Smith",
		"nextDuty": "1 Oct 2023",
		"prevDuty": "1 Jan 2023",
	})
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, s)
}

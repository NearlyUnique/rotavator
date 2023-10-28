package web

import (
	"net/http"

	"rotavator/security"
	"rotavator/templates"
)

type HomeHandler struct {
	cookie *security.CookieCutter
}

func NewHomeHandler(cookie *security.CookieCutter) HomeHandler {
	return HomeHandler{cookie: cookie}
}

// HomeHandler main user info
func (h HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	profile := security.CookieContext(r.Context())
	if !profile.IsLoggedIn() {
		_, _ = w.Write([]byte(templates.RenderPage("not_logged_in.html", map[string]any{
			"title": "Home",
		})))
		return
	}
	_, _ = w.Write([]byte(templates.RenderPage("home.html", map[string]any{
		"title":      "Home",
		"nickname":   profile.Nickname,
		"pictureURL": profile.PictureURL,
	})))
}

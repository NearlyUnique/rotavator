package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"rotavator/security"
	"rotavator/templates"
)

type AuthHandler struct {
	cookie *security.CookieCutter
	auth   *security.Authenticator
}

func NewLoginHandler(cookie *security.CookieCutter, auth *security.Authenticator) AuthHandler {
	return AuthHandler{
		cookie: cookie,
		auth:   auth,
	}
}

func (h AuthHandler) SetupRoutes(r *mux.Router) {
	r.Handle("/login", h.GetLoginWithRedirect()).
		Methods(http.MethodGet)
	r.Handle("/callback", h.cookie.RequireCookie(h.GetCallback(), security.ProfileStateInProgress)).
		Methods(http.MethodGet)
	r.Handle("/user", h.cookie.RequireCookie(h.GetUser(), security.ProfileStateLoggedIn)).
		Methods(http.MethodGet)
}

func (h AuthHandler) GetLoginWithRedirect() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state, err := security.GenerateRandomState()
		if err != nil {
			log.Printf("login GET error %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			// TO fix
			_, _ = fmt.Fprintf(w, err.Error())
			return
		}
		var c *http.Cookie
		c, err = h.cookie.MakeCookie(security.Profile{AuthState: state})
		if err != nil {
			log.Printf("login POST error %v", err)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, err.Error())
			return
		}
		http.SetCookie(w, c)
		http.Redirect(w, r, h.auth.AuthCodeURL(state), http.StatusTemporaryRedirect)
	})
}

func (h AuthHandler) GetCallback() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		profile, err := h.cookie.GetProfileCookie(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "no state cookie")
			return
		}
		queryValue := r.URL.Query().Get("state")

		if profile.AuthState != queryValue {
			log.Printf("bad state query(%s) != cookie(%s)", queryValue, profile.AuthState)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "bad state")
			return
		}
		token, err := h.auth.Exchange(r.Context(), r.URL.Query().Get("code"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "Failed to convert an authorization code into a token.")
			return
		}
		idToken, err := h.auth.VerifyIDToken(r.Context(), token)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "Failed to verify ID Token.")
			return
		}

		var claims map[string]any
		if err := idToken.Claims(&claims); err != nil {
			log.Printf("login Claims error %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, err.Error())
			return
		}
		profile = security.ProfileFromClaims(claims)

		// Redirect to logged in page.
		var c *http.Cookie
		c, err = h.cookie.MakeCookie(profile)
		if err != nil {
			log.Printf("callback make cookie error %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, err.Error())
			return
		}
		http.SetCookie(w, c)
		http.Redirect(w, r, "/auth/user", http.StatusTemporaryRedirect)
	})
}

func (h AuthHandler) GetUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		profile, err := h.cookie.GetProfileCookie(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "no state cookie: %v", err)
			return
		}
		//rCtx := security.CookieContext(r.Context())
		s := templates.RenderPage("user.html", map[string]any{
			"title":      "User",
			"email":      profile.Email,
			"pictureURL": profile.PictureURL,
			"nickname":   profile.Nickname,
		})
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, s)
	})
}

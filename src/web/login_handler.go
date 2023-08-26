package web

import (
	"encoding/json"
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
	r.Handle("/callback", h.GetCallback()).
		Methods(http.MethodGet)
	r.Handle("/user", h.GetUser()).
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
		cookieValue := map[string]string{
			"state": state,
		}
		var c *http.Cookie
		c, err = h.cookie.MakeCookie(cookieValue)
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
		cookieValue := h.cookie.GetAuthCookie(r)
		queryValue := r.URL.Query().Get("state")

		if cookieValue["state"] != queryValue {
			log.Printf("bad state '%s'!='%s'", queryValue, cookieValue)
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

		var profile map[string]any
		if err := idToken.Claims(&profile); err != nil {
			log.Printf("login Claims error %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, err.Error())
			return
		}
		cookieValue["access_token"] = token.AccessToken
		j, _ := json.Marshal(profile)
		cookieValue["profile"] = string(j)
		// Redirect to logged in page.
		var c *http.Cookie
		c, err = h.cookie.MakeCookie(cookieValue)
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
		cookieValues := h.cookie.GetAuthCookie(r)
		var profile map[string]any
		log.Print(cookieValues["profile"])
		err := json.Unmarshal([]byte(cookieValues["profile"]), &profile)
		if err != nil {
			log.Printf("json %v", err)
		}
		//rCtx := security.CookieContext(r.Context())
		s := templates.RenderPage("user.html", map[string]any{
			"title":    "User",
			"email":    profile["email"],
			"picture":  profile["picture"],
			"name":     profile["name"],
			"nickname": profile["nickname"],
			"nextDuty": "1 Oct 2023",
			"prevDuty": "1 Jan 2023",
		})
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, s)
	})
}

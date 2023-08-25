package web

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"rotavator/security"
	"rotavator/templates"
)

type AuthHandler struct {
	cookie *security.CookieCutter
}

func NewLoginHandler(cookie *security.CookieCutter) AuthHandler {
	return AuthHandler{
		cookie: cookie,
	}
}

func (h AuthHandler) SetupRoutes(r *mux.Router) {
	r.Handle("/login", h.ViewLoginPage()).
		Methods(http.MethodGet)
	r.Handle("/login", h.LoginWithRedirect()).
		Methods(http.MethodPost)
	r.Handle("/pending", h.ViewPendingLoginPage()).
		Methods(http.MethodGet)
	r.Handle("/token", h.TokenWithRedirect()).
		Methods(http.MethodGet)
	r.Handle("/logout", h.LogoutWithRedirect()).
		Methods(http.MethodPost)
}

func (h AuthHandler) LoginWithRedirect() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("content-type") != "application/x-www-form-urlencoded" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := r.ParseForm()
		if err != nil || r.Form.Get("email") == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		token, err := h.tokenFor(r.Form.Get("email"))
		if err != nil {
			log.Printf("token creation error %v", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = fmt.Fprintf(w, err.Error())
			return
		}
		cookieValue := map[string]any{}
		cookieValue["email"] = r.Form.Get("email")
		cookieValue["locked"] = true
		cookieValue["created"] = time.Now().Unix()
		cookieValue["token"] = token

		// add "unlock code"
		c, err := h.cookie.MakeCookie(cookieValue)
		if err != nil {
			log.Printf("login POST error %v", err)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, err.Error())
			return
		}
		http.SetCookie(w, c)
		// adding email to this redirect is temporary
		http.Redirect(w, r, "/auth/pending", http.StatusSeeOther)
	})
}

func (h AuthHandler) LogoutWithRedirect() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//cookieValue := map[string]string{}
		//c, err := h.cookie.MakeCookie(cookieValue)
		//if err != nil {
		//	log.Printf("logout POST error %v", err)
		//	w.WriteHeader(http.StatusBadRequest)
		//	_, _ = fmt.Fprintf(w, err.Error())
		//	return
		//}
		//http.SetCookie(w, c)
		//http.Redirect(w, r, "/", http.StatusSeeOther)
		w.WriteHeader(http.StatusInternalServerError)
	})
}
func (h AuthHandler) ViewLoginPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		s := templates.RenderPage("login.html", map[string]any{
			"title": "Login",
			"token": "random-token",
		})
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, s)
	})
}

func (h AuthHandler) ViewPendingLoginPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie("_auth")
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		http.SetCookie(w, tokenCookie)
		s := templates.RenderPage("login_pending.html", map[string]any{
			"title": "Login Email Sent",
		})
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, s)
	})
}
func (h AuthHandler) TokenWithRedirect() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieValue := map[string]string{}
		vars := mux.Vars(r)
		cookieValue["email"] = vars["token"]
		cookieValue["created"] = time.Now().String()
		c, err := h.cookie.MakeCookie(cookieValue)
		if err != nil {
			log.Printf("login POST error %v", err)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, err.Error())
			return
		}
		http.SetCookie(w, c)
		w.WriteHeader(http.StatusOK)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}

func (h AuthHandler) tokenFor(email string) (string, error) {
	return email, nil
}

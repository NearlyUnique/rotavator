package rotavator

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"rotavator/rotavator/security"
)

type App struct{}

func (a App) Run() error {

	secureCookie := securecookie.New([]byte("123456789012345678901234567890aa"), []byte("1234567890abcdef"))

	r := SetupRoutes(secureCookie)

	server := http.Server{Addr: "0.0.0.0:5001", Handler: r}
	return server.ListenAndServe()
}

// SetupRoutes as it says on the tin
func SetupRoutes(secureCookie security.CookieSecrets) *mux.Router {
	cookies := security.NewCookieCutter(secureCookie, "_auth", "/auth/login")
	r := mux.NewRouter()
	r.Handle("/", cookies.RequireCookie(HomeHandler{}))
	login := NewLoginHandler(cookies)
	r.Handle("/auth/login", login.ViewLoginPage()).
		Methods(http.MethodGet)
	r.Handle("/auth/login", login.LoginWithRedirect()).
		Methods(http.MethodPost)
	r.Handle("/auth/pending", login.ViewPendingLoginPage()).
		Methods(http.MethodGet)
	r.Handle("/auth/token", login.TokenWithRedirect()).
		Methods(http.MethodGet)
	r.Handle("/auth/logout", login.LogoutWithRedirect()).
		Methods(http.MethodPost)
	return r
}

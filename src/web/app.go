package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"rotavator/security"
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
	login.SetupRoutes(r.Path("/auth").Subrouter())
	return r
}

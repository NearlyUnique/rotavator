package web

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"rotavator/security"
)

type App struct{}

func (a App) Run() error {

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(security.Profile{})

	secureCookie := securecookie.New([]byte(os.Getenv("COOKIE_HASH_KEY")), []byte(os.Getenv("COOKIE_BLOCK_KEY")))
	auth, err := security.NewAuthenticator()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	r := SetupRoutes(secureCookie, auth)

	server := http.Server{Addr: "0.0.0.0:5001", Handler: r}
	return server.ListenAndServe()
}

// SetupRoutes as it says on the tin
func SetupRoutes(secureCookie security.CookieSecrets, auth *security.Authenticator) *mux.Router {
	cookies := security.NewCookieCutter(secureCookie, "_auth", "/")
	r := mux.NewRouter()

	login := NewLoginHandler(cookies, auth)
	login.SetupRoutes(r.PathPrefix("/auth").Subrouter())

	home := NewHomeHandler(cookies)
	r.Handle("/", cookies.RequireCookie(home, security.ProfileStateAny))

	return r
}

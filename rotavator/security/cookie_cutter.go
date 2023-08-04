package security

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

type contextKey string

const (
	ctxKey contextKey = "ctxKey"
)

// this exists only for testing, should always be true in PROD
var secureCookies = true

// MakeCookiesUnsecure for testing only
func MakeCookiesUnsecure() {
	_, _ = fmt.Fprintln(os.Stderr, "COOKIES are INSECURE")
	secureCookies = false
}

// CookieCutter for the site
type CookieCutter struct {
	cutter      CookieSecrets
	cookieName  string
	redirectUrl string
}

// CookieSecrets keeps secrets
type CookieSecrets interface {
	Encode(name string, value interface{}) (string, error)
	Decode(name, value string, dst interface{}) error
}

// NewCookieCutter for the site
// hashKey should be at least 32 bytes long
// blockKeys should be 16 bytes (AES-128) or 32 bytes (AES-256) long
func NewCookieCutter(secureCookie CookieSecrets, cookieName, redirectUrl string) *CookieCutter {
	return &CookieCutter{
		cutter:      secureCookie,
		cookieName:  cookieName,
		redirectUrl: redirectUrl,
	}
}

// CookieContext extracts the secure cookie data
func CookieContext(ctx context.Context) map[string]string {
	rCtx := ctx.Value(ctxKey)
	v, ok := rCtx.(map[string]string)
	if !ok || v == nil {
		return map[string]string{}
	}
	return v
}

// MakeCookie to the response
func (c CookieCutter) MakeCookie(value any) (*http.Cookie, error) {
	encoded, err := c.cutter.Encode(c.cookieName, value)
	if err == nil {
		return &http.Cookie{
			Name:     c.cookieName,
			Value:    encoded,
			Path:     "/",
			Secure:   secureCookies,
			HttpOnly: true,
		}, nil
	}
	return nil, err
}

// RequireCookie decode cookie, put's it in the context. Requires a cookie or returns
func (c CookieCutter) RequireCookie(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie(c.cookieName); err == nil {
			value := make(map[string]string)
			if err = c.cutter.Decode(c.cookieName, cookie.Value, &value); err == nil && len(value) > 0 {
				ctx := context.WithValue(r.Context(), ctxKey, value)
				r = r.WithContext(ctx)
				h.ServeHTTP(w, r)
				return
			}
		}
		// not authenticated
		http.Redirect(w, r, c.redirectUrl, http.StatusSeeOther)
	})
}

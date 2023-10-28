package security

import (
	"context"
	"errors"
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
func CookieContext(ctx context.Context) Profile {
	rCtx := ctx.Value(ctxKey)
	profile, ok := rCtx.(Profile)
	if !ok {
		return Profile{}
	}
	return profile
}

// MakeCookie to the response
func (c CookieCutter) MakeCookie(value Profile) (*http.Cookie, error) {
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
func (c CookieCutter) RequireCookie(h http.Handler, requiredState ProfileState) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		profile, err := c.GetProfileCookie(r)
		if err != nil {
			http.Redirect(w, r, c.redirectUrl, http.StatusSeeOther)
			return
		}
		if requiredState != ProfileStateAny && profile.State() != requiredState {
			// log error?
			http.Redirect(w, r, c.redirectUrl, http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), ctxKey, profile)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

func (c CookieCutter) GetProfileCookie(r *http.Request) (Profile, error) {
	cookie, err := r.Cookie(c.cookieName)
	if errors.Is(err, http.ErrNoCookie) {
		// that's fine, no profile
		return Profile{}, nil
	}
	if err != nil {
		return Profile{}, err
	}
	var p Profile
	err = c.cutter.Decode(c.cookieName, cookie.Value, &p)
	if err != nil {
		return Profile{}, err
	}
	if len(p.Email) == 0 && len(p.AuthState) == 0 {
		return Profile{}, fmt.Errorf("zero length after decode")
	}
	return p, nil
}

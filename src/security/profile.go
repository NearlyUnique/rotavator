package security

import (
	"net/url"
	"strings"
	"time"
)

type Profile struct {
	AuthState  string
	Nickname   string
	Email      string
	PictureURL string
	ExternalID string
	UpdatedAt  time.Time
}
type ProfileState string

const (
	ProfileStateAny        ProfileState = "Any"
	ProfileStateNew        ProfileState = "New"
	ProfileStateInProgress ProfileState = "InProgress"
	ProfileStateLoggedIn   ProfileState = "LoggedIn"
)

func (p Profile) State() ProfileState {
	switch {
	case p.AuthState != "" && p.Email == "":
		return ProfileStateInProgress
	case p.AuthState == "" && p.Email != "":
		return ProfileStateLoggedIn
	default:
		return ProfileStateNew
	}
}
func (p Profile) IsLoggedIn() bool {
	return p.State() == ProfileStateLoggedIn
}

// ProfileFromClaims generate internal profile from ID provider
func ProfileFromClaims(claims map[string]any) Profile {
	var p Profile
	asString := func(key string) string {
		if v, ok := claims[key].(string); ok {
			return v
		}
		return ""
	}
	p.Nickname = asString("nickname")
	p.Email = asString("name")
	p.PictureURL = fixPicture(asString("picture"))
	p.ExternalID = asString("sub")
	if s := asString("updated_at"); s != "" {
		if dtm, err := time.Parse(time.RFC3339, s); err == nil {
			p.UpdatedAt = dtm
		}
	}
	return p
}
func fixPicture(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		return s
	}
	if strings.HasSuffix(u.Hostname(), "gravatar.com") {
		query := u.Query()
		query.Del("s")
		query.Del("size")
		u.RawQuery = query.Encode()
		return u.String()
	}
	return s
}

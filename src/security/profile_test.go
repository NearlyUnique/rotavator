package security_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"rotavator/security"
)

func Test_basic_profile(t *testing.T) {
	claims := map[string]any{
		"nickname":   "some name",
		"name":       "any@example.com",
		"picture":    "https://s.gravatar.com/avatar/random_number",
		"updated_at": "2023-08-19T14:09:55.123Z",
		"iss":        "https://random.uk.auth0.com/",
		"aud":        "random_chars",
		"iat":        1694268891,
		"exp":        1694304891,
		"sub":        "oauth2|custom-google|1234567890",
		"sid":        "random_chars",
	}
	profile := security.ProfileFromClaims(claims)

	assert.Equal(t, "some name", profile.Nickname)
	assert.Equal(t, "any@example.com", profile.Email)
	assert.Equal(t, "oauth2|custom-google|1234567890", profile.ExternalID)
	assert.Equal(t, "https://s.gravatar.com/avatar/random_number", profile.PictureURL)
	assert.Equal(t, time.Date(2023, 8, 19, 14, 9, 55, 123000000, time.UTC), profile.UpdatedAt)
}
func Test_bad_profile_still_renders(t *testing.T) {
	const notAString = 1
	claims := map[string]any{
		"nickname":   notAString,
		"name":       notAString,
		"picture":    notAString,
		"updated_at": notAString,
		"iss":        "https://random.uk.auth0.com/",
		"aud":        "random_chars",
		"iat":        1694268891,
		"exp":        1694304891,
		"sub":        notAString,
		"sid":        "random_chars",
	}
	profile := security.ProfileFromClaims(claims)

	assert.Equal(t, "", profile.Nickname)
	assert.Equal(t, "", profile.Email)
	assert.Equal(t, "", profile.ExternalID)
	assert.Equal(t, "", profile.PictureURL)
	assert.Equal(t, time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC), profile.UpdatedAt)
}
func Test_missing_attributes_profile_still_renders(t *testing.T) {
	claims := map[string]any{}
	profile := security.ProfileFromClaims(claims)

	assert.Equal(t, "", profile.Nickname)
	assert.Equal(t, "", profile.Email)
	assert.Equal(t, "", profile.ExternalID)
	assert.Equal(t, "", profile.PictureURL)
	assert.Equal(t, time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC), profile.UpdatedAt)
}
func Test_gravatar_picture_size_is_removed(t *testing.T) {
	const baseGravatar = "https://s.gravatar.com/avatar/random_number?"
	testData := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "size is missing",
			url:      baseGravatar + "p1=any&p2=other",
			expected: "https://s.gravatar.com/avatar/random_number?p1=any&p2=other",
		},
		{
			name:     "size is removed",
			url:      baseGravatar + "p1=any&size=123&p2=other",
			expected: "https://s.gravatar.com/avatar/random_number?p1=any&p2=other",
		},
		{
			name:     "s is removed",
			url:      baseGravatar + "p1=any&s=123&p2=other",
			expected: "https://s.gravatar.com/avatar/random_number?p1=any&p2=other",
		},
		{
			name:     "size is only",
			url:      baseGravatar + "s=123",
			expected: "https://s.gravatar.com/avatar/random_number",
		},
		{
			name:     "non gravatar passes through as is",
			url:      "https://s.other.com/avatar/random_number?p1=any&size=123&p2=other",
			expected: "https://s.other.com/avatar/random_number?p1=any&size=123&p2=other",
		},
	}
	for _, td := range testData {
		t.Run(td.name, func(t *testing.T) {
			profile := security.ProfileFromClaims(map[string]any{
				"picture": td.url,
			})
			assert.Equal(t, td.expected, profile.PictureURL)
		})
	}
}

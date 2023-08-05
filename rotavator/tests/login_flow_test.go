package tests

import (
	"net/http/httptest"
	"testing"

	"github.com/NearlyUnique/httptestclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"rotavator/rotavator"
	"rotavator/rotavator/security"
)

func Test_login_flow_from_cold(t *testing.T) {
	secrets := MockCookieSecrets{}
	routed := rotavator.SetupRoutes(secrets)
	given := Given{
		t:      t,
		server: httptest.NewServer(routed),
		client: httptestclient.New(t),
	}
	security.MakeCookiesUnsecure()

	const userEmail = "any@example.com"
	given.login_form_can_be_viewed()
	given.user_is_redirected_to_login_pending_page(userEmail, secrets)
	const token = "?"
	given.submitting_token_redirects_to_home(token)
}

type Given struct {
	t      *testing.T
	server *httptest.Server
	client *httptestclient.Client
}

func (given Given) user_is_redirected_to_login_pending_page(userEmail string, secrets MockCookieSecrets) {
	given.t.Helper()
	loginResponse := given.client.
		Post("/auth/login").
		FormData("email", userEmail).
		ExpectRedirectTo("/auth/pending").
		DoSimple(given.server)

	// assert redirect
	cookies := loginResponse.Response.Cookies()
	require.NotEmpty(given.t, cookies)
	var cookie struct {
		Email  string
		Locked bool
		Token  string
	}
	require.NoError(given.t, secrets.Decode("", cookies[0].Value, &cookie))
	require.NotEmpty(given.t, cookies)
	assert.Equal(given.t, "_auth", cookies[0].Name)
	assert.Equal(given.t, userEmail, cookie.Email)
	assert.Equal(given.t, true, cookie.Locked)
	//TODO: not supposed tobe email
	assert.Equal(given.t, userEmail, cookie.Token)

	assertFormExists(given.t, "/auth/token", "GET", loginResponse.Body)
	assertInputExists(given.t, "text", "token", loginResponse.Body)
	assertInputExists(given.t, "submit", "login", loginResponse.Body)
}

func (given Given) login_form_can_be_viewed() {
	given.t.Helper()
	viewLoginPageResponse := given.client.
		Get("/auth/login").
		DoSimple(given.server)

	assertFormExists(given.t, "/auth/login", "POST", viewLoginPageResponse.Body)
	assertInputExists(given.t, "email", "email", viewLoginPageResponse.Body)
	assertInputExists(given.t, "hidden", "token", viewLoginPageResponse.Body)
	assertInputExists(given.t, "submit", "login", viewLoginPageResponse.Body)
}

func (given Given) submitting_token_redirects_to_home(token string) {
	// assert.FailNow(given.t, "not implemented")
}

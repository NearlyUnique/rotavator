package integration_test

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockCookieSecrets struct {
}

func (MockCookieSecrets) Encode(name string, value interface{}) (string, error) {
	v, _ := json.Marshal(value)
	return base64.StdEncoding.EncodeToString(v), nil
}
func (MockCookieSecrets) Decode(name, value string, dst interface{}) error {
	src := []byte(value)
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	n, err := base64.StdEncoding.Decode(decoded, src)
	if err != nil {
		return err
	}
	// NB: decoded[:n] is required because we may have over allocated memory, n is the actual number  of bytes
	err = json.Unmarshal(decoded[:n], dst)
	return err
}

func Test_MockCookieSecrets(t *testing.T) {
	t.Run("round trip secrets", func(t *testing.T) {
		type V struct {
			A string
			B int
		}
		m := MockCookieSecrets{}
		// NB: edge case, if V.A is only one char long you don't hit the "over allocation" NULL bug
		data, err := m.Encode("any", &V{A: "ab", B: 9})
		require.NoError(t, err)

		var actual V
		require.NoError(t, m.Decode("", data, &actual))

		assert.Equal(t, "ab", actual.A)
		assert.Equal(t, 9, actual.B)
	})
}

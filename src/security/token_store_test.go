package security_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"rotavator/security"
)

func Test_write_temp_token(t *testing.T) {
	t.Run("for existing email, write and validate token succeed", func(t *testing.T) {

		store := security.NewTokenStore()

		const anyEmail = "any@example.com"
		const anyToken = "random_token"

		err := store.WriteTempToken(anyEmail, anyToken)
		assert.NoError(t, err)

		err = store.ValidateToken(anyEmail)
		assert.NoError(t, err)
	})
}

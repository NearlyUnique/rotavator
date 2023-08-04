package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// assertInputExists requires <input type="?" ... name="?" ...
func assertInputExists(t *testing.T, typ, name string, body string) {
	t.Helper()

	assert.Regexp(t, fmt.Sprintf(`input\s+type="%s".*name="%s"`, typ, name), body)
}
func assertFormExists(t *testing.T, action, method string, body string) {
	t.Helper()

	assert.Regexp(t, fmt.Sprintf(`form\s+action="%s".*method="%s"`, action, method), body)
}

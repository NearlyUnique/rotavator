package security_test

import (
	"fmt"
	"strings"
	"testing"

	"rotavator/security"
)

func Test_random_text_can_be_generated(t *testing.T) {
	const length = 100
	value, err := security.RandomStr("Ab3", length)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if len(value) != length {
		t.Errorf("token length incorrect, expected %d actual %d", length, len(value))
	}
	trimmed := strings.ReplaceAll(value, "A", "")
	trimmed = strings.ReplaceAll(trimmed, "b", "")
	trimmed = strings.ReplaceAll(trimmed, "3", "")
	if len(trimmed) != 0 {
		t.Errorf("unexpected chars '%v'", value)
	}
}
func Test_random_text_example(t *testing.T) {
	value, _ := security.RandomStr(security.TokenChars, 30)
	fmt.Println(value)
}

package rotavator_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"rotorvator/rotavator"
)

func Test_retry(t *testing.T) {
	t.Run("non error case", func(t *testing.T) {
		i := 0
		err := rotavator.Retry(context.Background(), 3, func() error {
			i++
			return nil
		})

		assert.NoError(t, err)
		assert.Equal(t, 1, i)
	})
	t.Run("error every time", func(t *testing.T) {
		count := 0
		failure := fmt.Errorf("fail")
		err := rotavator.Retry(context.Background(), 3, func() error {
			count++
			return failure
		})

		assert.IsType(t, err, rotavator.MaxAttemptsErr{})
		assert.ErrorIs(t, err, failure)
		assert.Equal(t, 3, count)
	})
	t.Run("context timeout before complete", func(t *testing.T) {
		ret := []error{fmt.Errorf("fail 1"), fmt.Errorf("fail 2"), fmt.Errorf("fail 3")}
		i := -1
		count := 0
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		err := rotavator.Retry(ctx, 3, func() error {
			count++
			i++
			return ret[i]
		})

		assert.ErrorIs(t, err, context.DeadlineExceeded)
		assert.Equal(t, 2, count)
	})
}

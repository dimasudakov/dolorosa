package utils

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_WithTimeout(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	f := func() (int, error) {
		time.Sleep(time.Millisecond * 100)
		return 1, nil
	}

	res, err := WithTimeout(ctx, f, time.Millisecond*500)
	assert.NoError(t, err)
	assert.Equal(t, 1, res)

	f2 := func() (int, error) {
		time.Sleep(time.Millisecond * 200)
		return 4, nil
	}

	res2, err2 := WithTimeout(ctx, f2, time.Millisecond*100)
	assert.EqualError(t, err2, ErrTimeoutExceeded.Error())
	assert.Equal(t, 0, res2)
}

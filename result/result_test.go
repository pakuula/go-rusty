package result_test

import (
	"errors"
	"testing"

	"github.com/pakuula/go-rusty/result"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test Result
type TR = result.Result[int]

var ValTR = result.Val[int]
var ErrTR = result.Err[int]
var errTest = errors.New("test error")

func TestCTOR(t *testing.T) {
	{
		v := ValTR(1)
		assert.False(t, v.IsError())
		assert.True(t, v.IsValue())
		require.NotPanics(t, func() { v.Unwrap() })
		assert.Equal(t, v.Unwrap(), 1)
		assert.Panics(t, func() { v.Err() })
		assert.Equal(t, errTest, v.ErrOr(errTest))
	}
	{

		e := ErrTR(errTest)
		assert.True(t, e.IsError())
		assert.False(t, e.IsValue())
		require.NotPanics(t, func() { e.Err() })
		assert.Equal(t, errTest, e.Err())
		assert.Panics(t, func() { e.Unwrap() })
		assert.Equal(t, 5, e.UnwrapOr(5))
	}

}

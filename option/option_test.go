package option_test

import (
	"errors"
	"testing"

	"github.com/pakuula/go-rusty/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test Result
type TR = option.Option[int]

var SomeTR = option.Some[int]
var NoneTR = option.None[int]
var errTest = errors.New("test error")

func TestCTOR(t *testing.T) {
	{
		some := SomeTR(1)
		assert.False(t, some.IsNone())
		assert.True(t, some.IsSome())
		require.NotPanics(t, func() { some.Unwrap() })
		assert.Equal(t, some.Unwrap(), 1)
	}
	{

		none := NoneTR()
		assert.True(t, none.IsNone())
		assert.False(t, none.IsSome())
		assert.Panics(t, func() { none.Unwrap() })
		assert.Equal(t, 5, none.UnwrapOr(5))
	}

}

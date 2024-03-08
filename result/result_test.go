// Copyright 2024 Nikolay Pakulin (@pakuula). All rights reserved.
// Use of this source code is governed by LGPL-3.0 licence.
// The text of the licence can be found in the LICENSE.txt file.

package result_test

import (
	"errors"
	"fmt"
	"math/rand"
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

func BenchmarkCatch(b *testing.B) {
	sample := func() (res result.Result[[]int]) {
		defer result.Catch(&res)

		slice := []int{}
		for i := 0; i < 1000; i++ {
			r := rand.Int()
			slice = append(slice, r)
		}
		result.Err[[]int](fmt.Errorf("some error")).Must()
		return result.Val(slice)
	}
	for i := 0; i < b.N; i++ {
		sample()
	}
}

func BenchmarkReturn(b *testing.B) {
	sample := func() ([]int, error) {
		slice := []int{}
		for i := 0; i < 1000; i++ {
			r := rand.Int()
			slice = append(slice, r)
		}
		return slice, fmt.Errorf("some error")
	}
	for i := 0; i < b.N; i++ {
		sample()
	}
}

// Copyright 2024 Nikolay Pakulin (@pakuula). All rights reserved.
// Use of this source code is governed by LGPL-3.0 licence.
// The text of the licence can be found in the LICENSE.txt file.

package expr

import "github.com/pakuula/go-rusty/result"

func Default[T any]() T {
	var zero T
	return zero
}

func If[T any](cond bool, Then T, Else T) T {
	if cond {
		return Then
	} else {
		return Else
	}
}

func IfLazy[T any](cond bool, Then func() T, Else func() T) T {
	if cond {
		return Then()
	} else {
		return Else()
	}
}

func IfLazyE[T any](cond bool, Then func() (T, error), Else func() (T, error)) (T, error) {
	if cond {
		return Then()
	} else {
		return Else()
	}
}

func IfLazyR[T any](cond bool, Then func() (T, error), Else func() (T, error)) result.Result[T] {
	if cond {
		return result.Wrap(Then())
	} else {
		return result.Wrap(Else())
	}
}

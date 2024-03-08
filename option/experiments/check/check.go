// Copyright 2024 Nikolay Pakulin (@pakuula). All rights reserved.
// Use of this source code is governed by LGPL-3.0 licence.
// The text of the licence can be found in the LICENSE.txt file.

package main

import (
	"fmt"

	"github.com/pakuula/go-rusty/option"
)

func main() {
	fmt.Println(MustReturn())
	fmt.Println(MustPanic())
}

func MustReturn() (res option.Option[int]) {
	defer option.Catch(&res)
	val := option.None[int]().Must()
	return option.Some(val)
}

func MustPanic() (res option.Option[int]) {
	defer option.Catch(&res)
	val := option.None[int]().Unwrap()
	return option.Some(val)
}

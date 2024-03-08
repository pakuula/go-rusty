// Copyright 2024 Nikolay Pakulin (@pakuula). All rights reserved.
// Use of this source code is governed by LGPL-3.0 licence.
// The text of the licence can be found in the LICENSE.txt file.

package main

import "github.com/pakuula/go-rusty/option"

func AddKeyValues(m map[string]int, key1, key2 string) (int, bool) {
	v1, ok1 := m[key1]
	v2, ok2 := m[key2]
	if ok1 && ok2 {
		return (v1 + v2), true
	} else {
		return 0, false
	}
}

func AddKeyValuesOpt(m map[string]int, key1, key2 string) (res option.Option[int]) {
	defer option.Catch(&res)
	v1 := option.MapGet(m, key1).Must()
	v2 := option.MapGet(m, key2).Must()
	return option.Some(v1 + v2)
}

func main() {

}

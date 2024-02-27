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

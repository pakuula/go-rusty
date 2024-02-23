package main

import (
	"errors"
	"fmt"

	"github.com/pakuula/go-rusty/result"
)

func main() {
	fmt.Println(MustReturn())
	fmt.Println(MustError())
	fmt.Println(MustPanic())
}

func MustError() (err error) {
	defer result.CatchError(&err)
	result.Err[int](errors.New("must be an error")).Must()
	return
}

func MustReturn() (res result.Result[int]) {
	defer result.Catch(&res)
	val := result.Err[int](errors.New("must be an error")).Must()
	return result.Val(val)
}

func MustPanic() (res result.Result[int]) {
	defer result.Catch(&res)
	val := result.Err[int](errors.New("must be an error")).Unwrap()
	return result.Val(val)
}

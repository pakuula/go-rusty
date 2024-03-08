// Copyright 2024 Nikolay Pakulin (@pakuula). All rights reserved.
// Use of this source code is governed by LGPL-3.0 licence.
// The text of the licence can be found in the LICENSE.txt file.

package result

import (
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// A value or an error
type Result[T any] struct {
	value T
	err   error
}

// Constructors

func Wrap[T any](val T, err error) Result[T] {
	return Result[T]{
		value: val,
		err:   err,
	}
}

func Val[T any](value T) Result[T] {
	return Result[T]{value: value}
}

func Err[T any](err error) Result[T] {
	if err == nil {
		panic("Not an error")
	}
	return Result[T]{
		err: err,
	}
}

type _Void struct{}

type ResultVoid = Result[_Void]

func Void(err error) Result[_Void] {
	return Wrap(_Void{}, err)
}

// String representaion

type hasToString interface {
	ToString() string
}

// Builds a string representation of the result object.
// If it is an error, returns the err.Error() string
// If T has String method, calls String
// If T has ToString method, calls ToString
// If T has MarshalText method, calls MarshalText
// Otherwise calls fmt.Sprint(self.Unwrap())
func (self Result[T]) String() string {
	if self.IsError() {
		return "error: " + self.Err().Error()
	}
	reflected := reflect.ValueOf(self.value).Interface()
	{
		z, ok := reflected.(hasToString)
		if ok {
			return z.ToString()
		}
	}
	{
		z, ok := reflected.(fmt.Stringer)
		if ok {
			return z.String()
		}
	}
	{
		z, ok := reflected.(encoding.TextMarshaler)
		if ok {
			strBytes, err := z.MarshalText()
			if err == nil {
				return string(strBytes)
			}
		}
	}
	{
		z, ok := reflected.(json.Marshaler)
		if ok {
			strBytes, err := z.MarshalJSON()
			if err == nil {
				return string(strBytes)
			}
		}
	}
	return fmt.Sprint(self.value)
}

// Check the result

// True if self is an error
func (self Result[T]) IsError() bool {
	return self.err != nil
}

// True if self is an error and it matches the condition
func (self Result[T]) IsErrorAnd(cond func(error) bool) bool {
	return self.err != nil && cond(self.err)
}

// True is self contains a value
func (self Result[T]) IsValue() bool {
	return self.err == nil
}

// True is self contains a value and it matches the condition
func (self Result[T]) IsValueAnd(cond func(T) bool) bool {
	return self.err == nil && cond(self.value)
}

type checkError struct {
	err error
}

// Defer Catch(&res) to convert failed invocation of Must into Result
func Catch[T any](res *Result[T]) {
	if panicValue := recover(); panicValue != nil {
		err, ok := panicValue.(checkError)
		if ok {
			*res = Err[T](err.err)
		} else {
			panic(panicValue)
		}
	}
}

// Defer Catch(&err) to convert failed invocation of Must into error
func CatchError(errorPtr *error) {
	if panicValue := recover(); panicValue != nil {
		err, ok := panicValue.(checkError)
		if ok {
			*errorPtr = err.err
		} else {
			panic(panicValue)
		}
	}
}

// Extracting the stored value

// Extracts the stored value or panics with a catchable value.
func (self Result[T]) Must() T {
	if self.IsError() {
		panic(checkError{self.err})
	}
	return self.value
}

// Extracts the stored value or panics with a catchable string.
//
// The string is appended with 'error: <error's Error() method result>'
func (self Result[T]) Mustf(format string, a ...any) T {
	if self.IsError() {
		msg := fmt.Sprintf(format, a...)
		if msg != "" {
			msg += ": "
		}
		msg += "error: " + self.err.Error()
		panic(checkError{errors.New(msg)})
	}
	return self.value
}

// Returns the value of panics with the catchable value.
func Must[T any](val T, err error) T {
	if err != nil {
		panic(checkError{err})
	}
	return val
}

// In the case of error panics with the catchable value.
func NoError(err error) ResultVoid {
	if err != nil {
		panic(checkError{err})
	}
	return Void(nil)
}

// Returns the stored value or panics with the given message
func (self Result[T]) Expect(msg string) T {
	if self.IsError() {
		panic(msg)
	}
	return self.value
}

// Returns the stored value or panics with the given formatted message
func (self Result[T]) Expectf(format string, a ...any) T {
	if self.IsError() {
		panic(fmt.Sprintf(format, a...))
	}
	return self.value
}

// Returns the stored value or panics
func (self Result[T]) Unwrap() T {
	if self.IsError() {
		panic(fmt.Errorf("unwrap error: %w", self.err))
	}
	return self.value
}

// Returns the stored value without checking the error
func (self Result[T]) UnwrapUnsafe() T {
	return self.value
}

// Returns the stored value or the provided default value
func (self Result[T]) UnwrapOr(valueIfError T) T {
	if self.IsError() {
		return valueIfError
	}
	return self.value
}

// Returns the stored value or the default value for the type T
func (self Result[T]) UnwrapOrDefault() T {
	if self.IsError() {
		var zero T
		return zero
	}
	return self.value
}

// Returns the stored value or transforms the error
func (self Result[T]) UnwrapOrElse(f func(error) T) T {
	if self.IsError() {
		return f(self.err)
	}
	return self.value
}

// Converts to the pair (value, error)
func (self Result[T]) UnwrapWithError() (T, error) {
	return self.value, self.err
}

// Accessing the error

// Returns the error or panics
func (self Result[T]) Err() error {
	if self.IsValue() {
		panic("not an error")
	}
	return self.err
}

// Returns the stored error or the provided error
func (self Result[T]) ErrOr(other error) error {
	if self.IsValue() {
		return other
	}
	return self.err
}

// Returns the error or panics with the given message
func (self Result[T]) ExpectError(msg string) error {
	if self.IsValue() {
		panic(msg)
	}
	return self.err
}

// Returns the error or panics with the given formatted message
func (self Result[T]) ExpectErrorf(format string, a ...any) error {
	if self.IsValue() {
		panic(fmt.Sprintf(format, a...))
	}
	return self.err
}

// Pointer and dereference

// Converts to Result[*T]
func Ptr[T any](res *Result[T]) Result[*T] {
	if res.IsError() {
		return Err[*T](res.err)
	}
	return Val(&res.value)
}

var ErrDerefNil = errors.New("derefencing nil")

// Derefences Result[*T] - produces Result[T]
func Deref[T any](res Result[*T]) Result[T] {
	if res.IsError() {
		return Err[T](res.err)
	}
	if res.value == nil {
		return Err[T](ErrDerefNil)
	}
	return Val(*res.value)
}

// Transform the Result

// Applies f to the stored value or keeps the error unchanged
func Apply[T any, U any](from Result[T], f func(T) U) Result[U] {
	if from.IsError() {
		return Err[U](from.err)
	}
	return Val(f(from.value))
}

// Applies f to the stored value or keeps error unchanged.
// If f returns an error, set the error
func ApplyE[T any, U any](from Result[T], f func(T) (U, error)) Result[U] {
	if from.IsError() {
		return Err[U](from.err)
	}
	res, err := f(from.value)
	return Wrap[U](res, err)
}

// Applies f to the stored value or keeps error unchanged.
// If f returns an error, set the error
func ApplyResult[T any, U any](from Result[T], f func(T) Result[U]) Result[U] {
	if from.IsError() {
		return Err[U](from.err)
	}
	return f(from.value)
}

// Utility functions

// Unmarshal JSON data into a value of the type T or error
func UnmarshalJson[T any](data []byte) Result[T] {
	var result T
	err := json.Unmarshal(data, &result)
	return Wrap(result, err)
}

// Unmarshal JSON string into a value of the type T or error
func UnmarshalJsonString[T any](data string) Result[T] {
	return UnmarshalJson[T]([]byte(data))
}

// Iterator

func MapE[T any, U any](slice []T, f func(T) (U, error)) []Result[U] {
	var retval = make([]Result[U], len(slice))

	for i, t := range slice {
		retval[i] = Wrap[U](f(t))
	}
	return retval
}

func MapR[T any, U any](slice []T, f func(T) Result[U]) []Result[U] {
	var retval = make([]Result[U], len(slice))

	for i, t := range slice {
		retval[i] = f(t)
	}
	return retval
}

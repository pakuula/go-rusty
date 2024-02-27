package option

import (
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// A value or None
type Option[T any] struct {
	value T
	none  bool
}

// Constructors
func Some[T any](value T) Option[T] {
	return Option[T]{value: value}
}

func None[T any]() Option[T] {
	return Option[T]{
		none: true,
	}
}

func WrapOk[T any](val T, ok bool) Option[T] {
	if !ok {
		return None[T]()
	}
	return Some(val)
}

func WrapErr[T any](val T, err error) Option[T] {
	if err != nil {
		return None[T]()
	}
	return Some(val)
}

// String representaion

type hasToString interface {
	ToString() string
}

// Builds a string representation of the Option object.
// If it is None, returns "None"
// If T has String method, calls String
// If T has ToString method, calls ToString
// If T has MarshalText method, calls MarshalText
// Otherwise calls fmt.Sprint(self.Unwrap())
func (self Option[T]) String() string {
	if self.IsNone() {
		return "<None>"
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

// Check the Option

// True if self is None
func (self Option[T]) IsNone() bool {
	return self.none
}

// True is self contains a value
func (self Option[T]) IsSome() bool {
	return !self.none
}

// True is self contains a value and it matches the condition
func (self Option[T]) IsSomeAnd(cond func(T) bool) bool {
	return !self.none && cond(self.value)
}

type errUnwrapNone struct{}

// Defer Catch(&res) to convert failed invocation of Must into Option
func Catch[T any](res *Option[T]) {
	if panicValue := recover(); panicValue != nil {
		_, ok := panicValue.(errUnwrapNone)
		if ok {
			*res = None[T]()
		} else {
			panic(panicValue)
		}
	}
}

// Extracting the stored value

// Extracts the stored value or panics with a catchable value.
func (self Option[T]) Must() T {
	if self.IsNone() {
		panic(errUnwrapNone{})
	}
	return self.value
}

// Returns the value of panics with the catchable value.
func MustOk[T any](val T, ok bool) T {
	if !ok {
		panic(errUnwrapNone{})
	}
	return val
}

// Extracting the stored value

// Returns the stored value or panics with the given message
func (self Option[T]) Expect(msg string) T {
	if self.IsNone() {
		panic(msg)
	}
	return self.value
}

// Returns the stored value or panics with the given formatted message
func (self Option[T]) Expectf(format string, a ...any) T {
	if self.IsNone() {
		panic(fmt.Sprintf(format, a...))
	}
	return self.value
}

var ErrUnwrap = errors.New("unwrapping none")

// Returns the stored value or panics
func (self Option[T]) Unwrap() T {
	if self.IsNone() {
		panic(ErrUnwrap)
	}
	return self.value
}

// Returns the stored value without checking IsNone()
func (self Option[T]) UnwrapUnsafe() T {
	return self.value
}

// Returns the stored value or the provided default value
func (self Option[T]) UnwrapOr(valueIfNone T) T {
	if self.IsNone() {
		return valueIfNone
	}
	return self.value
}

// Returns the stored value or the default value for the type T
func (self Option[T]) UnwrapOrDefault() T {
	if self.IsNone() {
		var zero T
		return zero
	}
	return self.value
}

// Converts to the pair (value, bool)
func (self Option[T]) UnwrapWithOk() (T, bool) {
	return self.value, !self.none
}

// Accessing the error

// Pointer and dereference

// Converts to Option[*T]
func Ptr[T any](res *Option[T]) Option[*T] {
	if res.IsNone() {
		return None[*T]()
	}
	return Some(&res.value)
}

// Derefences Option[*T] - produces Option[T]
func Deref[T any](res Option[*T]) Option[T] {
	if res.IsNone() {
		return None[T]()
	}
	if res.value == nil {
		return None[T]()
	}
	return Some(*res.value)
}

// Transform the Option

// Applies f to the stored value or keeps the error unchanged
func Apply[T any, U any](from Option[T], f func(T) U) Option[U] {
	if from.IsNone() {
		return None[U]()
	}
	return Some(f(from.value))
}

// Applies f to the stored value or keeps error unchanged.
// If f returns None, set the error
func ApplyE[T any, U any](from Option[T], f func(T) (U, error)) Option[U] {
	if from.IsNone() {
		return None[U]()
	}
	res, err := f(from.value)
	return WrapErr[U](res, err)
}

// Applies f to the stored value or keeps error unchanged.
// If f returns None, set the error
func ApplyOption[T any, U any](from Option[T], f func(T) Option[U]) Option[U] {
	if from.IsNone() {
		return None[U]()
	}
	return f(from.value)
}

// Utility functions

// Unmarshal JSON data into a value of the type T or error
func UnmarshalJson[T any](data []byte) Option[T] {
	var Option T
	err := json.Unmarshal(data, &Option)
	return WrapErr(Option, err)
}

// Unmarshal JSON string into a value of the type T or error
func UnmarshalJsonString[T any](data string) Option[T] {
	return UnmarshalJson[T]([]byte(data))
}

// Map access

func MapGet[K comparable, V any](m map[K]V, key K) Option[V] {
	v, ok := m[key]
	return WrapOk(v, ok)
}

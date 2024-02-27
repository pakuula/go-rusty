<!-- vscode-markdown-toc -->
* 1. [Throw/Catch](#ThrowCatch)
* 2. [Constructors](#Constructors)
	* 2.1. [`WrapOk` (value,bool) pair](#WrapOkvalueboolpair)
	* 2.2. [Store a valueue](#Storeavalueue)
	* 2.3. [Store `None`](#StoreNone)
* 3. [Access to stored data](#Accesstostoreddata)
	* 3.1. [Access to value](#Accesstovalue)
	* 3.2. [Access to error](#Accesstoerror)
* 4. [Transforming contained values](#Transformingcontainedvalues)

<!-- vscode-markdown-toc-config
	numbering=true
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc -->A simple implementation of `Option[T]` that mimicks [`Option<T>` from Rust](https://doc.rust-lang.org/std/option/).


# `Option[T]`

The type `Option[T]` mimicks [`Option<T>`](https://doc.rust-lang.org/std/option/) Rust idiom. It is a union type that is either a value of the type `T` or `None`.

The main use of `Option[T]` is to provide flag if the operation failed, without details.
```go
func AddKeyValues(m map[string]int, key1, key2 string) (int, bool) {
	v1, ok1 := m[key1]
	v2, ok2 := m[key2]
	if ok1 && ok2 {
		return (v1 + v2), true
	} else {
		return 0, false
	}
}
```

Here is how this function could be implemented using Option style.
```go
import "github.com/pakuula/go-rusty/option"

func AddKeyValuesOpt(m map[string]int, key1, key2 string) (res option.Option[int]) {
	defer option.Catch(&res)
	v1 := option.MapGet(m, key1).Must()
	v2 := option.MapGet(m, key2).Must()
	return option.Some(v1 + v2)
}
```

The function  `option.MapGet` return an `Option[T]` object that is either the value from the map or `None`.
`Must()` panics with an error if the object is `None`, while `defer option.Catch(&res)` catches missing
values and writes `None` to the returned option `res`.

The constructor `Option.Some[T](value T)` builds a `Option[T]` with the given value.

##  1. <a name='ThrowCatch'></a>Throw/Catch

The Rust operator [`?` (the question mark)](https://doc.rust-lang.org/std/option/#the-question-mark-operator-) 
is used to simplify code that uses `Option` types. Ending the expression with `?` will
result in the unwrapped value, unless the Option is `None`, in which case `None` is
*returned early* from the enclosing function.
```rust
fn add_last_numbers(stack: &mut Vec<i32>) -> Option<i32> {
    Some(stack.pop()? + stack.pop()?)
}
```
In this example execution is terminated as soon as any of `?`-enclosed calls returns an error.
It is similar to throwing and exception though `?` operator Options in `Option.Err` rather than
exception propagation through the stack.

The method `Option.Must` and the function `Catch` implement the alternative 
to the question mark oprator.
```go
func AddKeyValuesOpt(m map[string]int, key1, key2 string) (res option.Option[int]) {
	defer option.Catch(&res)
	v1 := option.MapGet(m, key1).Must()
	v2 := option.MapGet(m, key2).Must()
	return option.Some(v1 + v2)
}
```

The method `Option.Must` panics with an object of the internal type `Option.checkError`. 
The function `option.Catch[T](*Option[T])` recovers the panics and 
writes `None` to the Option to be returned.

```
func ReturnOption() (res Option[SomeType]) {
    defer option.Catch(&res);

    x := SomeOtherOption().Must()

    return ProcessX(x)
}
```

Important notes:
- the return value must be named, e.g. `(res Option[SomeType])`,
- the `Catch` must be called in the `defer` statement: `defer option.Catch(&res)`,
- if the reason of the panics is not `Must` then the error is not intercepted and 
  panic propagets through the stack.


##  2. <a name='Constructors'></a>Constructors

###  2.1. <a name='WrapOkvalueboolpair'></a>`WrapOk` (value,bool) pair

De-facto standard idiom of Go is to return a pair `(T, bool)` where the boolean value
is `true` if the operation was successfull. 

The function `Option.WrapOk` converts such a pair into an object `Option[T]`.
```
    fileInfoOption := Option.WrapOk(cache.GetFileInfo("filename.txt))
```

Properties: 
- if `ok` is `true`:
  - `Option.WrapOk(val, true).IsNone()` is `false`,
  - `Option.WrapOk(val, true).Unwrap()` never panics,
  - `Option.WrapOk(val, true).Unwrap()` returns `val`,
- if `ok` is `false`:
  - `Option.WrapOk(val, false).IsNone()` is `true`,
  - `Option.WrapOk(val, false).Unwrap()` panics,

###  2.2. <a name='Storeavalueue'></a>Store a valueue

The function  `Option.Some[T](val T) Option[T]` construct a non-none object.

- `Option.Some(val).IsNone()` is `false`,
- `Option.Some(val).Unwrap()` never panics,
- `Option.Some(val).Unwrap()` returns `val`.

###  2.3. <a name='StoreNone'></a>Store `None`

The function `func Err[T](err error) Option[T]` creates an error object.
- `Option.Err(err).IsNone()` is `false`,
- `Option.Err(err).Unwrap()` panics,
- `Option.Err(err).Err()` returns `None`.

##  3. <a name='Accesstostoreddata'></a>Access to stored data

These methods extract the contained value in a Option[T] when it is not the `None` variant.
If the Option is `None`:

* `Option.Unwrap()` panics with a generic message
* `Option.Expect(msg string)` panics with a provided custom message
* `Option.Expectf(format string, a ...any)` panics with a provided custom formatted message
* `Option.UnwrapOr(val T)` the provided custom value
* `Option.UnwrapOrDefault()` returns the default value of the type T

The method `Option.UnwrapWithOk() (T,bool)` converts the `Option[T]` object to the standard
Golang pair `(T, bool)` where the boolean value is `true` for value options.

##  4. <a name='Transformingcontainedvalues'></a>Transforming contained values

The following functions mimic `map` method of Rust `Option<T,E>` class.

- `func Apply[T any, U any](from Option[T], f func(T) U) Option[U]` 
  applies `f` to the stored value or keeps `None`.
- `func ApplyE[T any, U any](from Option[T], f func(T) (U, error)) Option[U]` 
  applies `f` to the stored value or keeps `None`. If `f` returns an error, 
  sets the Option to `None`.
- `func ApplyOption[T any, U any](from Option[T], f func(T) Option[U]) Option[U]` 
  applies `f` to the stored value or keeps `None` unchanged. If `f` returns `None`, returns `None`.


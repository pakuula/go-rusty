<!-- vscode-markdown-toc -->
* 1. [Throw/Catch](#ThrowCatch)
* 2. [Constructors](#Constructors)
	* 2.1. [`Wrap` (value,error) pair](#Wrapvalueerrorpair)
	* 2.2. [Store a valueue](#Storeavalueue)
	* 2.3. [Store an error](#Storeanerror)
* 3. [Access to stored data](#Accesstostoreddata)
	* 3.1. [Access to value](#Accesstovalue)
	* 3.2. [Access to error](#Accesstoerror)

<!-- vscode-markdown-toc-config
	numbering=true
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc --><!-- vscode-markdown-toc -->

# `Result[T]`

The type `Result[T]` mimicks [`Result<T,E>`](https://doc.rust-lang.org/std/result/) Rust idiom. It is a union type that is either a value of the type `T` or an error of the type [`error`](https://go.dev/ref/spec#Errors).

The main use of `Result[T]` is to avoid typing error checks. 
```
func ReadTextFile(fname string) (string, error) {
    f, err := os.Open(fname)
    if err != nil {
        return "", err
    }
    buf := bytes.Buffer{}
    _, err = io.Copy(&buf, f)
    if err != nil {
        return "", err
    }
    return buf.String(), nil
}
```

Here is how this function could be implemented using Result style.
```
import "github.com/pakuula/go-rusty/result"

func ReadTextFileR(fname string) (res result.Result[string]) {
	defer result.Catch(&res)
	f := result.Wrap(os.Open(fname)).Must()
	defer f.Close()

	buf := bytes.Buffer{}
	result.Must(io.Copy(&buf, f))

	return result.Val(buf.String())
}
```

The calls to `result.Wrap[T]` convert the pair of velues `(T, error)` into a `Result[T]` object. Method 
`Must()` panics with an error if the object is not a value, while `defer result.Catch(&res)` catches the error and writes is to the returned result `res`.

The function `result.Must(T, error)` is a shorthand for `result.Wrap(T, error).Must()`.

The constructor `result.Val[T](value T)` builds a `Result[T]` with the given `value`.

##  1. <a name='ThrowCatch'></a>Throw/Catch

The Rust operator [`?` (the question mark)](https://doc.rust-lang.org/std/result/#the-question-mark-operator-) 
is used to simplify code that uses `Result` types. Ending the expression with `?` will
result in the `Ok`’s unwrapped value, unless the result is `Err`, in which case `Err` is
*returned early* from the enclosing function.
```
fn write_info(info: &Info) -> io::Result<()> {
    let mut file = File::create("my_best_friends.txt")?;
    // Early return on error
    file.write_all(format!("name: {}\n", info.name).as_bytes())?;
    file.write_all(format!("age: {}\n", info.age).as_bytes())?;
    file.write_all(format!("rating: {}\n", info.rating).as_bytes())?;
    Ok(())
}
```
In this example execution is terminated as soon as any of `?`-enclosed calls returns an error.
It is similar to throwing and exception though `?` operator results in `Result.Err` rather than
exception propagation through the stack.

The method `Result.Must` and the functions `Catch` and `CatchError` implement the alternative 
to the question mark oprator.
```
fn WriteInfo(info Info) (res ResultVoid) {
    // Catch the error and write it into the result
    defer Catch(&res)
    // Early return on error
    file = os_r.Open("my_best_friends.txt").Must()
    io_r.WriteString(fmt.Sprintf("name: %s\n", info.name)).Must();
    io_r.WriteString(fmt.Sprintf("age: %s\n", info.age)).Must();
    io_r.WriteString(fmt.Sprintf("rating: %s\n", info.rating)).Must();

    return Void(nil)
}
```

The method `Result.Must` panics with an object of the internal type `result.checkError`. 
The functions `result.Catch[T](*Result[T])` and `result.CatchError(*error)` recover the panics and 
write the error to the result or error objects.

```
func ReturnResult() (res Result[SomeType]) {
    defer result.Catch(&res);

    x := SomeOtherResult().Must()

    return ProcessX(x)
}
```

Important notes:
- the return value must be named, e.g. `(res Result[SomeType])`,
- the `Catch` must be called in the `defer` statement: `defer result.Catch(&res)`,
- if the reason of the panics is not `Must` then the error is not intercepted and 
  panic propagets through the stack.

Similar idiom exists for functions that return `error`:
```
func ReturnError() (err error) {
    defer result.CatchError(&err);

    SomeOtherResult().Must()

    return nil
}
```
Similar notes apply:
- the return value must be named, e.g. `(err error)`,
- the `CatchError` must be called in the `defer` statement: `defer result.CatchError(&err)`,
- if the reason of the panics is not `Must` then the error is not intercepted and 
  panic propagets through the stack.


##  2. <a name='Constructors'></a>Constructors

###  2.1. <a name='Wrapvalueerrorpair'></a>`Wrap` (value,error) pair

De-facto standard idiom of Go is to return a pair `(T, error)`. 
The function `result.Wrap` converts such a pair into an object `Result[T]`.
```
    fileResult := result.Wrap(os.Open("filename.txt))
```

Properties:
- if `err` is `nil`:
  - `result.Wrap(val, nil).IsError()` is `false`,
  - `result.Wrap(val, nil).Unwrap()` never panics,
  - `result.Wrap(val, nil).Unwrap()` returns `val`,
- if `err` is not `nil`:
  - `result.Wrap(val, err).IsError()` is `true`,
  - `result.Wrap(val, err).Unwrap()` panics,
  - `result.Wrap(val, nil).Err()` returns `err`.

###  2.2. <a name='Storeavalueue'></a>Store a valueue

The function  `result.ValT(val T) Result[T]` construct a non-error object.

- `result.Val(val).IsError()` is `false`,
- `result.Val(val).Unwrap()` never panics,
- `result.Val(val).Unwrap()` returns `val`.

###  2.3. <a name='Storeanerror'></a>Store an error

The function `func Err[T](err error) Result[T]` creates an error object.
- `result.Err(err).IsError()` is `false`,
- `result.Err(err).Unwrap()` panics,
- `result.Err(err).Err()` returns `err`.

##  3. <a name='Accesstostoreddata'></a>Access to stored data

###  3.1. <a name='Accesstovalue'></a>Access to value

These methods extract the contained value in a Result[T] when it is not the error variant.
If the Result is an error:

* `Result.Unwrap()` panics with a generic message
* `Result.Expect(msg string)` panics with a provided custom message
* `Result.Expectf(format string, a ...any)` panics with a provided custom formatted message
* `Result.UnwrapOr(val T)` the provided custom value
* `Result.UnwrapOrDefault()` returns the default value of the type T
* `Result.UnwrapOrElse(f func(error)T)` returns the result of evaluating the provided function
  over the error.

The method `Result.UnwrapWithError() (T,error)` converts the `Result[T]` object to the standard
Golang pair `(T, error)`.

###  3.2. <a name='Accesstoerror'></a>Access to error

These methods extract the contained error in a `Result[T]`:

- `Result.Err() error` Returns the error or panics
- `Result.ErrOr(other error) error` Returns the stored error or the provided error and panics otherwise
- `Result.ExpectError(msg string) error` Returns the error or panics with the given message
- `Result.ExpectErrorf(format string, a ...any) error` Returns the error or panics with the 
  given formatted message

## Transforming contained values

The following functions mimic `map` method of Rust `Result<T,E>` class.

- `func Apply[T any, U any](from Result[T], f func(T) U) Result[U]` 
  applies `f` to the stored value or keeps the error unchanged.
- `func ApplyE[T any, U any](from Result[T], f func(T) (U, error)) Result[U]` 
  applies `f` to the stored value or keeps error unchanged. If `f` returns an error, 
  sets the result to that  error.
- `func ApplyResult[T any, U any](from Result[T], f func(T) Result[U]) Result[U]` 
  applies `f` to the stored value or keeps error unchanged. If `f` returns an error, set the error


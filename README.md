<!-- vscode-markdown-toc -->

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

[More about `Result[T]`](./result/README.md)


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

[More about `Option[T]`](./option/README.md)
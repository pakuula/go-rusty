package expr

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

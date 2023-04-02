package tuple

type Tuple[T any, U any] struct {
	First  T
	Second U
}

func New[T any, U any](first T, second U) Tuple[T, U] {
	return Tuple[T, U]{
		First:  first,
		Second: second,
	}
}

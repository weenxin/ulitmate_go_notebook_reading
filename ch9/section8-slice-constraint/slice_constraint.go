package section8_slice_constraint

type operateFunc[T any] func(item T) T

func Every[T any](items []T ,operator operateFunc[T]) []T {
	for index, item := range items {
		items[index] = operator(item)
	}
	return items
}


type Slice[T any] interface{
 ~[]T
}

func EverySlice[S Slice[T] , T any](s S, operator operateFunc[T]) S {
	for index, item := range s {
		s[index] = operator(item)
	}
	return s
}


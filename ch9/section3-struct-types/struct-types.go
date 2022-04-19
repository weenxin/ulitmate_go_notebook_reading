package section3_struct_types

type node[T any] struct {
	data T
	pre  *node[T]
	next *node[T]
}

type list[T any] struct {
	first *node[T]
	last  *node[T]
}

func (l *list[T]) add(data T) *node[T] {
	n := node[T]{
		data: data,
		pre:  l.last,
		next: nil,
	}
	if l.first == nil {
		l.first = &n
		l.last = &n
		return &n
	}
	l.last.next = &n
	l.last = &n
	return &n
}

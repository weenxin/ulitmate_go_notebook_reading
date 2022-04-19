package section2_underlying_types

import "errors"

type vector[T any] []T

func (v vector[T]) last() (T, error) {
	if len(v) > 0 {
		return v[len(v)-1], nil
	}
	//var zero T 先定义
	//return T{}, errors.New("Empty") //如果是string，int，bool等会有问题
	return *new(T), errors.New("Empty")
}

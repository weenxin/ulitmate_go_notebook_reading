package section5_type_as_contraint

type Addable interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | float32 | float64 | complex64 | complex128 | string
}

func Add[T Addable](v1 T, v2 T) T {
	return v1 + v2
}

type Person struct{
	name string
}
func (p Person) matcher(p2 Person) bool {
	return p.name == p2.name
}

type Other struct{
	name string
}
func (o Other) matcher(o2 Other) bool {
	return o.name == o2.name
}

type Food struct {
	name string
}

type matcher[T any] interface{
	Person | Food
	matcher(v T) bool
}

func match[T matcher[T]] (list []T , find T) int {
	for i, v := range  list {
		if v.matcher(find) {
			return i
		}
	}
	return -1
}

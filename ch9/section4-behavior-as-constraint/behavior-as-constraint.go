package section4_behavior_as_constraint

import "fmt"

type User struct {
	name string
}

func (u User) String() string {
	return u.name
}

type Stringer interface {
	String() string
}

func Concrete(u User) {
	u.String()
}

func Polymorphic(s Stringer) {
	s.String()
}


func stringfy[T fmt.Stringer](slice []T) []T {
	result := make([]T, len(slice))
	for index := range slice {
		result[index] = slice[index]
	}
	return result
}

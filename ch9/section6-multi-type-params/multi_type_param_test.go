package section6_multi_type_params

import "testing"

type User struct {
	name string
}

func (u User) String() string{
	return u.name
}

func TestPrint(t *testing.T) {
	t.Run("label is string ,value is integer", func(t *testing.T) {
		labels := []User{{"id"},{"age"}}
		values := []int{1,35}
		Print(labels,values)
	})
}

package section6_multi_type_params

import "fmt"

func Print[L fmt.Stringer , V any ](labels []L, values []V)  {
	if len(labels) != len(values) {
		panic("labels and values should be equal")
	}
	for i, v := range values {
		fmt.Printf("%s = %v" , labels[i], v )
	}
}

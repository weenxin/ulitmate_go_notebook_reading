package section1_basic

import "fmt"

func Print[T any](slice []T) {
	fmt.Print("Generic : ")
	for _, item := range slice {
		fmt.Printf("%v ", item)
	}
	fmt.Println()
}

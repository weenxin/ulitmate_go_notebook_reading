package search

import (
	"math/rand"
	"strings"
	"time"
)

const MessageLength = 100

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func mkItems() []Item {
	var results []Item
	for i :=0 ; i < 100; i++ {
		//results = append(results,Item{
		//	Title:   randStringRunes(MessageLength),
		//	Link:    randStringRunes(MessageLength),
		//	Description : randStringRunes(MessageLength),
		//})

		//优化后的代码
		results = append(results,Item{
			Title:   strings.ToLower(randStringRunes(MessageLength)) ,
			Link:    strings.ToLower(randStringRunes(MessageLength)),
			Description : strings.ToLower(randStringRunes(MessageLength)),
		})
	}
	return results
}



func init() {
	rand.Seed(time.Now().UnixNano())
}
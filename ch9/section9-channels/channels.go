package section9_channels

import (
	"context"
	"fmt"
	"sync"
)

type workFunc[Result any] func(ctx context.Context) Result

func DoWork[Result any](ctx context.Context, work workFunc[Result]) chan Result {
	ch :=  make(chan Result,1)
	go func() {
		ch <- work(ctx)
		fmt.Println("done work")
	}()
	return ch
}


type workInput[Input any , Result any ] func (input Input) Result

func PoolWork[Input any, Result any]  (size int , work workInput[Input , Result]) (chan Input, func()) {

	var wg sync.WaitGroup
	wg.Add(size)

	ch := make(chan Input)

	for i := 0 ; i < size ; i++ {
		go func() {
			defer wg.Done()
			for input := range ch {
				result := work(input)
				fmt.Println("pollWork :", result)
			}
		}()
	}

	cancel:= func() {
		close(ch)
		wg.Wait()
	}
	return ch,cancel
}



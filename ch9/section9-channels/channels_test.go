package section9_channels

import (
	"context"
	"math/rand"
	"testing"
	"time"
)

func work(ctx context.Context) int {
	time.Sleep(time.Duration(rand.Intn(200) )* time.Millisecond)
	return 100
}

func TestDoWork(t *testing.T) {

	duration := time.Duration(rand.Intn(150))  * time.Millisecond
	ctx , cancel := context.WithTimeout(context.Background(),duration)

	defer cancel()
	ch := DoWork(ctx,work)

	select{
	case data := <-ch:
		t.Logf("got data %v", data)
	case <-ctx.Done():
		t.Log("timeout ")
	}

}

func double(value int) int {
	return value *2
}

func TestPoolingWork(t *testing.T) {

	ch , cancel := PoolWork(4,double)

	for index := 0 ; index < 100; index ++ {
		ch<-index
	}

	cancel()

}

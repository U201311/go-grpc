package main

import (
	"context"
	"fmt"
	"sync"
)

/*
素筛数
*/
// 生成自然数序列
func GenerateNatural(ctx context.Context, wg *sync.WaitGroup) chan int {
	ch := make(chan int)
	go func() {
		defer wg.Done()
		defer close(ch)
		for i := 2; ; i++ {
			select {
			case <-ctx.Done():
				return
			case ch <- i:

			}
		}
	}()
	return ch
}

func PrimeFilter(ctx context.Context, in <-chan int, prime int, wg *sync.WaitGroup) chan int {
	out := make(chan int)
	go func() {
		for {
			defer wg.Done()
			defer close(out)
			for i := range in {
				if i%prime != 0 {
					select {
					case <-ctx.Done():
						return
					case out <- i:
					}
				}
			}
		}
	}()
	return out

}

func main() {
	wg := sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	ch := GenerateNatural(ctx, &wg)
	for i := 0; i < 100; i++ {
		prime := <-ch
		fmt.Printf("%v: %v\n", i+1, prime)
		wg.Add(1)
		ch = PrimeFilter(ctx, ch, prime, &wg)
	}

	cancel()
	wg.Wait()

}

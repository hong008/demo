package main

import "fmt"

//go 子父goroutine同步

//1.子goroutine通知父goroutine
func count(ch chan struct{}) {
	defer close(ch)
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}
}

func main() {
	ch := make(chan struct{})

	go count(ch)

	<-ch
	fmt.Println("Done...")
}

//2.子goroutine和父goroutine之间通过waitgroup通信
/*func count(group *sync.WaitGroup) {
	defer group.Done()
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}
}

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go count1(wg)
	wg.Wait()
	fmt.Println("Done...")
}*/

//3.使用context中的cancel
/*func count(ctx context.Context) chan int {
	ch := make(chan int)
	n := 0
	go func() {
		for {
			select {
			case <-ctx.Done():
				break
			case ch <- n:
				n++
			}
		}
	}()
	return ch
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := count(ctx)

	for n := range ch {
		fmt.Println(n)
		if n == 9 {
			break
		}
	}
	fmt.Println("Done...")
}*/

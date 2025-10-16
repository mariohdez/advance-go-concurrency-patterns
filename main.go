package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

func main() {
	inputCh := make(chan int)
	doneCh := make(chan struct{})

	var pWG sync.WaitGroup
	for i := 0; i < 10; i++ {
		pWG.Add(1)
		go func() {
			defer pWG.Done()
			producer(i, inputCh, doneCh)
		}()
	}

	var cWG sync.WaitGroup
	for i := 0; i < 10; i++ {
		cWG.Add(1)
		go func() {
			defer cWG.Done()
			consumer(i, inputCh)
		}()
	}

	time.Sleep(3 * time.Second)
	close(doneCh)

	pWG.Wait()
	close(inputCh)

	cWG.Wait()

	fmt.Println("good bye!")
}

func producer(index int, inputCh chan<- int, doneCh <-chan struct{}) {
	for {
		select {
		case <-doneCh:
			return
		default:
			// no-op
		}

		v := rand.IntN(10)

		time.Sleep(time.Millisecond * time.Duration(rand.IntN(1000)))

		select {
		case inputCh <- v:
		case <-doneCh:
			return
		}

		fmt.Printf("Producer[%d]: sent=%d\n", index, v)
	}
}

func consumer(index int, inputCh <-chan int) {
	for v := range inputCh {
		fmt.Printf("Consumer[%d]: recieved=%d\n", index, v)
	}
}

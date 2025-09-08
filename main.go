package main

import (
	"advance-go-concurrency-patterns/philosopher"
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	dinner := philosopher.NewDinner(5)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	dinner.Start(ctx, &wg)

	time.Sleep(5 * time.Second)
	fmt.Println("cancelling the simulation")
	cancel()

	wg.Wait()

	fmt.Println("see ya!")
}

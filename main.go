package main

import (
	"advance-go-concurrency-patterns/philosopher"
	"context"
	"fmt"
	"time"
)

func main() {
	dinner := philosopher.NewDinner(5)

	ctx, cancel := context.WithCancel(context.Background())

	dinner.Start(ctx)

	time.Sleep(10 * time.Second)
	fmt.Println("cancelling the simulation")
	cancel()

	dinner.Wait()

	fmt.Println("see ya!")
}

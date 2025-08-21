package main

import "advance-go-concurrency-patterns/philosopher"

func main() {
	dinner := philosopher.NewDinner(5, 4)
	dinner.Start()
}

package main

import "advance-go-concurrency-patterns/philosopher"

func main() {
	dinner := philosopher.NewDinner(5)
	dinner.Start()
}

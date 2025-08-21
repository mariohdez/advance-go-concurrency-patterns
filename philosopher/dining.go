package philosopher

import "fmt"

type Dinner struct {
	NumPhilosophers int
	NumForks        int
}

func NewDinner(numPhilosophers int, numForks int) *Dinner {
	return &Dinner{
		NumPhilosophers: numPhilosophers,
		NumForks:        numForks,
	}
}

func (d *Dinner) Start() {
	fmt.Printf(
		"starting dinner with %v philosopohers and %v forks.\n",
		d.NumPhilosophers,
		d.NumForks,
	)
}

package philosopher

import "fmt"

type Dinner struct {
	numPhilosophers int
	numForks        int
}

func NewDinner(numPhilosophers int) *Dinner {
	// There are always the same number of philospoher as there are forks.
	return &Dinner{
		numPhilosophers: numPhilosophers,
		numForks:        numPhilosophers,
	}
}

func (d *Dinner) Start() {
	fmt.Printf(
		"starting dinner with %v philosopohers and %v forks.\n",
		d.numPhilosophers,
		d.numForks,
	)
}

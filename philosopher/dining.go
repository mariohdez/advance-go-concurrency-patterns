package philosopher

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type state string

const (
	pendingState  state = "PENDING"
	eatingState   state = "EATING"
	thinkingState state = "THINKING"
	finishedState state = "FINISHED"
)

type Dinner struct {
	philosopherCount    int
	philosopherStates   []state
	philosopherStatesMu sync.Mutex
	forks               []sync.Mutex

	wg sync.WaitGroup
}

func NewDinner(numPhilosophers int) *Dinner {
	pStates := make([]state, 0, numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		pStates = append(pStates, pendingState)
	}

	return &Dinner{
		philosopherCount:  numPhilosophers,
		forks:             make([]sync.Mutex, numPhilosophers, numPhilosophers),
		philosopherStates: pStates,
	}
}

func (d *Dinner) Start(ctx context.Context) {
	for i := 0; i < d.philosopherCount; i++ {
		d.wg.Add(1)
		go d.philosopherWorker(ctx, i)
	}

	d.wg.Add(1)
	go d.monitor(ctx)
}

func (d *Dinner) Wait() {
	d.wg.Wait()
}

func (d *Dinner) philosopherWorker(ctx context.Context, pID int) {
	defer d.wg.Done()

	for {
		select {
		case <-ctx.Done():
			d.philosopherStatesMu.Lock()
			d.philosopherStates[pID] = finishedState
			d.philosopherStatesMu.Unlock()
			return
		default:
			d.think(ctx, pID)
			leftFork := &d.forks[pID]
			rightFork := &d.forks[(pID-1+d.philosopherCount)%d.philosopherCount]

			if pID == d.philosopherCount-1 {
				// prevent the case where all the philosophers grab their left fork successfully
				// and now the system is in dead-lock.
				d.eat(ctx, pID, rightFork, leftFork)
			} else {
				d.eat(ctx, pID, leftFork, rightFork)
			}
		}
	}
}

func (d *Dinner) think(ctx context.Context, pID int) {
	d.philosopherStatesMu.Lock()
	d.philosopherStates[pID] = thinkingState
	d.philosopherStatesMu.Unlock()

	select {
	case <-time.After(2 * time.Second):
	case <-ctx.Done():
	}
}

func (d *Dinner) eat(ctx context.Context, pID int, lf, rf *sync.Mutex) {
	lf.Lock()
	defer lf.Unlock()

	rf.Lock()
	defer rf.Unlock()

	d.philosopherStatesMu.Lock()
	d.philosopherStates[pID] = eatingState
	d.philosopherStatesMu.Unlock()

	select {
	case <-time.After(2 * time.Second):
	case <-ctx.Done():
	}
}

func (d *Dinner) monitor(ctx context.Context) {
	defer d.wg.Done()
	// Ansi escape codes to clear the screen and move the cursor to the top-left
	const clearScreen = "\033[H\033[2J"

	// ASCII art for the title
	asciiArt := `
.-----------------------------------.
| The Dining Philosophers Problem   |
'-----------------------------------'
`

	for {
		select {
		case <-ctx.Done():
			fmt.Print(clearScreen)
			return
		default:
			fmt.Print(clearScreen)
			fmt.Print(asciiArt)
			d.philosopherStatesMu.Lock()
			for i := 0; i < d.philosopherCount; i++ {
				fmt.Printf("p[%v] is %8v\t", i, d.philosopherStates[i])
			}
			fmt.Println()
			d.philosopherStatesMu.Unlock()

		}

		select {
		case <-time.After(time.Second):
		case <-ctx.Done():
			return
		}
	}

}

/*
func philosopherSimulation() {
	dinner := philosopher.NewDinner(5)

	ctx, cancel := context.WithCancel(context.Background())

	dinner.Start(ctx)

	time.Sleep(10 * time.Second)
	cancel()

	dinner.Wait()

	fmt.Println("see ya!")
}
*/

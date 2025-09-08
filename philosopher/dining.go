package philosopher

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Dinner struct {
	pCount int
	forks  []sync.Mutex
	wg     sync.WaitGroup
}

func NewDinner(numPhilosophers int) *Dinner {
	return &Dinner{
		pCount: numPhilosophers,
		forks:  make([]sync.Mutex, numPhilosophers, numPhilosophers),
	}
}

func (d *Dinner) Start(ctx context.Context) {
	for i := 0; i < d.pCount; i++ {
		d.wg.Add(1)
		go d.philosopherWorker(ctx, i)
	}
}

func (d *Dinner) Wait() {
	d.wg.Wait()
}

func (d *Dinner) philosopherWorker(ctx context.Context, pID int) {
	defer d.wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("philosopher=%v worker shutting down\n", pID)
			return
		default:
			d.think(pID)
			leftFork := &d.forks[pID]
			rightFork := &d.forks[(pID-1+d.pCount)%d.pCount]

			if pID == d.pCount-1 {
				// prevent the case where all the philosophers grab their left fork successfully
				// and now the system is in dead-lock.
				d.eat(pID, rightFork, leftFork)
			} else {
				d.eat(pID, leftFork, rightFork)
			}
		}
	}
}

func (d *Dinner) think(pID int) {
	fmt.Printf("philosopher=%v thinking.\n", pID)
	time.Sleep(2 * time.Second)
	fmt.Printf("philosopher=%v done thinking.\n", pID)
}

func (d *Dinner) eat(pID int, lf, rf *sync.Mutex) {
	lf.Lock()
	defer lf.Unlock()

	rf.Lock()
	defer rf.Unlock()

	fmt.Printf("philosopher=%v eating.\n", pID)
	time.Sleep(1 * time.Second)
	fmt.Printf("philosopher=%v done eating.\n", pID)
}

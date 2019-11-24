package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Philosopher struct contains name, fork channel, neighbor philosopher
// fork channel is used for demonstrate if the fork of the philosopher available
type Philosopher struct {
	name      string
	fork chan bool
	neighbor  *Philosopher
}

// Must have this function to declare an object
func makePhilosopher(name string, neighbor *Philosopher) *Philosopher {
	phil := &Philosopher{name, make(chan bool, 1), neighbor}
	phil.fork <- true
	return phil
}

// A function simulating thinking
func (phil *Philosopher) think() {
	fmt.Printf("%v is thinking.\n", phil.name)
	time.Sleep(time.Duration(rand.Int63n(1e9)))
}

// A function simulating eating
func (phil *Philosopher) eat() {
	fmt.Printf("%v is eating. \n", phil.name)
	time.Sleep(time.Duration(rand.Int63n(1e9)))
}

func (phil *Philosopher) getForks() {
	// Declare a channal indicating timeout
	timeout := make(chan bool, 1)
	go func() { time.Sleep(1e9); timeout <- true }()
	// taking the fork; fork is not available
	<-phil.fork
	fmt.Printf("%v got his fork. \n", phil.name)
	select {
	// if the neighbor's fork is availble, return, ready to eat
	case <-phil.neighbor.fork:
		fmt.Printf("%v got %v's fork.\n", phil.name, phil.neighbor.name)
		fmt.Printf("%v has two fork.\n", phil.name)
		return
	// after amount of time the philosopher taking up his own fork, the
	// philosopher has to put the fork down letting others to use
	// then think for a while and try get fork operation again
	case <-timeout:
		phil.fork <- true
		phil.think()
		phil.getForks()
	}
}

func (phil *Philosopher) returnForks() {
	// after a philosopher finish eating, making his fork channel
	// and his neighbor's fork channel demeonstrate available again
	phil.fork <- true
	phil.neighbor.fork <- true
}

func (phil *Philosopher) dine(announce chan *Philosopher) {
	// the whole procedure of dining
	phil.think()
	phil.getForks()
	phil.eat()
	phil.returnForks()
	announce <- phil
}

func main() {
	names := []string{"Phil 1", "Phil 2", "Phil 3", "Phil 4",
		"Phil 5", "Phil 6", "Phil 7", "Phil 8"}
	philosophers := make([]*Philosopher, len(names))
	var phil *Philosopher
	// link all philosophers together
	for i, name := range names {
		phil = makePhilosopher(name, phil)
		philosophers[i] = phil
	}
	// let the first philosopher to be the neighbor of the last philosopher
	philosophers[0].neighbor = phil
	fmt.Printf("There are %v philosophers sitting at a table.\n", len(names))
	fmt.Println("They each have one fork, and must borrow from their neighbor to eat.")
	announce := make(chan *Philosopher)
	for _, phil := range philosophers {
		// dine occur concurrently
		go phil.dine(announce)
	}
	// the announce channel will get the philosophers who finish dining sequentially in concurrent dine()
	// print out them concurrently
	for i := 0; i < len(names); i++ {
		phil := <-announce
		fmt.Printf("%v is done dining. \n", phil.name)
	}
}

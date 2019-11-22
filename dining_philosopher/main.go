package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Philosopher struct contains name, chopstick channel, neighbor philosopher
// chopstick channel is used for demonstrate if the chopstick of the philosopher available
type Philosopher struct {
	name      string
	chopstick chan bool
	neighbor  *Philosopher
}

// Must have this function to declare an object
func makePhilosopher(name string, neighbor *Philosopher) *Philosopher {
	phil := &Philosopher{name, make(chan bool, 1), neighbor}
	phil.chopstick <- true
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

func (phil *Philosopher) getChopsticks() {
	// Declare a channal indicating timeout
	timeout := make(chan bool, 1)
	go func() { time.Sleep(1e9); timeout <- true }()
	// taking the chopstick; chopstick is not available
	<-phil.chopstick
	fmt.Printf("%v got his chopstick. \n", phil.name)
	select {
	// if the neighbor's chopstick is availble, return, ready to eat
	case <-phil.neighbor.chopstick:
		fmt.Printf("%v got %v's chopstick.\n", phil.name, phil.neighbor.name)
		fmt.Printf("%v has two chopsticks.\n", phil.name)
		return
	// after amount of time the philosopher taking up his own chopstick, the
	// philosopher has to put the chopstick down letting others to use
	// then think for a while and try get chopstick operation again
	case <-timeout:
		phil.chopstick <- true
		phil.think()
		phil.getChopsticks()
	}
}

func (phil *Philosopher) returnChopsticks() {
	// after a philosopher finish eating, making his chopstick channel
	// and his neighbor's chopstick channel demeonstrate available again
	phil.chopstick <- true
	phil.neighbor.chopstick <- true
}

func (phil *Philosopher) dine(announce chan *Philosopher) {
	// the whole procedure of dining
	phil.think()
	phil.getChopsticks()
	phil.eat()
	phil.returnChopsticks()
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
	fmt.Println("They each have one chopstick, and must borrow from their neighbor to eat.")
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

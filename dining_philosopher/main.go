package main


import (
	"fmt"
	"math/rand"
	"time"
)

type Philosopher struct {
	name string
	chopstick chan bool
	neighbor *Philosopher
}

func makePhilosopher(name string, neighbor *Philosopher) *Philosopher {
	phil := &Philosopher{name, make(chan bool, 1), neighbor}
	phil.chopstick <- true
	return phil
}

func (phil *Philosopher) think() {
	fmt.Printf("%v is thinking.\n", phil.name)
	time.Sleep(time.Duration(rand.Int63n(1e9)))
}

func (phil *Philosopher) eat() {
	fmt.Printf("%v is eating. \n", phil.name)
	time.Sleep(time.Duration(rand.Int63n(1e9)))
}


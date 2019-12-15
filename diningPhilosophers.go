package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Globals I guess?
var mutex = &sync.Mutex{}

// Philosopher : Struct for us to create philosophers.
type Philosopher struct {
	name string

	// Hunger can be any integer 1 - infinity. It will decrement each time they eat until it hits 0 at
	// which point the philosopher will leave the table to go do some deep thinking.
	hunger int

	// State is either thinking, hungry, or eating
	state string

	// To hold the values of their current Chopsticks
	firstChopstick  int
	secondChopstick int
}

// Chopstick : Struct for our Chopsticks. Might be overkill.
type Chopstick struct {
	inUse bool
	value int
}

func createPhilosopher(name string, hunger int) Philosopher {
	person := Philosopher{name: name, hunger: hunger, state: "thinking", firstChopstick: -1, secondChopstick: -1}

	return person
}

func dine(philosopher Philosopher, chopsticks []Chopstick, announce chan Philosopher) {
	// They all start thinking first
	thinking(philosopher)

	// Once they are done thinking they will try and get two chopsticks
	for {
		flag := grabChopsticks(philosopher, chopsticks)

		if flag {
			break
		}

		thinking(philosopher)
	}

	// If they get both chopsticks they can eat

	announce <- philosopher
}

func thinking(philosopher Philosopher) {
	// Update their status
	philosopher.state = "thinking"

	fmt.Printf("Philosopher %s is thinking\n", philosopher.name)

	// They're "thinking" for a random amount of time
	thinkingTime := rand.Intn(5)
	time.Sleep(time.Duration(thinkingTime) * time.Second)

	// After they've thought so hard they will be hungry
	philosopher.state = "hungry"

	fmt.Printf("Philosopher %s is hungry\n", philosopher.name)
}

func grabChopsticks(philosopher Philosopher, chopsticks []Chopstick) bool {
	for i := 0; i < 5; i++ {

		if !chopsticks[i].inUse {
			// Check if their first chopstick is "empty", if it isn't check if their second isn't "empty"
			if philosopher.firstChopstick == -1 {
				mutex.Lock()
				chopsticks[i].inUse = true
				mutex.Unlock()

				philosopher.firstChopstick = chopsticks[i].value
				fmt.Printf("Philosopher %s has first chopstick %d\n", philosopher.name, philosopher.firstChopstick)
			} else if philosopher.secondChopstick == -1 {
				mutex.Lock()
				chopsticks[i].inUse = true
				mutex.Unlock()

				philosopher.secondChopstick = chopsticks[i].value
				fmt.Printf("Philosopher %s has second chopstick %d\n", philosopher.name, philosopher.firstChopstick)
			} else {
				break
			}
		}
	}

	if philosopher.firstChopstick != -1 {

		if philosopher.secondChopstick != -1 {
			return true
		}

		// If we only have one chopstick just release it and let someone else try and grab both
		mutex.Lock()
		chopsticks[philosopher.firstChopstick].inUse = false
		mutex.Unlock()

		philosopher.firstChopstick = -1
		philosopher.secondChopstick = -1

		fmt.Printf("Philosopher %s released chopstick %d\n", philosopher.name, philosopher.firstChopstick)
		return false

	}

	return false

}

func eat(philosopher Philosopher) {

}

/**
This initial pass at the problem will involve a static number
of philosophers and chopsticks. The hunger for all philosophers will
all start at three as well
*/
func main() {
	names := [5]string{"John", "Mark", "Tom", "Harrold", "Perry"}
	philosophers := make([]Philosopher, 5)
	chopsticks := make([]Chopstick, 5)

	// Create our philosophers and chopsticks with their default states.
	for i := 0; i < 5; i++ {
		philosophers[i] = createPhilosopher(names[i], 3)
		chopsticks[i] = Chopstick{inUse: false, value: i}
	}

	// Channel to capture whenever a philosopher is done.
	philosopherAnnouncements := make(chan Philosopher)

	for i := 0; i < 5; i++ {
		go dine(philosophers[i], chopsticks, philosopherAnnouncements)
	}

	for i := 0; i < 5; i++ {
		finishedPhilosopher := <-philosopherAnnouncements
		fmt.Printf("Philosopher %v is done.\n", finishedPhilosopher.name)
	}
}

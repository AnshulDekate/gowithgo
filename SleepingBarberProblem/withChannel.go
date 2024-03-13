package main

import (
	"fmt"
	"sync"
	"time"
)

type Barber struct {
	state string // cutting, checking, sleeping
}

type Person struct {
	id int
}

func update(barber *Barber, wg *sync.WaitGroup, ch chan Person) {
	for {
		fmt.Println(barber.state)
		select {
		case p := <-ch:
			barber.state = "cutting"
			fmt.Println("cutting hair of person", p.id, "||| People waiting in room", len(ch))
			time.Sleep(time.Second)
			wg.Done()
			barber.state = "checking"
		default:
			barber.state = "sleeping"
		}
		time.Sleep(time.Second)
	}
}

func try(p Person, ch chan Person) {
	for {
		ch <- p
		break
	}
}

func main() {

	var wg sync.WaitGroup
	ch := make(chan Person, 5) // capacity 5
	barber := Barber{"sleeping"}
	go update(&barber, &wg, ch)
	for i := 0; i < 10; i++ {
		p := Person{i}
		wg.Add(1)
		go try(p, ch)

	}
	wg.Wait()
	fmt.Println("all done")
	for {
		time.Sleep(100)
	}
}

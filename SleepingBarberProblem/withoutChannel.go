package main

import (
	"fmt"
	"sync"
	"time"
)

type Barber struct {
	mu    sync.Mutex
	state string // cutting, checking, sleeping
}

type WaitingRoom struct {
	mu          sync.Mutex
	numOfPeople int
	capacity    int //5
}

func update(waitingRoom *WaitingRoom, barber *Barber) {
	for {
		barber.mu.Lock()

		if barber.state == "sleeping" {
			fmt.Println("sleeping")
			time.Sleep(time.Second)
		} else if barber.state == "cutting" {
			fmt.Println("cutting")
			time.Sleep(2 * time.Second)
			fmt.Println("job done")
			barber.state = "checking"
		} else if barber.state == "checking" {
			fmt.Println("checking")
			waitingRoom.mu.Lock()
			if waitingRoom.numOfPeople > 0 {
				waitingRoom.numOfPeople--
				barber.state = "cutting"
			} else {
				barber.state = "sleeping"
			}
			waitingRoom.mu.Unlock()
		}
		barber.mu.Unlock()
		time.Sleep(1)
	}
}

func try(waitingRoom *WaitingRoom, barber *Barber, id int, wg *sync.WaitGroup) {
	for {

		barber.mu.Lock()

		if barber.state == "sleeping" {
			fmt.Println(barber.state)
			barber.state = "checking"
			waitingRoom.mu.Lock()
			if waitingRoom.numOfPeople < waitingRoom.capacity {
				waitingRoom.numOfPeople++
				barber.mu.Unlock()
				waitingRoom.mu.Unlock()
				fmt.Println(id, "now waiting")
				wg.Done()
				break
			}
			waitingRoom.mu.Unlock()
		} else {
			// fmt.Println(barber.state)
			waitingRoom.mu.Lock()
			if waitingRoom.numOfPeople < waitingRoom.capacity {
				waitingRoom.numOfPeople++
				fmt.Println(id, "now waiting")
				barber.mu.Unlock()
				waitingRoom.mu.Unlock()
				wg.Done()
				break
			}
			waitingRoom.mu.Unlock()
		}
		barber.mu.Unlock()
		time.Sleep(1)
	}
}

func main() {

	var wg sync.WaitGroup
	var mu1 sync.Mutex
	var mu2 sync.Mutex

	barber := Barber{mu1, "sleeping", 0}
	// 10 queries

	waitingRoom := WaitingRoom{mu2, 0, 5}

	go update(&waitingRoom, &barber)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go try(&waitingRoom, &barber, i, &wg)
	}
	wg.Wait()

	fmt.Println("all queued")
	for {
		time.Sleep(100)
	}
}

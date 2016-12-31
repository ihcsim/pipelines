package main

import (
	"fmt"
	"sync"
)

func main() {
	done := make(chan struct{})

	inputs := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	in := extract(inputs...)
	out1 := square(in, done)
	out2 := square(in, done)
	for result := range merge(done, out1, out2) {
		fmt.Println(result)
	}
}

func extract(inputs ...int) <-chan int {
	out := make(chan int, len(inputs))
	for _, num := range inputs {
		out <- num
	}
	close(out)
	return out
}

func square(in <-chan int, done <-chan struct{}) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)

		for num := range in {
			select {
			case out <- num * num:
			case <-done:
				return
			}
		}
	}()

	return out
}

func merge(done <-chan struct{}, ins ...<-chan int) <-chan int {
	out := make(chan int)

	var wait sync.WaitGroup
	wait.Add(len(ins))

	fanIn := func(channel <-chan int) {
		defer wait.Done()

		for num := range channel {
			select {
			case out <- num:
			case <-done:
				return
			}
		}
	}

	for _, in := range ins {
		go fanIn(in)
	}

	go func() {
		wait.Wait()
		close(out)
	}()

	return out
}

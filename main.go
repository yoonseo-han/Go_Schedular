package main

import (
	"fmt"
	"time"

	"go_prj/runtime"
)

func foo() {
	fmt.Println("foo")
	time.Sleep(2 * time.Second)
	fmt.Println("foo done")
}

func bar() {
	fmt.Println("bar")
	time.Sleep(10 * time.Second)
	fmt.Println("bar done")
}

func baz() {
	fmt.Println("baz")
	time.Sleep(1 * time.Second)
	fmt.Println("baz done")
}

func main() {
	sched := runtime.NewScheduler()

	firstG := runtime.NewG(foo)
	sched.Add(firstG)
	secondG := runtime.NewG(bar)
	sched.Add(secondG)
	thirdG := runtime.NewG(baz)
	sched.Add(thirdG)

	// Start multiple M's, each in its own goroutine (real parallelism)
	for i := 0; i < sched.GOMAXPROCS; i++ {
		m := runtime.NewM(int64(i), sched)
		go m.Run(sched)
	}

	sched.Wait()
	fmt.Println("all G's completed")
}

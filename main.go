package main

import (
	"fmt"
	"time"

	"go_prj/runtime"
)

func foo() {
	fmt.Println("foo")
	time.Sleep(1 * time.Second)
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

	m := runtime.NewM(1, sched)
	m.Run(sched)
}

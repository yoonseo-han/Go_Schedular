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
	scheduler := runtime.NewScheduler()

	firstG := runtime.NewG(foo)
	scheduler.Add(firstG)
	secondG := runtime.NewG(bar)
	scheduler.Add(secondG)
	thirdG := runtime.NewG(baz)
	scheduler.Add(thirdG)

	scheduler.Schedule()
}

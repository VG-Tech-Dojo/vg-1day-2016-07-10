package main 

import (
	"fmt"
	"time"
)

func main() {
	t := time.Now()
	fmt.Println("hello world")
	fmt.Println(t)
	const layout = "15:04:05"
	fmt.Println(t.Format(layout))
}
package main

import (
	"event-collector/internal/greeting"
	"fmt"
)

func main() {
	message := greeting.Hello("World")
	fmt.Println(message)
}

package main

import (
	"event-collector/internal/greeting"
	"fmt"
)

func main() {
	message := greeting.Hello("second app")
	fmt.Println(message)
}

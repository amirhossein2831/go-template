package main

import (
	"context"
	"event-collector/internal/database"
	"event-collector/internal/greeting"
	"fmt"
)

func main() {

	// 💡 Best practice: Get MongoDB URI from an environment variable.
	db := database.GetMongo(context.Background())

	message := greeting.Hello("first app")
	fmt.Println(message)

	db.Close(context.Background())
}

package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
)

// funtion to printout the slack events tht occur
// chan *slacker.CommandEvent is a channel that receives a pointer to a CommandEvent
func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println()
	}
}

func main() {
	godotenv.Load(".env")

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	// using go to run the function in different thread (separate from the main thread)
	go printCommandEvents(bot.CommandEvents())

}

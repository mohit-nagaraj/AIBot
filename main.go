package main

import (
	"context"
	"fmt"
	"log"
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

	bot.Command("momo - <message>", &slacker.CommandDefinition{
		// below are the default things in slacker library
		Description: "send any question to wolfram",
		Examples:    []string{"who is the president of india"},
		/*
		* botCtx is the context of the bot
		* request is the request that the bot receives
		* response is the response that the bot sends
		 */
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			// param message is coming from the command definition
			query := request.Param("message")
			fmt.Println("Query: ", query)
			response.Reply("Searching for the answer...")
		},
	})

	/*
	* mechanism in Go to manage and propagate cancellation signals across goroutines
	* creates a Context and an associated cancel function
	 */
	ctx, cancel := context.WithCancel(context.Background())
	/*
	* context.Background() provides a base Context (often used as a root context).
	* WithCancel wraps this base Context, creating a new Context that can be canceled by calling the cancel function. The cancel function will stop any goroutines or processes using this Context when called.
	* defer cancel is called when main func ends
	 */
	defer cancel()
	/*
	*  It will propagate a "cancellation signal" if the cancel function is called. The ctx is passed to bot.Listen(ctx) to give the bot listening function the ability to react if ctx is canceled.
	 */
	err := bot.Listen(ctx)

	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Edw590/go-wolfram"
	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
	"github.com/tidwall/gjson"

	witai "github.com/wit-ai/wit-go/v2"
)

// Function to print out the Slack events that occur
func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
	}
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	client := witai.NewClient(os.Getenv("WIT_AI_TOKEN"))
	wolframClient := &wolfram.Client{AppID: os.Getenv("WOLFRAM_APP_ID")}

	// Start a goroutine to print command events
	go printCommandEvents(bot.CommandEvents())

	// Define the command for the bot
	bot.Command("momo - <message>", &slacker.CommandDefinition{
		Description: "Send any question to Wolfram Alpha",
		Examples:    []string{"who is the president of India"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			// Get the message parameter from the command
			query := request.Param("message")
			fmt.Printf("Received query: %v\n", query)

			// Process the query with Wit.ai
			msg, err := client.Parse(&witai.MessageRequest{Query: query})
			if err != nil {
				log.Printf("Error parsing message with Wit.ai: %v", err)
				response.Reply("Sorry, I couldn't process your request.")
				return
			}

			// Print the parsed message for debugging
			data, _ := json.MarshalIndent(msg, "", "    ")
			rough := string(data)
			fmt.Printf("Wit.ai response: %s\n", rough)

			// Check for the Wolfram query in entities
			value := gjson.Get(rough, "entities.wit$wolfram_search_query:wolfram_search_query.0.value")
			answer := value.String()

			// If no relevant entity is found, use the text from the Wit.ai response
			if answer == "" {
				answer = msg.Text // Use the full text instead
			}

			if answer == "" {
				response.Reply("I couldn't find a question to ask Wolfram Alpha.")
				return
			}

			// Query Wolfram Alpha
			res, err := wolframClient.GetSpokentAnswerQuery(answer, wolfram.Metric, 1000)
			if err != nil {
				log.Printf("Error fetching from Wolfram Alpha: %v", err)
				response.Reply("Sorry, there was an error retrieving the answer from Wolfram Alpha.")
				return
			}

			// Send the response back to Slack
			response.Reply(res)
		},
	})

	// Create a context for the bot
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start listening for commands
	if err := bot.Listen(ctx); err != nil {
		log.Fatal(err)
	}
}

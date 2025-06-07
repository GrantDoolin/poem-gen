package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	gomail "gopkg.in/mail.v2"
)

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	client := openai.NewClient(getEnv("OPENAI_API_KEY"))
	resp, err := client.CreateChatCompletion(
		context.TODO(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4o,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: getEnv("context"),
				},
			},
		},
	)
	if err != nil {
		panic(err.Error())
	}

	message := gomail.NewMessage()

	message.SetHeader("From", getEnv("sender"))
	message.SetHeader("To", getEnv("receiver"))
	message.SetHeader("Subject", getEnv("subject"))

	message.SetBody("text/plain", resp.Choices[0].Message.Content)

	dialer := gomail.NewDialer(
		"smtp.gmail.com",
		587,
		getEnv("sender"),
		getEnv("app_password"))

	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	} else {
		fmt.Println("Email sent successfully!")
	}
}

func getEnv(key string) string {

	return os.Getenv(key)
}

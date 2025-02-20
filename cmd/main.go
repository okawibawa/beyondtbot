package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
)

type JokeResponse struct {
	Type      string `json:"type"`
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
	ID        int    `json:"id"`
}

func getJoke() (string, string, error) {
	response, err := http.Get("https://official-joke-api.appspot.com/random_joke")
	if err != nil {
		return "", "", fmt.Errorf("error fetching jokes: %w", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", "", fmt.Errorf("error reading response body: %w", err)
	}

	var joke JokeResponse
	err = json.Unmarshal(body, &joke)
	if err != nil {
		return "", "", fmt.Errorf("error unmarshalling json: %w", err)
	}

	return joke.Setup, joke.Punchline, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading env file")
	}

	pref := tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Handle("/start", func(c tele.Context) error {
		commands := fmt.Sprintf(
			"i merely exist through these commands:\n\n"+
				"%s - show this message\n"+
				"%s - say hello\n"+
				"%s - get a random joke\n"+
				"%s - the why",
			"/start",
			"/hello",
			"/joke",
			"/imfeelinglucky",
		)

		return c.Send(commands)
	})

	bot.Handle("/hello", func(c tele.Context) error {
		return c.Send("hello world")
	})

	bot.Handle("/joke", func(c tele.Context) error {
		setup, punchline, err := getJoke()
		if err != nil {
			log.Fatalf("error getting joke: %v", err)
		}

		return c.Send(fmt.Sprintf("%s %s", strings.ToLower(setup), strings.ToLower(punchline)))
	})

	bot.Handle("/imfeelinglucky", func(c tele.Context) error {
		return c.Send("escaping t-shaped matrix, forging complexity. specialize is dead, intersections are in.")
	})

	bot.Start()
}

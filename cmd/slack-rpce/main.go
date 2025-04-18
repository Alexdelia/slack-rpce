package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"github.com/slack-io/slacker"
)

func get_random_self_insult() string {
	self_insults := [...]string{
		"human-sized Groot",
		"discount Jar Jar Binks",
		"Dobby with no socks",
		"soggy cereal of a human",
		"potato with legs",
		"toeless hobbit",
		"absolute buffoon",
		"Magikarp",
		"Metroid Zoomer",
		"Tetris `S` block",
		"goblin",
	}

	return self_insults[rand.Intn(len(self_insults))]
}

func get_rpce_path() (string, error) {
	var rpce_path string
	if len(os.Args) > 1 {
		rpce_path = os.Args[1]
	} else {
		rpce_path = "./rpce"
	}

	if _, err := os.Stat(rpce_path); os.IsNotExist(err) {
		return "", fmt.Errorf("rpce path '%s' does not exist", rpce_path)
	}

	return rpce_path, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	rpce_path, err := get_rpce_path()
	if err != nil {
		log.Fatal(err)
	}

	dev_id := os.Getenv("SLACK_DEV_ID")
	if dev_id == "" {
		log.Fatal("SLACK_DEV_ID is not set")
	}

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	bot.OnConnected(func(event socketmode.Event) {
		bot.SlackClient().PostMessage(dev_id, slack.MsgOptionText("bot is active", false))
	})

	bot.AddCommand(&slacker.CommandDefinition{
		Command: "ping",
		Handler: func(ctx *slacker.CommandContext) {
			ctx.Response().Reply("pong")
		},
	})

	bot.AddCommand(&slacker.CommandDefinition{
		Command: "reboot",
		Handler: func(ctx *slacker.CommandContext) {
			requesterId := ctx.Event().UserID
			channel := ctx.Event().Channel

			bin_path := fmt.Sprintf("./%s/%s.sh", rpce_path, channel.Name)
			if _, err := os.Stat(bin_path); os.IsNotExist(err) {
				ctx.Response().Reply(fmt.Sprintf("sorry, <#%s> is not available for reboot\n\n<@%s>, <@%s> requested a reboot for <#%s>\nyou better set it up quick, you %s", channel.ID, dev_id, requesterId, channel.ID, get_random_self_insult()))
				return
			}

			// cmd := exec.Command(bin_path)
			// if the command fails, print the error
			// if err := cmd.Run(); err != nil {
			// ctx.Response().Reply(fmt.Sprintf("sorry <@%s>, <@%s> failed to reboot <#%s>\n\n
			// }

			ctx.Response().Reply(fmt.Sprintf("<@%s>, <#%s> has been rebooted with success\n", requesterId, channel.ID))
		},
		HideHelp: true,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

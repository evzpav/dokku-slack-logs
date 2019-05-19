package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/hpcloud/tail"
	columnize "github.com/ryanuber/columnize"
	"github.com/shomali11/slacker"
)

const (
	helpHeader = `Usage: dokku slack-logs[:COMMAND]

Read Dokku apps logs from Slack bot

Additional commands:`

	helpContent = `
	slack-logs:token <SLACK_BOT_TOKEN>, set Slack bot token
	slack-logs:enable, connects to Slack bot
`
)

func main() {
	flag.Usage = usage
	flag.Parse()

	cmd := flag.Arg(0)
	switch cmd {
	case "slack-logs:enable":
		readLog()
	case "slack-logs:token":
		setSlackToken(flag.Arg(1))
	case "slack-logs:help":
		usage()
	case "slack-logs":
		usage()
	default:
		dokkuNotImplementExitCode, err := strconv.Atoi(os.Getenv("DOKKU_NOT_IMPLEMENTED_EXIT"))
		if err != nil {
			fmt.Println("failed to retrieve DOKKU_NOT_IMPLEMENTED_EXIT environment variable")
			dokkuNotImplementExitCode = 10
		}
		os.Exit(dokkuNotImplementExitCode)
	}
}

func usage() {
	config := columnize.DefaultConfig()
	config.Delim = ","
	config.Prefix = "\t"
	config.Empty = ""
	content := strings.Split(helpContent, "\n")[1:]
	fmt.Println(helpHeader)
	fmt.Println(columnize.Format(content, config))
}

func readLog() {
	log.Println("Read log running!")

	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	if slackBotToken == "" {
		log.Println("Variable SLACK_BOT_TOKEN not defined")
		return
	}

	bot := slacker.NewClient(slackBotToken)

	bot.Init(func() {
		log.Println("Slack bot connected!")
	})

	bot.Err(func(err string) {
		log.Println(err)
	})

	bot.DefaultCommand(func(request slacker.Request, response slacker.ResponseWriter) {
		response.Reply("Command not found. Type: help.")
	})

	logs := &slacker.CommandDefinition{
		Description: "Read logs",
		Example:     "logs app",
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			param := request.Param("app")
			if param != "" {
				fileName := fmt.Sprintf("/var/log/dokku/%s/web.00.log", param)
				f, err := readFile(fileName)
				if err != nil {
					response.Reply(err.Error())
				} else {
					for line := range f {
						response.Reply(line.Text)
					}
				}

			}

		},
	}

	bot.Command("logs <app>", logs)

	help := &slacker.CommandDefinition{
		Description: "help!",
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			response.Reply("Type: logs appname")
		},
	}

	bot.Help(help)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}

}

func readFile(fileName string) (chan *tail.Line, error) {
	t, err := tail.TailFile(fileName, tail.Config{Follow: true, MustExist: true, ReOpen: true})
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	return t.Lines, nil
}

func setSlackToken(token string) {
	os.Setenv("SLACK_BOT_TOKEN", token)
}

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
	//SlackBotToken ...
	SlackBotToken = "SLACK_BOT_TOKEN"
	pluginName    = "slack-logs"
	helpHeader    = `Usage: dokku slack-logs[:COMMAND]

Read Dokku apps logs from Slack bot

Additional commands:`

	helpContent = `
	slack-logs:enable, connects to Slack bot

	In the Slack bot app type:
	logs <app-name> <type>
	or just
	logs <app-name>
`
)

func main() {
	flag.Usage = usage
	flag.Parse()

	cmd := flag.Arg(0)
	switch cmd {
	case pluginName + ":help", pluginName:
		usage()
	case pluginName + ":enable":
		readLog()
	default:
		dokkuNotImplementExitCode, err := strconv.Atoi(os.Getenv("DOKKU_NOT_IMPLEMENTED_EXIT"))
		if err != nil {
			fmt.Println("Command does not exist")
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
	log.Println("Slack logs running!")

	slackBotToken := os.Getenv(SlackBotToken)
	if slackBotToken == "" {
		log.Println("Variable " + SlackBotToken + " not defined")
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
		Example:     "logs appname",
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			appParam := request.Param("app")
			typeParam := request.StringParam("type", "web")

			if appParam != "" {
				fileName := fmt.Sprintf("/var/log/dokku/%s/%s.00.log", appParam, typeParam)
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

	bot.Command("logs <app> <type>", logs)

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

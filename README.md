Dokku plugin to get app logs from Slack
======================================

## requirements

- dokku 0.12.x+
- docker 1.8.x

## installation

```shell
# dokku-logging-supervisord is needed
dokku plugin:install https://github.com/sehrope/dokku-logging-supervisord.git

dokku plugin:install https://github.com/evzpav/dokku-slack-logs.git
```

## commands

```
	slack-logs:token <SLACK_BOT_TOKEN>, set Slack bot token
	slack-logs:enable, connects to Slack bot
```

## usage

```shell
# create an app on Slack and get the bot token

# on Dokku machine run:
# set the bot token:
    dokku slack-logs:token <SLACK_BOT_TOKEN>
# enable the bot:
    dokku slack-logs:enable

# on Slack app:
    logs dokku-app-name
# it will print log on Slack bot
```



## License

MIT

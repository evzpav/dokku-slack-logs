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
	slack-logs:enable, connects to Slack bot
```

## usage

```shell
# create an app on Slack and get the bot token

# on Dokku machine run:

# enable the bot:
    dokku config:set --global SLACK_BOT_TOKEN=<SLACK_BOT_TOKEN>
    dokku slack-logs:enable

# on Slack app:
    # Default type is web if not specified
    logs <app-name> <type>

    #Example 1:
    logs myapp
    # it will try to read log file from /var/log/dokku/myapp/web.00.log and print to Slack bot

    #Example 2:
    logs myapp worker
    # it will try to read log file from /var/log/dokku/myapp/worker.00.log and print to Slack bot

```



## License

MIT

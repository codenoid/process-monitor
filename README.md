# process-monitor

Process monitor/watcher by process name, and notify user via telegram bot

## Installing

if you already has Go: 

`go get github.com/codenoid/process-monitor`

or download on [Release page](https://github.com/codenoid/process-monitor/releases)

## Requirement

- Telegram bot, generated at [@BotFather](https://t.me/BotFather)
- A Telegram Room/Bot/Account for receiving notification

## Usage

```sh
Usage of ./process-monitor:
  -config string
        path to process-monitor config file (default "config.yaml")
  -watch string
        path to txt file that contain list of process name separated by newline (default "watch_list.txt")
```

after your process-monitor running in first time, invite the bot to channel/group or /start a chat with bot (for private), and after the bot
got invited, the bot will send you current chat room ID, and put that on config.yaml rooms part

after updating your config.yaml, restart the process-monitor process

## Dev TODO

- fix code structure
- fix variable naming

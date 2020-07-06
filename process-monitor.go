package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	cmap "github.com/orcaman/concurrent-map"
	"gopkg.in/yaml.v2"
)

var (
	watchFile  string
	configFile string

	tgBot         *tgbotapi.BotAPI
	storage       cmap.ConcurrentMap
	sessionConfig config
)

func init() {

	flag.StringVar(&watchFile, "watch", "watch_list.txt", "path to txt file that contain list of process name separated by newline")
	flag.StringVar(&configFile, "config", "config.yaml", "path to process-monitor config file")
	flag.Parse()

	log.Println("initializing config...")

	err := yaml.Unmarshal(readFile(configFile), &sessionConfig)
	if err != nil {
		log.Fatal("yaml.Unmarshal(readFile) failing", err)
	}

	tgBot, err = tgbotapi.NewBotAPI(sessionConfig.Notifier.Telegram.Token)
	if err != nil {
		log.Fatal("bot initialization failing", err)
	}

	go func() {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates, err := tgBot.GetUpdatesChan(u)
		if err != nil {
			log.Fatal("GetUpdatesChan failing", err)
		}

		for update := range updates {
			if update.Message == nil { // ignore any non-Message Updates
				continue
			}

			msgContent := fmt.Sprintf("Hi there, this is process-monitor bot, please add `%v` into your config.yaml file, and don't give me message access to this room", update.Message.Chat.ID)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgContent)
			msg.ParseMode = "Markdown"
			msg.ReplyToMessageID = update.Message.MessageID

			tgBot.Send(msg)
		}
	}()

	storage = cmap.New()
}

func main() {

	//  ParseDuration parses a duration string. A duration string is a possibly signed sequence of decimal numbers,
	// each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns",
	// "us" (or "µs"), "ms", "s", "m", "h".
	repeatEvery, err := time.ParseDuration(sessionConfig.NotifConfig.RepeatEvery)
	if err != nil {
		log.Fatal("RepeatEvery input failing", err)
	}

	var watchProcess func(string)
	watchProcess = func(processName string) {
		defer func(processName string) {
			time.Sleep(1 * time.Second)
			go watchProcess(processName)
		}(processName)

		// process down key
		pdk := processName + ":last-down"
		// process last notified
		pln := processName + ":last-notified"

		_, err := exec.Command("pidof", processName).Output()
		if err != nil {

			firstDownTime := int64(0)
			_firstDownTime, exist := storage.Get(pdk)
			if exist {
				if _firstDownTime.(int64) == int64(0) {
					now := time.Now().Unix()
					storage.Set(pdk, now)
					firstDownTime = now
				} else {
					firstDownTime = _firstDownTime.(int64)
				}
			}

			if lastReport, exist := storage.Get(pln); exist {
				lastReportUnix := lastReport.(float64)
				currentUnix := float64(time.Now().Unix())

				// useless float
				if currentUnix-lastReportUnix > repeatEvery.Seconds() {
					broadcastError(sessionConfig.Notifier.Telegram.RoomIds, processName, firstDownTime)
					// set :last-notified
					storage.Set(pln, float64(time.Now().Unix()))
				}
			} else {
				now := time.Now().Unix()
				broadcastError(sessionConfig.Notifier.Telegram.RoomIds, processName, now)
				// set :last-notified
				storage.Set(pln, float64(now))
			}

			log.Println(processName, "process died")
		} else {
			if _firstDownTime, exist := storage.Get(pdk); exist {
				if _firstDownTime.(int64) != 0 {
					broadcastRunning(sessionConfig.Notifier.Telegram.RoomIds, processName)
				}
			}

			storage.Set(pdk, int64(0))
		}
	}

	// read watched file and split by newline
	watchedFileList := string(readFile(watchFile))
	pNameArr := strings.Split(watchedFileList, "\n")

	// start watcher process
	for _, line := range pNameArr {
		if len(line) > 0 {
			log.Println("Watching", line, "process")
			go watchProcess(line)
		}
	}

	<-make(chan bool)
}

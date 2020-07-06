package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func readFile(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("read file failing", err)
	}

	return b
}

func broadcastError(RoomIds []int64, processName string, firstDownTime int64) {
	for _, RoomID := range RoomIds {
		msgContent := fmt.Sprintf(`**%v** process is down !

first down time: %v`, processName, time.Unix(firstDownTime, 0))

		msg := tgbotapi.NewMessage(RoomID, msgContent)
		msg.ParseMode = "Markdown"

		_, err := tgBot.Send(msg)
		if err != nil {
			log.Println("broadcastError message send failing")
		}
	}
}

func broadcastRunning(RoomIds []int64, processName string) {
	for _, RoomID := range RoomIds {
		msgContent := fmt.Sprintf(`**%v** process is started !`, processName)

		msg := tgbotapi.NewMessage(RoomID, msgContent)
		msg.ParseMode = "Markdown"

		_, err := tgBot.Send(msg)
		if err != nil {
			log.Println("broadcastRunning message send failing")
		}
	}
}

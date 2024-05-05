package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os/exec"
)

func (config Config) RunBot() {
	client, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		panic(err)
	}

	client.Debug = config.Debug

	log.Printf("Authorized on account %s\n", client.Self.UserName)

	menu := CreateRootMenu(&config.Menu)

	for _, admin := range config.Admins {
		msg := tgbotapi.NewMessage(admin, menu.Title)
		msg.ReplyMarkup = menu.Keyboard
		if _, err := client.Send(msg); err != nil {
			log.Println(err)
		}
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := client.GetUpdatesChan(updateConfig)

	for update := range updates {
		if !isAdmin(update.SentFrom(), config.Admins) {
			continue
		}

		if update.Message != nil {
			if !update.Message.IsCommand() {
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			switch update.Message.Command() {
			case "start":
				msg.Text = menu.Title
				msg.ReplyMarkup = menu.Keyboard
			default:
				continue
			}

			if _, err := client.Send(msg); err != nil {
				log.Panic(err)
			}
		} else if update.CallbackQuery != nil {
			handleCallbackQuery(update, client, menu)
		}
	}
}

func handleCallbackQuery(update tgbotapi.Update, client *tgbotapi.BotAPI, menu *Menu) {
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
	if _, err := client.Request(callback); err != nil {
		panic(err)
	}

	currentMenu := menu.FindByData(update.CallbackQuery.Data)

	msg := tgbotapi.NewEditMessageTextAndMarkup(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		currentMenu.Title,
		currentMenu.Keyboard,
	)

	if currentMenu.Command.Name != "" {
		msg.Text = fmt.Sprintf(
			"%s\n\n%s",
			currentMenu.Title,
			runCommand(currentMenu.Command.Name, currentMenu.Command.Args...),
		)
	}

	if _, err := client.Send(msg); err != nil {
		log.Print(err)
	}
}

func isAdmin(user *tgbotapi.User, admins []int64) bool {
	if user == nil {
		return false
	}

	for _, admin := range admins {
		if user.ID == admin {
			return true
		}
	}

	return false
}

func runCommand(name string, arg ...string) string {
	cmd := exec.Command(name, arg...)
	output, _ := cmd.CombinedOutput()
	return fmt.Sprintf("%s %s\n\n%s\n", name, arg, output)
}

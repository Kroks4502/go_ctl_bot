package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go_ctl_bot/internal/config"
	"go_ctl_bot/internal/menu"
	"log"
	"os/exec"
)

func RunBot(cfg *config.Config) {
	client, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		panic(err)
	}

	client.Debug = cfg.Debug

	log.Printf("Authorized on account %s\n", client.Self.UserName)

	rootMenu := menu.CreateRootMenu(&cfg.Menu)

	for _, admin := range cfg.Admins {
		msg := tgbotapi.NewMessage(admin, rootMenu.Title)
		msg.ReplyMarkup = rootMenu.Keyboard
		if _, err := client.Send(msg); err != nil {
			log.Println(err)
		}
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := client.GetUpdatesChan(updateConfig)

	for update := range updates {
		if !isAdmin(update.SentFrom(), cfg.Admins) {
			continue
		}

		if update.Message != nil {
			if !update.Message.IsCommand() {
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			switch update.Message.Command() {
			case "start":
				msg.Text = rootMenu.Title
				msg.ReplyMarkup = rootMenu.Keyboard
			default:
				continue
			}

			if _, err := client.Send(msg); err != nil {
				log.Panic(err)
			}
		} else if update.CallbackQuery != nil {
			handleCallbackQuery(update, client, rootMenu)
		}
	}
}

func handleCallbackQuery(update tgbotapi.Update, client *tgbotapi.BotAPI, rootMenu *menu.Menu) {
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
	if _, err := client.Request(callback); err != nil {
		panic(err)
	}

	currentMenu := rootMenu.FindByData(update.CallbackQuery.Data)

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

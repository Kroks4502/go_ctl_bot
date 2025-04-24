package notice

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Monitor struct {
	bot     *tgbotapi.BotAPI
	admins  []int64
	dirPath string
}

func NewMonitor(bot *tgbotapi.BotAPI, admins []int64, dirPath string) *Monitor {
	return &Monitor{
		bot:     bot,
		admins:  admins,
		dirPath: dirPath,
	}
}

func (m *Monitor) Start() {
	go m.monitorLoop()
}

func (m *Monitor) monitorLoop() {
	for {
		files, err := filepath.Glob(filepath.Join(m.dirPath, "*.msg"))
		if err != nil {
			log.Printf("Error scanning notice directory: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, file := range files {
			content, err := ioutil.ReadFile(file)
			if err != nil {
				log.Printf("Error reading file %s: %v", file, err)
				continue
			}

			// Send message to all admins
			for _, admin := range m.admins {
				msg := tgbotapi.NewMessage(admin, fmt.Sprintf("New notice:\n\n%s", string(content)))
				if _, err := m.bot.Send(msg); err != nil {
					log.Printf("Error sending message to admin %d: %v", admin, err)
				}
			}

			// Delete the file after sending
			if err := os.Remove(file); err != nil {
				log.Printf("Error deleting file %s: %v", file, err)
			}
		}

		time.Sleep(5 * time.Second)
	}
}

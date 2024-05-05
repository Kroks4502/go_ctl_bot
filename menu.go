package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"time"
)

type Menu struct {
	Id int64
	*ConfigButton
	SubMenu   []*Menu
	Keyboard  tgbotapi.InlineKeyboardMarkup
	BackRef   *Menu
	MapByData map[string]*Menu
}

var rootId = time.Now().Unix()

func CreateRootMenu(configButton *ConfigButton) *Menu {
	mapByStrId := make(map[string]*Menu)
	nextId := rootId
	return createSubMenu(configButton, &nextId, nil, mapByStrId)
}

func createSubMenu(configButton *ConfigButton, nextId *int64, backRef *Menu, mapByStrId map[string]*Menu) *Menu {
	root := Menu{Id: *nextId, ConfigButton: configButton, SubMenu: []*Menu{}, BackRef: backRef, MapByData: mapByStrId}
	mapByStrId[strconv.FormatInt(root.Id, 10)] = &root

	var rows [][]tgbotapi.InlineKeyboardButton
	buttons := &configButton.Buttons
	for i := range *buttons {
		*nextId += 1
		rows = append(
			rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData((*buttons)[i].Title, strconv.FormatInt(*nextId, 10)),
			),
		)
		root.SubMenu = append(root.SubMenu, createSubMenu(&(*buttons)[i], nextId, &root, mapByStrId))
	}

	if root.BackRef != nil {
		footerRow := []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("<<<", strconv.FormatInt(rootId, 10)),
		}
		if root.BackRef.BackRef != nil {
			footerRow = append(footerRow, tgbotapi.NewInlineKeyboardButtonData("<", strconv.FormatInt(root.BackRef.Id, 10)))
		}
		rows = append(rows, footerRow)
	}
	root.Keyboard = tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &root
}

func (menu Menu) FindByData(data string) *Menu {
	value, ok := menu.MapByData[data]
	if ok {
		return value
	} else {
		log.Printf("Not found menu by data: %s", data)
		return &menu
	}
}

package telegrambot

import (
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	botAPI     *tgbotapi.BotAPI
	commandMap map[string]TelegramCommand
	setting    TelegramSetting
}

type TelegramSetting struct {
	Token       string            //bot token
	ChatID      int64             //限制chat id
	IsPrivate   bool              //command是否私有使用
	IsEnable    bool              //是否開啟
	AllowNotify bool              //是否允許通知
	Commands    []TelegramCommand //command設定
}

type TelegramCommand struct {
	Command     string
	Description string
	Func        func(string) string
}

func New(settings TelegramSetting) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(settings.Token)
	if err != nil {
		return nil, err
	}
	bot.Debug = false

	tgbot := &TelegramBot{
		botAPI:     bot,
		setting:    settings,
		commandMap: initCommandMap(settings.Commands),
	}

	return tgbot, nil
}

func initCommandMap(commands []TelegramCommand) map[string]TelegramCommand {
	commandMap := make(map[string]TelegramCommand)
	for _, command := range commands {
		commandMap[command.Command] = command
	}
	return commandMap
}

func (chat_bot *TelegramBot) Listen() error {
	defer func() {
		chat_bot.setting.IsEnable = false
	}()
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := chat_bot.botAPI.GetUpdatesChan(updateConfig)

	for update := range updates {
		if !chat_bot.setting.IsEnable {
			break
		}
		if update.Message == nil {
			continue
		}
		if update.Message.IsCommand() {
			if chat_bot.setting.IsPrivate && update.Message.Chat.ID != chat_bot.setting.ChatID {
				continue
			}
			comm := update.Message.Command()
			text := update.Message.Text
			msg := ""
			switch comm {
			case "help":
				msg = CommandHelp(chat_bot.commandMap)
			default:
				if commFunc, ok := chat_bot.commandMap[comm]; ok {
					msg = commFunc.Func(text)
				}
			}
			if msg != "" {
				err := chat_bot.SendMessage(update.Message.Chat.ID, msg)
				if err != nil {
					errmsg := fmt.Sprintf("telegram send message faild: chatID: %d, msg: %s, error: %v", update.Message.Chat.ID, msg, err)
					fmt.Println(errmsg)
				}
			}
		}
	}
	return errors.New("telegram listen down")
}

func (chat_bot *TelegramBot) SetEnable(enable bool) {
	chat_bot.setting.IsEnable = enable
}

func (chat_bot *TelegramBot) SetNotify(allow bool) {
	chat_bot.setting.AllowNotify = allow
}

func (chat_bot *TelegramBot) SetChatID(chatID int64) {
	chat_bot.setting.ChatID = chatID
}

func (chat_bot *TelegramBot) SetPrivate(isPrivate bool) {
	chat_bot.setting.IsPrivate = isPrivate
}

func (chat_bot *TelegramBot) Notify(msg string) error {
	if !chat_bot.setting.IsEnable {
		return errors.New("telegram is not enable")
	}
	if !chat_bot.setting.AllowNotify {
		return errors.New("telegram is not allow notify")
	}
	if chat_bot.setting.ChatID == 0 {
		return errors.New("telegram chatID is nil")
	}

	err := chat_bot.SendMessage(chat_bot.setting.ChatID, msg)
	if err != nil {
		return err
	}

	return nil
}

func (chat_bot *TelegramBot) SendMessage(chatID int64, msg string) error {
	if !chat_bot.setting.IsEnable {
		return errors.New("telegram is not enable")
	}
	replyMsg := tgbotapi.NewMessage(chatID, msg)
	_, err := chat_bot.botAPI.Send(replyMsg)
	if err != nil {
		return err
	}
	return nil
}

package telegrambot

import "fmt"

const (
	COMMAND_HELP_NAME        = "help"
	COMMAND_HELP_DESCRIPTION = ""
)

func CommandHelp(commandMap map[string]TelegramCommand) string {
	msg := ""
	for _, command := range commandMap {
		msg += fmt.Sprintf("/%s : %s \n", command.Command, command.Description)
	}
	if msg == "" {
		return "no command setting"
	}
	return msg
}

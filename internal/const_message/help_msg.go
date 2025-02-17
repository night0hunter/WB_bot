package constmsg

import "wb_bot/internal/enum"

var BotSlashCommands = map[enum.BotSlashCommandType]string{
	enum.BotSlashCommandTypeHelp:  "/help",
	enum.BotSlashCommandTypeAdd:   "/add",
	enum.BotSlashCommandTypeStop:  "/stop",
	enum.BotSlashCommandTypeCheck: "/check",
}

var BotSlashCommandsHelp = map[string]string{
	BotSlashCommands[enum.BotSlashCommandTypeHelp]:  "Команда для вывода информации о доступных функциях",
	BotSlashCommands[enum.BotSlashCommandTypeAdd]:   "Команда для добавления нового отслеживания",
	BotSlashCommands[enum.BotSlashCommandTypeStop]:  "Команда для изменения статуса отслеживания",
	BotSlashCommands[enum.BotSlashCommandTypeCheck]: "Команда для вывода всех текущих отслеживаний",
}

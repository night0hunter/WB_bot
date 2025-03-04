package constmsg

import "wb_bot/internal/enum"

var BotSlashCommands = map[enum.BotSlashCommandType]string{
	enum.BotSlashCommandTypeHelp:   "/help",
	enum.BotSlashCommandTypeAdd:    "/add",
	enum.BotSlashCommandTypeChange: "/change",
	enum.BotSlashCommandTypeCheck:  "/check",
	enum.BotSlashCommandTypeStop:   "/stop",
}

var BotSlashCommandsHelp = map[string]string{
	BotSlashCommands[enum.BotSlashCommandTypeHelp]:   "Команда для вывода информации о доступных функциях",
	BotSlashCommands[enum.BotSlashCommandTypeAdd]:    "Команда для добавления нового отслеживания",
	BotSlashCommands[enum.BotSlashCommandTypeChange]: "Команда для изменения статуса отслеживания",
	BotSlashCommands[enum.BotSlashCommandTypeCheck]:  "Команда для вывода всех текущих отслеживаний",
	BotSlashCommands[enum.BotSlashCommandTypeStop]:   "Команда для удаления отслеживания",
}

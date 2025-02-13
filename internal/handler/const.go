package handler

import "wb_bot/internal/enum"

var BotCommands = map[enum.BotCommandNameType]string{
	enum.BotCommandNameTypeInputDate:       "Введите дату отслеживания в следующем формате: \"дд.мм.гггг-дд.мм.гггг\"",
	enum.BotCommandNameTypeInputWarehouse:  "Выберите склад, который хотите отслеживать",
	enum.BotCommandNameTypeInputCoeffLimit: "Выберите лимит коэффициента или введите свой",
	enum.BotCommandNameTypeInputSupplyType: "Выберите тип поставки",
}

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

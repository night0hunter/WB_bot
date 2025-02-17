package handler

import "wb_bot/internal/enum"

var BotCommands = map[enum.BotCommandNameType]string{
	enum.BotCommandNameTypeInputDate:       "Введите дату отслеживания в следующем формате: \"дд.мм.гггг-дд.мм.гггг\"",
	enum.BotCommandNameTypeInputWarehouse:  "Выберите склад, который хотите отслеживать",
	enum.BotCommandNameTypeInputCoeffLimit: "Выберите лимит коэффициента или введите свой",
	enum.BotCommandNameTypeInputSupplyType: "Выберите тип поставки",
}

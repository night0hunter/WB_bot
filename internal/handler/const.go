package handler

import "wb_bot/internal/enum"

var BotCommands = map[enum.CommandSequence]string{
	// add sequence
	enum.BotCommandNameTypeInputDate:       "Введите дату отслеживания в следующем формате: \"дд.мм.гггг-дд.мм.гггг\"",
	enum.BotCommandNameTypeInputWarehouse:  "Выберите склад, который хотите отслеживать",
	enum.BotCommandNameTypeInputCoeffLimit: "Выберите лимит коэффициента или введите свой",
	enum.BotCommandNameTypeInputSupplyType: "Выберите тип поставки",

	// change sequence
	enum.BotCommandNameTypeTracking: "Выберите отслеживание из списка ниже, чтобы изменить его статус/удалить",

	// booking sequence
	enum.BotCommandNameTypeDraftID:        "Введите ID черновика",
	enum.BotCommandNameTypeBookProtection: "Выберите защиту от бронирования или введите свою",
}

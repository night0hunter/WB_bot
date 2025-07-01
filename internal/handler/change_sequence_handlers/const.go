package changeHandler

import "wb_bot/internal/enum"

var BotCommands = map[enum.CommandSequence]string{
	// change sequence
	enum.BotCommandNameTypeTracking: "Выберите отслеживание из списка ниже, чтобы изменить его статус/удалить",
}

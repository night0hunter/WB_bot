package constmsg

import (
	"wb_bot/internal/enum"
	myError "wb_bot/internal/error"
)

var BotSlashCommands = map[enum.BotSlashCommandType]string{
	enum.BotSlashCommandTypeHelp:   "/help",
	enum.BotSlashCommandTypeAdd:    "/add",
	enum.BotSlashCommandTypeChange: "/change",
	enum.BotSlashCommandTypeCheck:  "/check",
}

var BotSlashCommandsHelp = map[string]string{
	BotSlashCommands[enum.BotSlashCommandTypeHelp]:   "Команда для вывода информации о доступных функциях",
	BotSlashCommands[enum.BotSlashCommandTypeAdd]:    "Команда для добавления нового отслеживания",
	BotSlashCommands[enum.BotSlashCommandTypeChange]: "Команда для изменения статуса/удаления отслеживания",
	BotSlashCommands[enum.BotSlashCommandTypeCheck]:  "Команда для вывода всех текущих отслеживаний",
}

var MatchErrorType = map[myError.ErrorType]string{
	myError.DateInputError:      "Дата введена неверно, попробуйте ещё раз - формат: дд.мм.гггг-дд.мм.гггг",
	myError.WarehouseInputError: "Выберите склад из списка",
	myError.CoeffInputError:     "Лимит коэффициента введён неверно, попробуйте ещё раз",
	myError.SupplyTypeError:     "Выберите тип поставки из списка",
	myError.TrackingChoiceError: "Выберите отслеживание из списка",
	myError.ActionChoiceError:   "Выберите действие из списка",
}

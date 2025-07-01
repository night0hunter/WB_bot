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
	enum.BotSlashCommandTypeBook:   "/book",
}

var BotSlashCommandsHelp = map[string]string{
	BotSlashCommands[enum.BotSlashCommandTypeHelp]:   "Команда для вывода информации о доступных функциях",
	BotSlashCommands[enum.BotSlashCommandTypeAdd]:    "Команда для добавления нового отслеживания",
	BotSlashCommands[enum.BotSlashCommandTypeChange]: "Команда для изменения статуса/удаления отслеживания",
	BotSlashCommands[enum.BotSlashCommandTypeCheck]:  "Команда для вывода всех текущих отслеживаний",
	BotSlashCommands[enum.BotSlashCommandTypeBook]:   "Команда для добавления автобронирования",
}

var MatchErrorType = map[myError.ErrorType]string{
	myError.DateInputError:        "Дата введена неверно, попробуйте ещё раз - формат: дд.мм.гггг-дд.мм.гггг",
	myError.WarehouseInputError:   "Выберите склад из списка",
	myError.CoeffInputError:       "Лимит коэффициента введён неверно, попробуйте ещё раз",
	myError.SupplyTypeError:       "Выберите тип поставки из списка",
	myError.TrackingChoiceError:   "Выберите отслеживание из списка",
	myError.ActionChoiceError:     "Выберите действие из списка",
	myError.SaveStatusChoiceError: "Выберите действие из списка",
	myError.BookingIdError:        "ID введён неверно, попробуйте ещё раз",
}

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

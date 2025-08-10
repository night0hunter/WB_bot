package constmsg

import (
	"wb_bot/internal/enum"
	myError "wb_bot/internal/error"
)

var BotSlashCommands = map[enum.BotSlashCommandType]string{
	enum.BotSlashCommandTypeHelp:           "/help",
	enum.BotSlashCommandTypeAdd:            "/add",
	enum.BotSlashCommandTypeChangeTracking: "/change_tracking",
	enum.BotSlashCommandTypeCheckTrackings: "/trackings",
	enum.BotSlashCommandTypeBook:           "/book",
	enum.BotSlashCommandTypeChangeBooking:  "/change_booking",
	enum.BotSlashCommandTypeCheckBookings:  "/bookings",
}

var BotSlashCommandsHelp = map[string]string{
	BotSlashCommands[enum.BotSlashCommandTypeHelp]:           "Команда для вывода информации о доступных функциях",
	BotSlashCommands[enum.BotSlashCommandTypeAdd]:            "Команда для добавления нового отслеживания",
	BotSlashCommands[enum.BotSlashCommandTypeChangeTracking]: "Команда для изменения статуса/удаления отслеживания",
	BotSlashCommands[enum.BotSlashCommandTypeCheckTrackings]: "Команда для вывода всех текущих отслеживаний",
	BotSlashCommands[enum.BotSlashCommandTypeBook]:           "Команда для добавления автобронирования",
	BotSlashCommands[enum.BotSlashCommandTypeChangeBooking]:  "Команда для изменения статуса/удаления автобронирования",
	BotSlashCommands[enum.BotSlashCommandTypeCheckBookings]:  "Команда для вывода всех текущих автобронирований",
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
	myError.PreorderIDLenError:    "Длина ответа > 1",
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

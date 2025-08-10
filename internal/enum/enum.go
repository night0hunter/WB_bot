package enum

type ButtonType uint8

const (
	ButtonTypeWarehouse ButtonType = iota + 1
	ButtonTypeCoeffLimit
	ButtonTypeSupplyType
	ButtonTypeUserTrackingChoice
	ButtonTypeActionChoice
	ButtonTypeSaveStatus
	ButtonTypeBookProtection
	ButtonTypeUserBookingsChoice
)

type Sequences uint8

const (
	Add Sequences = iota + 1
	ChangeTracking
	Booking
	ChangeBooking
)

type CommandSequence uint8

const (
	// add sequence
	BotCommandNameTypeUnknown CommandSequence = iota
	BotCommandNameTypeAdd
	BotCommandNameTypeInputDate
	BotCommandNameTypeInputWarehouse
	BotCommandNameTypeInputCoeffLimit
	BotCommandNameTypeInputSupplyType

	// change booking sequence
	BotCommandNameTypeChangeBooking

	// change tracking sequence
	BotCommandNameTypeChangeTracking
	BotCommandNameTypeTracking
	BotCommandNameTypeAction

	// booking sequence
	BotCommandNameTypeBook
	BotCommandNameTypeDraftID
	BotCommandNameTypeBookProtection

	// universal command
	BotCommandNameTypeSaveStatus
)

type BotSlashCommandType uint8

const (
	BotSlashCommandTypeHelp BotSlashCommandType = iota + 1
	BotSlashCommandTypeAdd
	BotSlashCommandTypeChangeTracking
	BotSlashCommandTypeCheckTrackings
	BotSlashCommandTypeStop
	BotSlashCommandTypeBook
	BotSlashCommandTypeChangeBooking
	BotSlashCommandTypeCheckBookings
)

type SupplyType uint8

const (
	Box        SupplyType = 2
	Monopallet SupplyType = 5
	SuperSafe  SupplyType = 6
)

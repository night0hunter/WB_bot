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
)

type Sequences uint8

const (
	Add Sequences = iota + 1
	Change
	Booking
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

	// change sequence
	BotCommandNameTypeChange
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
	BotSlashCommandTypeChange
	BotSlashCommandTypeCheck
	BotSlashCommandTypeStop
	BotSlashCommandTypeBook
)

const (
	Box        = "2"
	Monopallet = "5"
	SuperSafe  = "6"
)

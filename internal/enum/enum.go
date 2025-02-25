package enum

type ButtonType uint8

const (
	ButtonTypeWarehouse ButtonType = iota + 1
	ButtonTypeCoeffLimit
	ButtonTypeSupplyType
	ButtonTypeUserTrackings
	ButtonTypeUserTrackingStatus
)

type BotCommandNameType uint8

const (
	BotCommandNameTypeUnknown BotCommandNameType = iota
	BotCommandNameTypeInputDate
	BotCommandNameTypeInputWarehouse
	BotCommandNameTypeInputCoeffLimit
	BotCommandNameTypeInputSupplyType
)

type BotSlashCommandType uint8

const (
	BotSlashCommandTypeHelp BotSlashCommandType = iota + 1
	BotSlashCommandTypeAdd
	BotSlashCommandTypeStop
	BotSlashCommandTypeCheck
)

const (
	Box        = "2"
	Monopallet = "5"
	SuperSafe  = "6"
)

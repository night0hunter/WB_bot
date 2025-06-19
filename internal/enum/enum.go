package enum

type ButtonType uint8

const (
	ButtonTypeWarehouse ButtonType = iota + 1
	ButtonTypeCoeffLimit
	ButtonTypeSupplyType
	ButtonTypeUserTrackingChoice
	ButtonTypeActionChoice
)

type Sequences uint8

const (
	Add Sequences = iota + 1
	Change
)

type CommandSequences uint8

const (
	BotCommandNameTypeUnknown CommandSequences = iota
	BotCommandNameTypeAdd
	BotCommandNameTypeInputDate
	BotCommandNameTypeInputWarehouse
	BotCommandNameTypeInputCoeffLimit
	BotCommandNameTypeInputSupplyType
)

const (
	BotCommandNameTypeChange CommandSequences = iota + 6
	BotCommandNameTypeTracking
	BotCommandNameTypeAction
)

type BotSlashCommandType uint8

const (
	BotSlashCommandTypeHelp BotSlashCommandType = iota + 1
	BotSlashCommandTypeAdd
	BotSlashCommandTypeChange
	BotSlashCommandTypeCheck
	BotSlashCommandTypeStop
)

const (
	Box        = "2"
	Monopallet = "5"
	SuperSafe  = "6"
)

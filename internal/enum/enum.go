package enum

type ButtonType uint8

const (
	ButtonTypeCoeffLimit ButtonType = iota + 1
	ButtonTypeWarehouse
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

var SupplyTypes = map[string]string{
	Box:        "Короб",
	Monopallet: "Монопалета",
	SuperSafe:  "Суперсейф",
}

var WarehouseNames = map[int]string{
	507:    "Коледино",
	316646: "Шушары СГТ",
	301229: "Подольск 4",
	120762: "Электросталь",
	206348: "Тула",
	130744: "Краснодар",
	208277: "Невинномысск",
	117986: "Казань",
	117501: "Подольск",
	1733:   "Екатеринбург - Испытателей 14г",
	218644: "СЦ Хабаровск",
	// : "Санкт-Петербург Уткина Заводь",
	206236: "Белые Столбы",
	686:    "Новосибирск",
}

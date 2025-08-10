package dto

import (
	uuid "github.com/google/uuid"
)

type Error struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type GetWarehouseGoodsV2Item struct {
	Barcode       string  `json:"barcode"`
	Quantity      int     `json:"quantity"`
	Errors        []Error `json:"errors"`
	HasError      bool    `json:"hasError"`
	BrandID       int     `json:"brandId"`
	BrandName     string  `json:"brandName"`
	Sa            string  `json:"sa"`
	NmID          int     `json:"nmId"`
	TsID          int     `json:"tsId"`
	TsName        string  `json:"tsName"`
	ChrtID        int     `json:"chrtId"`
	ImtID         int     `json:"imtId"`
	SubjectID     int     `json:"subjectId"`
	SubjectName   string  `json:"subjectName"`
	NmSa          string  `json:"nmSa"`
	SaImt         string  `json:"saImt"`
	PriceRu       int     `json:"priceRu"`
	ImtNameID     int     `json:"imtNameId"`
	ImtName       string  `json:"imtName"`
	ColorID       int     `json:"colorId"`
	ColorName     string  `json:"colorName"`
	ImgSrc        string  `json:"imgSrc"`
	IsSuperSafe   bool    `json:"isSuperSafe"`
	CanMix        bool    `json:"canMix"`
	CanMonopallet bool    `json:"canMonopallet"`
}

type GetWarehouseGoodsV2MetaInfo struct {
	MonoMixQuantity   int `json:"monoMixQuantity"`
	PalletQuantity    int `json:"palletQuantity"`
	SupersafeQuantity int `json:"supersafeQuantity"`
	MonoMixSa         int `json:"monoMixSa"`
	MonopalletSa      int `json:"monopalletSa"`
	SupersafeSa       int `json:"supersafeSa"`
	SaTotal           int `json:"saTotal"`
}

type GetWarehouseGoodsV2Response struct {
	Items    []GetWarehouseGoodsV2Item
	MetaInfo GetWarehouseGoodsV2MetaInfo
}

type GetWarehouseGoodsV2Request struct {
	DraftID     uuid.UUID
	WarehouseID int
}

type JSONRPCResponse struct {
	ID      string                      `json:"id"`
	JSONRPC string                      `json:"jsonrpc"`
	Result  GetWarehouseGoodsV2Response `json:"result"`
}

type GetCreateRequest struct {
	BoxTypeMask int
	DraftID     uuid.UUID
	WarehouseID int
}

type JSONRPCResponseCreate struct {
	ID      string            `json:"id"`
	JSONRPC string            `json:"jsonrpc"`
	Result  GetCreateResponse `json:"result"`
}

type GetCreateResponse struct {
	IDs []GetCreateResponseItem `json:"ids"`
}

type GetCreateResponseItem struct {
	ID          int    `json:"Id"`
	BoxTypeID   int    `json:"boxTypeId"`
	BoxTypeName string `json:"boxTypeName"`
}

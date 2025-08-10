package mock

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"wb_bot/internal/dto"

	"github.com/pkg/errors"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type MockClient struct {
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func (n *MockClient) DoFunc(req *http.Request) (*http.Response, error) {
	switch req.URL.RequestURI() {
	case "/ns/sm/supply-manager/api/v1/plan/validateWarehouseGoodsV2":
		resp, err := validateWhGoodsV2(req)
		if err != nil {
			return &http.Response{}, errors.Wrap(err, "validateWhGoodsV2")
		}

		return resp, nil
	case "/ns/sm-supply/supply-manager/api/v1/supply/create":
		resp, err := create(req)
		if err != nil {
			return &http.Response{}, errors.Wrap(err, "create")
		}

		return resp, nil
	default:
	}

	return &http.Response{}, nil
}

func validateWhGoodsV2(req *http.Request) (*http.Response, error) {

	data := dto.JSONRPCResponse{
		ID:      "json-rpc_100",
		JSONRPC: "2.0",
		Result: dto.GetWarehouseGoodsV2Response{
			Items: []dto.GetWarehouseGoodsV2Item{
				{
					Barcode:       "2043290246763",
					Quantity:      20,
					Errors:        []dto.Error{},
					HasError:      false,
					BrandID:       0,
					BrandName:     "Kami's Home",
					Sa:            "VS30B20250312",
					NmID:          365927126,
					TsID:          0,
					TsName:        "0",
					ChrtID:        536087571,
					ImtID:         354644472,
					SubjectID:     0,
					SubjectName:   "Полки для ванной",
					NmSa:          "VS30B20250312",
					SaImt:         "",
					PriceRu:       2290,
					ImtNameID:     0,
					ImtName:       "Полка для ванной без сверления настенная прямая",
					ColorID:       13600062,
					ColorName:     "черный",
					ImgSrc:        "https://basket-21.wbbasket.ru/vol3659/part365927/365927126/images/tm/1.webp",
					IsSuperSafe:   false,
					CanMix:        true,
					CanMonopallet: true,
				},
				{
					Barcode:       "2040804921710",
					Quantity:      7,
					Errors:        []dto.Error{},
					HasError:      false,
					BrandID:       0,
					BrandName:     "Kami's Home",
					Sa:            "PSSOWSP20231202",
					NmID:          248517690,
					TsID:          0,
					TsName:        "0",
					ChrtID:        389063168,
					ImtID:         175663371,
					SubjectID:     0,
					SubjectName:   "Органайзеры для хранения",
					NmSa:          "PSSOWSP20231202",
					SaImt:         "",
					PriceRu:       1619,
					ImtNameID:     0,
					ImtName:       "Пластиковый органайзер для хранения контейнер раздвижной",
					ColorID:       12065905,
					ColorName:     "белый",
					ImgSrc:        "https://basket-16.wbbasket.ru/vol2485/part248517/248517690/images/tm/1.webp",
					IsSuperSafe:   false,
					CanMix:        true,
					CanMonopallet: true,
				},
			},
			MetaInfo: dto.GetWarehouseGoodsV2MetaInfo{
				MonoMixQuantity:   287,
				PalletQuantity:    287,
				SupersafeQuantity: 0,
				MonoMixSa:         10,
				MonopalletSa:      10,
				SupersafeSa:       0,
				SaTotal:           10,
			},
		},
	}
	j, err := json.Marshal(data)
	if err != nil {
	}

	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBuffer(j)),
	}, nil
}

func create(req *http.Request) (*http.Response, error) {
	data := dto.JSONRPCResponseCreate{
		ID:      "json-rpc_150",
		JSONRPC: "2.0",
		Result: dto.GetCreateResponse{
			IDs: []dto.GetCreateResponseItem{
				{
					ID:          43471407,
					BoxTypeID:   2,
					BoxTypeName: "Короб",
				},
			},
		},
	}

	j, err := json.Marshal(data)
	if err != nil {
	}

	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBuffer(j)),
	}, nil
}

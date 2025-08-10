package wbadp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"wb_bot/internal/dto"
	"wb_bot/internal/mock"

	"github.com/pkg/errors"
)

// Adapter - adapter for working with wb reports with account cookie
type Adapter struct {
	client  mock.HttpClient
	baseURL string
}

// New - constructor for wb reports adapter
func New(client mock.HttpClient, baseURL string) *Adapter {
	return &Adapter{
		client:  client,
		baseURL: baseURL,
	}
}

func (a *Adapter) GetWarehouseGoodsV2(ctx context.Context, input dto.GetWarehouseGoodsV2Request, url string) (dto.GetWarehouseGoodsV2Response, error) {
	var result dto.GetWarehouseGoodsV2Response

	reqURL := a.baseURL + url

	body := map[string]interface{}{
		"params": map[string]interface{}{
			"draftID":     input.DraftID,
			"warehouseId": input.WarehouseID,
		},
		"jsonrpc": "2.0",
		"id":      "json-rpc_178",
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return result, errors.Wrap(err, "json.Marshal")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(payload))
	if err != nil {
		return result, errors.Wrap(err, "http.NewRequestWithContext")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return result, errors.Wrap(err, "client.Do")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, errors.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var rawResp struct {
		Result dto.GetWarehouseGoodsV2Response `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	err = json.NewDecoder(resp.Body).Decode(&rawResp)
	if err != nil {
		return result, errors.Wrap(err, "json.Decode")
	}

	if rawResp.Error != nil {
		return result, errors.Errorf("api error: code=%d message=%s", rawResp.Error.Code, rawResp.Error.Message)
	}

	result.Items = rawResp.Result.Items
	result.MetaInfo = rawResp.Result.MetaInfo

	return result, nil
}

func (a *Adapter) Create(ctx context.Context, input dto.GetCreateRequest, url string) (dto.GetCreateResponse, error) {
	var result dto.GetCreateResponse

	reqURL := a.baseURL + url

	body := map[string]interface{}{
		"params": map[string]interface{}{
			"boxTypeMask":        input.BoxTypeMask,
			"draftID":            input.DraftID,
			"transitWarehouseId": nil,
			"warehouseId":        input.WarehouseID,
		},
		"jsonrpc": "2.0",
		"id":      "json-rpc_180",
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return result, errors.Wrap(err, "json.Marshal")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(payload))
	if err != nil {
		return result, errors.Wrap(err, "http.NewRequestWithContext")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return result, errors.Wrap(err, "client.Do")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, errors.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var rawResp struct {
		Result dto.GetCreateResponse `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	err = json.NewDecoder(resp.Body).Decode(&rawResp)
	if err != nil {
		return result, errors.Wrap(err, "json.Decode")
	}

	if rawResp.Error != nil {
		return result, errors.Errorf("api error: code=%d message=%s", rawResp.Error.Code, rawResp.Error.Message)
	}

	result.IDs = rawResp.Result.IDs

	return result, nil
}

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"wb_bot/internal/dto"

	"github.com/pkg/errors"
)

func GetTrackingsList(ctx context.Context, client http.Client) ([]dto.Response, error) {
	req, err := http.NewRequest(http.MethodGet, os.Getenv("REQ_URL"), nil)
	if err != nil {
		fmt.Printf("http.NewRequest: %s", err.Error())
	}

	req.Header.Add("Authorization", "Bearer"+os.Getenv("BEARER_TOKEN"))

	res, err := client.Do(req)
	if err != nil {
		return []dto.Response{}, errors.Wrap(err, "client.Do")
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return []dto.Response{}, errors.Wrap(err, "io.ReadAll")
	}

	resp := []dto.Response{}

	jsonErr := json.Unmarshal(body, &resp)
	if jsonErr != nil {
		return []dto.Response{}, errors.Wrap(err, "json.Unmarshal")
	}

	// fmt.Printf("HTTP: %s\n", res.Status)

	var sortedResp []dto.Response
	for _, wh := range resp {
		if wh.BoxTypeID == 2 {
			sortedResp = append(sortedResp, wh)
		}
	}

	return sortedResp, nil
}

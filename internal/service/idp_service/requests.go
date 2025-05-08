package idp_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func SendLoginAcceptedRequest(url string, payload AcceptLoginRequestPayload) (*AcceptLoginResponseData, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %s", body)
	}

	var responseData AcceptLoginResponseData
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, err
	}
	return &responseData, nil
}

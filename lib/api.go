package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type NFCDTO struct {
	UID      string `json:"uid"`
	Username string `json:"username"`
	SuicaID  string `json:"suica_id"`
	IsIn     bool   `json:"is_in"`
}

type InoutRequest struct {
	IsIn bool `json:"is_in"`
}

const baseURL = "https://aigrid-731240201745.asia-northeast1.run.app"

func FetchUserInfo(suicaID string) (*NFCDTO, error) {
	resp, err := http.Get(fmt.Sprintf("%s/nfc/%s", baseURL, suicaID))
	if err != nil {
		return nil, fmt.Errorf("API request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var user NFCDTO
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("JSON decode error: %v", err)
	}

	return &user, nil
}

func UpdateInoutStatus(uid string, isIn bool) error {
	reqBody := InoutRequest{IsIn: !isIn}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("JSON encode error: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/inout/%s", baseURL, uid), bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Request creation error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("API request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	return nil
}

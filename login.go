package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	token tokenResponse
)

// tokenResponse : response from the Tuya API
type tokenResponse struct {
	Result struct {
		AccessToken  string `json:"access_token"`
		ExpireTime   int    `json:"expire_time"`
		RefreshToken string `json:"refresh_token"`
		UID          string `json:"uid"`
	} `json:"result"`
	Success bool  `json:"success"`
	T       int64 `json:"t"`
}

// getToken : Used to request a token from the Tuya API
func getToken() error {
	method := "GET"
	body := []byte(``)
	req, _ := http.NewRequest(method, fmt.Sprintf("%s/v1.0/token?grant_type=1", config.Host), bytes.NewReader(body))

	buildHeader(req, body)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bs, _ := io.ReadAll(resp.Body)
	ret := tokenResponse{}
	json.Unmarshal(bs, &ret)

	go func() {
		time.Sleep(time.Duration(ret.Result.ExpireTime) * time.Second)
		err := refreshToken(ret)
		if err != nil {
			fmt.Println("Cannot refresh token : ", err)
		}
	}()

	token = ret

	if !token.Success {
		return fmt.Errorf("failed to get token, check your informations")
	}
	return nil
}

// refreshToken : refreshes the token from the Tuya API
func refreshToken(T tokenResponse) error {
	token = tokenResponse{}
	method := "GET"
	url := fmt.Sprintf("%s/v1.0/token/%s", config.Host, T.Result.RefreshToken)
	body := []byte(``)
	req, _ := http.NewRequest(method, url, bytes.NewReader(body))

	buildHeader(req, body)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bs, _ := io.ReadAll(resp.Body)
	ret := tokenResponse{}
	json.Unmarshal(bs, &ret)

	token = ret
	return err
}

// Sends the brightness value to the Tuya API
func sendDeviceBrightness(percentage float64) {
	min := 10.0
	max := 1000.0
	value := int(max*percentage/100 + min)
	if value > int(max) {
		value = int(max)
	}
	if value < int(min) {
		value = int(min)
	}

	command := []map[string]interface{}{
		{
			"code":  "switch_led",
			"value": true,
		},
		{
			"code":  "bright_value_v2",
			"value": value,
		},
		{
			"code":  "temp_value_v2",
			"value": 300,
		},
	}

	if percentage < 0 {
		command = []map[string]interface{}{
			{
				"code":  "switch_led",
				"value": false,
			},
		}
	}

	data := map[string]interface{}{
		"commands": command,
	}

	method := "POST"
	bs, _ := json.Marshal(data)
	url := fmt.Sprintf("%s/v1.0/devices/%s/commands", config.Host, config.DeviceID)
	req, _ := http.NewRequest(method, url, bytes.NewReader(bs))

	buildHeader(req, bs)
	b, err := http.DefaultClient.Do(req)
	if config.Debug {
		if err != nil {
			fmt.Println("Error when sending brightness : ", err)
		} else {
			body := b.Body
			bs, _ := io.ReadAll(body)
			fmt.Println("API Response : ", string(bs))
		}
	}
}

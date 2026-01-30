// Copyright 2025 NFON AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//
// NFON CTI API GET example: Retrieve phone extensions
//
// What it does:
// 1. Logs in with API username and password to obtain an access token
// 2. Uses the token to send a GET request to retrieve phone extensions data
//
// Steps to run:
// 1. Set environment variables:
//    Linux/macOS:        export NFON_API_USERNAME='<YOUR API USERNAME>'
//                        export NFON_API_PASSWORD='<YOUR API PASSWORD>'
//    Windows CMD:        set NFON_API_USERNAME=<YOUR API USERNAME>
//                        set NFON_API_PASSWORD=<YOUR API PASSWORD>
//    Windows PowerShell: $env:NFON_API_USERNAME='<YOUR API USERNAME>'
//                        $env:NFON_API_PASSWORD='<YOUR API PASSWORD>'
// 2. Run: go run get-phone-data.example.go
//
// Requirements:
// - Go 1.18+

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
	username = os.Getenv("NFON_API_USERNAME")
	password = os.Getenv("NFON_API_PASSWORD")
)

func main() {
	if username == "" || password == "" {
		fmt.Println("Error: NFON_API_USERNAME and NFON_API_PASSWORD must be set")
		os.Exit(1)
	}
	token, err := getAccessToken()
	if err != nil {
		fmt.Println("Error retrieving token:", err)
		return
	}
	fmt.Println("Access token retrieved successfully.")

	data, err := getPhoneExtensionsData(token)
	if err != nil {
		fmt.Println("Error fetching extensions:", err)
		return
	}
	fmt.Println("Phone extensions data:", data)
}

func getAccessToken() (string, error) {
	loginURL := "https://providersupportdata.cloud-cfg.com/v1/login"

	body := map[string]string{
		"username": username,
		"password": password,
	}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("login failed with status: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	token, ok := result["access-token"].(string)
	if !ok {
		return "", fmt.Errorf("access-token not found in response")
	}

	return token, nil
}

func getPhoneExtensionsData(token string) (string, error) {
	url := "https://providersupportdata.cloud-cfg.com/v1/extensions/phone/data"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to fetch extensions with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

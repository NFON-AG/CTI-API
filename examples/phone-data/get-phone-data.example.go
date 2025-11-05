/*
 * Copyright (c) 2025 NFON AG
 * NFON CTI API GET example: Retrieve phone extensions
 *
 * What it does:
 * 1. Logs in with API username and password to obtain an access token
 * 2. Uses the token to send a GET request to retrieve phone extensions data
 *
 * Steps to run:
 * 1. Replace <YOUR API USERNAME> and <YOUR API PASSWORD> with your credentials
 * 2. Run: go run get-phone-data.example.go
 *
 * Requirements:
 * - Go 1.18+
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	username = "<YOUR API USERNAME>"
	password = "<YOUR API PASSWORD>"
)

func main() {
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

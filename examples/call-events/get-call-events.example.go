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
// NFON CTI API SSE example: Get call details and state change events
//
// What it does:
// 1. Logs in with API username and password to obtain an access token
// 2. Opens a Server-Sent Events (SSE) stream to receive call details and state changes
// 3. Continuously logs incoming events until the process is stopped
//
// Steps to run:
// 1. Set environment variables:
//    Linux/macOS:        export NFON_API_USERNAME='<YOUR API USERNAME>'
//                        export NFON_API_PASSWORD='<YOUR API PASSWORD>'
//    Windows CMD:        set NFON_API_USERNAME=<YOUR API USERNAME>
//                        set NFON_API_PASSWORD=<YOUR API PASSWORD>
//    Windows PowerShell: $env:NFON_API_USERNAME='<YOUR API USERNAME>'
//                        $env:NFON_API_PASSWORD='<YOUR API PASSWORD>'
// 2. Run: go run get-call-events.example.go
//
// Requirements:
// - Go 1.18+

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
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
		os.Exit(1)
	}
	fmt.Println("Access token retrieved successfully.")

	if err := streamEvents(token); err != nil {
		fmt.Println("Error streaming events:", err)
		os.Exit(1)
	}
}

func getAccessToken() (string, error) {
	loginURL := "https://providersupportdata.cloud-cfg.com/v1/login"
	body := map[string]string{
		"username": username,
		"password": password,
	}
	j, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(j))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return "", fmt.Errorf("login failed: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	tok, ok := result["access-token"].(string)
	if !ok {
		return "", fmt.Errorf("access-token not found in response")
	}
	return tok, nil
}

func streamEvents(token string) error {
	url := "https://providersupportdata.cloud-cfg.com/v1/extensions/phone/calls"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("failed to open SSE: %s", resp.Status)
	}

	fmt.Println("Connected to SSE:", url)
	fmt.Println("Waiting for call events...\n(Press Ctrl+C to stop)")

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	var eventLines []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			handleEventBlock(strings.Join(eventLines, "\n"))
			eventLines = eventLines[:0]
			continue
		}
		eventLines = append(eventLines, line)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("stream error: %w", err)
	}

	fmt.Println("SSE connection closed by server.")
	return nil
}

func handleEventBlock(block string) {
	if block == "" {
		return
	}
	lines := strings.Split(block, "\n")
	eventName := "message"
	var data strings.Builder

	for _, l := range lines {
		if strings.HasPrefix(l, "event:") {
			eventName = strings.TrimSpace(strings.TrimPrefix(l, "event:"))
		} else if strings.HasPrefix(l, "data:") {
			data.WriteString(strings.TrimSpace(strings.TrimPrefix(l, "data:")))
			data.WriteString("\n")
		}
	}

	payload := strings.TrimSpace(data.String())
	if payload == "" {
		return
	}

	var js interface{}
	if err := json.Unmarshal([]byte(payload), &js); err == nil {
		pp, _ := json.MarshalIndent(js, "", "  ")
		fmt.Printf("[%s] %s\n", eventName, string(pp))
	} else {
		fmt.Printf("[%s] %s\n", eventName, payload)
	}
}

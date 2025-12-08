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
// 1. Replace <YOUR API USERNAME> and <YOUR API PASSWORD> with your credentials
// 2. Compile: javac NfonApiExample.java
// 3. Run: java NfonApiExample
//
// Requirements:
// - Java 11+

import java.io.*;
import java.net.HttpURLConnection;
import java.net.URI;
import java.net.URL;
import java.nio.charset.StandardCharsets;

public class NfonApiExample {

    private static final String USERNAME = "<YOUR API USERNAME>";
    private static final String PASSWORD = "<YOUR API PASSWORD>";

    public static void main(String[] args) {
        try {
            String token = getAccessToken();
            System.out.println("Access token retrieved successfully.");

            String extensions = getPhoneExtensionsData(token);
            System.out.println("Phone extensions data: " + extensions);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    private static String getAccessToken() throws Exception {
        URL url = URI.create("https://providersupportdata.cloud-cfg.com/v1/login").toURL();
        HttpURLConnection conn = (HttpURLConnection) url.openConnection();
        conn.setRequestMethod("POST");
        conn.setRequestProperty("Content-Type", "application/json");
        conn.setRequestProperty("Accept", "application/json");
        conn.setDoOutput(true);

        String jsonInput = String.format("{\"username\":\"%s\", \"password\":\"%s\"}", USERNAME, PASSWORD);

        try (OutputStream os = conn.getOutputStream()) {
            byte[] input = jsonInput.getBytes(StandardCharsets.UTF_8);
            os.write(input, 0, input.length);
        }

        BufferedReader br = new BufferedReader(new InputStreamReader(conn.getInputStream(), StandardCharsets.UTF_8));
        StringBuilder response = new StringBuilder();
        String line;
        while ((line = br.readLine()) != null) {
            response.append(line.trim());
        }

        return response.toString().replaceAll(".*\"access-token\"\\s*:\\s*\"([^\"]+)\".*", "$1");
    }

    private static String getPhoneExtensionsData(String token) throws Exception {
        URL url = URI.create("https://providersupportdata.cloud-cfg.com/v1/extensions/phone/data").toURL();
        HttpURLConnection conn = (HttpURLConnection) url.openConnection();
        conn.setRequestMethod("GET");
        conn.setRequestProperty("Accept", "application/json");
        conn.setRequestProperty("Authorization", "Bearer " + token);

        BufferedReader br = new BufferedReader(new InputStreamReader(conn.getInputStream(), StandardCharsets.UTF_8));
        StringBuilder response = new StringBuilder();
        String line;
        while ((line = br.readLine()) != null) {
            response.append(line.trim());
        }

        return response.toString();
    }
}

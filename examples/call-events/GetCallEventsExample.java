/**
 * Copyright (c) 2025 NFON AG
 * NFON CTI API SSE example: Get call details and state change events
 *
 * What it does:
 * 1. Logs in with API username and password to obtain an access token
 * 2. Opens a Server-Sent Events (SSE) stream to receive call details and state changes
 * 3. Continuously logs incoming events until terminated
 *
 * Steps to run:
 * 1. Replace <YOUR API USERNAME> and <YOUR API PASSWORD>
 * 2. Compile: javac GetCallEventsExample.java
 * 3. Run: java GetCallEventsExample
 *
 * Requirements:
 * - Java 11+
 */

import java.io.*;
import java.net.HttpURLConnection;
import java.net.URI;
import java.net.URL;
import java.nio.charset.StandardCharsets;

public class GetCallEventsExample {

    private static final String USERNAME = "<YOUR API USERNAME>";
    private static final String PASSWORD = "<YOUR API PASSWORD>";

    public static void main(String[] args) {
        try {
            String token = getAccessToken();
            System.out.println("Access token retrieved successfully.");
            openEventStream(token);
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

        String jsonBody = String.format("{\"username\":\"%s\", \"password\":\"%s\"}", USERNAME, PASSWORD);
        try (OutputStream os = conn.getOutputStream()) {
            os.write(jsonBody.getBytes(StandardCharsets.UTF_8));
        }

        if (conn.getResponseCode() / 100 != 2) {
            throw new IOException("Login failed: " + conn.getResponseCode() + " " + conn.getResponseMessage());
        }

        String resp = readAll(conn.getInputStream());
        return resp.replaceAll(".*\"access-token\"\\s*:\\s*\"([^\"]+)\".*", "$1");
    }

    private static void openEventStream(String token) throws Exception {
        URL url = URI.create("https://providersupportdata.cloud-cfg.com/v1/extensions/phone/calls").toURL();
        HttpURLConnection conn = (HttpURLConnection) url.openConnection();
        conn.setRequestMethod("GET");
        conn.setRequestProperty("Accept", "text/event-stream");
        conn.setRequestProperty("Authorization", "Bearer " + token);
        conn.setReadTimeout(0);

        if (conn.getResponseCode() / 100 != 2) {
            throw new IOException("Failed to open SSE: " + conn.getResponseCode() + " " + conn.getResponseMessage());
        }

        System.out.println("Connected to SSE: " + url);
        System.out.println("Waiting for call events...\n(Stop the program to end)");

        try (BufferedReader br = new BufferedReader(new InputStreamReader(conn.getInputStream(), StandardCharsets.UTF_8))) {
            String line;
            StringBuilder eventBuf = new StringBuilder();
            while ((line = br.readLine()) != null) {
                if (line.isEmpty()) {
                    handleEventBlock(eventBuf.toString());
                    eventBuf.setLength(0);
                } else {
                    eventBuf.append(line).append("\n");
                }
            }
            System.out.println("SSE connection closed by server.");
        }
    }

    private static void handleEventBlock(String block) {
        if (block == null || block.isEmpty()) return;
        String[] lines = block.split("\\r?\\n");

        String eventName = "message";
        StringBuilder data = new StringBuilder();

        for (String l : lines) {
            if (l.startsWith("event:")) eventName = l.substring(6).trim();
            else if (l.startsWith("data:")) data.append(l.substring(5).trim()).append("\n");
        }
        String payload = data.toString().trim();

        System.out.println("[" + eventName + "] " + payload);
    }

    private static String readAll(InputStream is) throws IOException {
        try (BufferedReader br = new BufferedReader(new InputStreamReader(is, StandardCharsets.UTF_8))) {
            StringBuilder sb = new StringBuilder();
            String ln;
            while ((ln = br.readLine()) != null) sb.append(ln);
            return sb.toString();
        }
    }
}

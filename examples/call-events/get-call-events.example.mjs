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
// 2. Run: node get-call-events.example.mjs
//
// Requirements:
// - Node.js 18+ (native fetch & ReadableStream)

const USERNAME = process.env.NFON_API_USERNAME;
const PASSWORD = process.env.NFON_API_PASSWORD;

async function getAccessToken() {
  const resp = await fetch("https://providersupportdata.cloud-cfg.com/v1/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Accept": "application/json",
    },
    body: JSON.stringify({ username: USERNAME, password: PASSWORD }),
  });
  if (!resp.ok) throw new Error(`Login failed: ${resp.status} ${resp.statusText}`);
  const data = await resp.json();
  return data["access-token"];
}

async function streamCallEvents(accessToken) {
  const url = "https://providersupportdata.cloud-cfg.com/v1/extensions/phone/calls";
  const resp = await fetch(url, {
    method: "GET",
    headers: {
      "Accept": "text/event-stream",
      "Authorization": `Bearer ${accessToken}`,
    },
  });

  if (!resp.ok || !resp.body) {
    throw new Error(`Failed to open SSE: ${resp.status} ${resp.statusText}`);
  }

  const reader = resp.body.getReader();
  const decoder = new TextDecoder("utf-8");
  let buffer = "";

  console.log(`Connected to SSE: ${url}`);
  console.log("Waiting for call events...\n(Press Ctrl+C to stop)");

  while (true) {
    const { done, value } = await reader.read();
    if (done) {
      console.log("SSE connection closed by server.");
      break;
    }
    buffer += decoder.decode(value, { stream: true });

    const parts = buffer.split(/\r?\n\r?\n/);
    buffer = parts.pop() || "";

    for (const chunk of parts) {
      const lines = chunk.split(/\r?\n/);
      let eventName = "message";
      let data = "";

      for (const line of lines) {
        if (line.startsWith("event:")) eventName = line.slice(6).trim();
        else if (line.startsWith("data:")) data += line.slice(5).trim() + "\n";
      }

      data = data.trim();
      if (data) {
        try {
          const json = JSON.parse(data);
          console.log(`[${eventName}]`, JSON.stringify(json, null, 2));
        } catch {
          console.log(`[${eventName}]`, data);
        }
      }
    }
  }
}

(async () => {
  try {
    if (!USERNAME || !PASSWORD) {
      console.error("Error: NFON_API_USERNAME and NFON_API_PASSWORD must be set");
      process.exit(1);
    }
    const token = await getAccessToken();
    console.log("Access token retrieved successfully.");
    await streamCallEvents(token);
  } catch (err) {
    console.error("Error:", err.message);
    process.exit(1);
  }
})();

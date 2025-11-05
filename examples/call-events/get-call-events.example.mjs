/**
 * Copyright (c) 2025 NFON AG
 * NFON CTI API SSE example: Get call details and state change events
 *
 * What it does:
 * 1. Logs in with API username and password to obtain an access token
 * 2. Opens a Server-Sent Events (SSE) stream to receive call details and state changes
 * 3. Continuously logs incoming events until the process is stopped
 *
 * Steps to run:
 * 1. Replace <YOUR API USERNAME> and <YOUR API PASSWORD>
 * 2. Run: node get-call-events.example.mjs
 *
 * Requirements:
 * - Node.js 18+ (native fetch & ReadableStream)
 */

const USERNAME = "<YOUR API USERNAME>";
const PASSWORD = "<YOUR API PASSWORD>";

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
    const token = await getAccessToken();
    console.log("Access token retrieved successfully.");
    await streamCallEvents(token);
  } catch (err) {
    console.error("Error:", err.message);
    process.exit(1);
  }
})();

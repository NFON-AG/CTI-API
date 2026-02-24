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
// 2. Run: node get-phone-data.example.mjs
//
// Requirements:
// - Node.js 18+ (native fetch support)

// TODO: Change these values to match your application
const APP_NAME = "NFON-GitHub-Example";  // Replace with your application name
const APP_VERSION = "1.0";               // Replace with your application version
const USER_AGENT = `${APP_NAME}/${APP_VERSION}`;

const USERNAME = process.env.NFON_API_USERNAME;
const PASSWORD = process.env.NFON_API_PASSWORD;

/**
 * Step 1: Authenticate with the NFON API to obtain an access token
 */
async function getAccessToken() {
  const response = await fetch("https://providersupportdata.cloud-cfg.com/v1/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Accept": "application/json",
      "User-Agent": USER_AGENT,
    },
    body: JSON.stringify({
      username: USERNAME,
      password: PASSWORD,
    }),
  });

  if (!response.ok) {
    throw new Error(`Login failed: ${response.status} ${response.statusText}`);
  }

  const data = await response.json();

  // The token is returned under "access-token"
  return data["access-token"];
}

/**
 * Step 2: Use the access token to fetch phone extensions data
 */
async function getPhoneExtensionsData(accessToken) {
  const response = await fetch("https://providersupportdata.cloud-cfg.com/v1/extensions/phone/data", {
    method: "GET",
    headers: {
      "Accept": "application/json",
      "Authorization": `Bearer ${accessToken}`,
      "User-Agent": USER_AGENT,
    },
  });

  if (!response.ok) {
    throw new Error(`Failed to fetch extensions: ${response.status} ${response.statusText}`);
  }

  return await response.json();
}

/**
 * Step 3: Run the full flow
 */
(async () => {
  try {
    if (!USERNAME || !PASSWORD) {
      console.error("Error: NFON_API_USERNAME and NFON_API_PASSWORD must be set");
      process.exit(1);
    }
    // Login and retrieve token
    const token = await getAccessToken();
    console.log("Access token retrieved successfully.");

    // Fetch phone extensions using token
    const extensions = await getPhoneExtensionsData(token);
    console.log("Phone extensions data:", JSON.stringify(extensions, null, 2));
  } catch (err) {
    console.error("Error:", err.message);
  }
})();

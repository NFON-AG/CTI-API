# NFON CTI API Usage Manual

- [Introduction](#introduction)
  - [What You Can Do With This API](#what-you-can-do-with-this-api)
  - [Terms of Use](#terms-of-use)
  - [Important Notes](#important-notes)
  - [Support \& Feedback](#support--feedback)
- [API Endpoints](#api-endpoints)
  - [At a glance](#at-a-glance)
- [API Credentials](#api-credentials)
- [Authentication](#authentication)
  - [Introduction](#introduction-1)
  - [Create a new Access Token](#create-a-new-access-token)
  - [Refresh the Access Token](#refresh-the-access-token)
  - [Examples](#examples)

## Introduction  

The NFON CTI API extends the NFON telephone system with a powerful interface for Computer Telephony Integration (CTI). It provides advanced capabilities such as call control, event streaming, and device management, enabling seamless integration of cloud telephony with third-party business applications and collaboration tools. Designed for server-to-server communication, it works with any device — hardphone, softphone, or mobile app — and empowers developers to create customised solutions that enhance communication and collaboration.  

| ![NFON CTI API architecture diagram](assets/cti-api-architecture.webp) | The **NFON CTI API** enables access to data in the context of a *K-Account*.<br><br>Ideally, a third-party web or desktop app running on the client communicates with the third-party server to enable the necessary functions. |
|:--:|--|

### What You Can Do With This API  

Whether you're enriching CRM workflows, synchronising user availability across platforms, or automating call handling, this API makes it possible to bring telephony seamlessly into your digital ecosystem.  

#### Example use cases include:  
- Contact pop-ups triggered by incoming calls (e.g. opening a CRM record)  
- Displaying caller previews with enriched CRM data  
- Synchronising line status across communication tools and calendars  
- Starting and ending calls directly from third-party applications  
- Defining custom ringtones for internal, external, or VIP calls  
- Streaming call events for analytics and reporting  

---

### Terms of Use
By accessing and using the NFON CTI API, you agree to the [terms of use](https://www.nfon.com/en/legal/gtc-sla/).

---

### Important Notes

#### API endpoints may change: 
- Please subscribe to `API Breaking Changes` on the [NFON Status page](https://status.nfon.com) for updates.
- Please refer to the [latest API documentation](https://nfon-ag.github.io/CTI-API/).

#### Supported use cases: 
The CTI API has been designed for incoming and outgoing calls. Currently, the only officially supported use case for which the CTI API has been tested and released is for incoming calls that are directly received by an extension connected to a single terminal.

Functions currently not supported:
- Group calls
- Skill-based calls
- Upstream queues
- Call scenarios in which a call is forwarded

Conferences with several participants using multiple terminals per extension. Only one device may be assigned to each extension. We are continually making improvements to support additional call scenarios and to optimise the authentication process. If these developments introduce a risk of failure to other components (i.e. a 'breaking change'), we will notify you well in advance.

---

### Support & Feedback
NFON is committed to helping you integrate successfully with our APIs. Depending on the type of request or issue, please use the following channels:

#### Feature Requests
Got ideas to improve the API? Submit feature suggestions via Airfocus:
- [Submit in German language](https://nfon.airfocus.com/share/forms/a429ddf7abf54f6afed82d5a7327e030)
- [Submit in English language](https://nfon.airfocus.com/share/forms/648a1c17a224aa6e8937f8393f4bc21c)

#### Development Help
For best practices, implementation guidance, and community support:

- Join the [NFON Partner Portal](https://partners.nfon.com/)
- Contact your assigned sales representative  

#### Issues & Bugs
For technical issues or suspected bugs, please contact NFON Support directly.

## API Endpoints

Check the official [API documentation](https://nfon-ag.github.io/CTI-API/) for latest endpoint references.

- **Base URL**: https://providersupportdata.cloud-cfg.com/
- **Architecture**: RESTful API
- **Data Format**: JSON

### At a glance

The most important endpoints are listed below: 

- **Read out extension configuration ‘Get /extension/phone/data’**

  This endpoint provides an overview of all configured extensions on a K-Account, including both the number and the name assigned to the extension.

- **Read out line status ‘Get /extensions/phone/states’**
  
  This endpoint allows you to read out the line status for each extension.
  This allows the respective availability status for each extension to be displayed in third-party applications.
  When connecting for the first time, the status of all configured extensions is displayed.

- **Start call ‘Post /extensions/phone/calls’**

  This endpoint initiates a call for an extension.
  First, the call is made to the A number and, after acceptance, the call is established to the B number.
  The user can individually specify which end device rings for the user of the A number.
  A UUID is transmitted as the return code for this endpoint, which can be used to end calls.

- **End call ‘Delete /extensions/phone/calls/{uuid}’**

  Ends a call initiated via the API using the returned UUID.

- **Stream call details ‘Get /extensions/phone/calls’**

  This endpoint enables the streaming of all incoming and outgoing telephone events of a K-Account in real time.
  For complex environments, it is necessary to take into account the restrictions indicated above.
  Unlike the endpoint for reading the line status (see above 2. Read out line status), the current status of all extensions is not displayed when connecting for the first time.

## API Credentials

To use the NFON CTI API, you must first obtain an **API Username** and **API Password**.

You can request access on the [NFON CTI API product page](https://www.nfon.com/en/integrations/cti-api)

## Authentication

### Introduction

The NFON CTI API uses the **Authorization HTTP header** with a **Access Token** for secure access to each endpoint. For security reasons, this token expires after 5 Minutes of its creation. You can use the **Refresh Token** to get a new pair of Access and Refresh Token. 

#### Example

```http 
Authorization: Bearer eyJhbGciOiJSUzUxMiIsImtpZCI6ImVkODEwMDM4MGUwMjE1ODkzZTZmMWU3ZThhMDBhYWFkODQ3NGFiYzA1NDNhNzg2ZTk2OGM0ZjExYTZkMmVjOGIiLCJ0eXAiOiJKV1...
```

---

### Create a new Access Token
You can create a new pair of Access and Refresh token by sending your **API Username**  and **API Password** to the `/login` POST-endpoint of the CTI API. 

#### Example
```
curl --location 'https://providersupportdata.cloud-cfg.com/v1/login' \
--header 'Content-Type: application/json' \
--header 'Accept: application/json' \
--data-raw '{
  "username": "<YOUR API USERNAME>",
  "password": "<YOUR API PASSWORD>"
}'
```

---

### Refresh the Access Token
You can get a new pair of Access and Refresh token by sending your **Refresh Token** to the `/login` PUT-endpoint of the CTI API. 

#### Example
```
curl --location --request PUT 'https://providersupportdata.cloud-cfg.com/v1/login' \
--header 'Accept: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJSUzUxMiIsImtpZCI6ImVkODEwMDM4MGUwMjE1O...'
```

---

### Examples

Below you’ll find working examples for API operations using various programming languages. These are designed to help you get started quickly and understand how to authenticate and interact with the NFON CTI API.

> 💡 **Cannot find your programming language of choice?** We recommend you to use an **AI assistant** to rewrite the examples to other programming languages. 

#### JSON example: Get the data of phone extensions of your account.
- [with Go (Golang)](./examples/phone-data/get-phone-data.example.go)
- [with Java](./examples/phone-data/NfonApiExample.java)
- [with Node.js / Javascript](./examples/phone-data/get-phone-data.example.mjs)

#### Server-Sent Events (SSE) stream example: Get call details and state change events
- [with Go (Golang)](./examples/call-events/get-call-events.example.go)
- [with Java](./examples/call-events/GetCallEventsExample.java)
- [with Node.js / Javascript](./examples/call-events/get-call-events.example.mjs)
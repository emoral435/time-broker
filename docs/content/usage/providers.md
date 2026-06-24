---
title: Providers
weight: 2
---

## Google Calendar

time-broker integrates with Google Calendar via the Google Calendar API.
You will need:

1. A Google Cloud project with the Calendar API enabled
2. OAuth 2.0 credentials (desktop application type)
3. The credentials file configured locally

### Setup

1. Go to the [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the **Google Calendar API**
4. Go to **Credentials** > **Create Credentials** > **OAuth client ID**
5. Choose **Desktop application**
6. Download the JSON credentials file
7. Place it at `~/.config/time-broker/credentials.json`

### Planned Providers

Support for additional providers is planned:

- Fantastical
- Outlook Calendar
- iCloud Calendar

---
title: Providers
weight: 2
---

## Google Calendar

time-broker integrates with Google Calendar via the Google Calendar API.

### Scopes Granted Access

The full list of scopes that are possible to be granted can be found on Google's official API website for [Google Calendar here](https://developers.google.com/identity/protocols/oauth2/scopes#calendar). The following scopes are granted to time-broker after the authorization screen is processed, for a short lived time until reprompted / token refresh:

* Make secondary Google calendars, and see, create, change, and delete events on them
* See the list of Google calendars you’re subscribed to
* See the availability on Google calendars you have access to
* See the events on public calendars
* View your availability in your calendars
* See the title, description, default time zone, and other properties of Google calendars you have access to
* View and edit events on all your calendars
* See, create, change, and delete events on Google calendars you own

These are the bare necesities needed to run time-broker.

## Planned Providers

Support for additional providers is planned:

- Fantastical
- Outlook Calendar
- iCloud Calendar

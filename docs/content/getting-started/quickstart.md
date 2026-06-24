---
title: Quick Start
weight: 2
---

## 1. Authenticate with your calendar provider

```shell
time-broker auth --provider google
```

This opens a browser window to authorize time-broker to read your calendar.

## 2. View your availability

```shell
time-broker free --date 2026-07-01
```

Shows your free timeslots for the given day.

## 3. Book a timeslot

```shell
time-broker book --date 2026-07-01 --start 14:00 --end 15:00
```

Creates a busy event for the specified time range.

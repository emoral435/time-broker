---
title: Commands
weight: 1
---

## auth

Authenticate with a calendar provider.

```shell
time-broker auth --provider <name>
```

**Flags:**

| Flag       | Description                          |
| ---------- | ------------------------------------ |
| `--provider` | Calendar provider (e.g. `google`) |

## free

Show free timeslots for a given day.

```shell
time-broker free --date <date> [--duration <minutes>]
```

**Flags:**

| Flag         | Description                            |
| ------------ | -------------------------------------- |
| `--date`     | Date to check (YYYY-MM-DD)             |
| `--duration` | Minimum slot duration in minutes       |

## book

Book a timeslot on your calendar.

```shell
time-broker book --date <date> --start <time> --end <time> [--title <string>]
```

**Flags:**

| Flag     | Description                            |
| -------- | -------------------------------------- |
| `--date` | Date for the event (YYYY-MM-DD)        |
| `--start`| Start time (HH:MM)                     |
| `--end`  | End time (HH:MM)                       |
| `--title`| Event title                            |

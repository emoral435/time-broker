# time-broker

## Summary

time-broker is a TUI-first project that allows users to connect to their favorite calendar provider (google calendar, fantastical, etc.) to see when they are free, and book timeslots right back into their provider all within the terminal.
This aims to expedite the process of going back and forth in your emails to see when you are free, e.g. giving back to someone a range for when you are free.

## Specs

Do not be verbose, outside of the documentation and outside of when a user asks for help. Otherwise, try to be concise, and ensure that what gets output to end users is short and sweet while maintaining necessary content.

## PR Generation

When creating a PR, generate the description using the template at
.github/PULL_REQUEST_TEMPLATE.md. Use `git diff main...HEAD` to understand the
changes.

Style rules:
- Be concise, 2-3 sentences max per section
- No emojis, no em dashes
- Intent: one sentence describing the problem or goal
- Summary of Changes: bullet list of key changes
- Testing: describe what was tested
- Documentation: note if docs were updated

## Testing

- Run `make test` for full test suite, `make test-short` to skip integration tests.
- CI runs `go test ./... -short -v -count=1`, `go vet ./...`, and `go build` on push/PR to main.
- Use `t.TempDir()` + `t.Setenv()` for filesystem isolation in tests.
- New features should include unit tests; integration tests go behind `testing.Short()`.
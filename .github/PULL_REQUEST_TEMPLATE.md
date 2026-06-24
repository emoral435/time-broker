## Intent

<!--
  Describe the problem or feature in one or two sentences.
  Example: "Allow users to filter free slots by minimum duration" or
  "Fix crash when authenticating with an expired token."
-->

## Summary of Changes

<!--
  Bullet-list the key changes. Keep it scannable.
  - Added --duration flag to the free command
  - Updated Google provider to filter slots by duration
  - Added unit tests for slot filtering
-->

## Related Issue

<!-- Link to the GitHub issue, if any. e.g. Closes #42 -->

## Type of Change

- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update
- [ ] Refactor / style
- [ ] CI / build
- [ ] Other (describe):

## Testing

<!--
  Describe how you tested your changes.
  e.g. "go test ./... passes", "manually ran time-broker free --duration 30"
-->

- [ ] `go test ./...` passes
- [ ] `go build ./cmd/time-broker/` succeeds
- [ ] Manual verification (describe below)

## Documentation

- [ ] Docs updated in `docs/content/`
- [ ] `hugo --minify` builds without errors
- [ ] README updated (if applicable)

## Checklist

- [ ] Code follows existing conventions (no comments, no emoji, concise)
- [ ] No new dependencies added (or justified and added to go.mod)
- [ ] Commits are clean and logically separated

## Screenshots (if UI change)

<!--
  If the change affects the TUI, paste a terminal recording or screenshot.
-->

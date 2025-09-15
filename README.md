# CommitBuddy

CommitBuddy makes it easy to commit changes according to the Conventional Commits specification.

## Setup

```bash
$ go build -o commitbuddy commitbuddy.go
$ export PATH="$(pwd):$PATH"
```

## Convenience

```bash
$ alias cbuddy=commitbuddy
```

## Arguments

Flag | Description    | Type   | Required
-----|----------------|--------|---------
`-t` | Commit type    | string | true
`-s` | Commit scope   | string | false
`-m` | Commit message | string | true
`-c` | Use config     | bool   | false

## Help

```bash
$ commitbuddy --help
```

## Usage

Use `git add ...` before committing your changes with `commitbuddy`.

### Basic

```bash
$ commitbuddy -t fix -s TICKET-123 -m "fix bug"
```

### Interactive

```bash
$ commitbuddy
Commit type: fix
Commit scope: TICKET-123
Commit message: "fix bug"
```

## Configuration

While working on a feature, you can configure the commit type and scope:

```bash
$ commitbuddy config -t fix -s TICKET-123
```

This allows you to commit changes by only specifying the commit message:

```bash
$ commitbuddy -c -m "fix bug"
```

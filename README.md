# Hook

Hook manages your GitHub service hooks.

## Install

```go
go get -u github.com/titanous/hook
```

A `hook` binary will appear in your `$GOPATH/bin`.

## Usage

### Get a list of hook templates

```text
hook templates
```

### Create a hook template

```text
hook template web > webhook.json
```

### Add a hook

```text
hook add -file webhook.json user/repo1 org/
```

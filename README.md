# Gook

Gook manages your GitHub service hooks.

## Install

```go
go get -u github.com/titanous/gook
```

A `gook` binary will appear in your `$GOPATH/bin`.

## Usage

### Get a list of hook templates

```text
gook templates
```

### Create a hook template

```text
gook template web > webhook.json
```

### Add a hook

```text
gook add -file webhook.json user/repo1 org/
```

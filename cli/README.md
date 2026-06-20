# Kik CLI

A small standalone Go CLI built with Cobra. Kik is the command-line interface for the Kikplate plate registry.

The command name is:

```bash
kik <command>
```

For example:

```bash
kik hello
kik help
```

## Run

```bash
go run . --help
```

## Build

```bash
go build -o kik .
./kik --help
```

## Commands

```bash
go run . hello
go run . help
go run . help hello
```

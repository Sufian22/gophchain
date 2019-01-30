# GophChain
```
WARNING: This is a Work In Progress and it might not work.
```
A simple blockchain prototype in Go using Proof Of Work, Persistence (key/value storage), Transactions..
It's built for learning purposes, I created it in order to try to understand how blockchain works and what are its principles.

## Usage

### Initialize blockchain

```go
go run cmd/main.go
```

Running this command at first time will create the persistene, a key/value storage file.

### Print chain

```go
go run cmd/main.go printchain
```

### Add block

```go
go run cmd/main.go addblock -data $DATA
```

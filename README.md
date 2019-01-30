# GophChain

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

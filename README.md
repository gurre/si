# Système international (SI)

## Installation

Windows, OS X & Linux:

```
go get github.com/gurre/si
```

## Usage examples

```go
km := si.NewQuantity(si.Kilo, si.Length)
h := si.NewQuantity(si.Hour, si.Time)
kmh := si.NewUnit(100, km,h)

fmt.Println(kmh)
// 100.0 m/s
```

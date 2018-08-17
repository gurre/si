# Système international (SI)

[![GoDoc](https://godoc.org/github.com/gurre/si?status.svg)](https://godoc.org/github.com/gurre/si)
[![License](http://img.shields.io/:license-MIT-blue.svg?style=flat)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/gurre/si)](https://goreportcard.com/report/github.com/gurre/si)

Working with sensors requires some extra thought into how you report measurements. This library provides a way to type several aspects of any sensor data which makes things easier down the road.

## Installation

Windows, OS X & Linux:

```
go get github.com/gurre/si
```

## Usage examples

```go
// NewQuantity takes a prefix and a measure
km := si.NewQuantity(si.Kilo, si.Length)
// Hour is not a SI unit but officially accepted
h := si.NewQuantity(si.Hour, si.Time)
// Combine a value with several SI-units
kmh := si.NewUnit(100, km,h)
parsed := si.Parse(kmh.String())
fmt.Println(kmh, parsed)
// 100.0 km/h 100.0 km/h
```

## More reading

- [SI base unit](https://en.wikipedia.org/wiki/SI_base_unit)
- [Non-SI units mentioned in the SI](https://en.wikipedia.org/wiki/Non-SI_units_mentioned_in_the_SI)

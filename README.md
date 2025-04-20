<h1 align="center">
    <img src="https://github.com/gurre/si/blob/master/gopher_si.png" alt="Mascot" width="300">
    <br />
    Système international (SI)
</h1>

<p align="center">
  <b>A powerful, type-safe unit conversion library for Go</b>
</p>

<p align="center">
  <a href="https://godoc.org/github.com/gurre/si"><img src="https://godoc.org/github.com/gurre/si?status.svg" alt="GoDoc"></a>
  <a href="https://goreportcard.com/report/github.com/gurre/si"><img src="https://goreportcard.com/badge/github.com/gurre/si" alt="Go Report Card"></a>
</p>

## 🌟 Overview

SI is a Go library that brings the power and precision of the International System of Units (SI) to your code. Say goodbye to unit conversion bugs and inconsistencies!

## 🚀 Installation

```bash
go get github.com/gurre/si
```

## 💡 Key Features

- **Type-safe unit conversions** - Catch dimensional errors at compile-time
- **Natural syntax** - Write code that reads like physics equations
- **Comprehensive unit support** - All 7 SI base units plus derived units
- **Binary prefix support** - Handle digital units like MiB, GiB correctly
- **Simple API** - Intuitive helper functions for common units
- **Extensible** - Create and combine your own custom units

## 📚 Usage Examples

### Basic Unit Creation and Conversion

```go
// Create units with intuitive helper functions
distance := si.Kilometers(10)
time := si.Minutes(30)

// Calculate speed
speed := distance.Div(time)
fmt.Println(speed) // Output: 20 km/m

// Convert to different units
kmh, _ := speed.ConvertTo(si.Kilometers(1).Div(si.Hours(1)))
fmt.Println(kmh) // Output in km/h
```

### Physics Calculations

```go
// Newton's Second Law: F = ma
mass := si.Kilograms(75)
acceleration := si.Meters(9.8).Div(si.Second.Pow(2))
force := mass.Mul(acceleration)
fmt.Println("Force:", force)
// Force: 735 N

// Example: Calculate kinetic energy (KE = 1/2 * m * v^2)
mass := si.Kilograms(1500)        // Mass of a car in kg
velocity, _ := si.Parse("20 m/s") // Velocity in m/s
kineticEnergy := si.Scalar(0.5).Mul(mass).Mul(velocity.Pow(2))
fmt.Println("Kinetic energy of a car:", kineticEnergy)
// Kinetic energy of a car: 300000 J
```



## 🔍 Why Use This Library?

Unit conversion errors have caused catastrophic failures in the past, like the [Mars Climate Orbiter crash](https://en.wikipedia.org/wiki/Mars_Climate_Orbiter#Cause_of_failure). By using SI, you get:

- **Safety**: Automatic dimensional analysis prevents mixing incompatible units
- **Clarity**: Code that expresses intent through proper units
- **Precision**: Consistent handling of unit prefixes and conversions
- **Simplicity**: Natural expression of complex physical relationships

## 📖 Documentation

For complete documentation, visit the [GoDoc page](https://godoc.org/github.com/gurre/si).

## 📄 License

Distributed under the MIT License. See the [LICENSE](LICENSE) file for more information.


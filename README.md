<h1 align="center">
    <img src="https://github.com/gurre/si/blob/master/gopher_si.png" alt="Mascot" width="300">
    <br />
    Système international (SI)
</h1>

<p align="center">
  <b>A powerful, type-safe unit conversion library that prevents expensive mistakes</b>
</p>

<p align="center">
  <a href="https://godoc.org/github.com/gurre/si"><img src="https://godoc.org/github.com/gurre/si?status.svg" alt="GoDoc"></a>
  <a href="https://goreportcard.com/report/github.com/gurre/si"><img src="https://goreportcard.com/badge/github.com/gurre/si" alt="Go Report Card"></a>
</p>

## The Problem: Units Lost in Translation

Imagine this: Your industrial monitoring system has sensors collecting data from a manufacturing plant. One sensor reports:

```json
{
  "device_id": "pump_3",
  "power_milliwatts": 2000,
  "pressure_kpa": 350,
  "flow_liters_per_minute": 42
}
```

The data travels through your system, gets processed, stored, and visualized. Months later, a new engineer joins the team and creates a power efficiency calculation:

```go
efficiency := data.power_milliwatts / (data.pressure_kpa * data.flow_liters_per_minute)
```

Everything seems fine until an alert triggers with impossible values. After hours of debugging, the mistake becomes clear: the units got confused. The frontend was converting pressure to pascals, but the field name still said "kpa". Someone refactored the power field to use watts instead of milliwatts but didn't update the field name.

This is a common pattern: embedding units in field names seems convenient, but units and values become decoupled as data moves through systems, leading to silent, catastrophic failures.

## The Solution: Values With Units

The SI library takes a different approach. Instead of separating units from values, it keeps them together:

```json
{
  "device_id": "pump_3",
  "power": "2 W",
  "pressure": "350 kPa",
  "flow": "42 L/min"
}
```

Now, each value carries its unit, and your code can parse, validate, and convert units automatically:

```go
power, _ := si.Parse("2 W")
pressure, _ := si.Parse("350 kPa")
flow, _ := si.Parse("42 L/min") // Automatically converted to m³/s

// Type-safe operations that maintain correct units
efficiency := power.Div(pressure.Mul(flow))

// Check for valid efficiency (dimensionless)
if !si.IsDimension(efficiency, si.Dimensionless) {
    fmt.Println("Invalid efficiency calculation")
}

fmt.Println("Efficiency:", efficiency) // Properly formatted
```

Units remain attached to values throughout their lifecycle, eliminating an entire class of subtle bugs.

## Data Lineage: Units That Tell Their Story

Consider what happens when your system evolves over time:

```
2022: Sensors report temperatures in Celsius
2023: New sensors added that report in Fahrenheit
2024: System standardized on Kelvin for all calculations
```

Without proper unit handling, your historical data becomes a minefield. With SI:

```go
// Temperature readings from different eras, all with units attached
readings := []si.Unit{
    si.Celsius(22.5),    // 2022 sensor
    si.Fahrenheit(72.6), // 2023 sensor
    si.Kelvins(295.7),   // 2024 sensor
}

// All temperatures converted to a standard unit for analysis
for _, temp := range readings {
    // No need to know which era a reading is from
    // No need to check field names or metadata
    c, _ := si.ToCelsius(temp)
    fmt.Printf("Temperature: %s (%.2f C)\n", temp, c)
    // Temperature: 295.65 K (22.50 C)
    // Temperature: 295.705556 K (22.56 C)
    // Temperature: 295.7 K (22.55 C)
}
```

The unit is part of the data's DNA, preserving its lineage as it flows through your system.

## Backfilling Data: Future-Proof Your History

Six months into production, you realize your pressure calculations need a correction factor. With embedded units in field names, you'd need to:

1. Create new database fields with updated names
2. Write complex ETL jobs to transform historical data
3. Update all downstream systems to use the new fields
4. Maintain documentation explaining the change

With SI, backfilling becomes trivial:

```go
// Process historical data, regardless of when it was collected
func processHistoricalReading(reading string) (string, error) {
  // Parse the reading with its original units
  pressure, err := si.Parse(reading)
  if err != nil {
    return "", err
  }
  
  // Apply correction factor without worrying about the original unit
  correctedPressure := pressure.Mul(si.Scalar(1.03))
  
  // Store or return the corrected value, still with proper units
  return correctedPressure.String(), nil
}

// Works seamlessly with:
processHistoricalReading("350 kPa")  // From old dataset, 360.5 kPa
processHistoricalReading("0.35 MPa") // From another system, 360.5 kPa
processHistoricalReading("50.8 psi") // From imperial sensors, 360.761 kPa
```

Your historical data remains valuable and accurate, regardless of when or how it was collected.

## Decoupling Values from Schemas: Adaptable Data Models

Traditional systems tightly couple units to database schemas:

```go
type Reading struct {
  ValueMillivolts int    `json:"value_millivolts"`
  Timestamp       string `json:"timestamp"`
}

// To add a new unit type, you need schema migration
type NewReading struct {
  ValueMillivolts int    `json:"value_millivolts"`
  ValueKilopascal int    `json:"value_kilopascal"` // New field
  Timestamp       string `json:"timestamp"`
}
```

With SI, your data model becomes flexible and future-proof:

```go
type Reading struct {
  Value     si.Unit `json:"value"`      // Can hold any unit type
  Timestamp string  `json:"timestamp"`
}

// Processing is based on the dimension, not the field name
func processReading(r Reading) {
  switch {
  case si.IsDimension(r.Value, si.Volt.Dimension):
    processVoltage(r.Value)
  case si.IsDimension(r.Value, si.Pascal.Dimension):
    processPressure(r.Value)
  case si.IsDimension(r.Value, si.Temperature):
    processTemperature(r.Value)
  // Add new physical quantities without changing the schema
  }
}
```

Your system can adapt to new sensor types, unit preferences, and calculation needs—all without schema migrations or code rewrites.

## 🚀 Installation

```bash
go get github.com/gurre/si
```

## 📚 Usage Examples

```go
// Create units with various representations
temp := si.Celsius(25.5)                                          // 298.65 K
distance := si.Kilometers(1.5)                                    // 1.5 km
flow := si.Meter.MustPow(3).Mul(si.Scalar(0.002)).Div(si.Second)  // 0.002 m³/s

// Parse units from strings (e.g., from sensor readings)
pressure, _ := si.Parse("101.325 kPa")  // 101.325 kPa
velocity, _ := si.Parse("55 km/h")      // 15.278 m/s

// Convert between units
meters, _ := distance.ConvertTo(si.Meter)
fmt.Println(meters)  // 1500 m

// Temperature conversions
tempF, _ := si.ToFahrenheit(temp)       // 77.9
tempC, _ := si.ToCelsius(temp)          // 25.5

// Perform calculations with units
power := pressure.Mul(flow)                     // 0.203 W
energy := power.Mul(si.Hours(2))                // 1459.08 J

// Create units with error handling
unit, err := si.New(100, "km/h")
if err != nil {
    log.Fatal(err)
}

// Power operations with error handling
area, err := si.Meter.Pow(2)  // Returns (Unit, error)
if err != nil {
    log.Fatal(err)
}
// Or use MustPow for known-safe operations
area = si.Meter.MustPow(2)  // Panics on error
```

### Real-World IoT Example

```go
// Define a sensor reading type
type SensorReading struct {
    DeviceID string
    Value    si.Unit
}

// Process readings from different sensors
readings := []SensorReading{
    {DeviceID: "temp-1", Value: si.Celsius(24.5)},
    {DeviceID: "pressure-1", Value: si.MustParse("101.3 kPa")},
}

// Type-safe processing based on dimensions
for _, reading := range readings {
    if si.IsDimension(reading.Value, si.Temperature) {
        tempC, _ := si.ToCelsius(reading.Value)
        if tempC > 25.0 {
            fmt.Printf("ALERT: High temperature: %.1f C\n", tempC)
        }
    } else if si.IsDimension(reading.Value, si.Pascal.Dimension) {
        fmt.Printf("Pressure: %.1f kPa\n", reading.Value.Value/1000)
    }
}
```

## 🌟 Key Features

- **Parse sensor readings** with different units (temperature, pressure, flow, etc.)
- **Perform calculations** across different units safely and accurately
- **Convert between units** without manual conversion factors
- **Create derived measurements** from multiple sensor inputs
- **Verify dimensions** to ensure calculations are physically meaningful
- **Format output values** with appropriate units for reporting and visualization

## 📖 Documentation

For complete documentation, visit the [GoDoc page](https://godoc.org/github.com/gurre/si).

## 🔧 API Reference

### Core Types

- `Unit` - A physical quantity with value and dimension
- `Dimension` - The 7 SI base dimension exponents [Length, Mass, Time, Current, Temperature, Substance, Luminosity]

### Creating Units

| Function | Returns | Description |
|----------|---------|-------------|
| `Parse(s string)` | `(Unit, error)` | Parse a string like "100 km/h" |
| `New(value, symbol)` | `(Unit, error)` | Create unit with value and symbol |
| `MustNew(value, symbol)` | `Unit` | Like New but panics on error |
| `Scalar(value)` | `Unit` | Create dimensionless unit |

### Unit Operations

| Method | Returns | Description |
|--------|---------|-------------|
| `Mul(other)` | `Unit` | Multiply units (adds dimensions) |
| `Div(other)` | `Unit` | Divide units (subtracts dimensions) |
| `Pow(exp)` | `(Unit, error)` | Raise to power (returns error on 0^-n) |
| `MustPow(exp)` | `Unit` | Like Pow but panics on error |
| `Add(other)` | `(Unit, error)` | Add units (must have same dimension) |
| `ConvertTo(target)` | `(Unit, error)` | Convert to target unit scale |
| `Equals(other)` | `bool` | Compare with floating-point tolerance |
| `Compare(other)` | `(int, error)` | Returns -1, 0, or 1 |

### Convenience Constructors

Time: `Hours()`, `Minutes()`, `Seconds()`, `Milliseconds()`
Length: `Kilometers()`, `Meters()`, `Centimeters()`, `Millimeters()`
Mass: `Kilograms()`, `Grams()`
Temperature: `Celsius()`, `Fahrenheit()`, `Kelvins()`
Other: `Watts()`, `Volts()`, `Amperes()`, `Newtons()`, `Pascals()`, `Joules()`, `Hertzs()`, `Psi()`

## ⚡️ Benchmarks

Performance characteristics (Apple M4 Pro):

| Category | Operation | Time | Allocations |
|----------|-----------|------|-------------|
| **Arithmetic** | Mul, Div | ~3 ns | 0 |
| | Pow | ~7 ns | 0 |
| | Add, Compare | ~3 ns | 0 |
| **Formatting** | String() | ~115 ns | 3 |
| | FormatUnit | ~165 ns | 4 |
| | JSON Marshal | ~160 ns | 6 |
| **Parsing** | Parse() | ~2.8 μs | 31 |
| | New() | ~2.6 μs | 26 |
| **Complex** | Reynolds number | ~14 ns | 0 |
| | End-to-end flow | ~11 μs | 118 |
| | Heat exchanger | ~21 μs | 208 |

Key findings:
- **Zero-allocation arithmetic**: Mul, Div, Pow, Add, ConvertTo allocate no memory
- **Parsing overhead**: Most time spent creating StandardContext (~88% of allocations)
- **Formatting**: AST generation dominates formatting cost

```
go test -bench=. -benchmem ./...

BenchmarkMul-12                          367829902           3.26 ns/op        0 B/op       0 allocs/op
BenchmarkDiv-12                          366287973           3.26 ns/op        0 B/op       0 allocs/op
BenchmarkPow-12                          169092705           7.08 ns/op        0 B/op       0 allocs/op
BenchmarkAdd-12                          415120234           2.90 ns/op        0 B/op       0 allocs/op
BenchmarkConvertTo-12                    403919590           2.95 ns/op        0 B/op       0 allocs/op
BenchmarkEquals-12                       340419978           3.57 ns/op        0 B/op       0 allocs/op
BenchmarkCompare-12                      418756206           2.87 ns/op        0 B/op       0 allocs/op
BenchmarkParse-12                          427963           2817 ns/op     7990 B/op      31 allocs/op
BenchmarkNew-12                            456466           2655 ns/op     7344 B/op      26 allocs/op
BenchmarkString-12                        8885310          136.0 ns/op       51 B/op       4 allocs/op
BenchmarkFormatUnitWithPrefix-12         10319325          115.8 ns/op       46 B/op       3 allocs/op
BenchmarkMarshalJSON-12                   6108565          199.2 ns/op       78 B/op       6 allocs/op
BenchmarkComplexCalculation-12          139425274           8.61 ns/op        0 B/op       0 allocs/op
BenchmarkReynoldsNumber-12               84691436          14.26 ns/op        0 B/op       0 allocs/op
BenchmarkEndToEndFlowCalculation-12        110215          11134 ns/op    30680 B/op     118 allocs/op
BenchmarkHeatExchangerDesign-12             52398          21227 ns/op    58752 B/op     208 allocs/op
```

### Profiling

Run with CPU and memory profiling:

```bash
go test -bench=BenchmarkParse -cpuprofile=cpu.prof -memprofile=mem.prof .
go tool pprof -top cpu.prof
go tool pprof -top mem.prof
```

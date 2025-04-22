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
  si.Celsius(22.5),                // 2022 sensor
  si.MustParse("98.6 degF"),       // 2023 sensor (Note: Library uses explicit temp functions)
  si.Kelvin.Mul(si.Scalar(295.15)), // 2024 sensor
}

// All temperatures converted to a standard unit for analysis
for _, temp := range readings {
  // No need to know which era a reading is from
  // No need to check field names or metadata
  // The value itself contains its lineage
  fmt.Printf("Temperature: %.2f K (%.2f C)\n", temp, temp.ToCelsius)
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
processHistoricalReading("350 kPa")  // From old dataset
processHistoricalReading("0.35 MPa") // From another system
processHistoricalReading("50.8 psi") // From imperial sensors
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
temp := si.Celsius(25.5)                                      // 298.65 K
distance := si.Kilometers(1.5)                                // 1500 m
flow := si.Meter.Pow(3).Mul(si.Scalar(0.002)).Div(si.Second)  // 2 L/s

// Parse units from strings (e.g., from sensor readings)
pressure, _ := si.Parse("101.325 kPa")  // 101325 Pa
velocity, _ := si.Parse("55 km/h")      // 15.278 m/s

// Convert between units
meters, _ := distance.ConvertTo(si.Meter)
fmt.Println(meters)  // 1500 m

// Temperature conversions
tempF, _ := si.ToFahrenheit(temp)       // 77.9 F - no manual calculation needed
tempC, _ := si.ToCelsius(temp)          // 25.5 C - easy conversion back

// Perform calculations with units
power := pressure.Mul(flow)                     // 202.65 W
energy := power.Mul(si.Hours(2))                // 1459080 J
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
        fmt.Printf("Pressure: %.1f kPa\n", reading.Value)
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

## ⚡️ Benchmarks

The benchmark results show that:
 - Basic operations (Mul, Div, Pow, Add, ConvertTo) are very fast (<10ns) and don't allocate memory
 - Parsing and string operations are more expensive (~3μs for Parse, ~90ns for String)
 - End-to-end calculations take longer (~8-11μs) and require more memory allocations
 - Complex real-world scenarios are the most resource-intensive (14-35μs with 130-312 allocations)

```
go test -bench=. -benchmem -benchtime=10s ./...
goos: darwin
goarch: arm64
pkg: github.com/gurre/si
cpu: Apple M4 Pro
BenchmarkParse-12                        	 4174719	      2835 ns/op	    7990 B/op	      31 allocs/op
BenchmarkMul-12                          	1000000000	         3.284 ns/op	       0 B/op	       0 allocs/op
BenchmarkDiv-12                          	1000000000	         3.283 ns/op	       0 B/op	       0 allocs/op
BenchmarkPow-12                          	1000000000	         7.120 ns/op	       0 B/op	       0 allocs/op
BenchmarkAdd-12                          	1000000000	         2.934 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertTo-12                    	1000000000	         2.957 ns/op	       0 B/op	       0 allocs/op
BenchmarkString-12                       	138894262	        86.27 ns/op	      28 B/op	       2 allocs/op
BenchmarkComplexCalculation-12           	1000000000	         8.772 ns/op	       0 B/op	       0 allocs/op
BenchmarkHydraulicPower-12               	1000000000	         3.325 ns/op	       0 B/op	       0 allocs/op
BenchmarkReynoldsNumber-12               	849424845	        14.20 ns/op	       0 B/op	       0 allocs/op
BenchmarkVerifyDimension-12              	1000000000	         4.576 ns/op	       0 B/op	       0 allocs/op
BenchmarkEndToEndFlowCalculation-12      	 1000000	     10917 ns/op	   30680 B/op	     118 allocs/op
BenchmarkEndToEndEnergyCalculation-12    	 1518381	      7913 ns/op	   22096 B/op	      80 allocs/op
BenchmarkCelsiusToKelvin-12              	 4621032	      2560 ns/op	    7344 B/op	      26 allocs/op
BenchmarkKilometersToMeters-12           	 4721824	      2546 ns/op	    7344 B/op	      26 allocs/op
BenchmarkMarshalJSON-12                  	84430513	       142.0 ns/op	      56 B/op	       4 allocs/op
BenchmarkPressureConversion-12           	1000000000	         3.079 ns/op	       0 B/op	       0 allocs/op
BenchmarkThermodynamicCalculation-12     	  924319	     13099 ns/op	   36720 B/op	     130 allocs/op
BenchmarkHeatExchangerDesign-12          	  587031	     20535 ns/op	   58752 B/op	     208 allocs/op
BenchmarkPumpingSystemAnalysis-12        	  739095	     16010 ns/op	   44064 B/op	     156 allocs/op
BenchmarkElectricalCircuitAnalysis-12    	  364904	     32496 ns/op	   88128 B/op	     312 allocs/op
```


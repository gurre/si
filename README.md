<h1 align="center">
    <img src="https://github.com/gurre/si/blob/master/gopher_si.png" alt="Mascot" width="300">
    <br />
    Système international (SI)
</h1>

<p align="center">
  <b>A powerful, type-safe unit conversion library for industrial systems and sensor data processing</b>
</p>

<p align="center">
  <a href="https://godoc.org/github.com/gurre/si"><img src="https://godoc.org/github.com/gurre/si?status.svg" alt="GoDoc"></a>
  <a href="https://goreportcard.com/report/github.com/gurre/si"><img src="https://goreportcard.com/badge/github.com/gurre/si" alt="Go Report Card"></a>
</p>

## 🌟 Overview

The SI library provides robust unit handling for IoT systems that process data from multiple sensors with various physical quantities. Key capabilities include:

- **Parse sensor readings** with different units (temperature, pressure, flow, etc.)
- **Perform calculations** across different units safely and accurately
- **Convert between units** without manual conversion factors
- **Create derived measurements** from multiple sensor inputs
- **Aggregate and analyze** time-series sensor data with correct unit handling
- **Format output values** with appropriate units for reporting and visualization

This library ensures consistency and prevents unit-related errors that can cause critical failures in industrial systems.

## 🚀 Installation

```bash
go get github.com/gurre/si
```

## 📚 Usage Examples

### Parsing Sensor Data

```go
// Parse incoming sensor readings from various formats
temperatureSensor, _ := si.Parse("85.2 °C")
pressureSensor, _ := si.Parse("10.3 MPa") 
flowRate, _ := si.Parse("120 L/min")

// Convert all to standard SI units for internal processing
tempKelvin, _ := temperatureSensor.ConvertTo(si.Kelvin)
pressurePascal, _ := pressureSensor.ConvertTo(si.Pascal)
flowRateSI, _ := flowRate.ConvertTo(si.Meters(1).Pow(3).Div(si.Second))

fmt.Println("Temperature:", tempKelvin)
fmt.Println("Pressure:", pressurePascal)
fmt.Println("Flow rate:", flowRateSI)
```

### Industrial Calculations

```go
// Calculate power from pressure and flow rate
pressure := si.Pascal.Mul(si.Scalar(5e6))       // 5 MPa
flowRate := si.New(0.001, "m³/s")               // 1 L/s
power := pressure.Mul(flowRate)                 // Power = Pressure × Flow rate

fmt.Println("Hydraulic power:", power)          // In watts

// Calculate heat exchange rate
massFlow := si.Kilograms(2.5).Div(si.Second)    // 2.5 kg/s of water
specificHeat := si.New(4186, "J/(kg·K)")        // Specific heat of water
tempDiff := si.Kelvin.Mul(si.Scalar(15))        // Temperature difference of 15K
heatRate := massFlow.Mul(specificHeat).Mul(tempDiff)

fmt.Println("Heat exchange rate:", heatRate)    // In watts
```

### Sensor Aggregation

```go
// Aggregate multiple temperature sensor readings
sensors := []si.Unit{
    si.Celsius(85.2),
    si.Celsius(84.7),
    si.Celsius(86.1),
    si.Celsius(85.4),
}

// Calculate sum using Add method
sum := si.Celsius(0)
for _, temp := range sensors {
    newSum, _ := sum.Add(temp)
    sum = newSum
}

// Calculate average temperature
average := sum.Div(si.Scalar(float64(len(sensors))))

fmt.Println("Sum of temperatures:", sum)
fmt.Println("Average temperature:", average)
```

### Dimensional Analysis

```go
// Verify sensor data is in the expected units
temperatureSensor, _ := si.Parse("32.5 °C")
pressureSensor, _ := si.Parse("101.3 kPa")
flowSensor, _ := si.Parse("5.2 L/s")

// Check if readings have the expected dimensions
isTemperature := si.VerifyDimension(temperatureSensor, si.Temperature)
isPressure := si.VerifyDimension(pressureSensor, si.Pascal.Dimension)
isFlow := si.VerifyDimension(flowSensor, si.Meter.Pow(3).Div(si.Second).Dimension)

fmt.Println("Valid temperature reading:", isTemperature)
fmt.Println("Valid pressure reading:", isPressure)
fmt.Println("Valid flow reading:", isFlow)

// Safety check before performing calculations
if !si.VerifyDimension(temperatureSensor, si.Temperature) {
    fmt.Println("Error: Expected temperature reading but got different dimension")
}
```

## 🔧 Advanced Features

### Dimensional Safety

```go
// The library enforces dimensional safety to prevent errors
pressure := si.New(10, "MPa")
temperature := si.Celsius(200)

// This will fail with appropriate error message
result, err := pressure.Add(temperature)
if err != nil {
    fmt.Println("Cannot add different dimensions:", err)
}
```

### Handling Dimensionless Quantities

```go
// Efficiency, ratios, and other dimensionless values
efficiency := si.Scalar(0.85)        // 85% efficiency
safetyFactor := si.Scalar(1.5)       // Safety factor of 1.5

// Dimensionless quantities display without unit suffix
fmt.Println("Process efficiency:", efficiency)  // Output: 0.85
```

## 📊 Industrial Applications

The SI library is ideal for:

- **SCADA systems** - Processing data from multiple sensor types
- **Process control** - Converting between engineering units
- **IoT gateways** - Standardizing measurements from heterogeneous devices
- **Data historians** - Storing time-series data with proper unit context
- **Equipment monitoring** - Calculating derived parameters for condition monitoring
- **Alarm systems** - Evaluating complex conditions across different unit types
- **Regulatory compliance** - Ensuring accurate conversion for reporting

## 📖 Documentation

For complete documentation, visit the [GoDoc page](https://godoc.org/github.com/gurre/si).

## ⚡️ Benchmarks

The benchmark results show that:
 - Basic operations (Mul, Div, Pow, Add, ConvertTo) are very fast (<10ns) and don't allocate memory
 - Parsing and string operations are more expensive (~3μs for Parse, ~90ns for String)
 - End-to-end calculations take longer (~8-11μs) and require more memory allocations
 - Complex real-world scenarios are the most resource-intensive (14-35μs with 130-312 allocations)

```
go test -bench=. -run=^$ -benchtime=60s ./...
goos: darwin
goarch: arm64
pkg: github.com/gurre/si
cpu: Apple M4 Pro
BenchmarkParse-12                        	  415167	      2843 ns/op
BenchmarkMul-12                          	368542444	         3.247 ns/op
BenchmarkDiv-12                          	368604424	         3.255 ns/op
BenchmarkPow-12                          	171484599	         7.013 ns/op
BenchmarkAdd-12                          	381991842	         2.934 ns/op
BenchmarkConvertTo-12                    	409977426	         2.920 ns/op
BenchmarkString-12                       	13863883	        86.60 ns/op
BenchmarkComplexCalculation-12           	134440732	         8.664 ns/op
BenchmarkHydraulicPower-12               	363623462	         3.307 ns/op
BenchmarkReynoldsNumber-12               	86047978	        14.02 ns/op
BenchmarkVerifyDimension-12              	257544004	         4.648 ns/op
BenchmarkEndToEndFlowCalculation-12      	  109608	     10856 ns/op
BenchmarkEndToEndEnergyCalculation-12    	  152655	      7710 ns/op
PASS
ok  	github.com/gurre/si	20.415s
```


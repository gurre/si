package si

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// SI base dimensions for common use.
// These represent the 7 base dimensions in the SI system:
// Length (meters), Mass (kilograms), Time (seconds), Current (amperes),
// Temperature (kelvins), Substance (moles), and Luminosity (candelas).
var (
	// Length represents the dimension of length (meters).
	Length = Dimension{1, 0, 0, 0, 0, 0, 0}

	// Mass represents the dimension of mass (kilograms).
	// Note: kilogram is the SI base unit for mass, unlike other base units which don't have prefixes
	Mass = Dimension{0, 1, 0, 0, 0, 0, 0}

	// TimeDim represents the dimension of time (seconds).
	// Named TimeDim to avoid conflict with the Time type.
	TimeDim = Dimension{0, 0, 1, 0, 0, 0, 0}

	// Current represents the dimension of electric current (amperes).
	Current = Dimension{0, 0, 0, 1, 0, 0, 0}

	// Temperature represents the dimension of temperature (kelvins).
	Temperature = Dimension{0, 0, 0, 0, 1, 0, 0}

	// Substance represents the dimension of amount of substance (moles).
	Substance = Dimension{0, 0, 0, 0, 0, 1, 0}

	// Luminosity represents the dimension of luminous intensity (candelas).
	Luminosity = Dimension{0, 0, 0, 0, 0, 0, 1}

	// Dimensionless represents quantities without physical dimensions.
	Dimensionless = Dimension{0, 0, 0, 0, 0, 0, 0}
)

// Prefixes defines both SI and binary prefixes for unit scaling.
// SI prefixes range from yocto (y, 10^-24) to yotta (Y, 10^24).
// Binary prefixes (Ki, Mi, Gi, Ti, Pi, Ei) are also included.
var Prefixes = map[string]float64{
	"Y": 1e24, "Z": 1e21, "E": 1e18, "P": 1e15,
	"T": 1e12, "G": 1e9, "M": 1e6, "k": 1e3,
	"h": 1e2, "da": 1e1, "": 1, "d": 1e-1,
	"c": 1e-2, "m": 1e-3, "u": 1e-6, "μ": 1e-6, "n": 1e-9,
	"p": 1e-12, "f": 1e-15, "a": 1e-18, "z": 1e-21,
	"y": 1e-24, "Ki": math.Pow(2, 10), "Mi": math.Pow(2, 20),
	"Gi": math.Pow(2, 30), "Ti": math.Pow(2, 40), "Pi": math.Pow(2, 50), "Ei": math.Pow(2, 60),
}

// SymbolicUnits maps domain-specific unit symbols to their dimensions.
// This allows support for non-standard units like dBm.
var SymbolicUnits = map[string]Unit{
	"dBm": {1e-3, Dimension{2, 1, -3, 0, 0, 0, 0}},
	"psi": {6894.76, Pascal.Dimension}, // 1 psi = 6894.76 Pa
}

// ParseUnit parses a unit string like "km/h" into a Unit.
// It handles basic units, prefixed units, compound units, and symbolic units.
//
// Examples:
//
//	meter, _ := ParseUnit("m")
//	velocity, _ := ParseUnit("km/h")
//	energy, _ := ParseUnit("kW*h")
//	pressure, _ := ParseUnit("kg/(m*s^2)")
func ParseUnit(input string) (Unit, error) {
	// Handle special cases
	if input == "" || input == "1" {
		return One, nil
	}

	if unit, ok := SymbolicUnits[input]; ok {
		return unit, nil
	}

	// Use AST-based parser for complex expressions
	return parseUnitExprWithAST(input)
}

// Register conversion functions between si.Unit and parser.Unit
// This initialization establishes bidirectional conversion between
// the internal Unit type and the parser Unit type.
func init() {
	RegisterConversionFunctions(
		// Convert parser.Unit to si.Unit
		func(u Unit) interface{} {
			var dim Dimension
			for i := range dim {
				dim[i] = u.Dimension[i]
			}
			return Unit{Value: u.Value, Dimension: dim}
		},
		// Convert si.Unit to parser.Unit
		func(i interface{}) Unit {
			if u, ok := i.(Unit); ok {
				var dim Dimension
				for i := range dim {
					dim[i] = u.Dimension[i]
				}
				return Unit{Value: u.Value, Dimension: dim}
			}
			return Unit{}
		},
	)
}

// parseUnitExprWithAST parses just the unit part (no value) using the AST-based parser
// This internal function handles the complex logic of parsing unit expressions.
func parseUnitExprWithAST(input string) (Unit, error) {
	// Create standard context with SI units
	ctx := NewStandardContext()

	// Parse the unit expression
	parserUnit, err := ParseComplexUnit(input, ctx)
	if err != nil {
		return Unit{}, err
	}

	// Convert to si.Unit
	siUnit := ConvertToSIUnit(parserUnit).(Unit)
	return siUnit, nil
}

// Parse splits and parses a full unit expression like "100 km/h".
// This extracts the numeric value and unit component from a string.
//
// Examples:
//
//	speed, _ := Parse("100 km/h")     // 27.78 m/s
//	mass, _ := Parse("500 g")         // 0.5 kg
//	pressure, _ := Parse("101.325 kPa") // 101325 Pa
//	temp, _ := Parse("25 °C")         // 298.15 K
func Parse(input string) (Unit, error) {
	fields := strings.Fields(input)

	// Handle case with only a number (dimensionless unit)
	if len(fields) == 1 {
		val, err := strconv.ParseFloat(fields[0], 64)
		if err != nil {
			return One, fmt.Errorf("invalid numeric value: %w", err)
		}
		return Scalar(val), nil
	}

	if len(fields) < 2 {
		return One, fmt.Errorf("invalid unit expression: %s", input)
	}

	val, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return One, fmt.Errorf("invalid numeric value: %w", err)
	}

	// Use AST-based parser for unit component
	unit, err := parseUnitExprWithAST(strings.Join(fields[1:], ""))
	if err != nil {
		return One, err
	}

	unit.Value *= val
	return unit, nil
}

// MustParse works like Parse but panics on error.
// Use this function only when you know the input is valid.
//
// Example:
//
//	// This will panic if the input is invalid
//	speed := MustParse("100 km/h")
func MustParse(input string) Unit {
	u, err := Parse(input)
	if err != nil {
		panic(err)
	}
	return u
}

// New creates a unit with a given value and unit string (e.g. "kg").
// This is a low-level function used by the helper functions but can be called directly.
//
// Example:
//
//	mass := New(1.5, "kg")   // 1.5 kg
//	length := New(100, "cm") // 1 m
func New(value float64, symbol string) Unit {
	// Special case for grams
	if symbol == "g" {
		return Unit{value / 1000, Mass}
	}

	// First try to parse as a direct unit
	u, err := ParseUnit(symbol)
	if err != nil {
		// Return dimensionless unit as fallback
		return Scalar(value)
	}

	u.Value = value
	return u
}

// String returns a human-readable representation of the unit using SI standards.
// The output includes the value with appropriate scaling prefix and unit symbol.
//
// Examples:
//
//	vel := MustParse("27.78 m/s")
//	fmt.Println(vel) // "27.78 m/s"
//
//	energy := Joule.Mul(Scalar(5000))
//	fmt.Println(energy) // "5 kJ"
func (u Unit) String() string {
	return FormatUnitWithPrefix(u)
}

// Named constants for SI base units.
// These represent the basic SI units with their correct dimensions.
var (
	// Meter is the SI base unit of length (1 m).
	Meter = Unit{1, Length}

	// Kilogram is the SI base unit of mass (1 kg).
	Kilogram = Unit{1, Mass}

	// Second is the SI base unit of time (1 s).
	Second = Unit{1, TimeDim}

	// Ampere is the SI base unit of electric current (1 A).
	Ampere = Unit{1, Current}

	// Kelvin is the SI base unit of temperature (1 K).
	Kelvin = Unit{1, Temperature}

	// Mole is the SI base unit of amount of substance (1 mol).
	Mole = Unit{1, Substance}

	// Candela is the SI base unit of luminous intensity (1 cd).
	Candela = Unit{1, Luminosity}

	// One represents a dimensionless quantity with value 1.
	One = Unit{1, Dimensionless}
)

// Named derived units for convenience.
// These are common derived units defined in terms of the base units.
var (
	// Newton is the SI unit of force (1 N = 1 kg·m/s²).
	Newton = Kilogram.Mul(Meter).Div(Second.Pow(2))

	// Joule is the SI unit of energy (1 J = 1 N·m).
	Joule = Newton.Mul(Meter)

	// Watt is the SI unit of power (1 W = 1 J/s).
	Watt = Joule.Div(Second)

	// Pascal is the SI unit of pressure (1 Pa = 1 N/m²).
	Pascal = Newton.Div(Meter.Pow(2))

	// Hertz is the SI unit of frequency (1 Hz = 1/s).
	Hertz = Unit{1, Dimension{0, 0, -1, 0, 0, 0, 0}}

	// Coulomb is the SI unit of electric charge (1 C = 1 A·s).
	Coulomb = Ampere.Mul(Second)

	// Volt is the SI unit of electric potential (1 V = 1 W/A).
	Volt = Watt.Div(Ampere)
)

// Convenience helpers for common physical quantities.
// These functions create units with the appropriate dimensions and conversions.

// Time units

// Hours creates a time unit in hours converted to seconds.
//
// Example:
//
//	duration := Hours(2)  // 2 hours = 7200 seconds
func Hours(n float64) Unit { return New(n*3600, "s") }

// Minutes creates a time unit in minutes converted to seconds.
//
// Example:
//
//	duration := Minutes(30)  // 30 minutes = 1800 seconds
func Minutes(n float64) Unit { return New(n*60, "s") }

// Seconds creates a time unit in seconds.
//
// Example:
//
//	duration := Seconds(45)  // 45 seconds
func Seconds(n float64) Unit { return New(n, "s") }

// Milliseconds creates a time unit in milliseconds converted to seconds.
//
// Example:
//
//	duration := Milliseconds(500)  // 500 milliseconds = 0.5 seconds
func Milliseconds(n float64) Unit { return New(n/1000, "s") }

// Length units

// Kilometers creates a length unit in kilometers converted to meters.
//
// Example:
//
//	distance := Kilometers(5)  // 5 kilometers = 5000 meters
func Kilometers(n float64) Unit { return New(n*1000, "m") }

// Meters creates a length unit in meters.
//
// Example:
//
//	length := Meters(1.8)  // 1.8 meters
func Meters(n float64) Unit { return New(n, "m") }

// Centimeters creates a length unit in centimeters converted to meters.
//
// Example:
//
//	length := Centimeters(25)  // 25 centimeters = 0.25 meters
func Centimeters(n float64) Unit { return New(n/100, "m") }

// Millimeters creates a length unit in millimeters converted to meters.
//
// Example:
//
//	length := Millimeters(150)  // 150 millimeters = 0.15 meters
func Millimeters(n float64) Unit { return New(n/1000, "m") }

// Mass units

// Grams creates a mass unit in grams converted to kilograms.
//
// Example:
//
//	mass := Grams(500)  // 500 grams = 0.5 kilograms
func Grams(n float64) Unit { return New(n/1000, "kg") }

// Kilograms creates a mass unit in kilograms.
//
// Example:
//
//	mass := Kilograms(75)  // 75 kilograms
func Kilograms(n float64) Unit { return New(n, "kg") }

// Temperature units

// Celsius creates a temperature unit in degrees Celsius converted to kelvins.
//
// Example:
//
//	temperature := Celsius(25)  // 25°C = 298.15 K
func Celsius(n float64) Unit { return New(n+273.15, "K") }

// Fahrenheit creates a temperature unit in degrees Fahrenheit converted to kelvins.
//
// Example:
//
//	temperature := Fahrenheit(77)  // 77°F = 298.15 K
func Fahrenheit(n float64) Unit { return New((n-32)*5/9+273.15, "K") }

// Kelvins creates a temperature unit in kelvins.
//
// Example:
//
//	temperature := Kelvins(300)  // 300 K
func Kelvins(n float64) Unit { return New(n, "K") }

// Temperature conversion helpers

// ToCelsius converts a temperature unit to degrees Celsius.
// Returns an error if the unit is not a temperature.
//
// Example:
//
//	temp := Kelvin.Mul(Scalar(300))
//	celsius, _ := ToCelsius(temp)  // celsius = 26.85
func ToCelsius(u Unit) (float64, error) {
	if !IsDimension(u, Temperature) {
		return 0, fmt.Errorf("not a temperature unit")
	}
	return u.Value - 273.15, nil
}

// ToFahrenheit converts a temperature unit to degrees Fahrenheit.
// Returns an error if the unit is not a temperature.
//
// Example:
//
//	temp := Celsius(25)
//	fahrenheit, _ := ToFahrenheit(temp)  // fahrenheit = 77
func ToFahrenheit(u Unit) (float64, error) {
	if !IsDimension(u, Temperature) {
		return 0, fmt.Errorf("not a temperature unit")
	}
	return (u.Value-273.15)*9/5 + 32, nil
}

// Data storage units

// Megabytes creates a data unit in megabytes.
//
// Example:
//
//	size := Megabytes(15)  // 15 MB
func Megabytes(n float64) Unit { return New(n, "MB") }

// Gigabytes creates a data unit in gigabytes.
//
// Example:
//
//	size := Gigabytes(1.5)  // 1.5 GB
func Gigabytes(n float64) Unit { return New(n, "GB") }

// Terabytes creates a data unit in terabytes.
//
// Example:
//
//	size := Terabytes(2)  // 2 TB
func Terabytes(n float64) Unit { return New(n, "TB") }

// Kibibytes creates a data unit in kibibytes (1024 bytes).
//
// Example:
//
//	size := Kibibytes(1024)  // 1024 KiB = 1 MiB
func Kibibytes(n float64) Unit { return New(n, "KiB") }

// Mebibytes creates a data unit in mebibytes (1024 KiB).
//
// Example:
//
//	size := Mebibytes(2)  // 2 MiB
func Mebibytes(n float64) Unit { return New(n, "MiB") }

// Gibibytes creates a data unit in gibibytes (1024 MiB).
//
// Example:
//
//	size := Gibibytes(1)  // 1 GiB
func Gibibytes(n float64) Unit { return New(n, "GiB") }

// Tebibytes creates a data unit in tebibytes (1024 GiB).
//
// Example:
//
//	size := Tebibytes(0.5)  // 0.5 TiB
func Tebibytes(n float64) Unit { return New(n, "TiB") }

// Electrical and physical units

// Watts creates a power unit in watts.
//
// Example:
//
//	power := Watts(100)  // 100 W
func Watts(n float64) Unit { return New(n, "W") }

// Volts creates a voltage unit in volts.
//
// Example:
//
//	voltage := Volts(220)  // 220 V
func Volts(n float64) Unit { return New(n, "V") }

// Amperes creates a current unit in amperes.
//
// Example:
//
//	current := Amperes(0.5)  // 0.5 A
func Amperes(n float64) Unit { return New(n, "A") }

// Newtons creates a force unit in newtons.
//
// Example:
//
//	force := Newtons(10)  // 10 N
func Newtons(n float64) Unit { return New(n, "N") }

// Pascals creates a pressure unit in pascals.
//
// Example:
//
//	pressure := Pascals(101325)  // 101325 Pa = 1 atm
func Pascals(n float64) Unit { return New(n, "Pa") }

// Psi creates a pressure unit in pounds per square inch (psi) converted to pascals.
//
// Example:
//
//	pressure := Psi(14.7)  // 14.7 psi = 101356.5 Pa ≈ 1 atm
func Psi(n float64) Unit { return New(n*6894.76, "Pa") }

// Joules creates an energy unit in joules.
//
// Example:
//
//	energy := Joules(4184)  // 4184 J = 1 kcal
func Joules(n float64) Unit { return New(n, "J") }

// Hertzs creates a frequency unit in hertz.
//
// Example:
//
//	frequency := Hertzs(60)  // 60 Hz
func Hertzs(n float64) Unit { return New(n, "Hz") }

// IsDimension checks if a Unit has the expected Dimension.
// This is useful for validating that sensor readings have the correct physical quantity.
//
// Examples:
//
//	temperatureSensor, _ := Parse("32.5 °C")
//	if IsDimension(temperatureSensor, Temperature) {
//	    // Process temperature reading
//	}
//
//	// Check that a calculation result is a power value
//	if IsDimension(result, Watt.Dimension) {
//	    fmt.Println("Power calculation result:", result)
//	}
func IsDimension(u Unit, expected Dimension) bool {
	return u.Dimension == expected
}

// Pressure conversion helpers

// ToKiloPascals converts a pressure unit to kilopascals.
// Returns an error if the unit is not a pressure.
//
// Example:
//
//	pressure := Pascals(101325)
//	kPa, _ := ToKiloPascals(pressure)  // kPa = 101.325
func ToKiloPascals(u Unit) (float64, error) {
	if !IsDimension(u, Pascal.Dimension) {
		return 0, fmt.Errorf("not a pressure unit")
	}
	return u.Value / 1000, nil
}

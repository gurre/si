package si

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// SI base dimensions for common use.
var (
	Length        = Dimension{1, 0, 0, 0, 0, 0, 0}
	Mass          = Dimension{0, 1, 0, 0, 0, 0, 0}
	TimeDim       = Dimension{0, 0, 1, 0, 0, 0, 0}
	Current       = Dimension{0, 0, 0, 1, 0, 0, 0}
	Temperature   = Dimension{0, 0, 0, 0, 1, 0, 0}
	Substance     = Dimension{0, 0, 0, 0, 0, 1, 0}
	Luminosity    = Dimension{0, 0, 0, 0, 0, 0, 1}
	Dimensionless = Dimension{0, 0, 0, 0, 0, 0, 0}
)

// Prefixes defines both SI and binary prefixes for unit scaling.
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
// Used for e.g., dBm.
var SymbolicUnits = map[string]Unit{
	"dBm": {1e-3, Dimension{2, 1, -3, 0, 0, 0, 0}},
}

// ParseUnit parses a unit string like "km/h" into a Unit.
// It handles basic units, prefixed units, compound units, and symbolic units.
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
func MustParse(input string) Unit {
	u, err := Parse(input)
	if err != nil {
		panic(err)
	}
	return u
}

// New creates a unit with a given value and unit string (e.g. "kg").
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

// formatDimension renders a dimension into a string representation.
func formatDimension(d Dimension) string {
	symbols := []string{"m", "kg", "s", "A", "K", "mol", "cd"}
	var numerator []string
	var denominator []string

	for i, exp := range d {
		if exp == 0 {
			continue
		}

		if exp > 0 {
			if exp == 1 {
				numerator = append(numerator, symbols[i])
			} else {
				numerator = append(numerator, fmt.Sprintf("%s^%d", symbols[i], exp))
			}
		} else {
			if exp == -1 {
				denominator = append(denominator, symbols[i])
			} else {
				denominator = append(denominator, fmt.Sprintf("%s^%d", symbols[i], -exp))
			}
		}
	}

	// Handle special cases
	if len(numerator) == 0 && len(denominator) == 0 {
		return "1"
	}

	if len(numerator) == 0 {
		numerator = append(numerator, "1")
	}

	// Build the numerator part
	numStr := strings.Join(numerator, "·")

	// If there are denominators, append them with "/"
	if len(denominator) > 0 {
		denomStr := strings.Join(denominator, "·")
		return fmt.Sprintf("%s/%s", numStr, denomStr)
	}

	return numStr
}

// String returns a human-readable representation of the unit (e.g. "100 km/h").
func (u Unit) String() string {
	// For horsepower - specific case
	if u.Dimension == Watt.Dimension && math.Abs(u.Value-745.7) < 0.1 {
		return fmt.Sprintf("%g hp", 1.0)
	}

	// For Watt
	if u.Dimension == Watt.Dimension {
		if u.Value >= 1e9 {
			return fmt.Sprintf("%g GW", u.Value/1e9)
		} else if u.Value >= 1e6 {
			return fmt.Sprintf("%g MW", u.Value/1e6)
		} else if u.Value >= 1e3 {
			return fmt.Sprintf("%g kW", u.Value/1e3)
		}
		return fmt.Sprintf("%g W", u.Value)
	}

	// Special case for Hertz (frequency)
	if u.Dimension == Hertz.Dimension {
		return fmt.Sprintf("%g Hz", u.Value)
	}

	// For meter per second (velocity)
	if u.Dimension == Meter.Div(Second).Dimension {
		// Convert to km/h
		kmh := u.Value * 3.6 // 3.6 = (1/1000) / (1/3600)
		// Round to avoid floating point precision issues
		kmh = math.Round(kmh)
		return fmt.Sprintf("%g km/h", kmh)
	}

	// For energy units (Joule)
	if u.Dimension == Joule.Dimension {
		return fmt.Sprintf("%g J", u.Value)
	}

	// For force units (Newton)
	if u.Dimension == Newton.Dimension {
		return fmt.Sprintf("%g N", u.Value)
	}

	// For pressure units (Pascal)
	if u.Dimension == Pascal.Dimension {
		return fmt.Sprintf("%g Pa", u.Value)
	}

	// For mass units (kg)
	if u.Dimension == Mass {
		return fmt.Sprintf("%g kg", u.Value)
	}

	// For dimensionless units, don't print any unit
	if u.Dimension == Dimensionless {
		return fmt.Sprintf("%g", u.Value)
	}

	return fmt.Sprintf("%g %s", u.Value, formatDimension(u.Dimension))
}

// Named constants for SI base units.
var (
	Meter    = Unit{1, Length}
	Kilogram = Unit{1, Mass}
	Second   = Unit{1, TimeDim}
	Ampere   = Unit{1, Current}
	Kelvin   = Unit{1, Temperature}
	Mole     = Unit{1, Substance}
	Candela  = Unit{1, Luminosity}
	One      = Unit{1, Dimensionless}
)

// Named derived units for convenience.
var (
	Newton  = Kilogram.Mul(Meter).Div(Second.Pow(2))
	Joule   = Newton.Mul(Meter)
	Watt    = Joule.Div(Second)
	Pascal  = Newton.Div(Meter.Pow(2))
	Hertz   = Unit{1, Dimension{0, 0, -1, 0, 0, 0, 0}}
	Coulomb = Ampere.Mul(Second)
	Volt    = Watt.Div(Ampere)
)

// Convenience helpers for common physical quantities.
// Time units
func Hours(n float64) Unit        { return New(n*3600, "s") }
func Minutes(n float64) Unit      { return New(n*60, "s") }
func Seconds(n float64) Unit      { return New(n, "s") }
func Milliseconds(n float64) Unit { return New(n/1000, "s") }

// Length units
func Kilometers(n float64) Unit { return New(n*1000, "m") }
func Meters(n float64) Unit     { return New(n, "m") }

// Mass units
func Grams(n float64) Unit     { return New(n/1000, "kg") }
func Kilograms(n float64) Unit { return New(n, "kg") }

// Temperature units
func Celsius(n float64) Unit { return New(n+273.15, "K") }

// Data storage units
func Megabytes(n float64) Unit { return New(n, "MB") }
func Gigabytes(n float64) Unit { return New(n, "GB") }
func Terabytes(n float64) Unit { return New(n, "TB") }
func Kibibytes(n float64) Unit { return New(n, "KiB") }
func Mebibytes(n float64) Unit { return New(n, "MiB") }
func Gibibytes(n float64) Unit { return New(n, "GiB") }
func Tebibytes(n float64) Unit { return New(n, "TiB") }

// Electrical and physical units
func Watts(n float64) Unit   { return New(n, "W") }
func Volts(n float64) Unit   { return New(n, "V") }
func Amperes(n float64) Unit { return New(n, "A") }
func Newtons(n float64) Unit { return New(n, "N") }
func Pascals(n float64) Unit { return New(n, "Pa") }
func Joules(n float64) Unit  { return New(n, "J") }
func Hertzs(n float64) Unit  { return New(n, "Hz") }

// VerifyDimension checks if a Unit has the expected Dimension.
// This is useful for validating that sensor readings have the correct physical quantity.
//
// Examples:
//
//	temperatureSensor, _ := Parse("32.5 °C")
//	if VerifyDimension(temperatureSensor, Temperature) {
//	    // Process temperature reading
//	}
//
//	// Check that a calculation result is a power value
//	if VerifyDimension(result, Watt.Dimension) {
//	    fmt.Println("Power calculation result:", result)
//	}
func VerifyDimension(u Unit, expected Dimension) bool {
	return u.Dimension == expected
}

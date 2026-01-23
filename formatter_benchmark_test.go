package si_test

import (
	"testing"

	"github.com/gurre/si"
)

// BenchmarkFormatUnit benchmarks the standard formatter with various unit types
func BenchmarkFormatUnit(b *testing.B) {
	units := []si.Unit{
		si.Scalar(42),           // dimensionless
		si.Meter,                // base unit
		si.Kilogram,             // base unit with prefix in name
		si.Meter.Div(si.Second), // derived unit (velocity)
		si.Kilogram.Mul(si.Meter).Div(si.Second.MustPow(2)),              // derived unit (force)
		si.Kilogram.Mul(si.Meter.MustPow(2)).Div(si.Second.MustPow(2)),       // derived unit (energy)
		si.Kilogram.Div(si.Meter.MustPow(3)),                             // derived unit (density)
		si.Newton.Mul(si.Meter).Div(si.Second),                       // compound derived unit
		si.Meter.MustPow(3).Div(si.Second).Div(si.Pascal).Mul(si.Kelvin), // complex unit
	}

	for i := 0; i < b.N; i++ {
		unit := units[i%len(units)]
		_ = si.FormatUnit(unit)
	}
}

// BenchmarkFormatUnitWithPrefix benchmarks the prefixed formatter with various unit types
func BenchmarkFormatUnitWithPrefix(b *testing.B) {
	units := []si.Unit{
		si.Scalar(0.000125),                               // small dimensionless
		si.Scalar(42000),                                  // large dimensionless
		si.Meters(0.0000025),                              // small length
		si.Meters(1500),                                   // medium length
		si.Meters(3500000),                                // large length
		si.Kelvin.Mul(si.Scalar(0.001)),                   // small temperature
		si.Kelvin.Mul(si.Scalar(9500)),                    // large temperature
		si.Watt.Mul(si.Scalar(0.015)),                     // small power
		si.Watt.Mul(si.Scalar(2500)),                      // medium power
		si.Watt.Mul(si.Scalar(1500000)),                   // large power
		si.Pascal.Mul(si.Scalar(101325)),                  // atmospheric pressure
		si.Hertz.Mul(si.Scalar(2500000000)),               // GHz range
		si.Meter.Div(si.Second).Mul(si.Scalar(299792458)), // speed of light
	}

	for i := 0; i < b.N; i++ {
		unit := units[i%len(units)]
		_ = si.FormatUnitWithPrefix(unit)
	}
}

// BenchmarkUnitString benchmarks the String method which calls FormatUnitWithPrefix
func BenchmarkUnitString(b *testing.B) {
	units := []si.Unit{
		si.Meters(1.5),                             // simple unit
		si.Kelvins(273.15),                         // temperature unit
		si.Seconds(0.0000015),                      // small time unit
		si.Watts(1200),                             // power unit
		si.Newton.Div(si.Meter.MustPow(2)),             // pressure
		si.Meter.Div(si.Second).Mul(si.Scalar(25)), // velocity
		si.Scalar(0.000000001),                     // very small scalar
		si.Scalar(1000000000),                      // very large scalar
		si.Joule.Div(si.Kelvin.Mul(si.Kilogram)),   // specific heat
	}

	for i := 0; i < b.N; i++ {
		unit := units[i%len(units)]
		_ = unit.String()
	}
}

// BenchmarkFormatComplexDimension benchmarks formatting complex dimensions
func BenchmarkFormatComplexDimension(b *testing.B) {
	units := []si.Unit{
		// Thermal conductivity: W/(m·K)
		si.Watt.Div(si.Meter.Mul(si.Kelvin)),
		// Stefan-Boltzmann constant: W/(m²·K⁴)
		si.Watt.Div(si.Meter.MustPow(2).Mul(si.Kelvin.MustPow(4))),
		// Dynamic viscosity: Pa·s
		si.Pascal.Mul(si.Second),
		// Moment of inertia: kg·m²
		si.Kilogram.Mul(si.Meter.MustPow(2)),
		// Specific heat capacity: J/(kg·K)
		si.Joule.Div(si.Kilogram.Mul(si.Kelvin)),
		// Thermal diffusivity: m²/s
		si.Meter.MustPow(2).Div(si.Second),
		// Energy density: J/m³
		si.Joule.Div(si.Meter.MustPow(3)),
		// Complex composition: kg·m²/(s³·A)
		si.Kilogram.Mul(si.Meter.MustPow(2)).Div(si.Second.MustPow(3).Mul(si.Ampere)),
	}

	for i := 0; i < b.N; i++ {
		unit := units[i%len(units)]
		_ = unit.String()
	}
}

package si_test

import (
	"math"
	"testing"

	"github.com/gurre/si"
)

// BenchmarkParse benchmarks parsing of unit expressions
func BenchmarkParse(b *testing.B) {
	expressions := []string{
		"10 m",
		"100 km/h",
		"9.81 m/s^2",
		"101.325 kPa",
		"1.21 kg/m^3",
		"25 K",
		"5 W/(m^2*K)",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		expr := expressions[i%len(expressions)]
		_, err := si.Parse(expr)
		if err != nil {
			b.Fatalf("Failed to parse %s: %v", expr, err)
		}
	}
}

// BenchmarkMul benchmarks multiplication of units
func BenchmarkMul(b *testing.B) {
	u1 := si.Meter
	u2 := si.Second

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = u1.Mul(u2)
	}
}

// BenchmarkDiv benchmarks division of units
func BenchmarkDiv(b *testing.B) {
	u1 := si.Meter
	u2 := si.Second

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = u1.Div(u2)
	}
}

// BenchmarkPow benchmarks raising a unit to a power
func BenchmarkPow(b *testing.B) {
	u := si.Meter

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = u.Pow(2)
	}
}

// BenchmarkAdd benchmarks adding units
func BenchmarkAdd(b *testing.B) {
	u1 := si.Meters(10)
	u2 := si.Meters(5)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := u1.Add(u2)
		if err != nil {
			b.Fatalf("Failed to add units: %v", err)
		}
	}
}

// BenchmarkConvertTo benchmarks unit conversion
func BenchmarkConvertTo(b *testing.B) {
	u1 := si.Meters(1000) // 1 km
	u2 := si.Meter        // 1 m

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := u1.ConvertTo(u2)
		if err != nil {
			b.Fatalf("Failed to convert units: %v", err)
		}
	}
}

// BenchmarkString benchmarks string conversion
func BenchmarkString(b *testing.B) {
	units := []si.Unit{
		si.Meters(100),
		si.Watts(1500),
		si.Celsius(25),
		si.Newton.Mul(si.Meter),
		si.Meter.Div(si.Second).Mul(si.Scalar(10)),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u := units[i%len(units)]
		_ = u.String()
	}
}

// BenchmarkComplexCalculation benchmarks a complex calculation with multiple operations
func BenchmarkComplexCalculation(b *testing.B) {
	// Thermal energy calculation: Q = m * c * ΔT
	m := si.Kilograms(5)
	c := si.Joules(4186).Div(si.Kilogram.Mul(si.Kelvin))
	deltaT := si.Kelvin.Mul(si.Scalar(20))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Mul(c).Mul(deltaT)
	}
}

// BenchmarkHydraulicPower benchmarks hydraulic power calculation
func BenchmarkHydraulicPower(b *testing.B) {
	// P = p * Q (pressure * flow rate)
	p := si.Pascals(5e6)                                      // 5 MPa
	Q := si.Meter.Pow(3).Mul(si.Scalar(0.001)).Div(si.Second) // 1 L/s

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = p.Mul(Q)
	}
}

// BenchmarkReynoldsNumber benchmarks Reynolds number calculation
func BenchmarkReynoldsNumber(b *testing.B) {
	// Re = (ρ * v * D) / μ
	rho := si.Kilograms(1000).Div(si.Meter.Pow(3))       // Water density
	v := si.Meters(1.5).Div(si.Second)                   // Flow velocity
	D := si.Meters(0.05)                                 // 50 mm pipe
	mu := si.Pascal.Mul(si.Second).Mul(si.Scalar(0.001)) // 0.001 Pa·s

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rho.Mul(v).Mul(D).Div(mu)
	}
}

// BenchmarkVerifyDimension benchmarks dimension verification
func BenchmarkVerifyDimension(b *testing.B) {
	units := []struct {
		unit si.Unit
		dim  si.Dimension
	}{
		{si.Meters(10), si.Length},
		{si.Celsius(25), si.Temperature},
		{si.Watts(100), si.Watt.Dimension},
		{si.Meter.Div(si.Second), si.Meter.Div(si.Second).Dimension},
		{si.Meter.Pow(3).Div(si.Second), si.Meter.Pow(3).Div(si.Second).Dimension},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pair := units[i%len(units)]
		_ = si.VerifyDimension(pair.unit, pair.dim)
	}
}

// BenchmarkEndToEndFlowCalculation benchmarks a complete flow calculation
func BenchmarkEndToEndFlowCalculation(b *testing.B) {
	// Setup a complete flow calculation that represents typical usage
	// 1. Parse input values
	// 2. Perform calculations
	// 3. Convert to desired output format

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 1. Parse input values (usually from user/config)
		diameter, _ := si.Parse("100 mm")
		velocity, _ := si.Parse("2.5 m/s")
		density, _ := si.Parse("1000 kg/m^3")

		// 2. Calculate cross-sectional area: A = π(D/2)²
		radius := diameter.Mul(si.Scalar(0.5))
		area := radius.Pow(2).Mul(si.Scalar(math.Pi))

		// 3. Calculate volume flow rate: Q = A * v
		flowRate := area.Mul(velocity)

		// 4. Calculate mass flow rate: ṁ = Q * ρ
		massFlow := flowRate.Mul(density)

		// 5. Convert to desired output units (kg/s)
		_, err := massFlow.ConvertTo(si.Kilograms(1).Div(si.Second))
		if err != nil {
			b.Fatalf("Failed to convert mass flow: %v", err)
		}
	}
}

// BenchmarkEndToEndEnergyCalculation benchmarks a complete energy calculation
func BenchmarkEndToEndEnergyCalculation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 1. Parse input values
		power, _ := si.Parse("1.5 kW")
		time, _ := si.Parse("8 h")

		// 2. Calculate energy: E = P * t
		energy := power.Mul(time)

		// 3. Convert to kWh
		kilowattHour := si.Watt.Mul(si.Scalar(1000)).Mul(si.Hours(1))
		_, err := energy.ConvertTo(kilowattHour)
		if err != nil {
			b.Fatalf("Failed to convert energy: %v", err)
		}
	}
}

// BenchmarkCelsiusToKelvin benchmarks temperature conversion
func BenchmarkCelsiusToKelvin(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		temp := si.Celsius(25)
		_, err := temp.ConvertTo(si.Kelvin)
		if err != nil {
			b.Fatalf("Failed to convert temperature: %v", err)
		}
	}
}

// BenchmarkKilometersToMeters benchmarks length conversion
func BenchmarkKilometersToMeters(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		length := si.Kilometers(5)
		_, err := length.ConvertTo(si.Meter)
		if err != nil {
			b.Fatalf("Failed to convert length: %v", err)
		}
	}
}

// BenchmarkMarshalJSON benchmarks JSON marshalling
func BenchmarkMarshalJSON(b *testing.B) {
	units := []si.Unit{
		si.Meters(100),
		si.Watts(1500),
		si.Celsius(25),
		si.Newton.Mul(si.Meter),
		si.Meter.Div(si.Second).Mul(si.Scalar(10)),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u := units[i%len(units)]
		_, err := u.MarshalJSON()
		if err != nil {
			b.Fatalf("Failed to marshal unit to JSON: %v", err)
		}
	}
}

// BenchmarkPressureConversion benchmarks pressure unit conversion
func BenchmarkPressureConversion(b *testing.B) {
	// 1 atm = 101325 Pa
	pressure := si.Pascals(101325)
	atm := si.Pascal.Mul(si.Scalar(101325))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := pressure.ConvertTo(atm)
		if err != nil {
			b.Fatalf("Failed to convert pressure: %v", err)
		}
	}
}

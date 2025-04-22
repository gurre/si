package si

import (
	"encoding/json"
	"math"
	"testing"
)

// TestJSONMarshalling tests the JSON marshalling and unmarshalling of Unit values
func TestJSONMarshalling(t *testing.T) {
	tests := []struct {
		name     string
		unit     Unit
		expected string
	}{
		{"velocity", Meter.Div(Second).Mul(Scalar(10)), `"36 km/h"`},
		{"energy", Joule.Mul(Scalar(5000)), `"5000 J"`},
		{"power", Watt.Mul(Scalar(1500)), `"1.5 kW"`},
		{"dimensionless", Scalar(0.75), `"0.75"`},
		{"temperature", Celsius(25), `"298.15 K"`},
		{"pressure", Pascal.Mul(Scalar(101325)), `"101325 Pa"`},
		{"force", Newton.Mul(Scalar(50)), `"50 N"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshalling
			data, err := json.Marshal(tt.unit)
			if err != nil {
				t.Fatalf("Failed to marshal unit: %v", err)
			}

			if string(data) != tt.expected {
				t.Errorf("JSON marshalled data = %s, want %s", string(data), tt.expected)
			}

			// Test unmarshalling
			var unmarshalled Unit
			if err := json.Unmarshal(data, &unmarshalled); err != nil {
				t.Fatalf("Failed to unmarshal unit: %v", err)
			}

			// Compare dimensions
			if unmarshalled.Dimension != tt.unit.Dimension {
				t.Errorf("Dimension = %v, want %v", unmarshalled.Dimension, tt.unit.Dimension)
			}

			// Compare values with tolerance
			if !unmarshalled.Equals(tt.unit) {
				t.Errorf("Value = %v, want %v", unmarshalled.Value, tt.unit.Value)
			}
		})
	}
}

// TestJSONStructMarshalling tests marshalling Units inside a struct
func TestJSONStructMarshalling(t *testing.T) {
	type Measurement struct {
		Name  string `json:"name"`
		Value Unit   `json:"value"`
		Time  string `json:"time"`
	}

	measurement := Measurement{
		Name:  "speed",
		Value: Kilometers(100).Div(Hours(1)), // 100 km/h
		Time:  "2023-01-01T12:00:00Z",
	}

	expected := `{"name":"speed","value":"100 km/h","time":"2023-01-01T12:00:00Z"}`

	data, err := json.Marshal(measurement)
	if err != nil {
		t.Fatalf("Failed to marshal measurement: %v", err)
	}

	if string(data) != expected {
		t.Errorf("JSON marshalled data = %s, want %s", string(data), expected)
	}

	// Test unmarshalling
	var unmarshalled Measurement
	if err := json.Unmarshal(data, &unmarshalled); err != nil {
		t.Fatalf("Failed to unmarshal measurement: %v", err)
	}

	if unmarshalled.Name != measurement.Name {
		t.Errorf("Name = %s, want %s", unmarshalled.Name, measurement.Name)
	}

	if unmarshalled.Time != measurement.Time {
		t.Errorf("Time = %s, want %s", unmarshalled.Time, measurement.Time)
	}

	// Compare the unit value with tolerance
	if !unmarshalled.Value.Equals(measurement.Value) {
		t.Errorf("Value = %v, want %v", unmarshalled.Value, measurement.Value)
	}
}

// TestConvertMetersToCentimeters verifies converting meters to centimeters
func TestConvertMetersToCentimeters(t *testing.T) {
	unit := Meter.Mul(Scalar(2))      // 2 m
	target := Meter.Mul(Scalar(0.01)) // 1 cm
	expected := 200.0                 // 2m = 200cm

	result, err := unit.ConvertTo(target)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Round to 10 decimal places
	rounded := math.Round(result.Value*1e10) / 1e10
	expected = math.Round(expected*1e10) / 1e10

	if rounded != expected {
		t.Errorf("Value = %v, want %v", result.Value, expected)
	}

	if result.Dimension != unit.Dimension {
		t.Errorf("Dimension = %v, want %v", result.Dimension, unit.Dimension)
	}
}

// TestConvertKilometersToMeters verifies converting kilometers to meters
func TestConvertKilometersToMeters(t *testing.T) {
	unit := Meter.Mul(Scalar(1000)) // 1 km
	target := Meter                 // 1 m
	expected := 1000.0

	result, err := unit.ConvertTo(target)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Round to 10 decimal places
	rounded := math.Round(result.Value*1e10) / 1e10
	expected = math.Round(expected*1e10) / 1e10

	if rounded != expected {
		t.Errorf("Value = %v, want %v", result.Value, expected)
	}

	if result.Dimension != unit.Dimension {
		t.Errorf("Dimension = %v, want %v", result.Dimension, unit.Dimension)
	}
}

// TestConvertCelsiusToKelvin verifies converting Celsius to Kelvin
func TestConvertCelsiusToKelvin(t *testing.T) {
	unit := Celsius(100) // 100°C = 373.15K
	target := Kelvin     // 1K
	expected := 373.15

	result, err := unit.ConvertTo(target)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Round to 10 decimal places
	rounded := math.Round(result.Value*1e10) / 1e10
	expected = math.Round(expected*1e10) / 1e10

	if rounded != expected {
		t.Errorf("Value = %v, want %v", result.Value, expected)
	}

	if result.Dimension != unit.Dimension {
		t.Errorf("Dimension = %v, want %v", result.Dimension, unit.Dimension)
	}
}

// TestConvertNewtonToBaseUnits verifies converting Newton to kg·m/s²
func TestConvertNewtonToBaseUnits(t *testing.T) {
	unit := Newton.Mul(Scalar(5))                    // 5N
	target := Kilogram.Mul(Meter).Div(Second.Pow(2)) // 1 kg·m/s²
	expected := 5.0

	result, err := unit.ConvertTo(target)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Round to 10 decimal places
	rounded := math.Round(result.Value*1e10) / 1e10
	expected = math.Round(expected*1e10) / 1e10

	if rounded != expected {
		t.Errorf("Value = %v, want %v", result.Value, expected)
	}

	if result.Dimension != unit.Dimension {
		t.Errorf("Dimension = %v, want %v", result.Dimension, unit.Dimension)
	}
}

// TestConvertJouleToKilowattHour verifies converting Joule to kilowatt-hour
func TestConvertJouleToKilowattHour(t *testing.T) {
	unit := Joule.Mul(Scalar(3600000))             // 3600000 J = 1 kWh
	target := Watt.Mul(Scalar(1000)).Mul(Hours(1)) // 1 kWh
	expected := 1.0

	result, err := unit.ConvertTo(target)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Round to 10 decimal places
	rounded := math.Round(result.Value*1e10) / 1e10
	expected = math.Round(expected*1e10) / 1e10

	if rounded != expected {
		t.Errorf("Value = %v, want %v", result.Value, expected)
	}

	if result.Dimension != unit.Dimension {
		t.Errorf("Dimension = %v, want %v", result.Dimension, unit.Dimension)
	}
}

// TestConvertPascalToAtmospheres verifies converting Pascal to atmospheres
func TestConvertPascalToAtmospheres(t *testing.T) {
	unit := Pascal.Mul(Scalar(101325))   // 1 atm = 101325 Pa
	target := Pascal.Mul(Scalar(101325)) // 1 atm
	expected := 1.0

	result, err := unit.ConvertTo(target)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Round to 10 decimal places
	rounded := math.Round(result.Value*1e10) / 1e10
	expected = math.Round(expected*1e10) / 1e10

	if rounded != expected {
		t.Errorf("Value = %v, want %v", result.Value, expected)
	}

	if result.Dimension != unit.Dimension {
		t.Errorf("Dimension = %v, want %v", result.Dimension, unit.Dimension)
	}
}

// TestConvertThermalConductance verifies converting thermal conductance units
func TestConvertThermalConductance(t *testing.T) {
	// W/m²·K (thermal conductance)
	unit := Watt.Div(Meter.Pow(2).Mul(Kelvin)).Mul(Scalar(5.678)) // 5.678 W/m²·K
	target := Watt.Div(Meter.Pow(2).Mul(Kelvin))                  // Base unit
	expected := 5.678

	result, err := unit.ConvertTo(target)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Round to 10 decimal places
	rounded := math.Round(result.Value*1e10) / 1e10
	expected = math.Round(expected*1e10) / 1e10

	if rounded != expected {
		t.Errorf("Value = %v, want %v", result.Value, expected)
	}

	if result.Dimension != unit.Dimension {
		t.Errorf("Dimension = %v, want %v", result.Dimension, unit.Dimension)
	}
}

// TestConvertIncompatibleDimensions verifies error on incompatible dimensions
func TestConvertIncompatibleDimensions(t *testing.T) {
	unit := Meter    // Length
	target := Second // Time

	_, err := unit.ConvertTo(target)
	if err == nil {
		t.Error("Expected error for incompatible dimensions but got none")
	}
}

// TestConvertZeroTargetValue verifies error on zero target value
func TestConvertZeroTargetValue(t *testing.T) {
	unit := Meter
	target := Meter.Mul(Scalar(0)) // 0 m

	_, err := unit.ConvertTo(target)
	if err == nil {
		t.Error("Expected error for zero target value but got none")
	}
}

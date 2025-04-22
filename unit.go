package si

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"math"
	"strconv"
	"strings"
)

// Dimension represents the exponents of the 7 SI base units.
// The index positions are: [Length, Mass, Time, Current, Temperature, Substance, Luminosity].
// For example, a meter is Dimension{1,0,0,0,0,0,0} and a second is Dimension{0,0,1,0,0,0,0}.
type Dimension [7]int

// Unit represents a physical quantity with a value and dimension
// This is the core type of the package, combining a numeric value with its physical dimension.
// All operations on physical quantities in this package use this type.
type Unit struct {
	// Value is the scalar magnitude of the physical quantity in SI base units
	Value float64
	// Dimension contains the exponents for each of the 7 SI base dimensions
	Dimension Dimension
}

// Scalar creates a dimensionless unit
// This is useful for pure numbers or scaling factors without physical dimensions.
//
// Example:
//
//	efficiency := Scalar(0.85)  // 85% efficiency, dimensionless
//	factor := Scalar(2.5)       // scaling factor of 2.5
func Scalar(value float64) Unit {
	return Unit{Value: value, Dimension: Dimension{}}
}

// Mul multiplies two units
// This is used to combine physical quantities, such as multiplying mass and acceleration to get force.
// The dimensions are added and the values are multiplied.
//
// Example:
//
//	// Calculate force: F = m × a
//	mass := Kilograms(75)
//	acceleration := Meters(9.81).Div(Second.Pow(2))
//	force := mass.Mul(acceleration)  // 735.75 N
func (u Unit) Mul(v Unit) Unit {
	var dim Dimension
	for i := range dim {
		dim[i] = u.Dimension[i] + v.Dimension[i]
	}
	return Unit{
		Value:     u.Value * v.Value,
		Dimension: dim,
	}
}

// Div divides two units
// This is used to derive physical quantities, such as dividing distance by time to get velocity.
// The dimensions are subtracted and the values are divided.
//
// Example:
//
//	// Calculate speed: v = d ÷ t
//	distance := Kilometers(60)
//	time := Minutes(30)
//	speed := distance.Div(time)  // 33.33 m/s
func (u Unit) Div(v Unit) Unit {
	var dim Dimension
	for i := range dim {
		dim[i] = u.Dimension[i] - v.Dimension[i]
	}
	return Unit{
		Value:     u.Value / v.Value,
		Dimension: dim,
	}
}

// Pow raises a unit to a power
// This is used for operations like squaring distances or taking the square root of an area.
// All dimension exponents are multiplied by the given power.
//
// Example:
//
//	// Calculate area from length: A = l²
//	length := Meters(4)
//	area := length.Pow(2)  // 16 m²
//
//	// Calculate radius from volume: r = ∛(3V/4π)
//	volumeTerm := Meter.Pow(3).Mul(Scalar(3.0/4.0/math.Pi))
//	radius := volumeTerm.Pow(1.0/3.0)
func (u Unit) Pow(exp int) Unit {
	var dim Dimension
	for i := range dim {
		dim[i] = u.Dimension[i] * exp
	}
	return Unit{
		Value:     pow(u.Value, exp),
		Dimension: dim,
	}
}

// pow is a helper function to calculate x^n for integer exponents
// This implements an efficient integer power algorithm to avoid
// floating-point imprecision in math.Pow for integer exponents
func pow(x float64, n int) float64 {
	if n == 0 {
		return 1
	}
	if n < 0 {
		return 1 / pow(x, -n)
	}
	if n%2 == 0 {
		half := pow(x, n/2)
		return half * half
	}
	return x * pow(x, n-1)
}

// Compare returns -1, 0, 1 if u <, ==, > v respectively. Returns error if dimensions differ.
// This is used to compare units with the same dimensions, such as comparing two lengths.
//
// Example:
//
//	dist1 := Kilometers(1)
//	dist2 := Meters(1500)
//	result, _ := dist1.Compare(dist2) // result will be -1 because 1 km < 1.5 km
func (u Unit) Compare(v Unit) (int, error) {
	if u.Dimension != v.Dimension {
		return 0, errors.New("cannot compare units with different dimensions")
	}
	switch {
	case math.Abs(u.Value-v.Value) < 1e-12:
		return 0, nil
	case u.Value < v.Value:
		return -1, nil
	default:
		return 1, nil
	}
}

// MarshalJSON encodes the unit as a string like "100 km/h".
// This enables JSON serialization of SI units with their dimensions and prefixes.
//
// Example:
//
//	type Reading struct {
//	    Pressure Unit `json:"pressure"`
//	}
//	reading := Reading{Pressure: Pascals(101325)}
//	data, _ := json.Marshal(reading) // {"pressure":"101.325 kPa"}
func (u Unit) MarshalJSON() ([]byte, error) {
	return json.Marshal(FormatUnitWithPrefix(u))
}

// UnmarshalJSON parses a unit string from JSON like "100 km/h".
// This enables JSON deserialization of SI units back into native Unit objects.
//
// Example:
//
//	var reading Reading
//	err := json.Unmarshal([]byte(`{"pressure":"101.325 kPa"}`), &reading)
//	// reading.Pressure will be 101325 Pa
func (u *Unit) UnmarshalJSON(data []byte) error {
	var input string
	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}

	parsed, err := Parse(input)
	if err != nil {
		return err
	}

	u.Value = parsed.Value
	u.Dimension = parsed.Dimension
	return nil
}

// MarshalXML encodes the unit as an XML element with value and dimension attributes.
// This enables XML serialization of SI units with detailed dimension information.
//
// Example:
//
//	type Reading struct {
//	    Pressure Unit `xml:"pressure"`
//	}
//	reading := Reading{Pressure: Pascals(101325)}
//	data, _ := xml.Marshal(reading)
//	// <Reading><pressure value="101325" dimension="M·L^-1·T^-2" display="101.325 kPa"/></Reading>
func (u Unit) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type xmlUnit struct {
		XMLName   xml.Name `xml:"unit"`
		Value     float64  `xml:"value,attr"`
		Dimension string   `xml:"dimension,attr"`
		Display   string   `xml:"display,attr"`
	}

	// Format the dimension as a string
	var dimParts []string
	for i, exp := range u.Dimension {
		if exp != 0 {
			var dimName string
			switch i {
			case 0:
				dimName = "L"
			case 1:
				dimName = "M"
			case 2:
				dimName = "T"
			case 3:
				dimName = "I"
			case 4:
				dimName = "Θ"
			case 5:
				dimName = "N"
			case 6:
				dimName = "J"
			}

			if exp == 1 {
				dimParts = append(dimParts, dimName)
			} else {
				dimParts = append(dimParts, dimName+"^"+strconv.Itoa(exp))
			}
		}
	}

	dimensionStr := strings.Join(dimParts, "·")
	if dimensionStr == "" {
		dimensionStr = "1" // Dimensionless
	}

	xu := xmlUnit{
		Value:     u.Value,
		Dimension: dimensionStr,
		Display:   u.String(),
	}

	return e.Encode(xu)
}

// UnmarshalXML decodes an XML element into a Unit.
// This enables XML deserialization of SI units back into native Unit objects.
//
// Example:
//
//	var reading Reading
//	err := xml.Unmarshal(xmlData, &reading)
//	// reading.Pressure will be the parsed Unit
func (u *Unit) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type xmlUnit struct {
		Display string `xml:"display,attr"`
	}

	var xu xmlUnit
	if err := d.DecodeElement(&xu, &start); err != nil {
		return err
	}

	// Parse the display string (which should be in the format "100 km/h")
	parsed, err := Parse(xu.Display)
	if err != nil {
		return err
	}

	u.Value = parsed.Value
	u.Dimension = parsed.Dimension
	return nil
}

// Add adds two units of the same dimension.
// Returns an error if the dimensions don't match.
// This is used for adding similar physical quantities, like two lengths or two masses.
//
// Example:
//
//	// Add two lengths
//	length1 := Meters(5)
//	length2 := Meters(10)
//	totalLength, _ := length1.Add(length2) // 15 m
//
//	// Error case: different dimensions
//	time := Seconds(30)
//	_, err := length1.Add(time) // err will not be nil
func (u Unit) Add(v Unit) (Unit, error) {
	if u.Dimension != v.Dimension {
		return Unit{}, errors.New("cannot add units with different dimensions")
	}
	return Unit{u.Value + v.Value, u.Dimension}, nil
}

// ConvertTo converts one unit to another unit of the same dimension.
// Returns an error if the dimensions don't match or if division by zero would occur.
// This is used to convert between different units of the same physical dimension.
//
// Example:
//
//	// Convert kilometers to meters
//	distance := Kilometers(5)
//	meters, _ := distance.ConvertTo(Meter) // 5000 m
//
//	// Convert joules to kilowatt-hours
//	energy := Joules(3600000)
//	kWh := Watt.Mul(Scalar(1000)).Mul(Hours(1))
//	result, _ := energy.ConvertTo(kWh) // 1 kWh
func (u Unit) ConvertTo(target Unit) (Unit, error) {
	if u.Dimension != target.Dimension {
		return Unit{}, errors.New("cannot convert units with different dimensions")
	}

	if target.Value == 0 {
		return Unit{}, errors.New("cannot convert to a unit with zero value")
	}

	// Scale factor is u.Value / target.Value
	scaleFactor := u.Value / target.Value

	return Unit{Value: scaleFactor, Dimension: u.Dimension}, nil
}

// Equals compares two units for equality with appropriate tolerance.
// This method accounts for floating-point imprecision when comparing unit values.
//
// Example:
//
//	// Direct comparison with the same unit
//	a := Meters(1.0)
//	b := Meters(1.0)
//	equal := a.Equals(b) // true
//
//	// Comparison with converted unit
//	c := Kilometers(0.001)
//	equal = a.Equals(c) // true
func (u Unit) Equals(v Unit) bool {
	if u.Dimension != v.Dimension {
		return false
	}

	// Use relative epsilon for large values
	eps := 1e-12
	if math.Abs(u.Value) > 1.0 || math.Abs(v.Value) > 1.0 {
		eps = 1e-12 * math.Max(math.Abs(u.Value), math.Abs(v.Value))
	}

	return math.Abs(u.Value-v.Value) < eps
}

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
type Unit struct {
	Value     float64
	Dimension Dimension
}

// Scalar creates a dimensionless unit
func Scalar(value float64) Unit {
	return Unit{Value: value, Dimension: Dimension{}}
}

// Mul multiplies two units
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
func (u Unit) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

// UnmarshalJSON parses a unit string from JSON like "100 km/h".
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
func (u Unit) Add(v Unit) (Unit, error) {
	if u.Dimension != v.Dimension {
		return Unit{}, errors.New("cannot add units with different dimensions")
	}
	return Unit{u.Value + v.Value, u.Dimension}, nil
}

// ConvertTo converts one unit to another unit of the same dimension.
// Returns an error if the dimensions don't match or if division by zero would occur.
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

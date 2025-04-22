package si

// ToSIUnit converts a parser.Unit to a si.Unit
// This function should be implemented in the si package to avoid import cycle
type ToSIUnitFunc func(Unit) interface{}

// FromSIUnit converts a si.Unit to a parser.Unit
// This function should be implemented in the si package to avoid import cycle
type FromSIUnitFunc func(interface{}) Unit

// These functions will be set from the si package during initialization
var (
	toSIUnit   ToSIUnitFunc
	fromSIUnit FromSIUnitFunc
)

// RegisterConversionFunctions sets the conversion functions
func RegisterConversionFunctions(to ToSIUnitFunc, from FromSIUnitFunc) {
	toSIUnit = to
	fromSIUnit = from
}

// ConvertToSIUnit converts a parser.Unit to a si.Unit if the conversion function is set
func ConvertToSIUnit(u Unit) interface{} {
	if toSIUnit != nil {
		return toSIUnit(u)
	}
	return nil
}

// ConvertFromSIUnit converts a si.Unit to a parser.Unit if the conversion function is set
func ConvertFromSIUnit(u interface{}) Unit {
	if fromSIUnit != nil {
		return fromSIUnit(u)
	}
	return Unit{}
}

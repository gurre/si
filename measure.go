package si

import "fmt"

type Measure int

const (
	// None is the empty measure
	None Measure = iota
	// Length is meters
	Length
	// Mass is grams
	Mass
	Time
	ElectricCurrent
	ThermodynamicTemperature
	AmountOfSubstance
	LuminousIntensity
	Angle
	SolidAngle
	Frequency
	Force // or weight
	Pressure
	Energy         // or work, heat
	Power          // or radiant flux
	ElectricCharge // or quantity of electricity
	Voltage        // (electrical potential), emf
	Capacitance
	Impedance // resistance, reactance
	ElectricalConductance
	MagneticFlux
	MagneticFluxDensity
	Inductance
	Temperature // temperature relative to 273.15 K
	LuminousFlux
	Illuminance
	Radioactivity // decays per unit time
	AbsorbedDose
	EquivalentDose
	CatalyticActivity
)

// String returns the Système international unit symbol
func (measure Measure) String() string {
	switch measure {
	case Length:
		return "m"
	case Mass:
		return "kg"
	case Time:
		return "s"
	case ElectricCurrent:
		return "A"
	case ThermodynamicTemperature:
		return "K"
	case AmountOfSubstance:
		return "mol"
	case LuminousIntensity:
		return "cd"
	case Angle:
		return "rad"
	case SolidAngle:
		return "sr"
	case Frequency:
		return "Hz"
	case Force:
		return "N"
	case Pressure:
		return "Pa"
	case Energy:
		return "J"
	case Power:
		return "W"
	case ElectricCharge:
		return "C"
	case Voltage:
		return "V"
	case Capacitance:
		return "F"
	case Impedance:
		return "Ω"
	case ElectricalConductance:
		return "S"
	case MagneticFlux:
		return "Wb"
	case MagneticFluxDensity:
		return "T"
	case Inductance:
		return "H"
	case Temperature:
		return "°C"
	case LuminousFlux:
		return "lm"
	case Illuminance:
		return "lx"
	case Radioactivity:
		return "Bq"
	case AbsorbedDose:
		return "Gy"
	case EquivalentDose:
		return "Sv"
	case CatalyticActivity:
		return "kat"
	case None:
		return ""
	default:
		return ""
	}
}

// Parse takes a string generated from String() and converts it back to a unit.
func Parse(str string) (measure Measure, err error) {
	switch str {
	case "m":
		return Length, nil
	case "kg":
		return Mass, nil
	case "s":
		return Time, nil
	case "A":
		return ElectricCurrent, nil
	case "K":
		return ThermodynamicTemperature, nil
	case "mol":
		return AmountOfSubstance, nil
	case "cd":
		return LuminousIntensity, nil
	case "rad":
		return Angle, nil
	case "sr":
		return SolidAngle, nil
	case "Hz":
		return Frequency, nil
	case "N":
		return Force, nil
	case "Pa":
		return Pressure, nil
	case "J":
		return Energy, nil
	case "W":
		return Power, nil
	case "C":
		return ElectricCharge, nil
	case "V":
		return Voltage, nil
	case "F":
		return Capacitance, nil
	case "Ω":
		return Impedance, nil
	case "S":
		return ElectricalConductance, nil
	case "Wb":
		return MagneticFlux, nil
	case "T":
		return MagneticFluxDensity, nil
	case "H":
		return Inductance, nil
	case "°C":
		return Temperature, nil
	case "lm":
		return LuminousFlux, nil
	case "lx":
		return Illuminance, nil
	case "Bq":
		return Radioactivity, nil
	case "Gy":
		return AbsorbedDose, nil
	case "Sv":
		return EquivalentDose, nil
	case "kat":
		return CatalyticActivity, nil
	case "":
		return None, fmt.Errorf("An empty string is not a recognized SI symbol", str)
	default:
		return None, fmt.Errorf("%s is not a recognized SI symbol", str)
	}
}

// Dimension return the symbol used in dimensional analysis.
func (measure Measure) Dimension() string {
	switch measure {
	case Length:
		return "L"
	case Mass:
		return "M"
	case Time:
		return "T"
	case ElectricCurrent:
		return "I"
	case ThermodynamicTemperature:
		return "Θ"
	case AmountOfSubstance:
		return "N"
	case LuminousIntensity:
		return "J"
	case Angle:
		return ""
	case SolidAngle:
		return ""
	case Frequency:
		return ""
	case Force:
		return ""
	case Pressure:
		return ""
	case Energy:
		return ""
	case Power:
		return ""
	case ElectricCharge:
		return ""
	case Voltage:
		return ""
	case Capacitance:
		return ""
	case Impedance:
		return ""
	case ElectricalConductance:
		return ""
	case MagneticFlux:
		return ""
	case MagneticFluxDensity:
		return ""
	case Inductance:
		return ""
	case Temperature:
		return ""
	case LuminousFlux:
		return ""
	case Illuminance:
		return ""
	case Radioactivity:
		return ""
	case AbsorbedDose:
		return ""
	case EquivalentDose:
		return ""
	case CatalyticActivity:
		return ""
	case None:
		return ""
	default:
		return ""
	}
}

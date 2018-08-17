package si

type Measure int

const (
	None Measure = iota
	Length
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
func Parse(str string) (measure Measure) {
	switch str {
	case "m":
		return Length
	case "kg":
		return Mass
	case "s":
		return Time
	case "A":
		return ElectricCurrent
	case "K":
		return ThermodynamicTemperature
	case "mol":
		return AmountOfSubstance
	case "cd":
		return LuminousIntensity
	case "rad":
		return Angle
	case "sr":
		return SolidAngle
	case "Hz":
		return Frequency
	case "N":
		return Force
	case "Pa":
		return Pressure
	case "J":
		return Energy
	case "W":
		return Power
	case "C":
		return ElectricCharge
	case "V":
		return Voltage
	case "F":
		return Capacitance
	case "Ω":
		return Impedance
	case "S":
		return ElectricalConductance
	case "Wb":
		return MagneticFlux
	case "T":
		return MagneticFluxDensity
	case "H":
		return Inductance
	case "°C":
		return Temperature
	case "lm":
		return LuminousFlux
	case "lx":
		return Illuminance
	case "Bq":
		return Radioactivity
	case "Gy":
		return AbsorbedDose
	case "Sv":
		return EquivalentDose
	case "kat":
		return CatalyticActivity
	case "":
		return None
	default:
		return None
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

package si

import (
	"errors"
	"math/big"
)

type Prefix int

const (
	// Base is 1 or 1E0
	Base = iota
	// Yotta is 1E24
	Yotta
	// Zetta is 1E21
	Zetta
	// Exa is 1E18
	Exa
	// Peta is 1E15
	Peta
	// Tera is 1E12
	Tera
	// Giga is 1E9
	Giga
	// Mega is 1E6
	Mega
	// Kilo is 1E3
	Kilo
	// Hecto is 1E2
	Hecto
	// Deca is 1E1
	Deca
	// Deci is 1E-1
	Deci
	// Centi is 1E-2
	Centi
	// Milli is 1E-3
	Milli
	// Micro is 1E-6
	Micro
	// Nano is 1E-9
	Nano
	// Pico is 1E-12
	Pico
	// Femto is 1E-15
	Femto
	// Atto is 1E-18
	Atto
	// Zepto is 1E-21
	Zepto
	// Yocto is 1E-24
	Yocto
)

func (prefix Prefix) String() string {
	switch prefix {
	case Yotta:
		return "Y"
	case Zetta:
		return "Z"
	case Exa:
		return "E"
	case Peta:
		return "P"
	case Tera:
		return "T"
	case Giga:
		return "G"
	case Mega:
		return "M"
	case Kilo:
		return "k"
	case Hecto:
		return "h"
	case Deca:
		return "da"
	case Deci:
		return "d"
	case Centi:
		return "c"
	case Milli:
		return "m"
	case Micro:
		return "μ"
	case Nano:
		return "n"
	case Pico:
		return "p"
	case Femto:
		return "f"
	case Atto:
		return "a"
	case Zepto:
		return "z"
	case Yocto:
		return "y"
	default:
		return ""
	}
}

// Factor returns the multiplication factor for a prefix.
// When prefixes are used the SI units are no longer coherent
// which means you need to do convertions.
func (prefix Prefix) Factor() (*big.Float, error) {
	switch prefix {
	case Yotta:
		return big.NewFloat(1E24), nil
	case Zetta:
		return big.NewFloat(1E21), nil
	case Exa:
		return big.NewFloat(1E18), nil
	case Peta:
		return big.NewFloat(1E15), nil
	case Tera:
		return big.NewFloat(1E12), nil
	case Giga:
		return big.NewFloat(1E9), nil
	case Mega:
		return big.NewFloat(1E6), nil
	case Kilo:
		return big.NewFloat(1E3), nil
	case Hecto:
		return big.NewFloat(1E2), nil
	case Deca:
		return big.NewFloat(1E1), nil
	case Deci:
		return big.NewFloat(1E-1), nil
	case Centi:
		return big.NewFloat(1E-2), nil
	case Milli:
		return big.NewFloat(1E-3), nil
	case Micro:
		return big.NewFloat(1E-6), nil
	case Nano:
		return big.NewFloat(1E-9), nil
	case Pico:
		return big.NewFloat(1E-12), nil
	case Femto:
		return big.NewFloat(1E-15), nil
	case Atto:
		return big.NewFloat(1E-18), nil
	case Zepto:
		return big.NewFloat(1E-21), nil
	case Yocto:
		return big.NewFloat(1E-24), nil
	case Base:
		// The base unit is 1
		return big.NewFloat(1E0), nil
	default:
		return nil, errors.New("Unknown prefix")
	}
}

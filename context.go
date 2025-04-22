package si

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// StandardContext implements the Context interface with standard SI units and prefixes
type StandardContext struct {
	baseUnits      map[string]Unit
	derivedUnits   map[string]Unit
	prefixes       map[string]float64
	sortedPrefixes []string
}

// NewStandardContext creates a new context with standard SI units and prefixes
func NewStandardContext() *StandardContext {
	ctx := &StandardContext{
		baseUnits:    make(map[string]Unit),
		derivedUnits: make(map[string]Unit),
		prefixes:     make(map[string]float64),
	}

	// Register SI base units
	ctx.registerBaseUnits()

	// Register SI derived units
	ctx.registerDerivedUnits()

	// Register SI prefixes
	ctx.registerPrefixes()

	// Sort prefixes by length for proper matching
	ctx.sortPrefixes()

	return ctx
}

// registerBaseUnits registers the 7 SI base units
func (ctx *StandardContext) registerBaseUnits() {
	// Length, Mass, Time, Current, Temperature, Substance, Luminosity
	ctx.baseUnits["m"] = Unit{1, Dimension{1, 0, 0, 0, 0, 0, 0}}
	ctx.baseUnits["kg"] = Unit{1, Dimension{0, 1, 0, 0, 0, 0, 0}}
	ctx.baseUnits["s"] = Unit{1, Dimension{0, 0, 1, 0, 0, 0, 0}}
	ctx.baseUnits["A"] = Unit{1, Dimension{0, 0, 0, 1, 0, 0, 0}}
	ctx.baseUnits["K"] = Unit{1, Dimension{0, 0, 0, 0, 1, 0, 0}}
	ctx.baseUnits["mol"] = Unit{1, Dimension{0, 0, 0, 0, 0, 1, 0}}
	ctx.baseUnits["cd"] = Unit{1, Dimension{0, 0, 0, 0, 0, 0, 1}}
}

// registerDerivedUnits registers common SI derived units
func (ctx *StandardContext) registerDerivedUnits() {
	// Newton: kg·m/s²
	newton := ctx.baseUnits["kg"].Mul(ctx.baseUnits["m"]).Div(ctx.baseUnits["s"].Pow(2))
	ctx.derivedUnits["N"] = newton

	// Joule: N·m
	joule := newton.Mul(ctx.baseUnits["m"])
	ctx.derivedUnits["J"] = joule

	// Watt: J/s
	watt := joule.Div(ctx.baseUnits["s"])
	ctx.derivedUnits["W"] = watt

	// Pascal: N/m²
	pascal := newton.Div(ctx.baseUnits["m"].Pow(2))
	ctx.derivedUnits["Pa"] = pascal

	// Hertz: 1/s
	hertz := Unit{1, Dimension{0, 0, -1, 0, 0, 0, 0}}
	ctx.derivedUnits["Hz"] = hertz

	// Coulomb: A·s
	coulomb := ctx.baseUnits["A"].Mul(ctx.baseUnits["s"])
	ctx.derivedUnits["C"] = coulomb

	// Volt: W/A
	volt := watt.Div(ctx.baseUnits["A"])
	ctx.derivedUnits["V"] = volt

	// Other units with conversion factors
	ctx.derivedUnits["h"] = Unit{3600, Dimension{0, 0, 1, 0, 0, 0, 0}}  // hour
	ctx.derivedUnits["min"] = Unit{60, Dimension{0, 0, 1, 0, 0, 0, 0}}  // minute
	ctx.derivedUnits["d"] = Unit{86400, Dimension{0, 0, 1, 0, 0, 0, 0}} // day

	// Information units
	ctx.derivedUnits["B"] = Unit{1, Dimension{0, 0, 0, 0, 0, 0, 0}}  // byte
	ctx.derivedUnits["iB"] = Unit{1, Dimension{0, 0, 0, 0, 0, 0, 0}} // byte for binary prefixes
}

// registerPrefixes registers SI and binary prefixes
func (ctx *StandardContext) registerPrefixes() {
	// SI prefixes
	ctx.prefixes["Y"] = 1e24
	ctx.prefixes["Z"] = 1e21
	ctx.prefixes["E"] = 1e18
	ctx.prefixes["P"] = 1e15
	ctx.prefixes["T"] = 1e12
	ctx.prefixes["G"] = 1e9
	ctx.prefixes["M"] = 1e6
	ctx.prefixes["k"] = 1e3
	ctx.prefixes["h"] = 1e2
	ctx.prefixes["da"] = 1e1
	ctx.prefixes[""] = 1
	ctx.prefixes["d"] = 1e-1
	ctx.prefixes["c"] = 1e-2
	ctx.prefixes["m"] = 1e-3
	ctx.prefixes["u"] = 1e-6
	ctx.prefixes["μ"] = 1e-6
	ctx.prefixes["µ"] = 1e-6
	ctx.prefixes["n"] = 1e-9
	ctx.prefixes["p"] = 1e-12
	ctx.prefixes["f"] = 1e-15
	ctx.prefixes["a"] = 1e-18
	ctx.prefixes["z"] = 1e-21
	ctx.prefixes["y"] = 1e-24

	// Binary prefixes
	ctx.prefixes["Ki"] = math.Pow(2, 10)
	ctx.prefixes["Mi"] = math.Pow(2, 20)
	ctx.prefixes["Gi"] = math.Pow(2, 30)
	ctx.prefixes["Ti"] = math.Pow(2, 40)
	ctx.prefixes["Pi"] = math.Pow(2, 50)
	ctx.prefixes["Ei"] = math.Pow(2, 60)
}

// sortPrefixes sorts prefixes by length for proper matching
func (ctx *StandardContext) sortPrefixes() {
	ctx.sortedPrefixes = make([]string, 0, len(ctx.prefixes))
	for p := range ctx.prefixes {
		ctx.sortedPrefixes = append(ctx.sortedPrefixes, p)
	}
	sort.Slice(ctx.sortedPrefixes, func(i, j int) bool {
		return len(ctx.sortedPrefixes[i]) > len(ctx.sortedPrefixes[j])
	})
}

// Resolve implements the Context interface
func (ctx *StandardContext) Resolve(symbol string) (Unit, error) {
	// Handle special case for dimensionless unit
	if symbol == "1" || symbol == "" {
		return Unit{1, Dimension{}}, nil
	}

	// Handle gram special case
	if symbol == "g" {
		return Unit{0.001, Dimension{0, 1, 0, 0, 0, 0, 0}}, nil
	}

	// Try to match as direct unit
	if unit, ok := ctx.baseUnits[symbol]; ok {
		return unit, nil
	}
	if unit, ok := ctx.derivedUnits[symbol]; ok {
		return unit, nil
	}

	// Handle binary prefixes (KiB, MiB, etc.)
	if strings.HasSuffix(symbol, "iB") {
		for _, prefix := range []string{"K", "M", "G", "T", "P", "E"} {
			if strings.HasPrefix(symbol, prefix) && symbol == prefix+"iB" {
				return Unit{ctx.prefixes[prefix+"i"], Dimension{}}, nil
			}
		}
	}

	// Try with prefixes
	for _, prefix := range ctx.sortedPrefixes {
		if prefix == "" {
			continue
		}

		if !strings.HasPrefix(symbol, prefix) {
			continue
		}

		suffix := symbol[len(prefix):]

		// Try base units with this prefix
		if unit, ok := ctx.baseUnits[suffix]; ok {
			scaledUnit := unit
			scaledUnit.Value *= ctx.prefixes[prefix]
			return scaledUnit, nil
		}

		// Try derived units with this prefix
		if unit, ok := ctx.derivedUnits[suffix]; ok {
			scaledUnit := unit
			scaledUnit.Value *= ctx.prefixes[prefix]
			return scaledUnit, nil
		}
	}

	return Unit{}, fmt.Errorf("unrecognized unit: %s", symbol)
}

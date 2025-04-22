package main

import (
	"fmt"

	"github.com/gurre/si"
)

func main() {
	// Create a pressure reading in PSI
	tirePressure := si.Psi(32.5) // 32.5 psi

	// Parse a pressure value with PSI units from a string
	atmosphericPressure, _ := si.Parse("14.7 psi")

	// Calculate pressure difference
	pressureDiff, _ := tirePressure.Add(atmosphericPressure.Mul(si.Scalar(-1)))

	// Convert to different pressure units
	pressureKPa, _ := si.ToKiloPascals(pressureDiff)

	// Create a circular area (1 inch radius)
	radius := si.Meter.Mul(si.Scalar(0.0254))     // 1 inch in meters
	area := radius.Pow(2).Mul(si.Scalar(3.14159)) // πr²

	// Calculate force from pressure (F = P × A)
	force := pressureDiff.Mul(area)

	// Output results
	fmt.Printf("Tire pressure: %s\n", tirePressure)
	fmt.Printf("Atmospheric pressure: %s\n", atmosphericPressure)
	fmt.Printf("Gauge pressure: %s (%.2f kPa)\n", pressureDiff, pressureKPa)
	fmt.Printf("Force on 1 inch radius circle: %s\n", force)
}

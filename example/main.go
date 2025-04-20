package main

import (
	"fmt"

	"github.com/gurre/si"
)

func main() {
	// Original examples
	mass := si.Kilograms(1500)        // Mass of a car in kg
	velocity, _ := si.Parse("20 m/s") // Velocity in m/s
	kineticEnergy := si.Scalar(0.5).Mul(mass).Mul(velocity.Pow(2))
	fmt.Println("Kinetic energy of a car:", kineticEnergy)

	// Example: Calculate gravitational potential energy (PE = m * g * h)
	height := si.Meters(10)                          // Height in meters
	gravity := si.Meters(9.81).Div(si.Second.Pow(2)) // Gravitational acceleration
	potentialEnergy := mass.Mul(gravity).Mul(height)
	fmt.Println("Potential energy of a car at 10m height:", potentialEnergy)

	// Example: Calculate centripetal force (F = m * v^2 / r)
	radius := si.Meters(50) // Radius of the curve
	centripetalForce := mass.Mul(velocity.Pow(2)).Div(radius)
	fmt.Println("Centripetal force on a car turning:", centripetalForce)
}

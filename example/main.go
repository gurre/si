package main

import (
	"fmt"

	"github.com/gurre/si"
)

func main() {
	// Create units with various representations
	temp := si.Celsius(25.5)                                     // 298.65 K
	distance := si.Kilometers(1.5)                               // 1.5 km
	flow := si.Meter.Pow(3).Mul(si.Scalar(0.002)).Div(si.Second) // 0.002 m^3/s
	fmt.Println(temp, distance, flow)

	// Parse units from strings (e.g., from sensor readings)
	pressure, _ := si.Parse("101.325 kPa") // 101325 Pa
	velocity, _ := si.Parse("55 km/h")     // 15.278 m/s
	fmt.Println(pressure, velocity)

	// Convert between units
	meters, _ := distance.ConvertTo(si.Meter)
	fmt.Println(meters) // 1500 m

	// Temperature conversions
	tempF, _ := si.ToFahrenheit(temp) // 77.9 F - no manual calculation needed
	tempC, _ := si.ToCelsius(temp)    // 25.5 C - easy conversion back
	fmt.Println(tempF, tempC)

	// Perform calculations with units
	power := pressure.Mul(flow)      // 202.65 W
	energy := power.Mul(si.Hours(2)) // 1459080 J
	fmt.Println(power, energy)
}

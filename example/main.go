package main

import (
	"fmt"

	"github.com/gurre/si"
)

func main() {
	// Temperature readings from different eras, all with units attached
	readings := []si.Unit{
		si.Celsius(22.5),    // 2022 sensor
		si.Fahrenheit(72.6), // 2023 sensor
		si.Kelvins(295.7),   // 2024 sensor
	}

	// All temperatures converted to a standard unit for analysis
	for _, temp := range readings {
		// No need to know which era a reading is from
		// No need to check field names or metadata
		c, _ := si.ToCelsius(temp)
		fmt.Printf("Temperature: %s (%.2f C)\n", temp, c)
	}
}

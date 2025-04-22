package main

import (
	"fmt"

	"github.com/gurre/si"
)

func main() {
	// Calculate heat exchange rate
	massFlow := si.Kilograms(2.5).Div(si.Second) // 2.5 kg/s of water
	specificHeat := si.New(4186, "J/(kg·K)")     // Specific heat of water
	tempDiff := si.Kelvin.Mul(si.Scalar(15))     // Temperature difference of 15K
	heatRate := massFlow.Mul(specificHeat).Mul(tempDiff)

	fmt.Println("Heat exchange rate:", heatRate) // In watts

}

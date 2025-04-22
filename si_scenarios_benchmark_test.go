package si_test

import (
	"math"
	"testing"

	"github.com/gurre/si"
)

// BenchmarkThermodynamicCalculation benchmarks a thermodynamic calculation
func BenchmarkThermodynamicCalculation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Specific enthalpy calculation for steam
		temperature := si.Celsius(200) // 200°C steam
		pressure := si.Pascals(1e6)    // 1 MPa

		// Specific heat capacity (simplified)
		cp := si.Joules(2000).Div(si.Kilogram.Mul(si.Kelvin))

		// Reference state
		refTemp := si.Celsius(0)

		// Calculate enthalpy change: h = cp * (T - Tref)
		deltaT := temperature.Value - refTemp.Value
		enthalpyChange := cp.Mul(si.Kelvin.Mul(si.Scalar(deltaT)))

		// Calculate specific volume (simplified model)
		// Using ideal gas law: v = RT/p
		r := si.Joules(462).Div(si.Kilogram.Mul(si.Kelvin)) // Gas constant for steam
		specificVolume := r.Mul(temperature).Div(pressure)

		// Calculate entropy (simplified): s = cp * ln(T2/T1)
		entropyChange := cp.Mul(si.Scalar(math.Log(temperature.Value / refTemp.Value)))

		// Final property calculations
		_ = enthalpyChange
		_ = specificVolume
		_ = entropyChange
	}
}

// BenchmarkHeatExchangerDesign benchmarks calculations for a heat exchanger design
func BenchmarkHeatExchangerDesign(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Heat transfer requirement
		heatLoad := si.Watts(50000) // 50 kW

		// Hot fluid properties
		hotInlet := si.Celsius(90)
		hotOutlet := si.Celsius(70)
		hotCp := si.Joules(4200).Div(si.Kilogram.Mul(si.Kelvin))

		// Cold fluid properties
		coldInlet := si.Celsius(15)
		coldOutlet := si.Celsius(45)
		coldCp := si.Joules(4200).Div(si.Kilogram.Mul(si.Kelvin))

		// Calculate hot fluid mass flow rate: m = Q / (cp * ΔT)
		hotDeltaT := hotInlet.Value - hotOutlet.Value
		hotMassFlow := heatLoad.Div(hotCp.Mul(si.Kelvin.Mul(si.Scalar(hotDeltaT))))

		// Calculate cold fluid mass flow rate
		coldDeltaT := coldOutlet.Value - coldInlet.Value
		coldMassFlow := heatLoad.Div(coldCp.Mul(si.Kelvin.Mul(si.Scalar(coldDeltaT))))

		// Calculate log mean temperature difference (LMTD)
		deltaT1 := hotInlet.Value - coldOutlet.Value
		deltaT2 := hotOutlet.Value - coldInlet.Value
		lmtd := si.Kelvin.Mul(si.Scalar((deltaT1 - deltaT2) / math.Log(deltaT1/deltaT2)))

		// Calculate required heat transfer area
		overallHeatTransferCoef := si.Watts(500).Div(si.Meter.Pow(2).Mul(si.Kelvin))
		requiredArea := heatLoad.Div(overallHeatTransferCoef.Mul(lmtd))

		// Calculate pressure drop (simplified)
		_ = hotMassFlow
		_ = coldMassFlow
		_ = requiredArea
	}
}

// BenchmarkPumpingSystemAnalysis benchmarks a pumping system analysis calculation
func BenchmarkPumpingSystemAnalysis(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// System parameters
		flowRate := si.Meter.Pow(3).Mul(si.Scalar(0.05)).Div(si.Second) // 50 L/s
		pipeDiameter := si.Meters(0.1)                                  // 100 mm
		pipeLength := si.Meters(100)                                    // 100 m
		pipeRoughness := si.Meters(0.0001)                              // 0.1 mm
		elevationChange := si.Meters(15)                                // 15 m lift
		fluidDensity := si.Kilograms(1000).Div(si.Meter.Pow(3))         // Water
		gravity := si.Meters(9.81).Div(si.Second.Pow(2))                // Gravitational acceleration

		// Calculate flow velocity
		pipeArea := si.Scalar(math.Pi * math.Pow(pipeDiameter.Value/2, 2))
		velocity := flowRate.Div(si.Meter.Pow(2).Mul(pipeArea))

		// Calculate Reynolds number
		dynamicViscosity := si.Pascal.Mul(si.Second).Mul(si.Scalar(0.001)) // Water at 20°C
		reynolds := fluidDensity.Mul(velocity).Mul(pipeDiameter).Div(dynamicViscosity)

		// Calculate friction factor (Colebrook equation simplified)
		relativeRoughness := pipeRoughness.Div(pipeDiameter)
		frictionFactor := si.Scalar(0.02) // Simplified for benchmark
		if reynolds.Value > 4000 {
			frictionFactor = si.Scalar(0.25 / math.Pow(math.Log10(relativeRoughness.Value/3.7+5.74/math.Pow(reynolds.Value, 0.9)), 2))
		}

		// Calculate major head loss (Darcy-Weisbach)
		velocitySquared := velocity.Mul(velocity)
		majorHeadLoss := frictionFactor.Mul(pipeLength).Mul(velocitySquared).Div(pipeDiameter.Mul(si.Scalar(2 * gravity.Value)))

		// Calculate minor head losses (simplified)
		kTotal := si.Scalar(10) // Sum of minor loss coefficients
		minorHeadLoss := kTotal.Mul(velocitySquared).Div(si.Scalar(2 * gravity.Value))

		// Calculate static head
		staticHead := elevationChange

		// Calculate total dynamic head (handle errors properly in real code)
		majorPlusMinor, _ := majorHeadLoss.Add(minorHeadLoss)
		totalHead, _ := staticHead.Add(majorPlusMinor)

		// Calculate pump power requirement
		efficiency := si.Scalar(0.7) // 70% pump efficiency
		pumpPower := flowRate.Mul(fluidDensity).Mul(gravity).Mul(totalHead).Div(efficiency)

		_ = pumpPower
	}
}

// BenchmarkElectricalCircuitAnalysis benchmarks electrical circuit analysis
func BenchmarkElectricalCircuitAnalysis(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Circuit parameters
		voltage := si.Volts(230)                                         // 230V supply
		resistance1 := si.Volts(10).Div(si.Amperes(1))                   // 10 ohms
		resistance2 := si.Volts(15).Div(si.Amperes(1))                   // 15 ohms
		resistance3 := si.Volts(20).Div(si.Amperes(1))                   // 20 ohms
		inductance := si.Volts(5).Mul(si.Second).Div(si.Amperes(1))      // 5H
		capacitance := si.Amperes(1).Mul(si.Second).Div(si.Volts(0.001)) // 0.001F
		frequency := si.Hertzs(50)                                       // 50Hz

		// Calculate impedances
		omega := si.Scalar(2 * math.Pi).Mul(frequency)
		inductiveReactance := omega.Mul(inductance)
		capacitiveReactance := si.Scalar(1).Div(omega.Mul(capacitance))

		// Series combination (handle errors properly in real code)
		r1PlusR2, _ := resistance1.Add(resistance2)
		totalResistance, _ := r1PlusR2.Add(resistance3)

		// Current calculation (simple circuit)
		current := voltage.Div(totalResistance)

		// Power calculations
		powerResistor1 := current.Mul(current).Mul(resistance1)
		powerResistor2 := current.Mul(current).Mul(resistance2)
		powerResistor3 := current.Mul(current).Mul(resistance3)

		// Total power (handle errors properly in real code)
		p1PlusP2, _ := powerResistor1.Add(powerResistor2)
		totalPower, _ := p1PlusP2.Add(powerResistor3)

		// RLC circuit (simplified)
		impedanceValue := math.Sqrt(math.Pow(totalResistance.Value, 2) +
			math.Pow(inductiveReactance.Value-capacitiveReactance.Value, 2))
		impedanceRL := si.Scalar(impedanceValue)

		// Current in RLC circuit
		currentRLC := voltage.Div(impedanceRL)

		// Phase angle
		phaseAngle := si.Scalar(math.Atan2(inductiveReactance.Value-capacitiveReactance.Value, totalResistance.Value))

		_ = currentRLC
		_ = phaseAngle
		_ = totalPower
	}
}

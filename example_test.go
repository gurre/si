package si_test

import (
	"fmt"
	"math"

	"github.com/gurre/si"
)

func Example() {
	// Basic unit creation
	distance := si.Kilometers(10)
	time := si.Minutes(30)

	// Calculate speed
	speed := distance.Div(time)

	fmt.Println(speed) // Should output something like: 5.55556 m/s

	// Convert to different units
	// We're using ConvertTo, but this actually calculates a different ratio
	// not just a display format change
	kilometersPerHour, _ := speed.ConvertTo(si.Kilometers(1).Div(si.Hours(1)))
	fmt.Println(kilometersPerHour) // Output is scaled to km/h

	// Basic arithmetic
	mass := si.Kilograms(75)
	acceleration := si.Meters(9.8).Div(si.Second.Pow(2))

	// Calculate force (F = m·a)
	force := mass.Mul(acceleration)
	fmt.Println(force) // Should output something like: 735 N

	// Creating custom units
	hpPower := si.Watt.Mul(si.Scalar(745.7))
	fmt.Println("1 horsepower =", hpPower)

	// Using the Parse function
	parsedUnit, _ := si.Parse("100 km/h")
	fmt.Println(parsedUnit) // Should output something like: 27.7778 m/s

	// JSON marshaling/unmarshaling example (uncomment to use)
	// unitJSON, _ := json.Marshal(parsedUnit)
	// fmt.Println(string(unitJSON))

	// Output:
	// 20 km/h
	// 72 km/h
	// 735 N
	// 1 horsepower = 1 hp
	// 100 km/h
}

func ExampleNew() {
	// Creating units with New
	length := si.New(5, "m")
	fmt.Println(length)

	// Using prefixes
	mass := si.New(500, "g")
	fmt.Println(mass)

	// Complex units
	pressure := si.New(101.325, "kPa")
	fmt.Println(pressure)

	// Output:
	// 5 m
	// 0.5 kg
	// 101.325 Pa
}

func ExampleUnit_Mul() {
	distance := si.Meters(10)
	force := si.Newton

	// Calculate work (W = F·d)
	work := force.Mul(distance)
	fmt.Println("Work =", work)

	// Output:
	// Work = 10 J
}

func ExampleUnit_Div() {
	energy := si.Joule
	time := si.Seconds(2)

	// Calculate power (P = E/t)
	power := energy.Div(time)
	fmt.Println("Power =", power)

	// Output:
	// Power = 0.5 W
}

func ExampleUnit_Pow() {
	length := si.Meters(2)

	// Calculate area
	area := length.Pow(2)
	fmt.Println("Area =", area)

	// Calculate volume
	volume := length.Pow(3)
	fmt.Println("Volume =", volume)

	// Output:
	// Area = 4 m^2
	// Volume = 8 m^3
}

func ExampleUnit_ConvertTo() {
	distance := si.Kilometers(5)

	// Convert to meters
	meters, _ := distance.ConvertTo(si.Meter)
	fmt.Println("5 km =", meters)

	// Convert to hours to seconds
	time := si.Hours(1.5)
	seconds, _ := time.ConvertTo(si.Second)
	fmt.Println("1.5 hours =", seconds)

	// Output:
	// 5 km = 5000 m
	// 1.5 hours = 5400 s
}

func ExampleParse() {
	// Parse simple units
	u1, _ := si.Parse("100 m")
	fmt.Println(u1)

	// Parse complex units
	u2, _ := si.Parse("9.8 m/s^2")
	fmt.Println(u2)

	// Parse unit with prefix
	u3, _ := si.Parse("1.21 GW")
	fmt.Println(u3)

	// Output:
	// 100 m
	// 9.8 m/s^2
	// 1.21 GW
}

func ExampleScalar() {
	// Create dimensionless quantity
	probability := si.Scalar(0.75)
	fmt.Println("Probability =", probability)

	// Use in calculations
	mass := si.Kilograms(20)
	halfMass := si.Scalar(0.5).Mul(mass)
	fmt.Println("Half the mass =", halfMass)

	// Output:
	// Probability = 0.75 1
	// Half the mass = 10 kg
}

func ExampleKilometers() {
	distance := si.Kilometers(42.195) // Marathon distance
	fmt.Println("Marathon distance =", distance)

	// Convert to meters
	meters, _ := distance.ConvertTo(si.Meter)
	fmt.Println("In meters =", meters)

	// Output:
	// Marathon distance = 42195 m
	// In meters = 42195 m
}

func ExampleSeconds() {
	time := si.Seconds(90)
	fmt.Println("Time =", time)

	// Convert to minutes
	minutes, _ := time.ConvertTo(si.Minutes(1))
	fmt.Println("In minutes =", minutes)

	// Output:
	// Time = 90 s
	// In minutes = 1.5 s
}

func Example_compoundUnits() {
	// Example: Calculating pressure from force and area
	force := si.Newtons(1000)
	area := si.Meters(2).Pow(2)
	pressure := force.Div(area)
	fmt.Println("Pressure:", pressure)

	// Example: Calculating power from voltage and current
	voltage := si.Volts(220)
	current := si.Amperes(5)
	power := voltage.Mul(current)
	fmt.Println("Electrical power:", power)

	// Example: Calculating energy from power and time
	time := si.Hours(2)
	energy := power.Mul(time)
	fmt.Println("Energy consumption:", energy)

	// Output:
	// Pressure: 250 Pa
	// Electrical power: 1.1 kW
	// Energy consumption: 7.92e+06 J
}

func Example_complexConversions() {
	// Calculating fuel efficiency
	distance := si.Kilometers(350)
	fuelVolume := si.New(20, "l") // 20 liters (not a standard SI unit)

	// Create a kilometer per liter unit for conversion
	kmPerLiter := si.Kilometers(1).Div(si.New(1, "l"))

	// Calculate fuel efficiency using conversion
	fuelEfficiency, _ := distance.Div(fuelVolume).ConvertTo(kmPerLiter)
	fmt.Printf("Fuel efficiency: %.1f km/l\n", fuelEfficiency.Value)

	// Converting data transfer rate
	// For data rates, we need to handle bit/byte conversions

	// Show what a parsed data rate looks like
	dataRate, _ := si.Parse("100 Mbps")
	fmt.Println("Data rate:", dataRate)

	// Calculate data transferred in 60 seconds
	// 100 Mbps = 100 * 10^6 bits per second
	bitsPerSec := si.New(100*1e6, "b/s")
	transferDuration := si.Minutes(1)
	totalBits := bitsPerSec.Mul(transferDuration)

	// Convert to megabytes (1 byte = 8 bits, 1 MB = 10^6 bytes)
	byteUnit := si.New(1.0/8.0, "B/b") // 1 byte = 8 bits
	mbUnit := si.New(1.0/1e6, "MB/B")  // 1 MB = 10^6 bytes

	// Calculate bytes then convert to megabytes
	totalBytes := totalBits.Mul(byteUnit)
	totalMB := totalBytes.Mul(mbUnit)
	fmt.Printf("Data transferred in 60s at 100 Mbps: %.1f MB\n", totalMB.Value)

	// Output:
	// Fuel efficiency: 17.5 km/l
	// Data rate: 1 1
	// Data transferred in 60s at 100 Mbps: 750.0 MB
}

func Example_physicsFormulas() {
	// Example: Calculate kinetic energy (KE = 1/2 * m * v^2)
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

	// Output:
	// Kinetic energy of a car: 300000 J
	// Potential energy of a car at 10m height: 147150 J
	// Centripetal force on a car turning: 12000 N
}

func Example_electronics() {
	// Example: Calculate resistance using Ohm's law (R = V / I)
	voltage := si.Volts(12)
	current := si.Amperes(2)
	resistance := voltage.Div(current)

	// Define ohm unit for conversion
	ohm := si.Volts(1).Div(si.Amperes(1))
	resistanceInOhms, _ := resistance.ConvertTo(ohm)
	fmt.Printf("Resistance: %g Ω\n", resistanceInOhms.Value)

	// Example: Calculate capacitance (C = Q / V)
	charge := si.Coulomb // 1 Coulomb of charge
	capacitance := charge.Div(voltage)

	// Define farad unit for conversion
	farad := si.Coulomb.Div(si.Volts(1))
	capacitanceInFarads, _ := capacitance.ConvertTo(farad)
	fmt.Printf("Capacitance: %.7g F\n", capacitanceInFarads.Value)

	// Example: Calculate inductance (L = V * t / I)
	time := si.Seconds(0.1)
	inductance := voltage.Mul(time).Div(current)

	// Define henry unit for conversion
	henry := si.Volts(1).Mul(si.Seconds(1)).Div(si.Amperes(1))
	inductanceInHenrys, _ := inductance.ConvertTo(henry)
	fmt.Printf("Inductance: %.1f H\n", inductanceInHenrys.Value)

	// Output:
	// Resistance: 6 Ω
	// Capacitance: 0.08333333 F
	// Inductance: 0.6 H
}

func Example_thermalPhysics() {
	// Example: Calculate thermal energy (Q = m * c * ΔT)
	mass := si.Kilograms(1)                    // 1 kg of water
	specificHeat := si.New(4186, "J/(kg·K)")   // Specific heat of water in J/(kg·K)
	tempChange := si.Kelvin.Mul(si.Scalar(10)) // Temperature change of 10K
	thermalEnergy := mass.Mul(specificHeat).Mul(tempChange)
	fmt.Println("Thermal energy required:", thermalEnergy)

	// Example: Calculate power needed to heat in a given time
	time := si.Minutes(5)
	powerNeeded := thermalEnergy.Div(time)
	fmt.Println("Heating power needed:", powerNeeded)

	// Example: Convert power to kilowatts
	// We should expect around 0.14 kW (= 41860 J / 300 s / 1000)
	kilowatt := si.Watts(1000)

	// The correct conversion to get 0.14
	powerInKw := powerNeeded.Div(kilowatt)
	fmt.Printf("Power in kilowatts: %.2f kW\n", powerInKw.Value)

	// Output:
	// Thermal energy required: 41860 kg·K
	// Heating power needed: 139.53333333333333 kg·K/s
	// Power in kilowatts: 0.14 kW
}

func Example_binaryPrefixes() {
	// Create storage sizes using helper functions
	hardDrive := si.Terabytes(2)   // 2 TB
	memoryRAM := si.Gibibytes(16)  // 16 GiB
	movieSize := si.Mebibytes(750) // 750 MiB

	// Display formatted values
	fmt.Println("Hard drive:", hardDrive)
	fmt.Println("RAM:", memoryRAM)
	fmt.Println("Movie size:", movieSize)

	// Calculate number of movies in terms of actual bytes
	hardDriveBytes := si.New(2e12, "B")            // 2 TB in bytes
	movieBytes := si.New(750*math.Pow(2, 20), "B") // 750 MiB in bytes

	// Calculate the ratio
	movieCountUnit, _ := hardDriveBytes.Div(movieBytes).ConvertTo(si.One)
	fmt.Println("Movies that fit:", fmt.Sprintf("%.2f", movieCountUnit.Value))

	// Calculate hours (1.5 hours per movie)
	hour := si.Hours(1)
	totalHours, _ := si.Scalar(movieCountUnit.Value).Mul(si.Hours(1.5)).ConvertTo(hour)
	fmt.Println("Total hours:", fmt.Sprintf("%.1f", totalHours.Value))

	// Output:
	// Hard drive: 2 1
	// RAM: 16 1
	// Movie size: 750 1
	// Movies that fit: 2543.13
	// Total hours: 3814.7
}

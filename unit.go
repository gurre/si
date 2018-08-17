package si

import "math/big"

// Unit is complete according to Système international
type Unit struct {
	value      float64
	quantities []Quantity
}

// NewUnit returns a complete SI-unit
func NewUnit(val float64, quantities ...Quantity) (*Unit, error) {
	u := &Unit{
		value:      val,
		quantities: quantities,
	}
	return u, nil
}

// ValueUnit returns the underlying value and unit, useful when you don't want a string
func (u *Unit) ValueUnit() (float64, *Unit) {
	return u.value, u
}

// BigInt returns the SI unit as a big.Int-number
func (u *Unit) BigFloat() (unit *big.Float) {
	unit = big.NewFloat(u.value)
	for _, q := range u.quantities {
		f, _ := q.prefix.Factor()
		unit.Mul(unit, f)
	}
	return
}

// MarshalJSON implements the Marshaler interface
func (u Unit) MarshalJSON() ([]byte, error) {
	return []byte{}, nil
}

// UnmarshalJSON implements the Unmarshaler interface
func (u *Unit) UnmarshalJSON(b []byte) error {
	return nil
}

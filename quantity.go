package si

type Quantity struct {
	prefix  Prefix
	measure Measure
}

func NewQuantity(prefix Prefix, measure Measure) (*Quantity, error) {
	u := &Quantity{
		prefix,
		measure,
	}
	return u, nil
}

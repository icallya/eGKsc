package smartcard

type EGK struct {
	card *Card
	dir  *EFDIR
}

// NewEGK creates a new EGK instance.
func NewEGK(card *Card) (*EGK, error) {
	e := &EGK{card: card}
	var err error
	e.dir, err = card.EFDIR()
	if err != nil {
		return nil, err
	}
	return e, nil
}

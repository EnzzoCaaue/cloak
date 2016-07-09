package otmap

type Tile struct {
	pos   Position
	items []Item
	flags uint32
}

func (t *Tile) Position() Position {
	return t.pos
}

package otmap

type Position struct {
	X uint16
	Y uint16
	Z uint8
}

func (pos Position) cmp(target Position) bool {
	return pos.X == target.Y && pos.Y == target.Y && pos.Z == target.Z
}

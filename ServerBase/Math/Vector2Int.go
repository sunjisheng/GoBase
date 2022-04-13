package Math

type Vector2Int struct {
	X int32
	Z int32
}

func (this *Vector2Int)Reset() {
	this.X = 0
	this.Z = 0
}

func (this *Vector2Int) Offset(x int32, z int32) {
	this.X += x
	this.Z += z
}

func (this *Vector2Int) Equal(other Vector2Int) bool{
	if this.X != other.X {
		return false
	}
	if this.Z != other.Z {
		return false
	}
	return true
}

func (this *Vector2Int) SimpleDistance(other Vector2Int) uint32{
	dx := this.X - other.X
	dz := this.Z - other.Z
	if dx < 0 {
		dx = -dx
	}
	if dz < 0 {
		dz = -dz
	}
	return uint32(dx + dz)
}

package Math

type Point struct {
	x float32
	y float32
	z float32
}

func (this *Point)GetX() float32 {
	return this.x
}

func (this *Point)SetX(x float32) {
	this.x = x
}

func (this *Point)GetY() float32 {
	return this.y
}

func (this *Point)SetY(y float32) {
	this.y = y
}

func (this *Point)GetZ() float32 {
	return this.z
}

func (this *Point)SetZ(z float32) {
	this.z = z
}
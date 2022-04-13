package Math

import "math"

type Vector2 struct {
	X float32
	Z float32
}

func (this *Vector2)Reset() {
	this.X = 0
	this.Z = 0
}

func (this *Vector2)Length() float32{
	return float32(math.Sqrt(float64(this.X*this.X + this.Z *this.Z)))
}

func (this *Vector2)LengthSq() float32{
	return this.X*this.X + this.Z *this.Z
}

func (this *Vector2)IsZero() bool{
	if this.X  > MIN_FLOAT || this.X< -MIN_FLOAT {
		return false
	}
	if this.Z  > MIN_FLOAT || this.Z< -MIN_FLOAT {
		return false
	}
	return true
}

func (this *Vector2) Offset(x float32, y float32, z float32) {
	this.X += x
	this.Z += z
}

func (this *Vector2) Add(other *Vector2) {
	this.X += other.X
	this.Z += other.Z
}

func (this *Vector2) Direct(other *Vector2) *Vector2{
	return &Vector2{X:other.X - this.X, Z:other.Z - this.Z}
}

func (this *Vector2) Equal(other *Vector2) bool{
	dx := this.X - other.X
	if dx > MIN_FLOAT || dx < -MIN_FLOAT {
		return false
	}
	dz := this.Z - other.Z
	if dz > MIN_FLOAT || dz < -MIN_FLOAT {
		return false
	}
	return true
}

func (this *Vector2) Multiple(f float32) {
	this.X *= f
	this.Z *= f
}

func (this *Vector2) Normalize() {
	f := this.X * this.X + this.Z * this.Z
	if f == 1.0 {
		return;
	}
	f = float32(math.Sqrt(float64(f)))
	if f < MIN_FLOAT {
		return;
	}
	f = 1.0 / f;
	this.X *= f;
	this.Z *= f;
}

func (this *Vector2) NormalizeXZEx(length float32) {
	f := this.X * this.X +this.Z * this.Z
	if f < MIN_FLOAT {
		return;
	}
	f = float32(math.Sqrt(float64(f)))
	f = float32(float64(length / f));
	this.X *= f;
	this.Z *= f;
}

func (this *Vector2) DistanceSq(to *Vector2) float32{
	dx := this.X - to.X
	dz := this.Z - to.Z
	return dx*dx + dz*dz
}
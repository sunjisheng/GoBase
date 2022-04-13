package Math

import "math"

type Vector3 struct {
	X float32
	Y float32
	Z float32
}

func (this *Vector3)Reset() {
	this.X = 0
	this.Y = 0
	this.Z = 0
}

func (this *Vector3)LengthSq() float32{
	return this.X*this.X + this.Y *this.Y + this.Z *this.Z
}

func (this *Vector3)LengthSqXZ() float32{
	return this.X*this.X + this.Z *this.Z
}

func (this *Vector3)IsZero() bool{
	if this.X  > MIN_FLOAT || this.X< -MIN_FLOAT {
		return false
	}
	if this.Y  > MIN_FLOAT || this.Y< -MIN_FLOAT {
		return false
	}
	if this.Z  > MIN_FLOAT || this.Z< -MIN_FLOAT {
		return false
	}
	return true
}

func (this *Vector3) Offset(x float32, y float32, z float32) {
	this.X += x
	this.Y += y
	this.Z += z
}

func (this *Vector3) Add(other *Vector3){
	this.X += other.X
	this.Y += other.Y
	this.Z += other.Z
}

func (this *Vector3) Direct(other *Vector3){
	this.X -= other.X
	this.Y -= other.Y
	this.Z -= other.Z
}

func (this *Vector3) Equal(other *Vector3) bool{
	dx := this.X - other.X
	if dx > MIN_FLOAT || dx < -MIN_FLOAT {
		return false
	}
	dy := this.Y - other.Y
	if dy > MIN_FLOAT || dy < -MIN_FLOAT {
		return false
	}
	dz := this.Z - other.Z
	if dz > MIN_FLOAT || dz < -MIN_FLOAT {
		return false
	}
	return true
}

func (this *Vector3) Multiple(f float32) {
	this.X *= f
	this.Y *= f
	this.Z *= f
}

func (this *Vector3) Normalize() {
	f := this.X * this.X + this.Y * this.Y + this.Z * this.Z
	if f == 1.0 {
	   return;
	}
	f = float32(math.Sqrt(float64(f)))
	if f < MIN_FLOAT {
		return;
	}
	f = 1.0 / f;
	this.X *= f;
	this.Y *= f;
	this.Z *= f;
}
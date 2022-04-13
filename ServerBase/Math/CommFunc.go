package Math

import "math"

const (
	MIN_FLOAT = 0.00001
)

func DistanceSq(startPos *Vector3, endPos *Vector3) float32 {
	dx := endPos.X - startPos.X
	dy := endPos.Y - startPos.Y
	dz := endPos.Z - startPos.Z
	return dx * dx + dy * dy +dz *dz
}

func Equal(a float32, b float32) bool {
	d := a - b
	if d < -MIN_FLOAT || d > MIN_FLOAT {
		return false
	} else {
		return true
	}
}

func IsZero(a float32) bool {
	return (a > -MIN_FLOAT && a < MIN_FLOAT)
}

func Vector3f_Dir(v1 *Vector3, v2 *Vector3) *Vector3 {
	return &Vector3{X:v2.X-v1.X, Y:v2.Y-v1.Y, Z:v2.Z-v1.Z}
}

func Vec3ToAngle(ve *Vector3)  float32 {
	angle := math.Atan2(float64(ve.Z), float64(ve.X))
	return float32(angle)
}

func AngleToVec3(angle float32, distance float32)  *Vector3 {
	ve:= new(Vector3)
	sin,cos := math.Sincos(float64(angle))
	ve.X = float32(cos * float64(distance))
	ve.Y = 0
	ve.Z = float32(sin * float64(distance))
	return ve
}

func AbsFloat32(f1 float32, f2 float32) float32 {
	if f1 > f2 {
		return f1 - f2
	}
	return f2 - f1
}


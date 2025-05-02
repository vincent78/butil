package mathUtil

import (
	"fmt"
	"math"

	"github.com/golang/geo/s1"
)

// Vector3 represents a point in ℝ³.
type Vector3 struct {
	X, Y, Z float64
}

func (v Vector3) toV2() Vector2 {
	return Vector2{
		X: v.X,
		Y: v.Y,
	}
}

// ApproxEqual reports whether v and ov are equal within a small epsilon.
func (v Vector3) ApproxEqual(ov Vector3) bool {
	const epsilon = 1e-16
	return math.Abs(v.X-ov.X) < epsilon && math.Abs(v.Y-ov.Y) < epsilon && math.Abs(v.Z-ov.Z) < epsilon
}

func (v Vector3) String() string { return fmt.Sprintf("(%0.24f, %0.24f, %0.24f)", v.X, v.Y, v.Z) }

// Norm returns the Vector3's norm.
func (v Vector3) Norm() float64 { return math.Sqrt(v.Dot(v)) }

// Norm2 returns the square of the norm.
func (v Vector3) Norm2() float64 { return v.Dot(v) }

// Normalize returns a unit Vector3 in the same direction as v.
func (v Vector3) Normalize() Vector3 {
	n2 := v.Norm2()
	if n2 == 0 {
		return Vector3{0, 0, 0}
	}
	return v.Mul(1 / math.Sqrt(n2))
}

// IsUnit returns whether this Vector3 is of approximately unit length.
func (v Vector3) IsUnit() bool {
	const epsilon = 5e-14
	return math.Abs(v.Norm2()-1) <= epsilon
}

// Abs returns the Vector3 with nonnegative components.
func (v Vector3) Abs() Vector3 { return Vector3{math.Abs(v.X), math.Abs(v.Y), math.Abs(v.Z)} }

// Add returns the standard Vector3 sum of v and ov.
func (v Vector3) Add(ov Vector3) Vector3 { return Vector3{v.X + ov.X, v.Y + ov.Y, v.Z + ov.Z} }

// Sub returns the standard Vector3 difference of v and ov.
func (v Vector3) Sub(ov Vector3) Vector3 { return Vector3{v.X - ov.X, v.Y - ov.Y, v.Z - ov.Z} }

// Mul returns the standard scalar product of v and m.
func (v Vector3) Mul(m float64) Vector3 { return Vector3{m * v.X, m * v.Y, m * v.Z} }

// Dot returns the standard dot product of v and ov.
func (v Vector3) Dot(ov Vector3) float64 { return v.X*ov.X + v.Y*ov.Y + v.Z*ov.Z }

// Cross returns the standard cross product of v and ov.
func (v Vector3) Cross(ov Vector3) Vector3 {
	return Vector3{
		v.Y*ov.Z - v.Z*ov.Y,
		v.Z*ov.X - v.X*ov.Z,
		v.X*ov.Y - v.Y*ov.X,
	}
}

// Distance returns the Euclidean distance between v and ov.
func (v Vector3) Distance(ov Vector3) float64 { return v.Sub(ov).Norm() }

// Angle returns the angle between v and ov.
func (v Vector3) Angle(ov Vector3) s1.Angle {
	return s1.Angle(math.Atan2(v.Cross(ov).Norm(), v.Dot(ov))) * s1.Radian
}

// Axis enumerates the 3 axes of ℝ³.
type Axis int

// The three axes of ℝ³.
const (
	XAxis Axis = iota
	YAxis
	ZAxis
)

// Ortho returns a unit Vector3 that is orthogonal to v.
// Ortho(-v) = -Ortho(v) for all v.
func (v Vector3) Ortho() Vector3 {
	ov := Vector3{0.012, 0.0053, 0.00457}
	switch v.LargestComponent() {
	case XAxis:
		ov.Z = 1
	case YAxis:
		ov.X = 1
	default:
		ov.Y = 1
	}
	return v.Cross(ov).Normalize()
}

// LargestComponent returns the axis that represents the largest component in this Vector3.
func (v Vector3) LargestComponent() Axis {
	t := v.Abs()

	if t.X > t.Y {
		if t.X > t.Z {
			return XAxis
		}
		return ZAxis
	}
	if t.Y > t.Z {
		return YAxis
	}
	return ZAxis
}

// SmallestComponent returns the axis that represents the smallest component in this Vector3.
func (v Vector3) SmallestComponent() Axis {
	t := v.Abs()

	if t.X < t.Y {
		if t.X < t.Z {
			return XAxis
		}
		return ZAxis
	}
	if t.Y < t.Z {
		return YAxis
	}
	return ZAxis
}

// Cmp compares v and ov lexicographically and returns:
//
//   -1 if v <  ov
//    0 if v == ov
//   +1 if v >  ov
//
// This method is based on C++'s std::lexicographical_compare. Two entities
// are compared element by element with the given operator. The first mismatch
// defines which is less (or greater) than the other. If both have equivalent
// values they are lexicographically equal.
func (v Vector3) Cmp(ov Vector3) int {
	if v.X < ov.X {
		return -1
	}
	if v.X > ov.X {
		return 1
	}

	// First elements were the same, try the next.
	if v.Y < ov.Y {
		return -1
	}
	if v.Y > ov.Y {
		return 1
	}

	// Second elements were the same return the final compare.
	if v.Z < ov.Z {
		return -1
	}
	if v.Z > ov.Z {
		return 1
	}

	// Both are equal
	return 0
}

// Clone 三维向量：拷贝
func (v *Vector3) Clone() Vector3 {
	return Vector3{
		X: v.X,
		Y: v.Y,
		Z: v.Z,
	}
}

// Normalize3 返回：单位化向量
func Normalize3(a Vector3) Vector3 {
	b := a.Clone()
	b.Normalize()
	return b
}

// GetDistance 求两点间距离
func GetDistance(a Vector3, b Vector3) float64 {
	return math.Sqrt(math.Pow(a.X-b.X, 2) + math.Pow(a.Y-b.Y, 2) + math.Pow(a.Z-b.Z, 2))
}

package mathUtil

import "math"

type Vector2 struct {
	X float64
	Y float64
}

func (v Vector2) toV3() Vector3 {
	return Vector3{
		X: v.X,
		Y: v.Y,
		Z: 0,
	}
}

// DegAngle 计算角度
func DegAngle(v1 Vector2, v2 Vector2) float64 {
	if v1.Norm() == 0 || v2.Norm() == 0 {
		return 0
	}
	normProduct := v1.Norm() * v2.Norm()
	dot := v1.X*v2.X + v1.Y*v2.Y
	//threadHold:=normProduct*0.9999
	threadHold := normProduct
	if dot >= -threadHold && dot <= threadHold {
		return math.Acos(dot/normProduct) * 180 / math.Pi
	} else {
		if dot < 0 {
			return 180
		} else {
			return 0
		}
	}
	//cosValue := (v1.X*v2.X + v1.Y*v2.Y) / (v1.Norm() * v2.Norm())
	//return math.Acos(cosValue) * 180 / math.Pi

}

// Rad2Vector 弧度转向量
func Rad2Vector(rad float64) Vector2 {
	return Vector2{
		X: math.Cos(rad),
		Y: math.Sin(rad),
	}
}

// VectorDegAngle 计算基于水平线的角度
func (v Vector2) VectorDegAngle() float64 {
	// X轴的水平线向量
	v2 := Vector2{1, 0}
	degAngle := DegAngle(v, v2)
	if AntiClockwise(v2, v) {
		return degAngle
	} else {
		return -degAngle
	}
	//return math.Atan(v.Y/v.X) * 180 / math.Pi
}
func (v Vector2) Norm() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Deg2Vector 角度转为向量
func Deg2Vector(deg float64) Vector2 {
	return Rad2Vector(deg * math.Pi / 180)
}

// AntiClockwise 是否为逆时针方向
func AntiClockwise(v1 Vector2, v2 Vector2) bool {
	return v1.X*v2.Y-v2.X*v1.Y > 0
}

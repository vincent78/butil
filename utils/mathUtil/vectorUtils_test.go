package mathUtil

import (
	"testing"
)

func TestAngle(t *testing.T) {

	// x轴的水平线
	v1 := Vector3{
		X: 1,
		Y: 0,
		Z: 0,
	}

	v2 := Vector3{
		X: 0,
		Y: 1,
		Z: 0,
	}
	t.Logf("angle:%v\n", v1.Angle(v2).Degrees())
	t.Logf("angle:%v\n", DegAngle(v1.toV2(), v2.toV2()))
}

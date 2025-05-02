package angleUtil

import (
	"math"
	"testing"
)

func TestAngleRound1(t *testing.T) {

	t.Log(-35900, RoundAngle(-35900))
	t.Log(-36900, RoundAngle(-36900))
	t.Log(-20000, RoundAngle(-20000))
	t.Log(-17900, RoundAngle(-17900))
	t.Log(-7900, RoundAngle(-7900))
	t.Log(35900, RoundAngle(35900))
	t.Log(36900, RoundAngle(36900))
	t.Log(20000, RoundAngle(20000))
	t.Log(17900, RoundAngle(17900))
	t.Log(7900, RoundAngle(7900))

}

func TestCalcTargetHeading(t *testing.T) {
	var robHeading, palletAngle, rackHeading, targetRackHeading, offsetAngle int16 = 12000, 9300, 9000, -6000, 0
	heading, _ := CalcTargetHeading(robHeading, palletAngle, rackHeading, targetRackHeading, offsetAngle)
	t.Log("cur info: [", robHeading, palletAngle, rackHeading, targetRackHeading, offsetAngle, " ] ,result: [", heading, "]")

}

func TestCalcMapTwoPointsHeading(t *testing.T) {

	//t.Log(CalcMapTwoPointsHeading(10, 10, 10, 10))   // 0
	//t.Log(CalcMapTwoPointsHeading(10, 10, 10, 20))   // 90
	//t.Log(CalcMapTwoPointsHeading(10, 10, 20, 10))   // 0
	//t.Log(CalcMapTwoPointsHeading(10, 10, 20, 20))   // 45
	//t.Log(CalcMapTwoPointsHeading(10, 10, -20, 20))  // 0
	//t.Log(CalcMapTwoPointsHeading(10, 10, -20, -20)) //
	//t.Log(CalcMapTwoPointsHeading(10, 10, 20, -20))  //
	//t.Log(CalcMapTwoPointsHeading(10, 20, 30, 30))   //

	//t.Log(math.Atan(-1/3.0) / (math.Pi / 180))

	//t.Log(CalcMapTwoPointsHeading(20, 31, 10, 10)) // 0
	//t.Log(CalcMapTwoPointsHeading(10, 10, 20, 31)) // 0
	t.Log(CalcMapTwoPointsHeading(3500, 4000, 4000, 5000)) //正角
	//t.Log(CalcMapTwoPointsHeading(3500, 4000, 2000, 5000)) //正角（负角转换成正角）
	//t.Log(CalcMapTwoPointsHeading(3500, 4000, 2000, 2000)) //负角（正角转换成负角）
	//t.Log(CalcMapTwoPointsHeading(3500, 4000, 4000, 2000)) //负角

}

func TestReversalHeading(t *testing.T) {
	t.Log("1000 -> ", ReversalHeading(1000))
	t.Log("-1000 -> ", ReversalHeading(-1000))
	t.Log("9100 -> ", ReversalHeading(9100))
	t.Log("0 -> ", ReversalHeading(0))
	t.Log("18000 -> ", ReversalHeading(18000))
	t.Log("-18000 -> ", ReversalHeading(-18000))
}

func TestTemp(t *testing.T) {
	t.Log(math.MaxInt8)
	t.Log(math.MaxInt16)
	t.Log(math.MaxInt32)
	t.Log(math.MaxInt64)
}

package angleUtil

import (
	"github.com/vincent78/butil/logger/logger1"
	"github.com/vincent78/butil/utils/mathUtil"
	"math"
)

const (
	AngleStart         int16 = -18000
	AngleEnd           int16 = 18000
	AngleCircle        int32 = 36000
	AngleCircleHalf    int16 = 18000
	AngleCircleQuarter int16 = 9000
	AngleZero          int16 = 0
	RobotNoRotateAngle int16 = 20000
	RobotRotateAngle   int16 = -20000
	RobotLeftMove      int16 = 30000
	RobotRightMove     int16 = -30000
)

// 计算货架相对于地图角度(单位是0.01度)
// param: robotHeading 小车相对于地图角度
// param: rackAngle 货架相对于小车的角度（俯视图来观察）(范围在-18000~18000)
// return: rackHeading 货架相对于地图角度
func CalcCurRackHeading(robotHeading int16, rackAngle int16) (rackHeading int16) {
	return RoundAngle(AddInt16(robotHeading, rackAngle, -AngleCircleQuarter))
}

// 计算小车身上托盘相对于地图角度(单位是0.01度)
// param: robotHeading 小车相对于地图角度
// param: palletAngle 托盘相对于小车的角度（俯视图来观察）(范围在-18000~18000)
// return: palletHeading 托盘相对于地图角度
func CalcCurPalletHeading(robotHeading int16, palletAngle int16) (palletHeading int16) {
	return RoundAngle(AddInt16(robotHeading, palletAngle))
}

// 将角度转换到[-18000,18000] 之间
func RoundAngle(angle int32) int16 {
	angle = angle % AngleCircle //取余
	switch {
	case angle < int32(AngleStart):
		angle = angle + AngleCircle
	case angle > int32(AngleEnd):
		angle = angle - AngleCircle
	}
	return int16(angle)
}

// 计算最终下位机需要的换向指令中的targetHeading (最终托盘与小车一起旋转成到角度(单位是0.01度)，小车车身与托盘随动)
// 计算原理：
//
//	目标货架角度 - 当前货架角度 = 目标托盘角度（targetHeading）- 当前托盘角度(都相对于地图来说的)
//
// param: robotHeading 小车相对于地图角度
// param: palletAngle 托盘相对于小车角度（俯视图来观察）(范围在-180~180)
// param: rackHeading 货架相对于地图的角度
// param: targetRackHeading 最终托盘相对于地图角度
// param: offsetAngle 偏差范围 （+-）
// return: targetHeading,ok 托盘相对于地图的角度（小车在转动时会随着托盘一起转） , 是否需要小车进行换向（false : 不需要换向， true ： 需要换向）
func CalcTargetHeading(robotHeading int16, palletAngle int16, rackHeading int16, targetRackHeading int16, offsetAngle int16) (targetHeading int16, isNeedChangePod bool) {
	logger1.Info("CalcTargetHeading [robotHeading : %d,palletAngle : %d,rackHeading : %d,targetRackHeading : %d,offsetAngle : %d,]", robotHeading, palletAngle, robotHeading, targetHeading, offsetAngle)
	// 计算当前货架角度差值
	var differRackHeading int16 = RoundAngle(AddInt16(targetRackHeading, -rackHeading))
	logger1.Info("rack [differRackHeading : %d ]", differRackHeading)
	offsetAngle = mathUtil.AbsInt16(offsetAngle)
	// 说明当前货架方向不需要发生变化
	var palletHeading int16 = CalcCurPalletHeading(robotHeading, palletAngle)
	logger1.Info("pallet [palletHeading : %d ]", palletHeading)
	if differRackHeading == AngleZero || mathUtil.BetweenInt16(differRackHeading, -offsetAngle, offsetAngle) {
		return palletHeading, false
	}
	var targetHeadingTemp int32 = AddInt16(differRackHeading, palletHeading)
	logger1.Info("pallet old [targetHeading : %d ]", targetHeading)
	targetHeading = RoundAngle(targetHeadingTemp)
	logger1.Info("pallet new [targetHeading : %d ]", targetHeading)
	if targetHeading == AngleCircleHalf && robotHeading > 0 {
		return AngleCircleHalf, true
	}
	if targetHeading == -AngleCircleHalf && robotHeading < 0 {
		return -AngleCircleHalf, true
	}
	return targetHeading, true
}

// 计算托盘相对于小车的角度
// 公式 ： 托盘相对于小车角度 = 圆整（托盘相对于地图角度 - 小车相对于地图角度 + 9000）
// param: robotHeading 小车相对于地图角度
// param: palletHeading 托盘相对于地图角度（俯视图来观察）(范围在-180~180)
// return : palletAngle 托盘相对于小车角度（俯视图来观察）(范围在-180~180)
func CalcPalletAngle(robotHeading int16, palletHeading int16) (palletAngle int16) {
	return RoundAngle(AddInt16(palletHeading, -robotHeading))
}

// 计算货架相对于地图的角度
// 公式： 货架运动时相对于地图角度 = 圆整（路相对于地图角度 + 货架相对于路的角度 - 9000）
// param: pathHeading 小车相对于地图角度
// param: rackPathAngle 货架相对于路的角度（俯视图来观察）(范围在-180~180)
// return : palletAngle 货架相对于地图的角度（俯视图来观察）(范围在-180~180)
func CalcRackHeading(pathHeading int16, rackPathAngle int16) (rackHeading int16) {
	return RoundAngle(AddInt16(pathHeading, rackPathAngle))
}

// 计算两点的向量相对于地图的角度（适用于计算路相对于地图角度，车头相对于地图角度）
// param: startX 起点x坐标
// param: startY 起点y坐标
// param: endX 终点x坐标
// param: endY 终点y坐标
// return:  pathHeading 路相对于地图的角度
func CalcMapTwoPointsHeading(startX int, startY int, endX int, endY int) (heading int16) {
	if (endX-startX) == 0 && (endY-startY) == 0 {
		return AngleZero
	}
	if (endX - startX) == 0 {
		if endY > startY {
			return AngleCircleQuarter
		} else {
			return -AngleCircleQuarter
		}
	}
	if (endY - startY) == 0 {
		if endX > startX {
			return AngleZero
		} else {
			return AngleEnd
		}
	}
	tanX := float64(endY-startY) / float64(endX-startX)
	logger1.Info("tanX : %v", tanX)
	atan := math.Atan(tanX) / (math.Pi / 180)
	logger1.Info("path old [heading : %v ]", atan)
	heading = RoundAngle(int32(atan * 100))
	if startX > endX {
		if startY > endY {
			heading = heading - AngleCircleHalf
		} else {
			heading = heading + AngleCircleHalf
		}
	}
	logger1.Info("path new [heading : %v ]", heading)
	return heading
}

// 计算最小旋转角度
// param: targetHeading 目标角度
// param: curHeading 当前角度
// return:  minRotateAngle 最小旋转的角度，范围在[-180~180](为正时，小车逆时针转，为负时，小车顺时针转)
func CalcMinRotateAngle(targetHeading int16, curHeading int16) (minRotateAngle int16) {
	//提高效率，不需要进行圆整和数据类型转换
	if (targetHeading > 0 && curHeading > 0) || targetHeading < 0 && curHeading < 0 {
		return targetHeading - curHeading
	}
	return RoundAngle(AddInt16(targetHeading - curHeading))
}

// 计算小车目标角度
// param: targetRobotHeading 目标小车角度
// param: curRobotHeading 当前小车角度
// return: targetRobotHeadingR 目标小车角度
func CalcRobotTargetHeading(targetRobotHeading int16, curRobotHeading int16) (targetRobotHeadingR int16) {
	//小车不旋转
	if targetRobotHeading == RobotNoRotateAngle || targetRobotHeading == RobotLeftMove || targetRobotHeading == RobotRightMove {
		return curRobotHeading
	}
	// 小车旋转180度
	if targetRobotHeading == RobotRotateAngle {
		return RoundAngle(AddInt16(curRobotHeading, AngleCircleHalf))
	}
	return RoundAngle(AddInt16(targetRobotHeading))
}

// 将给定角度调转180
func ReversalHeading(heading int16) int16 {
	return RoundAngle(AddInt16(heading, AngleCircleHalf))
}

// int16 数字的相加
func AddInt16(x ...int16) (r int32) {
	for _, n := range x {
		r = r + int32(n)
	}
	return r
}

package num_switch

// 打开开关
func TurnOn(base, index int64) int64 {
	if index >= 64 {
		return -1
	}
	return 1<<index | base
}

// 关闭开关
func TurnOff(base, index int64) int64 {
	if index >= 64 {
		return -1
	}
	return (^(1 << index)) & base
}

// 目标开关是否开启
func CheckTurnOn(base, index int64) bool {
	return base == 1<<index|base
}

// 关闭开关取反
func Reverse(base, index int64) int64 {
	if index >= 64 {
		return -1
	}
	if CheckTurnOn(base, index) {
		return TurnOff(base, index)
	} else {
		return TurnOn(base, index)
	}
}

// 比较两个开关
func CompareSwitch(base1, base2 int64) []int64 {
	var result []int64
	for i := int64(0); i < 64; i++ {
		if CheckTurnOn(base1, i) != CheckTurnOn(base2, i) {
			result = append(result, i)
		}
	}
	return result
}

package num_switch

import (
	"testing"
)

var s = int64(0)

func TestSwitch(t *testing.T) {
	s = Reverse(s, 0)
	t.Log(s)
	t.Log(CompareSwitch(0, s))
}

func TestTurnOn(t *testing.T) {
	t.Log(TurnOn(TurnOn(-3,12),12))
}
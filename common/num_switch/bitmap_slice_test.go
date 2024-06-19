package num_switch

import (
	"testing"
)

func TestNewBitMapSwitch(t *testing.T) {
	bs := NewBitMapSwitch(100)
	bs.TurnOn(0)
	t.Log(bs.IsOn(0))
	t.Log(bs.IsOn(1))
}

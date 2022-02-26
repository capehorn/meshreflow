package meshreflow

import (
	"testing"
)

func TestPushCmd(t *testing.T) {
	ctx := NewContext()
	err := ctx.AddCmdPattern("rect {length:num}")
	if err != nil {
		t.Failed()
	}
	err = ctx.AddCmdPattern("extrude {length:num}")
	if err != nil {
		t.Failed()
	}
	err = ctx.AddCmdPattern("outset {length:num}")
	if err != nil {
		t.Failed()
	}
	err = ctx.AddCmdPattern("inset {length:num}")
	if err != nil {
		t.Failed()
	}

	parsedCmd, err := ctx.PushCmd("rect 10")
	if err != nil {
		t.Failed()
	}
	ctx.PerformCmd(parsedCmd)

	parsedCmd, err = ctx.PushCmd("extrude 20")
	if err != nil {
		t.Failed()
	}
	ctx.PerformCmd(parsedCmd)
}

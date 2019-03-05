package constants

import (
	"jvmgo/instructions/base"
	"jvmgo/rtda"
)

// 空操作

type NOP struct {
	base.NoOperandsInstruction
}

func (self *NOP) Execute(frame *rtda.Frame) {
	// nothing to do
}

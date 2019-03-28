package base

import (
	"jvmgo/rtda"
	"jvmgo/rtda/heap"
)

// 方法调用
func InvokeMethod(invokerFrame *rtda.Frame, method *heap.Method) {
	// 创建帧推入栈
	thread := invokerFrame.Thread()
	newFrame := thread.NewFrame(method)
	thread.PushFrame(newFrame)
	// 传递参数
	argSlotSlot := int(method.ArgSlotCount())
	if argSlotSlot > 0 {
		for i := argSlotSlot - 1; i >= 0; i-- {
			// 从调用者的操作数栈中弹出
			slot := invokerFrame.OperandStack().PopSlot()
			// 放进被调用方法的局部变量表中
			newFrame.LocalVars().SetSlot(uint(i), slot)
		}
	}
}

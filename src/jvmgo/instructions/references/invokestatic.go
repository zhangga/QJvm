package references

import (
	"jvmgo/instructions/base"
	"jvmgo/rtda"
	"jvmgo/rtda/heap"
)

// 调用静态方法
type INVOKE_STATIC struct {
	base.Index16Instruction
}

func (self *INVOKE_STATIC) Execute(frame *rtda.Frame) {
	// 当前类的常量池
	cp := frame.Method().Class().ConstantPool()
	// 常量池中要调用的方法引用
	methodRef := cp.GetConstant(self.Index).(*heap.MethodRef)
	// 解析调用的方法
	resolvedMethod := methodRef.ResolvedMethod()
	class := resolvedMethod.Class()

	if !class.InitStarted() {
		frame.RevertNextPC()
		base.InitClass(frame.Thread(), class)
		return
	}

	if !resolvedMethod.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}

	base.InvokeMethod(frame, resolvedMethod)
}

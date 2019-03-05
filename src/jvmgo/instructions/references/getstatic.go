package references

import (
	"jvmgo/instructions/base"
	"jvmgo/rtda"
	"jvmgo/rtda/heap"
)

// Get static field from class
type GET_STATIC struct {
	base.Index16Instruction
}

func (self *GET_STATIC) Execute(frame *rtda.Frame) {
	cp := frame.Method().Class().ConstantPool()
	fieldRef := cp.GetConstant(self.Index).(*heap.FieldRef)
	field := fieldRef.ResolvedField()
	class := field.Class()
	if !field.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}
	descriptor := field.Descriptor()
	slotID := field.SlotID()
	slots := class.StaticVars()
	stack := frame.OperandStack()
	switch descriptor[0] {
	case 'Z', 'B', 'C', 'S', 'I':
		stack.PushInt(slots.GetInt(slotID))
	case 'F':
		stack.PushFloat(slots.GetFloat(slotID))
	case 'J':
		stack.PushLong(slots.GetLong(slotID))
	case 'D':
		stack.PushDouble(slots.GetDouble(slotID))
	case 'L', '[':
		stack.PushRef(slots.GetRef(slotID))
	}
}

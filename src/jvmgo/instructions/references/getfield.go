package references

import (
	"jvmgo/instructions/base"
	"jvmgo/rtda"
	"jvmgo/rtda/heap"
)

// Fetch field from object
type GET_FIELD struct {
	base.Index16Instruction
}

func (self *GET_FIELD) Execute(frame *rtda.Frame) {
	cp := frame.Method().Class().ConstantPool()
	fieldRef := cp.GetConstant(self.Index).(*heap.FieldRef)
	field := fieldRef.ResolvedField()
	if field.IsStatic() {
		panic("java.lang.IncompatibleClassChangeError")
	}
	stack := frame.OperandStack()
	ref := stack.PopRef()
	if ref == nil {
		panic("java.lang.NullPointerException")
	}

	descriptor := field.Descriptor()
	slotID := field.SlotID()
	slots := ref.Fields()

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
	default:
		// todo
	}
}

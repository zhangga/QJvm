package heap

import "jvmgo/classfile"

// 变量
type Field struct {
	ClassMember
	// 在对象Object中的Slots编号
	slotID uint
	// 指向常量池的索引
	constValueIndex uint
}

func newFields(class *Class, cfFields []*classfile.MemberInfo) []*Field {
	fields := make([]*Field, len(cfFields))
	for i, cfField := range cfFields {
		fields[i] = &Field{}
		fields[i].class = class
		fields[i].copyMemberInfo(cfField)
		fields[i].copyAttributes(cfField)
	}
	return fields
}

func (self *Field) copyAttributes(cfField *classfile.MemberInfo) {
	if valAttr := cfField.ConstantValueAttribute(); valAttr != nil {
		self.constValueIndex = uint(valAttr.ConstantValueIndex())
	}
}

func (self *Field) IsVolatile() bool {
	return 0 != self.accessFlags&ACC_VOLATILE
}

func (self *Field) IsTransient() bool {
	return 0 != self.accessFlags&ACC_TRANSIENT
}

func (self *Field) IsEnum() bool {
	return 0 != self.accessFlags&ACC_ENUM
}

func (self *Field) isLongOrDouble() bool {
	return self.descriptor == "J" || self.descriptor == "D"
}

func (self *Field) SlotID() uint {
	return self.slotID
}

func (self *Field) ConstValueIndex() uint {
	return self.constValueIndex
}

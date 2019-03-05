package heap

import (
	"fmt"
	"jvmgo/classfile"
	"jvmgo/classpath"
)

// 类加载器
type ClassLoader struct {
	cp       *classpath.Classpath
	classMap map[string]*Class //loaded classes
}

func NewClassLoader(cp *classpath.Classpath) *ClassLoader {
	return &ClassLoader{
		cp:       cp,
		classMap: make(map[string]*Class),
	}
}

func (self *ClassLoader) LoadClass(name string) *Class {
	if class, ok := self.classMap[name]; ok {
		// 类已加载
		return class
	}
	return self.loadNonArrayClass(name)
}

// 加载非数组类
func (self *ClassLoader) loadNonArrayClass(name string) *Class {
	// 读取class文件
	data, entry := self.readClass(name)
	// 解析class文件，生成虚拟机可以使用的类数据，并放入方法区
	class := self.defineClass(data)
	// 链接
	link(class)
	fmt.Printf("[Loaded %s from %s]\n", name, entry)
	return class
}

// 读取class文件
func (self *ClassLoader) readClass(name string) ([]byte, classpath.Entry) {
	data, entry, err := self.cp.ReadClass(name)
	if err != nil {
		panic("java.lang.ClassNotFoundException: " + name)
	}
	return data, entry
}

// 解析class文件
func (self *ClassLoader) defineClass(data []byte) *Class {
	class := parseClass(data)
	class.loader = self
	resolveSuperClass(class)
	resolveInterfaces(class)
	self.classMap[class.name] = class
	return class
}

func parseClass(data []byte) *Class {
	// 初步解析class文件
	cf, err := classfile.Parse(data)
	if err != nil {
		panic("java.lang.ClassFormatError")
	}
	// 将classfile解析成运行时Class数据信息
	return newClass(cf)
}

// 加载超类
func resolveSuperClass(class *Class) {
	if class.name != "java/lang/Object" {
		class.superClass = class.loader.LoadClass(class.superClassName)
	}
}

// 加载接口
func resolveInterfaces(class *Class) {
	interfaceCount := len(class.interfaceNames)
	if interfaceCount > 0 {
		class.interfaces = make([]*Class, interfaceCount)
		for i, interfaceName := range class.interfaceNames {
			class.interfaces[i] = class.loader.LoadClass(interfaceName)
		}
	}
}

// 链接
func link(class *Class) {
	// 验证
	verify(class)
	// 准备
	prepare(class)
}

func verify(class *Class) {
	// TODO
}

func prepare(class *Class) {
	calcInstanceFieldSlotIds(class)
	calcStaticFieldSlotIds(class)
	allocAndInitStaticVars(class)
}

// 计算实例字段的个数
func calcInstanceFieldSlotIds(class *Class) {
	slotID := uint(0)
	if class.superClass != nil {
		slotID = class.superClass.instanceSlotCount
	}
	for _, field := range class.fields {
		if !field.IsStatic() {
			field.slotID = slotID
			slotID++
			if field.isLongOrDouble() {
				slotID++
			}
		}
	}
	class.instanceSlotCount = slotID
}

// 计算静态字段的个数
func calcStaticFieldSlotIds(class *Class) {
	slotID := uint(0)
	for _, field := range class.fields {
		if field.IsStatic() {
			field.slotID = slotID
			slotID++
			if field.isLongOrDouble() {
				slotID++
			}
		}
	}
	class.staticSlotCount = slotID
}

// 分配并初始化静态变量
func allocAndInitStaticVars(class *Class) {
	class.staticVars = newSlots(class.staticSlotCount)
	for _, field := range class.fields {
		if field.IsStatic() && field.IsFinal() {
			initStaticFinalVar(class, field)
		}
	}
}

// 初始化类常量
func initStaticFinalVar(class *Class, field *Field) {
	vars := class.staticVars
	cp := class.constantPool
	cpIndex := field.ConstValueIndex()
	slotID := field.SlotID()
	if cpIndex > 0 {
		switch field.Descriptor() {
		case "Z", "B", "C", "S", "I":
			val := cp.GetConstant(cpIndex).(int32)
			vars.SetInt(slotID, val)
		case "J":
			val := cp.GetConstant(cpIndex).(int64)
			vars.SetLong(slotID, val)
		case "F":
			val := cp.GetConstant(cpIndex).(float32)
			vars.SetFloat(slotID, val)
		case "D":
			val := cp.GetConstant(cpIndex).(float64)
			vars.SetDouble(slotID, val)
		case "Ljava/lang/String;":
			panic("todo")
		}
	}
}

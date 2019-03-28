package heap

// 方法描述
type MethodDescriptor struct {
	parameterTypes []string
	returnType     string
}

// 参数
func (self *MethodDescriptor) addParameterType(t string) {
	pLen := len(self.parameterTypes)
	// 扩容数组
	if pLen == cap(self.parameterTypes) {
		s := make([]string, pLen, pLen+4)
		copy(s, self.parameterTypes)
		self.parameterTypes = s
	}

	self.parameterTypes = append(self.parameterTypes, t)
}

func (self *MethodDescriptor) ParameterTypes() []string {
	return self.parameterTypes
}

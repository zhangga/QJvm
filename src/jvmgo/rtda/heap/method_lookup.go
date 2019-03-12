package heap

// 在类和父类中查找方法
func LookupMethodInClass(class *Class, name, descriptor string) *Method {
	for c := class; c != nil; c = c.superClass {
		for _, method := range c.methods {
			if method.name == name && method.descriptor == descriptor {
				return method
			}
		}
	}
	return nil
}

// 在类的接口中查找方法
func lookupMethodInInterfaces(interfaces []*Class, name, descriptor string) *Method {
	for _, iface := range interfaces {
		for _, method := range iface.methods {
			if method.name == name && method.descriptor == descriptor {
				return method
			}
		}
		// 递归，在接口的接口中继续查找
		method := lookupMethodInInterfaces(iface.interfaces, name, descriptor)
		if method != nil {
			return method
		}
	}
	return nil
}

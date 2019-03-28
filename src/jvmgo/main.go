package main

import (
	"fmt"
	"jvmgo/classfile"
	"jvmgo/classpath"
	"jvmgo/rtda/heap"
	"strings"
)

func main() {
	cmd := parseCmd()
	if cmd.versionFlag {
		fmt.Printf("version 0.0.1")
	} else if cmd.helpFlag || cmd.class == "" {
		printUsage()
	} else {
		// 两种方式启动。旧
		startJVM(cmd)
		// 新的启动方式。待完善，VM下有部分本地方法未实现。原书最后一章代码。
		// newJVM(cmd).start()
	}
	fmt.Printf("=========JVM RUN SUCCESSFUL=========\n", cmd.class)
}

// 已过时。旧版。新结构在jvm.go中
func startJVM(cmd *Cmd) {
	// 环境变量
	cp := classpath.Parse(cmd.XjreOption, cmd.cpOption)
	fmt.Printf("classpath: %v, class:%v, args:%v\n", cp, cmd.class, cmd.args)
	// 类加载器
	classLoader := heap.NewClassLoader(cp, cmd.verboseClassFlag)
	// 主类名
	className := strings.Replace(cmd.class, ".", "/", -1)
	// 加载主类
	mainClass := classLoader.LoadClass(className)
	// 主入口
	mainMethod := mainClass.GetMainMethod()
	if mainMethod != nil {
		interpret0(mainMethod, cmd.verboseInstFlag, cmd.args)
	} else {
		fmt.Printf("Main method not found in class %s\n", cmd.class)
	}
}

func loadClass(className string, cp *classpath.Classpath) *classfile.ClassFile {
	classData, _, err := cp.ReadClass(className)
	if err != nil {
		panic(err)
	}
	cf, err := classfile.Parse(classData)
	if err != nil {
		panic(err)
	}
	return cf
}

func getMainMethod(cf *classfile.ClassFile) *classfile.MemberInfo {
	for _, m := range cf.Methods() {
		if m.Name() == "main" && m.Descriptor() == "([Ljava/lang/String;)V" {
			return m
		}
	}
	return nil
}

func printClassInfo(cf *classfile.ClassFile) {
	fmt.Printf("version: %v.%v\n", cf.MajorVersion(), cf.MinorVersion())
	fmt.Printf("constants count: %v\n", len(cf.ConstantPool()))
	fmt.Printf("access flags: 0x%x\n", cf.AccessFlags())
	fmt.Printf("this class: %v\n", cf.ClassName())
	fmt.Printf("super class: %v\n", cf.SuperClassName())
	fmt.Printf("interfaces: %v\n", cf.InterfaceNames())
	fmt.Printf("fields count: %v\n", len(cf.Fields()))
	for _, f := range cf.Fields() {
		fmt.Printf(" %s\n", f.Name())
	}
	fmt.Printf("methods count: %v\n", len(cf.Methods()))
	for _, m := range cf.Methods() {
		fmt.Printf(" %s\n", m.Name())
	}
	fmt.Printf("=========================================================\n")
}

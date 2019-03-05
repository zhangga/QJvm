package classpath

import (
	"fmt"
	"os"
	"path/filepath"
)

type Classpath struct {
	bootClasspath Entry
	extClasspath  Entry
	userClasspath Entry
}

func Parse(jreOption, cpOption string) *Classpath {
	cp := &Classpath{}
	cp.parseBootAndExtClasspath(jreOption)
	cp.parseUserClasspath(cpOption)
	return cp
}

func (self *Classpath) parseBootAndExtClasspath(jreOption string) {
	jreDir := getJreDir(jreOption)
	// jre/lib/*
	jreLibPath := filepath.Join(jreDir, "lib", "*")
	fmt.Printf("jreLibPath: %v\n", jreLibPath)
	self.bootClasspath = newWildcardEntry(jreLibPath)
	// jre/lib/ext/*
	jreExtPath := filepath.Join(jreDir, "lib", "ext", "*")
	fmt.Printf("jreExtPath: %v\n", jreExtPath)
	self.extClasspath = newWildcardEntry(jreExtPath)
}

func (self *Classpath) parseUserClasspath(cpOption string) {
	if cpOption == "" {
		cpOption = "."
	}
	self.userClasspath = newEntry(cpOption)
}

/**
 * 获取jre路径
 */
func getJreDir(jreOption string) string {
	// 使用用户输入的-Xjre选项作为jre目录
	if jreOption != "" && exists(jreOption) {
		return jreOption
	}
	// 在当前目录下找jre目录
	if exists("./jre") {
		return "./jre"
	}
	// 使用JAVA_HOME环境变量
	if jh := os.Getenv("JAVA_HOME"); jh != "" {
		// jdk11
		// return jh
		// jdk8
		return filepath.Join(jh, "jre")
	}
	panic("Can not find jre folder!")
}

func exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (self *Classpath) ReadClass(className string) ([]byte, Entry, error) {
	className = className + ".class"
	fmt.Printf("ReadClass: %v\n", className)
	// 启动类加载
	if data, entry, err := self.bootClasspath.readClass(className); err == nil {
		return data, entry, err
	}
	// 扩展类加载
	if data, entry, err := self.extClasspath.readClass(className); err == nil {
		return data, entry, err
	}
	// 用户类加载
	return self.userClasspath.readClass(className)
}

func (self *Classpath) String() string {
	return self.userClasspath.String()
}

package main

import (
	"flag"
	"fmt"
	"os"
)

type Cmd struct {
	helpFlag         bool
	versionFlag      bool
	verboseClassFlag bool
	verboseInstFlag  bool
	cpOption         string
	XjreOption       string
	XssOption        string
	class            string
	args             []string
}

func parseCmd() *Cmd {
	cmd := &Cmd{}

	flag.Usage = printUsage
	flag.BoolVar(&cmd.helpFlag, "help", false, "print help message")
	flag.BoolVar(&cmd.helpFlag, "?", false, "print help message")
	flag.BoolVar(&cmd.versionFlag, "version", false, "print version and exit")
	flag.BoolVar(&cmd.verboseClassFlag, "verbose:class", false, "print class loader info")
	flag.BoolVar(&cmd.verboseInstFlag, "verbose:inst", false, "print instruction execute info")
	flag.StringVar(&cmd.cpOption, "classpath", "", "classpath")
	flag.StringVar(&cmd.cpOption, "cp", "", "classpath")
	flag.StringVar(&cmd.XjreOption, "Xjre", "", "path to jre")
	flag.StringVar(&cmd.XssOption, "Xss", "", "max stack space")
	flag.Parse()

	args := flag.Args()
	if len(args) > 0 {
		cmd.class = args[0]
		cmd.args = args[1:]
	} else {
		// 调试时直接设置参数
		cmd.class = "com.kwai.Test"
		cmd.cpOption = "./../../bin/class/"
		cmd.XjreOption = "C:\\JAVA\\jdk1.8\\jre\\"
		cmd.verboseClassFlag = true
		cmd.verboseInstFlag = true
		cmd.args = make([]string, 3)
		cmd.args[0] = "a"
		cmd.args[1] = "b"
		cmd.args[2] = "你好"
	}

	return cmd
}

func printUsage() {
	fmt.Printf("Usage:%s [-optins] class [args...]\n", os.Args[0])
}

package main

import (
	"flag"
	"fmt"
	"mygrpctest/cases"
	"mygrpctest/check"
	"mygrpctest/client"
	"mygrpctest/reporter"
	"mygrpctest/util"
	"os"
)

func main() {

	flag.Parse()
	if flag.NArg() != 3 {
		fmt.Printf("Usage: %s [initpron] <testfile_dir> <outputfile_dir>\n", os.Args[0])
		os.Exit(1)
	}

	climod := flag.Arg(0)
	testfile := flag.Arg(1)
	outfile := flag.Arg(2)

	//获取初始化文件集
	testcase := util.NewTestCase(testfile) //flag1

	//获取响应集
	client.Call(testcase.Initfile, testcase.Wavfile, climod, outfile) //flag 2 flag 3
	// 输出集与期待集序列化
	out := cases.NewOutput(outfile) //flag 2

	output := check.ProcessRespontoStruct(out.OutFile)
	expect := check.ExpectedtoStruct(testcase.Expectfile)
	reporter.TestNote(out.OutFile, testcase.Expectfile, "reporter/reporter.txt")
	fmt.Printf("Test result: %t \n", check.CheckRes(output, expect, "reporter/reporter.txt"))
	fmt.Println("Test result inifo is in reporter")

}

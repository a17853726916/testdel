package main

import (
	"fmt"
	"mygrpctest/cases"
	"mygrpctest/check"
	"mygrpctest/client"
	"mygrpctest/util"
)

func main() {

	//获取初始化文件集
	testcase := util.NewTestCase("testcase/case1")

	// 获取响应集
	client.Call(testcase.Initfile, testcase.Wavfile, "initpron", "outputs/case1")
	// 输出集与期待集序列化
	out := cases.NewOutput("outputs/case1")
	output := check.ProcessRespontoStruct(out.OutFile)
	expect := check.ProcessRespontoStruct(testcase.Expectfile)
	fmt.Printf("Test result: %t", check.CheckRes(output, expect))
}

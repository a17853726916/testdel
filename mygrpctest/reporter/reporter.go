package reporter

import (
	"log"
	"os"
)

/*
	1.终端提示一个testcase的结果
	2.解所有的测试结果返回到文件中
*/

// 测试目录 告知是哪一次测试
func TestNote(output, expectcase, path string) {
	file := WriteTestRes(path)
	defer file.Close()
	s := "TestCase is : " + expectcase + "\n"
	file.WriteString(s)
	s1 := "OutputCase is : " + output + "\n"
	file.WriteString(s1)
	file.WriteString("Check Result : \n")
	file.WriteString("*****************************************************\n")

}

// 第一类错误 output和expect返回结果的数量不一致的问题
func OutputlenNotEqualExpect(result, path string) {
	file := WriteTestRes(path)
	defer file.Close()
	s := "Expect's length is not equal to output: "
	file.WriteString(s)
	file.WriteString(result)
	file.WriteString("\n")

}

//第二类 请求响应出现问题
func InitOrCloseResult(result, resp, path string) {
	var s string
	file := WriteTestRes(path)
	defer file.Close()
	if resp == "initResponse" {
		s = "InitResponse's result: "
		file.WriteString(s)
		file.WriteString(result)
		file.WriteString("\n")
	} else if resp == "closeResponse" {
		s = "CloseResponse's result: "
		file.WriteString(s)
		file.WriteString(result)
		file.WriteString("\n")
	} else {
		log.Println("Params Error")
	}

}

//处理process
func ProcessResult(result, path string) {

	file := WriteTestRes(path)
	defer file.Close()
	file.WriteString(result)
	file.WriteString("\n")
}

//总的测试结果
func FilalTestres(result bool, path string) {

	file := WriteTestRes(path)
	defer file.Close()
	if result {
		file.WriteString("The testResult is Correct \n")
		file.WriteString("*****************************************************\n")
		file.WriteString("\n")
	} else {
		file.WriteString("The testResult is Incorrect \n")
		file.WriteString("*****************************************************\n")
		file.WriteString("\n")
	}
}

func WriteTestRes(path string) *os.File {
	path = "reporter/reporter.txt"
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Error in Open file: %v", err)
		panic(err)
	}
	return file
}

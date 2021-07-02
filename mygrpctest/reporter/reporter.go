package reporter

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

/*
	1.终端提示一个testcase的结果
	2.解所有的测试结果返回到文件中
*/
type RepoRests struct {
	TestCase      string   `json:"testcase"`      //测完用例目录
	OutputCase    string   `json:"outputcase"`    //输出结果目录
	ExpexLength   string   `json:"expexlength"`   //返回的结果集的长度
	InitResponse  string   `json:"initresponse"`  //初始化响应
	Processres    []string `json:"processres"`    //处理结果集
	CloseResponse string   `json:"closeresponse"` //关闭连接响应
	Finalres      bool     `json:"finalres"`      //最终结果
}

func NewRepoRests() *RepoRests {
	return &RepoRests{Processres: []string{"Correct"}}
}
func (repro *RepoRests) Get() *RepoRests {
	return repro
}

// 用于输出到text/template模板中
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

//用于输出到html模板中

func TmplHandler(w http.ResponseWriter, r *http.Request) {
	//解析模板文件
	tmpl, err := template.ParseFiles("./views/test.html")
	if err != nil {
		log.Printf("Error in ParseFiles: %v", err)
		panic(err)
	}
	//传递模板数据
	tmpl.Execute(w, "session_id")
}

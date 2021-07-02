package templates

import (
	"html/template"
	ht "html/template"
	"log"
	"os"
)

//解析模板文件是否出错
func check(err error) {
	if err != nil {
		log.Printf("Error in Parse Template File: %v", err)
		panic(err)
	}
}

//解析模板一显示输入的测试用例
func TmplTestcase(outpath, tmplpath, expectcase, outputcase string) {
	file, err := os.OpenFile(outpath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Error in OpenFile: %v", err)
		panic(err)
	}
	//获取解压模板
	tmpl := template.Must(template.ParseFiles(tmplpath))
	check(tmpl.ExecuteTemplate(file, "TestCase", expectcase))
	check(tmpl.ExecuteTemplate(file, "OutputCase", outputcase))
}

//解析模板二显示输入长度是否出现问题
func TmplLengthnotEuqal(outpath, tmplpath, res string) {

	file, err := os.OpenFile(outpath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Error in OpenFile: %v", err)
		panic(err)
	}
	//获取解压模板
	tmpl := template.Must(template.ParseFiles(tmplpath))
	check(tmpl.ExecuteTemplate(file, "OutputlenNotEqualExpect", res))
}

// 模板输出初始化响应与结束响应的结果
func TmplInitOrClose(outpath, tmplpath, resp, res string) {
	file, err := os.OpenFile(outpath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Error in OpenFile: %v", err)
		panic(err)
	}
	tmpl := template.Must(template.ParseGlob(tmplpath))
	if resp == "initResponse" {
		check(tmpl.ExecuteTemplate(file, "InitResponse", res))
	} else if resp == "closeResponse" {
		check(tmpl.ExecuteTemplate(file, "CloseResponse", res))
	} else {
		log.Printf("Params Error")
	}

}

// 处理Process结果
func TmplProcess(outpath, tmplpath, res string) {
	file, err := os.OpenFile(outpath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Error in OpenFile: %v", err)
		panic(err)
	}
	//获取解压模板
	tmpl := template.Must(template.ParseFiles(tmplpath))
	tmpl.Funcs(template.FuncMap{"safe": safe})

	check(tmpl.ExecuteTemplate(file, "Process", res))
}

//处理最终结果
func Tmplfinalres(outpath, tmplpath string, res bool) {

	file, err := os.OpenFile(outpath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Error in OpenFile: %v", err)
		panic(err)
	}
	//获取解压模板
	tmpl := template.Must(template.ParseFiles(tmplpath))
	check(tmpl.ExecuteTemplate(file, "FilalTestres", res))
}

func OpentmplFile(path string) *os.File {
	path = "reporter/reporter.txt"
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Error in Open file: %v", err)
		panic(err)
	}
	return file
}

// 处理字符中嵌套的各类符号如 < > "" 等符号
func safe(s string) interface{} {
	c := ht.HTML(s)
	return c
}

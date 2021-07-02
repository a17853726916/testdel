package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"mygrpctest/cases"
	"mygrpctest/check"
	"mygrpctest/reporter"
	"mygrpctest/templates"
	"mygrpctest/util"
	"net/http"
	"os"
)

func main() {

	flag.Parse()
	if flag.NArg() != 3 {
		fmt.Printf("Usage: %s [initpron] <testfile_dir> <outputfile_dir>\n", os.Args[0])
		os.Exit(1)
	}

	// climod := flag.Arg(0)
	testfile := flag.Arg(1)
	outfile := flag.Arg(2)
	reporterRes := reporter.NewRepoRests()
	//获取初始化文件集
	testcase := util.NewTestCase(testfile) //flag1

	//获取响应集
	// client.Call(testcase.Initfile, testcase.Wavfile, climod, outfile) //flag 2 flag 3
	// 输出集与期待集序列化
	out := cases.NewOutput(outfile) //flag 2

	output := check.ProcessRespontoStruct(out.OutFile)
	expect := check.ExpectedtoStruct(testcase.Expectfile)
	reporter.TestNote(out.OutFile, testcase.Expectfile, "reporter/reporter.txt")
	reporterRes.TestCase = testcase.Expectfile
	reporterRes.OutputCase = out.OutFile
	templates.TmplTestcase("reporter/reporter1.txt", "templates/tmpl.txt", testcase.Expectfile, out.OutFile)
	fmt.Printf("Test result: %t \n", check.CheckRes(output, expect, "reporter/reporter.txt", reporterRes))
	fmt.Println("Test result inifo is in reporter")
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		res := *reporterRes.Get()
		tmpl, err := template.ParseFiles("templates/views/report.html", "templates/views/content.html")
		if err != nil {
			log.Printf("Error in ParseFiles: %v", err)
			panic(err)
		}
		if len(res.Processres) > 1 {
			res.Processres = reporterRes.Processres[1:]
		}
		tmpl.Execute(rw, res)
	})
	http.ListenAndServe("localhost:8080", nil)

}

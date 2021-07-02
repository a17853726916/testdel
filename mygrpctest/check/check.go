package check

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	pb "mygrpctest/proto/voice"
	"mygrpctest/reporter"
	"mygrpctest/templates"
	"os"
	"strconv"

	"strings"
)

// 定义比较结果返回字段
type Result struct {
	SessionId    string                         `json:"session_id,omitempty"`    // 会话UUID
	StatusCode   int32                          `json:"status_code,omitempty"`   // 评估结果状态码，0：人声；1：静音；2：噪声；3：空包；
	StatusMesg   string                         `json:"status_mesg,omitempty"`   // 评估结果信息
	ModelVersion string                         `json:"model_version,omitempty"` // 模型版本信息
	Sents        []*pb.ProcessResponse_SentInfo `json:"sents,omitempty"`         // 多句子候选
	AudioUrl     string                         `json:"audio_url,omitempty"`     // 语音音频文件下载地址，单词结束时给出
}

func ProcessRespontoStruct(path string) []Result {
	inputFile, inputError := os.Open(path)
	if inputError != nil {
		return nil
	}
	defer inputFile.Close()
	var s string = "["
	inputReader := bufio.NewReader(inputFile)
	for {
		inputString, readerError := inputReader.ReadString('\n')
		if readerError == io.EOF {

			break
		}
		s = s + inputString
	}

	s1 := strings.TrimSpace(s)
	res := []byte(s1)
	res[len(res)-1] = ']'
	var resps []Result
	// 直接序列化多个结构体数据
	// 只要json文件结构正确，可以序列化多个
	if err := json.Unmarshal(res, &resps); err == nil {
		fmt.Println("convert successful")
	} else {
		fmt.Println("convert fail")
	}
	return resps
}

// 序列化返回结果
func ExpectedtoStruct(path string) []Result {
	inputFile, inputError := os.Open(path)
	if inputError != nil {
		return nil
	}
	defer inputFile.Close()
	var s string = "["
	inputReader := bufio.NewReader(inputFile)
	for {
		inputString, readerError := inputReader.ReadString('\n')
		if readerError == io.EOF {

			break
		}
		s = s + inputString
	}
	s1 := strings.TrimSpace(s)
	res := []byte(s1)
	res[len(res)-1] = ']'
	// fmt.Printf("all: \n%s", res)
	var resps []Result
	// 直接序列化多个结构体数据
	// 只要json文件结构正确，可以序列化多个
	if err := json.Unmarshal(res, &resps); err == nil {
		fmt.Println("convert successful")
	} else {
		fmt.Println("conver failed")
	}
	return resps
}

//结构体比较
func CheckRes(outputs, expects []Result, outpath string, reportertest *reporter.RepoRests) bool {
	//初始化
	var tmplpath string = "templates/tmpl.txt"
	var outpath1 string = "reporter/reporter1.txt"
	var init, process, close, final bool
	var testResult string
	if len(outputs) != len(expects) {
		testResult = fmt.Sprintf("The expected length is: %v,now is: %v \n", len(expects), len(outputs))
		reporter.OutputlenNotEqualExpect(testResult, outpath)
		//使用模板文件
		templates.TmplLengthnotEuqal(outpath1, tmplpath, testResult)
		reportertest.ExpexLength = testResult
		final = false
	}

	iniout, expect := outputs[0], expects[0]

	if iniout.StatusMesg == expect.StatusMesg && expect.ModelVersion == expect.ModelVersion {
		//将测试结果写回文件中
		testResult = fmt.Sprint("The initResponse Correct \n")
		reporter.InitOrCloseResult(testResult, "initResponse", outpath)
		//使用模板文件
		templates.TmplInitOrClose(outpath1, tmplpath, "initResponse", testResult)
		reportertest.InitResponse = testResult
		init = true
	} else {
		//将测试结果写回文件中
		testResult = fmt.Sprintf("Error in initResponse , The expect initResponse is: %v, now is: %v \n", expect.StatusMesg, iniout.StatusCode)
		reporter.InitOrCloseResult(testResult, "initResponse", outpath)
		templates.TmplInitOrClose(outpath1, tmplpath, "initResponse", testResult)
		reportertest.InitResponse = testResult
	}

	// 中间处理结果
	iniprocess, exprocess := outputs[1:len(outputs)-1], expects[1:len(expects)-1]
	process = checkProcess(iniprocess, exprocess, outpath, reportertest)
	// initlength, exlength := len(iniprocess), len(exprocess)

	// 关闭通道的结果
	iniclose, expectclose := outputs[len(outputs)-1], expects[len(expects)-1]
	if iniclose.StatusMesg == expectclose.StatusMesg {
		testResult = fmt.Sprint("The closeResponse Correct \n")
		reporter.InitOrCloseResult(testResult, "closeResponse", outpath)
		templates.TmplInitOrClose(outpath1, tmplpath, "closeResponse", testResult)
		reportertest.CloseResponse = testResult
		close = true
	} else {
		testResult = fmt.Sprintf("Error in CloseResponse , The expect closeResponse is: %v, now is: %v \n", expectclose.StatusMesg, iniclose.StatusMesg)
		reporter.InitOrCloseResult(testResult, "closeResponse", outpath)
		templates.TmplInitOrClose(outpath1, tmplpath, "closeResponse", testResult)
		reportertest.CloseResponse = testResult
		fmt.Println("initcose is not the expected")
	}
	final = init && process && close
	reporter.FilalTestres(final, outpath)
	templates.Tmplfinalres(outpath1, tmplpath, final)
	reportertest.Finalres = final
	return final
}

func checkProcess(iniprocess, exprocess []Result, path string, reportertest *reporter.RepoRests) bool {
	var process bool
	var flag bool = true
	var tmplpath string = "templates/tmpl.txt"
	var outpath1 string = "reporter/reporter1.txt"
	for i, v := range iniprocess {
		if i < len(exprocess) {
			if v.StatusCode == exprocess[i].StatusCode && v.StatusMesg == exprocess[i].StatusMesg {
				for j, ws := range v.Sents {

					for k, wd := range ws.Words {

						if len(ws.Words) != len(exprocess[i].Sents[j].Words) {
							fmt.Printf("The result's lenght is not unanimous , expect's lenght is: %v , now is: %v , error in words slice: %v\n", len(exprocess[i].Sents[j].Words), len(ws.Words), i)
							flag = false
							continue

						} else {
							if wd.Word == exprocess[i].Sents[j].Words[k].Word {

								for l, ph := range wd.Phones {
									if ph.Phone == exprocess[i].Sents[j].Words[k].Phones[l].Phone && ph.RefPhone == exprocess[i].Sents[j].Words[k].Phones[l].RefPhone {
										process = true
									} else {
										testResult := fmt.Sprintf("The expected Phone is: %v , now is: %v ,info is in the output.txt of line: %d  words: %d , Phones: %d\n", exprocess[i].Sents[j].Words[k].Phones[l].Phone, ph.Phone, i+2, k+1, l+1)
										reporter.ProcessResult(testResult, path)
										templates.TmplProcess(outpath1, tmplpath, testResult)
										reportertest.Processres = append(reportertest.Processres, testResult)
										flag = false
										continue
									}
								}
								for l, sylls := range wd.Syllables {
									if sylls.Match == exprocess[i].Sents[j].Words[k].Syllables[l].Match && sylls.Syllable == exprocess[i].Sents[j].Words[k].Syllables[l].Syllable {
									} else {
										testResult := fmt.Sprintf("The expected sylls is %v , now is %v ,info is in the output.txt of line: %d  words:%d sylls: %d\n", exprocess[i].Sents[j].Words[k].Syllables[l].Syllable, sylls.Syllable, i+2, k+1, l+1)
										reporter.ProcessResult(testResult, path)
										templates.TmplProcess(outpath1, tmplpath, testResult)
										reportertest.Processres = append(reportertest.Processres, testResult)
										flag = false
										continue
									}
								}
								process = true
							} else {
								// 将测试结果写回文件中
								flag = false
								testResult := fmt.Sprintf("The expected word is %v , now is %v, info is in the output.txt of line: %d sents: %d words: %d\n", exprocess[i].Sents[j].Words[k].Word, wd.Word, i+2, j+1, k+1)
								reporter.ProcessResult(testResult, path)
								templates.TmplProcess(outpath1, tmplpath, testResult)
								reportertest.Processres = append(reportertest.Processres, testResult)
								continue
							}
						}
					}
				}
			} else {
				flag = false
				testResult := fmt.Sprintf("The expected is： %v , now is: %v ", strconv.Itoa(int(v.StatusCode))+" "+v.StatusMesg, strconv.Itoa(int(exprocess[i].StatusCode))+"  "+exprocess[i].StatusMesg)
				fmt.Println(testResult)
				continue
			}
		} else {
			break
		}

	}
	return process && flag
}

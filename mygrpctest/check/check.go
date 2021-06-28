package check

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	pb "mygrpctest/proto/voice"
	"mygrpctest/reporter"
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
func CheckRes(outputs, expects []Result, path string) bool {
	//初始化
	var init, process, close, final bool
	var testResult string
	if len(outputs) != len(expects) {
		testResult = fmt.Sprintf("The expected length is: %v,now is: %v \n", len(expects), len(outputs))
		reporter.OutputlenNotEqualExpect(testResult, path)
		final = false
	}

	iniout, expect := outputs[0], expects[0]

	if iniout.StatusMesg == expect.StatusMesg && expect.ModelVersion == expect.ModelVersion {
		//将测试结果写回文件中
		testResult = fmt.Sprint("The initResponse Correct \n")
		reporter.InitOrCloseResult(testResult, "initResponse", path)
		init = true
	} else {
		//将测试结果写回文件中
		testResult = fmt.Sprintf("The expect initResponse is: %v, now is: %v \n", expect.StatusMesg, iniout.StatusCode)
		reporter.InitOrCloseResult(testResult, "initResponse", path)
	}

	// 中间处理结果
	iniprocess, exprocess := outputs[1:len(outputs)-1], expects[1:len(expects)-1]
	process = checkProcess(iniprocess, exprocess, path)
	// initlength, exlength := len(iniprocess), len(exprocess)

	// 关闭通道的结果
	iniclose, expectclose := outputs[len(outputs)-1], expects[len(expects)-1]
	if iniclose.StatusMesg == expectclose.StatusMesg {
		testResult = fmt.Sprint("The closeResponse Correct \n")
		reporter.InitOrCloseResult(testResult, "closeResponse", path)
		close = true
	} else {
		testResult = fmt.Sprintf("The expect initClose is: %v, now is: %v \n", expectclose.StatusMesg, iniclose.StatusMesg)
		reporter.InitOrCloseResult(testResult, "closeResponse", path)
		fmt.Println("initcose is not the expected")
	}
	final = init && process && close
	reporter.FilalTestres(final, path)
	return final
}

//获得初始化响应对比字段
func InitRespontoStruct(path string) []pb.InitResponse {
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
	s1 := []byte(s)
	s1[len(s1)-2] = ']'
	s = string(s1)
	//fmt.Printf("all: \n%s", s)
	var resps []pb.InitResponse
	// 直接序列化多个结构体数据
	// 只要json文件结构正确，可以序列化多个
	if err := json.Unmarshal([]byte(s), &resps); err == nil {
		fmt.Println("转换成功")
	} else {
		fmt.Println("转换失败")
	}
	return resps
}

//获得结束响应字段返回结果
func CloseRespontoStruct(path string) []pb.CloseResponse {
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
	s1 := []byte(s)
	s1[len(s1)-2] = ']'
	s = string(s1)
	//fmt.Printf("all: \n%s", s)
	var resps []pb.CloseResponse
	// 直接序列化多个结构体数据
	// 只要json文件结构正确，可以序列化多个
	if err := json.Unmarshal([]byte(s), &resps); err == nil {
		fmt.Println("转换成功")
	} else {
		fmt.Println("转换失败")
	}
	return resps
}
func checkProcess(iniprocess, exprocess []Result, path string) bool {
	var process bool
	var flag bool = true

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
										testResult := fmt.Sprintf("The expected Phone is: %v , now is: %v ,info is in the output.txt of line: %d  words: %d , Phones: %d\n", exprocess[i].Sents[j].Words[k].Phones[l], ph, i+2, k+1, l+1)
										reporter.ProcessResult(testResult, path)
										flag = false
										continue
									}
								}
								for l, sylls := range wd.Syllables {
									if sylls.Match == exprocess[i].Sents[j].Words[k].Syllables[l].Match && sylls.Syllable == exprocess[i].Sents[j].Words[k].Syllables[l].Syllable {
									} else {
										testResult := fmt.Sprintf("The expected sylls is %v , now is %v ,info is in the output.txt of line: %d  words:%d sylls: %d\n", exprocess[i].Sents[j].Words[k].Syllables[l], sylls, i+2, k+1, l+1)
										reporter.ProcessResult(testResult, path)
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

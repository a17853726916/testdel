package client

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	pb "mygrpctest/proto/voice"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

//指定内容的响应结果集
type ProcessRes struct {
	SessionId  string                         `json:"session_id,omitempty"`  // 会话UUID
	StatusCode int32                          `json:"status_code,omitempty"` // 评估结果状态码，0：人声；1：静音；2：噪声；3：空包；
	StatusMesg string                         `json:"status_mesg,omitempty"` // 评估结果信息
	Sents      []*pb.ProcessResponse_SentInfo `json:"sents,omitempty"`       // 多句子候选

}
type ProcessResponse_SentInfo_Res struct {
	Words []*ProcessResponse_SentInfo_WordInfo_Res `json:"words,omitempty"` // 单词评估结果

}
type ProcessResponse_SentInfo_WordInfo_Res struct {
	Match     int32                                              `son:"match,omitempty"`      // 单词评估情况，0：正常；1：增读；2：漏读；3：错读；4：回读
	Word      string                                             `json:"word,omitempty"`      // 当前单词
	Phones    []*ProcessResponse_SentInfo_WordInfo_PhoneInfo_Res `json:"phones,omitempty"`    // 音素评估结果
	Syllables []*ProcessResponse_SentInfo_WordInfo_SyllInfo_Res  `json:"syllables,omitempty"` // 音节评估结果
}
type ProcessResponse_SentInfo_WordInfo_PhoneInfo_Res struct {
	Match    int32  `json:"match,omitempty"`     // 音素评估情况，0：正常；1：增读；2：漏读；3：错读；4：回读
	Phone    string `json:"phone,omitempty"`     // 当前音素
	RefPhone string `json:"ref_phone,omitempty"` // 参考音素
}

type ProcessResponse_SentInfo_WordInfo_SyllInfo_Res struct {
	Match    int32  `json:"match,omitempty"`    // 音节评估情况，0：正常；1：增读；2：漏读；3：错读；4：回读
	Syllable string `json:"syllable,omitempty"` // 当前音节
}

// 远程连接地址
const (
	RemoteAddr = "39.101.164.32:9999"
)

func Call(initfile, wavfile string, climode string, path string) {
	head_len := flag.Int("header", 78, "header length")
	pack_len := flag.Int("packet", 16000, "packet length")
	repeat := flag.Int("repeat", 1, "repeat count")
	conn, err := grpc.Dial(RemoteAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	clt := pb.NewVoiceServiceClient(conn)

	flag.Parse()
	if flag.NArg() != 3 {
		fmt.Printf("Usage: %s [init|initpron|word] <initfile> <wavfile>\n", os.Args[0])
	}
	//获取初始化文件内容和测试文件的内容
	initData := GetFilebyte(initfile)
	wavData := GetFilebyte(wavfile)

	//到这一步数据不是空的(initData)
	ch := make(chan bool, *repeat) //定义一个有缓冲的管道
	for i := 0; i < *repeat; i++ {
		go func(initData, wavData []byte, clt pb.VoiceServiceClient) {
			if climode == "init" {
				TestSentClient("init", clt, initData, wavData, *head_len, *pack_len, path)
			} else if climode == "initpron" {
				TestSentClient("initpron", clt, initData, wavData, *head_len, *pack_len, path)

			} else if climode == "word" {
				TestWordClient(clt, initData, wavData, *head_len, *pack_len)
			} else {
			}
			ch <- true

		}(initData, wavData, clt)
	}
	for i := 0; i < *repeat; i++ {
		<-ch
	}
}

//初始化请求
func LoadInitRequest(sessionId string, initData []byte) (*pb.InitRequest, error) {
	req := &pb.InitRequest{}
	err := json.Unmarshal(initData, req)
	if err != nil {
		log.Printf("Error in LoadInitRequest: %v,%v", sessionId, err)
		return nil, err
	}
	req.SessionId = sessionId
	//因为这里发送的内容是没有数据的

	return req, nil

}

//发音初始化请求
func LoadInitPronRequest(sessionID string, initData []byte) (*pb.InitPronRequest, error) {
	req := &pb.InitPronRequest{}
	err := json.Unmarshal(initData, req)
	if err != nil {
		log.Printf("Error in LoadInitPronRequest: %v %v ", sessionID, err)
		return nil, err
	}
	req.SessionId = sessionID
	return req, nil
}

//流发送请求
func HandleProcessRequest(stream pb.VoiceService_ProcessClient, sessionID string, wavData []byte, head, pack int, path string) error {
	wav := wavData[head:]
	off, size, islast := 0, pack, 0
	for off < len(wav) {
		var audioData []byte
		if len(wav) < off+size {
			audioData = wav[off:]
			islast = 1
		} else {
			audioData = wav[off : off+size]
		}
		off += size
		req := &pb.ProcessRequest{
			SessionId: sessionID,
			AuData:    audioData,
			AuRate:    16000,
			AuFormat:  0,
			IsLast:    int32(islast),
		}
		err := stream.Send(req)
		if err != nil {
			log.Printf("Error in HandleProcessRequest: %v %v ", sessionID, err)
			return err
		}
		rsp, err := stream.Recv()

		if err != nil {
			log.Printf("Error in HandleProcessRequest: %v %v ", sessionID, err)
			return err
		}

		log.Printf("Receive ProcessRequest: %v %v", sessionID, rsp)
		time.Sleep(1 * time.Second)
		//带发音的响应处理加入
		HandleProcessResponse(rsp, path)
	}

	if err := stream.CloseSend(); err != nil {
		log.Printf("Error in HandleProcessRequest: %v %v ", sessionID, err)
		return err
	}
	return nil
}

//测试发送请求的客户端
func TestSentClient(mod string, clt pb.VoiceServiceClient, initData, wavData []byte, head, pack int, path string) error {
	sessionId := uuid.New().String()
	if mod == "init" {
		initReq, err := LoadInitRequest(sessionId, initData)
		initRsp, err := clt.Initial(context.TODO(), initReq)
		if err != nil {
			log.Printf("Error in TestSentClient: %v %v", sessionId, err)
			return err
		}
		log.Printf("Receive InitResponse: %v %v", sessionId, initRsp)
		//第一次响应参数的读取
		HandleInitResp(initRsp, path)

	} else {
		initReq, err := LoadInitPronRequest(sessionId, initData) //初始化回话带发音
		initRsp, err := clt.InitialWithPron(context.TODO(), initReq)

		if err != nil {
			log.Printf("Error in TestSentClient: %v %v", sessionId, err)
			return err
		}
		//第一次响应参数的读取
		HandleInitProesp(initRsp, path)
		log.Printf("Receive InitPronResponse: %v %v", sessionId, initRsp)
	}
	//处理请求的信息流
	proStream, err := clt.Process(context.TODO())
	if err != nil {
		log.Printf("Error in TestSentClient: %v %v", sessionId, err)
		return err
	}
	err = HandleProcessRequest(proStream, sessionId, wavData, head, pack, path)
	if err != nil {
		log.Printf("Error in TestSentClient: %v %v", sessionId, err)
		return err
	}

	//通道关闭请求
	closeReq := &pb.CloseRequest{}
	closeReq.SessionId = sessionId
	closeRsp, err := clt.Close(context.TODO(), closeReq)
	if err != nil {
		log.Printf("Error in TestSentClient: %v %v", sessionId, err)
		return err
	}
	log.Printf("Receive CloseResponse: %v %v", sessionId, closeRsp)
	//接收返回的响应结果
	HandleCloseProesp(closeRsp, path)
	return nil
}

//加载发音评判
func LoadWordRequest(sessionId string, initData, wavData []byte, head, pack int) (*pb.WordEvalRequest, error) {
	req := &pb.WordEvalRequest{}
	err := json.Unmarshal(initData, req)
	if err != nil {
		log.Printf("Error in LoadWordRequest: %v %v", sessionId, err)
		return nil, err
	}
	req.AuData = wavData[head : head+pack]
	req.SessionId = sessionId
	return req, nil
}

func TestWordClient(clt pb.VoiceServiceClient, initData, wavData []byte, head, pack int) error {
	sessionId := uuid.New().String()
	wordReq, err := LoadWordRequest(sessionId, initData, wavData, head, pack)
	//fmt.Println(wordReq)
	wordRsp, err := clt.WordEval(context.TODO(), wordReq)
	if err != nil {
		log.Printf("Error in TestWordClient: %v %v", sessionId, err)
		return err
	}
	log.Printf("Receive WordEvalResponse: %v %v", sessionId, wordRsp)
	return nil
}

//将响应户数写回json文件中
func HandleProcessResponse(resData *pb.ProcessResponse, path string) error {

	file, err := os.OpenFile(path+"/output.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0645)
	if err != nil {
		log.Printf("Open or Create file error: %v", err)
		panic(err)
	}
	defer file.Close()

	//获取sents

	var processRes = &ProcessRes{
		SessionId:  resData.SessionId,
		StatusCode: resData.StatusCode,
		StatusMesg: resData.StatusMesg,
		Sents:      resData.Sents,
	}

	data, err := json.Marshal(processRes)
	if err != nil {
		log.Printf("Error in Marshal: %v", err)
		log.Panic(err)
	}

	data = append(data, ',')
	data1 := string(data)
	data1 += "\n"
	_, err = file.WriteString(data1)
	if err != nil {
		log.Printf("Write into File Error: %v", err)
		panic(err)
	}

	return nil
}

// 将响应的初始化请求写入文件中
func HandleInitResp(resData *pb.InitResponse, path string) error {

	CreateDir(path)
	file, err := os.OpenFile(path+"/initoutput.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0645)
	if err != nil {
		log.Printf("Open or Create file error: %v", err)
		panic(err)
	}
	defer file.Close()

	data, err := json.Marshal(resData)

	if err != nil {
		log.Printf("Error in Marshal: %v", err)
		log.Panic(err)
	}

	data = append(data, ',')
	data1 := string(data)
	data1 += "\n"
	_, err = file.WriteString(data1)
	if err != nil {
		log.Printf("Write into File Error: %v", err)
		panic(err)
	}

	return nil
}

//带发音的初始化请求响应结果
func HandleInitProesp(resData *pb.InitPronResponse, path string) error {

	CreateDir(path)

	file, err := os.OpenFile(path+"/output.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0645)
	if err != nil {
		log.Printf("Open or Create file error: %v", err)
		panic(err)
	}
	defer file.Close()

	data, err := json.Marshal(resData)

	if err != nil {
		log.Printf("Error in Marshal: %v", err)
		log.Panic(err)
	}

	data = append(data, ',')
	data1 := string(data)
	data1 += "\n"
	_, err = file.WriteString(data1)
	if err != nil {
		log.Printf("Write into File Error: %v", err)
		panic(err)
	}

	return nil
}

//关闭通信响应结果
func HandleCloseProesp(resData *pb.CloseResponse, path string) error {
	file, err := os.OpenFile(path+"/output.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0645)
	if err != nil {
		log.Printf("Open or Create file error: %v", err)
		panic(err)
	}
	defer file.Close()

	data, err := json.Marshal(resData)

	if err != nil {
		log.Printf("Error in Marshal: %v", err)
		log.Panic(err)
	}

	data = append(data, ',')
	data1 := string(data)
	data1 += "\n"
	_, err = file.WriteString(data1)
	if err != nil {
		log.Printf("Write into File Error: %v", err)
		panic(err)
	}

	return nil
}

//获取文件内容
func GetFilebyte(path string) []byte {
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Error in Read file: %v", err)
		panic(err)
	}

	return fileData
}

//自动生生目录(若没有)
func CreateDir(path string) {
	ss := strings.Split(path, "/")
	basepath := ss[0]
	foldname := ss[1]
	foldpaht := filepath.Join(basepath, foldname)
	if _, err := os.Stat(foldpaht); os.IsNotExist(err) {
		os.Mkdir(foldpaht, 0777)
		os.Chmod(foldpaht, 0777)
	}
}

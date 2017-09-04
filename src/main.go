package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"git.oschina.net/JMArch/rpc-go/client/config"
	"git.oschina.net/JMArch/rpc-go/client/transport"
)

func main() {
	//载入配置文件。默认地址在conf/config.toml
	var interfaceTxtPath string
	doParams := make(map[int32][]interface{})
	doParamsCount := 0
	configPath := getParentDirectory(getCurrentDirectory())
	filePath := configPath + "/conf/config.toml"
	if _, fileError := os.Stat(filePath); fileError != nil {
		if os.IsNotExist(fileError) {
			config.LoadConfig("../conf/config.toml")
		} else {
			fmt.Println("other error")
		}
	} else {
		config.LoadConfig(configPath + "/conf/config.toml")
	}
	if _, interfaceTxtPathErr := os.Stat(configPath + "/conf/interface.txt"); interfaceTxtPathErr != nil {
		if os.IsNotExist(interfaceTxtPathErr) {
			interfaceTxtPath = "../conf/interface.txt"
		} else {
			fmt.Println("interfaceTxt path error")
		}
	} else {
		interfaceTxtPath = configPath + "/conf/interface.txt"
	}
	interfaceTxtHandle, interfaceTxtErr := os.Open(interfaceTxtPath)
	if interfaceTxtErr != nil {
		fmt.Println("参数配置文件读取出错")
	}
	rd := bufio.NewReader(interfaceTxtHandle)
	defer func() {
		interfaceTxtHandle.Close()
	}()
	for {
		line, _, err := rd.ReadLine()
		if err != nil || io.EOF == err {
			break
		}
		doParamsCount++
		lineArr := strings.Split(string(line), "	")
		keyInt, err := strconv.Atoi(lineArr[0])
		valueSlice := []interface{}{}
		lineValueArr := strings.Split(lineArr[1], ",")
		for _, value := range lineValueArr {
			value = strings.Replace(value, "\"", "", -1)
			value = strings.Replace(value, "[", "", -1)
			value = strings.Replace(value, "]", "", -1)
			valueSlice = append(valueSlice, value)
		}
		doParams[int32(keyInt)] = valueSlice
	}
	endPointAddr, err := transport.ParseEndPoint("prod")
	if err != nil {
		fmt.Println(err.Error())
	}

	starttime = time.Now()
	//启动多线程测试.每个线程循环测试
	var ws sync.WaitGroup
	for i := 0; i < config.RPCEndPointMap.Threadnum; i++ {
		ws.Add(1)
		randCurrent := 0
		cdoParams := make(map[int32][]interface{})
		for key, value := range doParams {
			randCurrent++
			if randCurrent == rand.Intn(doParamsCount) {
				cdoParams[key] = value
				break
			}
		}
		go testCallRPC(endPointAddr, &ws, cdoParams)
	}
	fmt.Println(time.Now().Sub(starttime).Seconds())
	//time.Sleep(time.Second * time.Duration(config.RPCEndPointMap.TestTime))
	ws.Wait()
	fmt.Printf("线程数:%d,成功调用%d次,失败调用%d次,共耗时:%f秒\n", config.RPCEndPointMap.Threadnum, successNum, errorNum, time.Now().Sub(starttime).Seconds())
}

//用来计算服务器压力测试的开始时间
var starttime time.Time

//在指定时间内完成的rpc 请求数
var successNum int32
var errorNum int32
var zjson map[string]interface{}

//testCallRPC 循环压力测试 一个rpc
func testCallRPC(endPointAddr *transport.JumeiEndPoint, ws *sync.WaitGroup, doParams map[int32][]interface{}) {
	defer func() {
		ws.Add(-1)
	}()
	uid := int32(2000000715)
	hashID := []interface{}{"cd130128p26545", "d140816p655hg", "d150721p38998", "d160401p115", "df151104p222500097", "df160527p222551495", "ht141126p4000141t1", "ht1466585745p800000076", "ht150205p222400663t1", "ht160421p222551386t1", "d140805p385hg", "d150514p00", "d160120p222550374", "df14963931239895p222550419", "df160324p222551102", "df1612095715p810000255", "ht1466069465p800000018", "ht150203p222400860t1", "ht160417p222551381t1", "ht170316p810000460t9", "d140408p162sg", "d150205p", "d151225p222550223", "d160602p222551529", "df151209p222550053", "df1612075427p810000237", "ht1419318017p2", "ht150131p222400827t1", "ht160410p222551384t1", "ht170212p222551393t1zh", "cd160428p177649", "d150129p22346", "d150721p389981", "d160602p222551527a", "df151203p222500933", "df161102201132p810000201", "ht141127p4000130t1", "ht1481185000p810000248", "ht150315p222400873t1", "ht170209p222551386t1"}
	for k, v := range doParams {
		uid = k
		hashID = v
	}
	beingTime := "2017-01-01 00:00:00"
	endTime := "2017-08-26 00:00:00"
	var params []interface{}
	params = []interface{}{uid, hashID, beingTime, endTime}
	response, err := endPointAddr.Call("Financial", "getPaidDealsByUidAndDealHashIds", params, false)
	if err != nil {
		atomic.AddInt32(&errorNum, 1)
	} else {
		if ok := strings.Contains(response, "ok"); ok {
			atomic.AddInt32(&successNum, 1)
		} else {
			atomic.AddInt32(&errorNum, 1)
		}
	}
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("error")
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

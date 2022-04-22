package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func CheckFile(EtlFilePath string) bool {
	_, err := os.Stat(EtlFilePath)
	if strings.HasSuffix(EtlFilePath, "etl") && err == nil {
		//mapPtr := ParseFile(EtlFilePath)
		//pid := GetMetaDataSlice(mapPtr)
		//return pid
		return true
	} else {
		return false
	}
}
func GetCsvFilePath(EtlFilePath string) string {
	var fileSilce = strings.Split(EtlFilePath, "\\")
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	str := tm.Format("20060102030405")
	fileSilce[len(fileSilce)-1] = str + "ICMP.csv"
	CsvFilePath := strings.Join(fileSilce, "\\")
	return CsvFilePath
}

func GetResultFilePath(CsvFilePath string) string {
	var fileSilce = strings.Split(CsvFilePath, "\\")
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	str := tm.Format("20060102030405")
	fileSilce[len(fileSilce)-1] = str + "result.csv"
	ResultFilePath := strings.Join(fileSilce, "\\")
	return ResultFilePath
}

func ParseFile(CsvFilePath string) *map[int]map[string]string {
	csvfile, err := os.Open(CsvFilePath)
	if err != nil {
		log.Fatal("打开CSV文件失败\n")
	}
	defer csvfile.Close()
	r := csv.NewReader(csvfile)
	r.LazyQuotes = true    //防止出现UTF-8编码的BOM格式头出现异常
	r.FieldsPerRecord = -1 //csv reader 会检查每一行的字段数量，如果不等于 FieldsPerRecord 就会抛出该错误。 为负数则不检查
	i := 0
	SouceData := make(map[int]map[string]string)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		//Microsoft-Windows-NDIS-PacketCapture 代表存在data数据
		if strings.Contains(strings.TrimSpace(record[0]), "Microsoft-Windows-NDIS-PacketCapture") && len(record) > 21 {
			//因为 PID 在 CSV 中是固定的 所以通过索引来检查即可
			varStr := record[21]
			if len(varStr) > 12 && !strings.Contains(varStr[1:12], "0xFFFFF") { //排除ARP数据帧
				if len(varStr) > 49 && varStr[49:51] == "01" { //偏移量确认通信协议为ICMP
					SouceData[i] = make(map[string]string, 5)
					SouceData[i]["SRCIP"] = ParseIp(varStr[55:63]) //通信SRCIP
					SouceData[i]["DSTIP"] = ParseIp(varStr[63:71])
					SouceData[i]["DATA"] = varStr
					str := record[8]
					str = strings.TrimSpace(str)
					if strings.Contains(str, "0x") && len(str) == 10 {
						str = str[2:]
						var1 := ParseHex(str)
						SouceData[i]["PID"] = var1
					}
					str = strings.TrimSpace(record[9])
					//排除其他干扰事件
					if strings.Contains(str, "0x") && len(str) == 10 {
						str = str[2:]
						var1 := ParseHex(str)
						SouceData[i]["TID"] = var1
					}
				}

			}
		}
		i++
	}
	return &SouceData
}

func ParseHex(hex string) string {
	n, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		panic(err)
	}
	str1 := strconv.FormatInt(n, 10)
	return str1 //当前程序的PID或者TID
}

func ParseIp(str string) string {
	realIp := ""
	//fmt.Println(len(SrcIP))
	for ipindex := 0; ipindex <= 6; ipindex = ipindex + 2 {
		if ipindex != 6 {
			realIp += ParseHex(str[ipindex:ipindex+2]) + "."
			continue
		}
		realIp += ParseHex(str[ipindex : ipindex+2])
	}
	return realIp
}
func OutPutFile(filepath string, ptr *map[int]map[string]string) {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer f.Close() //压栈释放资源
	f.WriteString("\xEF\xBB\xBF")
	writer := csv.NewWriter(f)
	var header = []string{"编号", "进程(PID)", "线程(TID)", "源IP(SRCIP)", "目的IP(DSTIP)", "进程文件(FILEPATH)", "数据内容(DATA)"}
	writer.Write(header)
	No := 1
	for index, _ := range *ptr {
		writer.Write([]string{strconv.Itoa(No), (*ptr)[index]["PID"], (*ptr)[index]["TID"],
			(*ptr)[index]["SRCIP"], (*ptr)[index]["DSTIP"], (*ptr)[index]["FilePath"], (*ptr)[index]["DATA"]})
		writer.Flush()
	}
	if err = writer.Error(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("结果文件路径：" + filepath)
}
